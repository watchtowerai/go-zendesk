package zendesk

import (
	"net/http"
	"testing"
)

func TestRedactChatComment(t *testing.T) {
	mockAPI := newMockAPI(http.MethodPut, "redact_chat_comment.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	req := &RedactChatCommentAttachmentRequest{
		TicketID:   123,
		MessageIDs: []string{"abcd-1234-efgh-5678"},
	}

	resp, err := client.RedactChatCommentAttachment(ctx, req)
	if err != nil {
		t.Fatalf("Failed to redact chat comment attachment: %s", err)
	}

	if resp.ChatEvent.ID != 1932802680168 {
		t.Fatalf("Expected chat event ID to be 1932802680168, got %d", resp.ChatEvent.ID)
	}
}

func TestRedactChatCommentAttachment(t *testing.T) {
	mockAPI := newMockAPI(http.MethodPut, "redact_chat_comment_attachment.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	req := &RedactChatCommentAttachmentRequest{
		TicketID:   123,
		MessageIDs: []string{"abcd-1234-efgh-5678"},
	}

	resp, err := client.RedactChatCommentAttachment(ctx, req)
	if err != nil {
		t.Fatalf("Failed to redact chat comment attachment: %s", err)
	}

	if resp.ChatEvent.ID != 1932802680168 {
		t.Fatalf("Expected chat event ID to be 1932802680168, got %d", resp.ChatEvent.ID)
	}
}
