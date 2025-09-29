package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
)

type RedactChatCommentRequest struct {
	TicketID  int64  `json:"id"`
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
}

type RedactChatCommentAttachmentRequest struct {
	TicketID   int64    `json:"id"`
	MessageIDs []string `json:"message_ids"`
}

type RedactChatCommentResponse struct {
	ChatEvent struct {
		ID    int64  `json:"id"`
		Type  string `json:"type"`
		Value struct {
			ChatID  string `json:"chat_id"`
			History []struct {
				ChatIndex int    `json:"chat_index"`
				Message   string `json:"message"`
				Type      string `json:"type"`
				Filename  string `json:"filename,omitempty"`
			} `json:"history"`
			VisitorID string `json:"visitor_id"`
		} `json:"value"`
	} `json:"chat_event"`
}

// RedactChatComment redacts a chat comment in a ticket.
//
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#redact-chat-comment
func (z *Client) RedactChatComment(ctx context.Context, req *RedactChatCommentRequest) (*RedactChatCommentResponse, error) {
	result := &RedactChatCommentResponse{}
	resp, err := z.put(ctx, fmt.Sprintf("/chat_redactions/%d.json", req.TicketID), req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RedactChatCommentAttachment redacts an attachments from chat comments in a ticket.
// NOTE: for chats, zendesk allows only one attachment per message
//
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#redact-chat-comment-attachment
func (z *Client) RedactChatCommentAttachment(ctx context.Context, req *RedactChatCommentAttachmentRequest) (*RedactChatCommentResponse, error) {
	result := &RedactChatCommentResponse{}
	resp, err := z.put(ctx, fmt.Sprintf("/chat_file_redactions/%d.json", req.TicketID), req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
