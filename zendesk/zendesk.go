package zendesk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	baseURLFormat = "https://%s.zendesk.com/api/v2"
)

var defaultHeaders = map[string]string{
	"User-Agent":   "nukosuke/go-zendesk",
	"Content-Type": "application/json",
}

var subdomainRegexp = regexp.MustCompile("^[a-z0-9][a-z0-9-]+[a-z0-9]$")

type (
	// Client of Zendesk API
	Client struct {
		baseURL    *url.URL
		httpClient *http.Client
		credential Credential
		headers    map[string]string
		maxSleep   time.Duration
		maxRetry   int
	}

	// BaseAPI encapsulates base methods for zendesk client
	BaseAPI interface {
		Get(ctx context.Context, path string) ([]byte, error)
		Post(ctx context.Context, path string, data interface{}) ([]byte, error)
		Put(ctx context.Context, path string, data interface{}) ([]byte, error)
		Delete(ctx context.Context, path string) error
	}

	// CursorPagination contains options for using cursor pagination.
	// Cursor pagination is preferred where possible.
	CursorPagination struct {
		// PageSize sets the number of results per page.
		// Most endpoints support up to 100 records per page.
		PageSize int `url:"page[size],omitempty"`

		// PageAfter provides the "next" cursor.
		PageAfter string `url:"page[after],omitempty"`

		// PageBefore provides the "previous" cursor.
		PageBefore string `url:"page[before],omitempty"`
	}

	// CursorPaginationMeta contains information concerning how to fetch
	// next and previous results, and if next results exist.
	CursorPaginationMeta struct {
		// HasMore is true if more results exist in the endpoint.
		HasMore bool `json:"has_more,omitempty"`

		// AfterCursor contains the cursor of the next result set.
		AfterCursor string `json:"after_cursor,omitempty"`

		// BeforeCursor contains the cursor of the previous result set.
		BeforeCursor string `json:"before_cursor,omitempty"`
	}
)

// NewClient creates new Zendesk API client
func NewClient(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{
		httpClient: httpClient,
		maxSleep:   5 * time.Second,
		maxRetry:   3,
	}
	client.headers = defaultHeaders
	return client, nil
}

// SetHeader saves HTTP header in client. It will be included all API request
func (z *Client) SetHeader(key string, value string) {
	z.headers[key] = value
}

// SetSubdomain saves subdomain in client. It will be used
// when call API
func (z *Client) SetSubdomain(subdomain string) error {
	if !subdomainRegexp.MatchString(subdomain) {
		return fmt.Errorf("%s is invalid subdomain", subdomain)
	}

	baseURLString := fmt.Sprintf(baseURLFormat, subdomain)
	baseURL, err := url.Parse(baseURLString)
	if err != nil {
		return err
	}

	z.baseURL = baseURL
	return nil
}

// SetEndpointURL replace full URL of endpoint without subdomain validation.
// This is mainly used for testing to point to mock API server.
func (z *Client) SetEndpointURL(newURL string) error {
	baseURL, err := url.Parse(newURL)
	if err != nil {
		return err
	}

	z.baseURL = baseURL
	return nil
}

// SetCredential saves credential in client. It will be set
// to request header when call API
func (z *Client) SetCredential(cred Credential) {
	z.credential = cred
}

// SetMaxRetrySleepDelay sets the maximum duration that a client will support sleeping
// if an API call returns a 429 error. Defaults to 5 seconds if not set.
func (z *Client) SetMaxRetrySleepDelay(duration time.Duration) {
	z.maxSleep = duration
}

// SetMaxRetry sets the maximum duration that a client will support sleeping
// if an API call returns a 429 error. Defaults to 3 if not set.
func (z *Client) SetMaxRetry(retries int) {
	if retries > 0 {
		z.maxRetry = retries
	}
}

// get fetches JSON data from API and returns its body as []bytes
func (z *Client) get(ctx context.Context, path string) ([]byte, error) {
	return z.execRequest(ctx, path, http.MethodGet, nil, []int{http.StatusOK})
}

// post send data to API and returns response body as []bytes
func (z *Client) post(ctx context.Context, path string, data interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return z.execRequest(ctx, path, http.MethodPost, bytes.NewReader(jsonBytes), []int{http.StatusOK, http.StatusCreated})
}

// put sends data to API and returns response body as []bytes
func (z *Client) put(ctx context.Context, path string, data interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return z.execRequest(ctx, path, http.MethodPut, bytes.NewReader(jsonBytes), []int{http.StatusOK, http.StatusNoContent})
}

// delete sends data to API and returns an error if unsuccessful
func (z *Client) delete(ctx context.Context, path string) error {
	_, err := z.execRequest(ctx, path, http.MethodDelete, nil, []int{http.StatusNoContent})
	return err
}

func (z *Client) execRequest(ctx context.Context, path string, verb string, reqBody io.Reader, successCodes []int) ([]byte, error) {
	var resp *http.Response
	var body []byte
	for attempts := 0; attempts < z.maxRetry; attempts++ {
		req, err := http.NewRequest(verb, z.baseURL.String()+path, reqBody)
		if err != nil {
			return nil, err
		}

		req = z.prepareRequest(ctx, req)
		resp, err = z.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		body, err = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests && attempts < z.maxRetry {
			retryStr := resp.Header.Get("Retry-After")
			retrySec, _ := strconv.Atoi(retryStr)
			if retrySec > 0 && time.Duration(retrySec) <= z.maxSleep {
				time.Sleep(time.Duration(retrySec) * time.Second)
				continue
			}
		}
		break
	}

	for _, code := range successCodes {
		if resp.StatusCode == code {
			return body, nil
		}
	}

	return nil, Error{
		body: body,
		resp: resp,
	}
}

// prepare request sets common request variables such as authn and user agent
func (z *Client) prepareRequest(ctx context.Context, req *http.Request) *http.Request {
	out := req.WithContext(ctx)
	z.includeHeaders(out)
	if z.credential != nil {
		if z.credential.Bearer() {
			out.Header.Add("Authorization", "Bearer "+z.credential.Secret())
		} else {
			out.SetBasicAuth(z.credential.Email(), z.credential.Secret())
		}
	}

	return out
}

// includeHeaders set HTTP headers from client.headers to *http.Request
func (z *Client) includeHeaders(req *http.Request) {
	for key, value := range z.headers {
		req.Header.Set(key, value)
	}
}

// addOptions build query string
func addOptions(s string, opts any) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// Get allows users to send requests not yet implemented
func (z *Client) Get(ctx context.Context, path string) ([]byte, error) {
	return z.get(ctx, path)
}

// Post allows users to send requests not yet implemented
func (z *Client) Post(ctx context.Context, path string, data interface{}) ([]byte, error) {
	return z.post(ctx, path, data)
}

// Put allows users to send requests not yet implemented
func (z *Client) Put(ctx context.Context, path string, data interface{}) ([]byte, error) {
	return z.put(ctx, path, data)
}

// Delete allows users to send requests not yet implemented
func (z *Client) Delete(ctx context.Context, path string) error {
	return z.delete(ctx, path)
}
