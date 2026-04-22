package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gin-quickstart/models"
	"gin-quickstart/services"

	"github.com/gin-gonic/gin"
)

// LLM history message constants (BE-M4).
const (
	msgProposalShown    = "Proposal shown to user. Waiting for confirmation."
	msgFormFilled       = "Form filled successfully."
	msgFieldFilled      = "Field filled successfully."
	msgUserConfirmed    = "[User confirmed the form proposal]"
)

type ChatHandler struct {
	conversations *services.ConversationService
	llama         *services.LlamaService
}

func NewChatHandler(apiKey string) *ChatHandler {
	return &ChatHandler{
		conversations: services.NewConversationService(),
		llama:         services.NewLlamaService(apiKey),
	}
}

type messageRequest struct {
	SessionID string `json:"sessionId"`
	Message   string `json:"message" binding:"required"`
	Language  string `json:"language"`
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req messageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := h.conversations.GetOrCreate(req.SessionID)

	h.conversations.AddMessage(session.ID, models.Message{
		Role:    models.RoleUser,
		Content: req.Message,
	})

	setSSEHeaders(c, session.ID)

	msgs := h.conversations.GetMessages(session.ID)
	messages := buildMessagesWithSystem(msgs, services.BuildSystemPrompt(req.Language))

	var assistantText strings.Builder

	toolCall, err := h.llama.StreamChat(c.Request.Context(), messages, services.BuildTools(), func(event services.ChatEvent) {
		if event.Type == "token" {
			assistantText.WriteString(event.Content)
		}
		writeEvent(c, event)
	})

	if err != nil {
		writeEvent(c, services.ChatEvent{Type: "error", Content: err.Error()})
		return
	}

	switch {
	case toolCall != nil && toolCall.Function.Name == "propose_form_data":
		var formData models.FormFillData
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &formData); err != nil {
			writeEvent(c, services.ChatEvent{Type: "error", Content: "failed to parse proposal data: " + err.Error()})
			return
		}
		h.conversations.SetPendingFormData(session.ID, &formData)
		addToolCallMessages(h, session.ID, toolCall, msgProposalShown)
		writeEvent(c, services.ChatEvent{Type: "form_confirm", Data: formData})

	case toolCall != nil && toolCall.Function.Name == "fill_demand_form":
		var formData models.FormFillData
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &formData); err != nil {
			writeEvent(c, services.ChatEvent{Type: "error", Content: "failed to parse form data: " + err.Error()})
			return
		}
		h.conversations.SetPendingFormData(session.ID, nil)
		addToolCallMessages(h, session.ID, toolCall, msgFormFilled)
		writeEvent(c, services.ChatEvent{Type: "form_fill", Data: formData})

	case assistantText.Len() > 0:
		h.conversations.AddMessage(session.ID, models.Message{
			Role:    models.RoleAssistant,
			Content: assistantText.String(),
		})
	}

	writeEvent(c, services.ChatEvent{Type: "done", Content: session.ID})
}

// buildMessagesWithSystem prepends the system prompt to the session messages (BE-M2).
func buildMessagesWithSystem(sessionMsgs []models.Message, systemPrompt string) []models.Message {
	messages := make([]models.Message, 0, len(sessionMsgs)+1)
	messages = append(messages, models.Message{Role: models.RoleSystem, Content: systemPrompt})
	return append(messages, sessionMsgs...)
}

// addToolCallMessages records the assistant tool call and its result in the conversation (BE-M3).
func addToolCallMessages(h *ChatHandler, sessionID string, toolCall *models.ToolCall, result string) {
	h.conversations.AddMessage(sessionID, models.Message{
		Role:      models.RoleAssistant,
		ToolCalls: []models.ToolCall{*toolCall},
	})
	h.conversations.AddMessage(sessionID, models.Message{
		Role:       models.RoleTool,
		ToolCallID: toolCall.ID,
		Content:    result,
	})
}

// setSSEHeaders writes the standard SSE response headers (BE-06).
func setSSEHeaders(c *gin.Context, sessionID string) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	if sessionID != "" {
		c.Header("X-Session-ID", sessionID)
	}
}

// writeEvent serialises and flushes a single SSE data event (BE-03).
func writeEvent(c *gin.Context, event services.ChatEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("ERROR: failed to marshal SSE event %q: %v", event.Type, err)
		return
	}
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	c.Writer.Flush()
}
