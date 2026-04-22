package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gin-quickstart/models"
	"gin-quickstart/services"

	"github.com/gin-gonic/gin"
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

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Session-ID", session.ID)

	messages := buildMessages(session, req.Language)

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
		session.PendingFormData = &formData
		h.conversations.AddMessage(session.ID, models.Message{
			Role:      models.RoleAssistant,
			ToolCalls: []models.ToolCall{*toolCall},
		})
		h.conversations.AddMessage(session.ID, models.Message{
			Role:       models.RoleTool,
			ToolCallID: toolCall.ID,
			Content:    "Proposal shown to user. Waiting for confirmation.",
		})
		writeEvent(c, services.ChatEvent{Type: "form_confirm", Data: formData})

	case toolCall != nil && toolCall.Function.Name == "fill_demand_form":
		var formData models.FormFillData
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &formData); err != nil {
			writeEvent(c, services.ChatEvent{Type: "error", Content: "failed to parse form data: " + err.Error()})
			return
		}
		session.PendingFormData = nil
		h.conversations.AddMessage(session.ID, models.Message{
			Role:      models.RoleAssistant,
			ToolCalls: []models.ToolCall{*toolCall},
		})
		h.conversations.AddMessage(session.ID, models.Message{
			Role:       models.RoleTool,
			ToolCallID: toolCall.ID,
			Content:    "Form filled successfully.",
		})
		writeEvent(c, services.ChatEvent{Type: "form_fill", Data: formData})

	case assistantText.Len() > 0:
		h.conversations.AddMessage(session.ID, models.Message{
			Role:    models.RoleAssistant,
			Content: assistantText.String(),
		})
	}

	writeEvent(c, services.ChatEvent{Type: "done", Content: session.ID})
}

func buildMessages(session *models.Session, language string) []models.Message {
	messages := []models.Message{
		{Role: models.RoleSystem, Content: services.BuildSystemPrompt(language)},
	}
	return append(messages, session.Messages...)
}

func writeEvent(c *gin.Context, event services.ChatEvent) {
	data, _ := json.Marshal(event)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	c.Writer.Flush()
}
