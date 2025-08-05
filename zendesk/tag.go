package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tag is an alias for string
type Tag string

// TagAPI an interface containing all tag related methods
type TagAPI interface {
	GetTicketTags(ctx context.Context, ticketID int64) ([]Tag, error)
	GetOrganizationTags(ctx context.Context, organizationID int64) ([]Tag, error)
	GetUserTags(ctx context.Context, userID int64) ([]Tag, error)
	AddTicketTags(ctx context.Context, ticketID int64, tags []Tag) ([]Tag, error)
	AddOrganizationTags(ctx context.Context, organizationID int64, tags []Tag) ([]Tag, error)
	AddUserTags(ctx context.Context, userID int64, tags []Tag) ([]Tag, error)
	ListTags(ctx context.Context, options TagListOptions) ([]Tag, Page, error)
	SearchTags(ctx context.Context, options SearchTagsOptions) ([]Tag, Page, error)
}

// GetTicketTags get ticket tag list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#show-tags
func (z *Client) GetTicketTags(ctx context.Context, ticketID int64) ([]Tag, error) {
	var result struct {
		Tags []Tag `json:"tags"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/tickets/%d/tags.json", ticketID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tags, err
}

type TagListOptions struct {
	PageOptions
}

// ListTags gets list of available tags
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#list-tags
func (z *Client) ListTags(ctx context.Context, options TagListOptions) ([]Tag, Page, error) {
	var result struct {
		Page
		Tags []struct {
			Name  string `json:"name"`
			Count int64  `json:"count"`
		} `json:"tags"`
	}

	u, err := addOptions("/tags.json", options)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, Page{}, err
	}

	var tags []Tag
	for _, t := range result.Tags {
		tags = append(tags, Tag(t.Name))
	}

	return tags, result.Page, nil
}

type SearchTagsOptions struct {
	PageOptions
	Name string `url:"name,omitempty"`
}

// SearchTags searches tags by name
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#search-tags
func (z *Client) SearchTags(ctx context.Context, options SearchTagsOptions) ([]Tag, Page, error) {
	var result struct {
		Page
		Tags []Tag `json:"tags"`
	}
	u, err := addOptions("/autocomplete/tags.json", options)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, Page{}, err
	}

	return result.Tags, result.Page, nil
}

// GetOrganizationTags get organization tag list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#show-tags
func (z *Client) GetOrganizationTags(ctx context.Context, organizationID int64) ([]Tag, error) {
	var result struct {
		Tags []Tag `json:"tags"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/organizations/%d/tags.json", organizationID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tags, err
}

// GetUserTags get user tag list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#show-tags
func (z *Client) GetUserTags(ctx context.Context, userID int64) ([]Tag, error) {
	var result struct {
		Tags []Tag `json:"tags"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/users/%d/tags.json", userID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tags, err
}

// AddTicketTags add tags to ticket
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#add-tags
func (z *Client) AddTicketTags(ctx context.Context, ticketID int64, tags []Tag) ([]Tag, error) {
	var data, result struct {
		Tags []Tag `json:"tags"`
	}
	data.Tags = tags

	body, err := z.put(ctx, fmt.Sprintf("/tickets/%d/tags", ticketID), data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}

// AddOrganizationTags add tags to organization
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#add-tags
func (z *Client) AddOrganizationTags(ctx context.Context, organizationID int64, tags []Tag) ([]Tag, error) {
	var data, result struct {
		Tags []Tag `json:"tags"`
	}
	data.Tags = tags

	body, err := z.put(ctx, fmt.Sprintf("/organizations/%d/tags", organizationID), data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}

// AddUserTags add tags to user
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#add-tags
func (z *Client) AddUserTags(ctx context.Context, userID int64, tags []Tag) ([]Tag, error) {
	var data, result struct {
		Tags []Tag `json:"tags"`
	}
	data.Tags = tags

	body, err := z.put(ctx, fmt.Sprintf("/users/%d/tags", userID), data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}
