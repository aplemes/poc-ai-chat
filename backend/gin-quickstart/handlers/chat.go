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

type confirmRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
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
		// Store for later confirmation
		session.PendingFormData = &formData
		// Add assistant tool call + tool result to history so AI has context
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

func (h *ChatHandler) ConfirmForm(c *gin.Context) {
	var req confirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := h.conversations.GetByID(req.SessionID)
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	formData := session.PendingFormData
	if formData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no pending form data to confirm"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Session-ID", session.ID)

	// Add confirmation to history so AI knows the form was accepted
	h.conversations.AddMessage(session.ID, models.Message{
		Role:    models.RoleUser,
		Content: "[User confirmed the form proposal]",
	})

	argsBytes, _ := json.Marshal(formData)
	toolCallID := fmt.Sprintf("confirm_%s_%d", session.ID[:8], len(session.Messages))
	h.conversations.AddMessage(session.ID, models.Message{
		Role: models.RoleAssistant,
		ToolCalls: []models.ToolCall{{
			ID:   toolCallID,
			Type: "function",
			Function: models.FuncCall{
				Name:      "fill_demand_form",
				Arguments: string(argsBytes),
			},
		}},
	})
	h.conversations.AddMessage(session.ID, models.Message{
		Role:       models.RoleTool,
		ToolCallID: toolCallID,
		Content:    "Form filled successfully.",
	})

	session.PendingFormData = nil

	writeEvent(c, services.ChatEvent{Type: "form_fill", Data: *formData})
	writeEvent(c, services.ChatEvent{Type: "done", Content: session.ID})
}

type fieldMessageRequest struct {
	SessionID string `json:"sessionId"`
	FieldName string `json:"fieldName" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Language  string `json:"language"`
}

func (h *ChatHandler) SendFieldMessage(c *gin.Context) {
	var req fieldMessageRequest
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

	messages := buildFieldMessages(session, req.FieldName, req.Language)
	tools := services.BuildFieldTool(req.FieldName)

	var assistantText strings.Builder

	toolCall, err := h.llama.StreamChat(c.Request.Context(), messages, tools, func(event services.ChatEvent) {
		if event.Type == "token" {
			assistantText.WriteString(event.Content)
		}
		writeEvent(c, event)
	})

	if err != nil {
		writeEvent(c, services.ChatEvent{Type: "error", Content: err.Error()})
		return
	}

	if toolCall != nil && toolCall.Function.Name == "fill_field" {
		var raw struct {
			Value json.RawMessage `json:"value"`
		}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &raw); err != nil {
			writeEvent(c, services.ChatEvent{Type: "error", Content: "failed to parse field value: " + err.Error()})
			return
		}
		h.conversations.AddMessage(session.ID, models.Message{
			Role:      models.RoleAssistant,
			ToolCalls: []models.ToolCall{*toolCall},
		})
		h.conversations.AddMessage(session.ID, models.Message{
			Role:       models.RoleTool,
			ToolCallID: toolCall.ID,
			Content:    "Field filled successfully.",
		})
		var fieldValue interface{}
		if req.FieldName == "busInterested" {
			var vals []string
			_ = json.Unmarshal(raw.Value, &vals)
			fieldValue = vals
		} else {
			var s string
			_ = json.Unmarshal(raw.Value, &s)
			fieldValue = s
		}
		writeEvent(c, services.ChatEvent{
			Type: "field_fill",
			Data: map[string]interface{}{"fieldName": req.FieldName, "value": fieldValue},
		})
	} else if assistantText.Len() > 0 {
		h.conversations.AddMessage(session.ID, models.Message{
			Role:    models.RoleAssistant,
			Content: assistantText.String(),
		})
	}

	writeEvent(c, services.ChatEvent{Type: "done", Content: session.ID})
}

func buildFieldMessages(session *models.Session, fieldName, language string) []models.Message {
	messages := []models.Message{
		{Role: models.RoleSystem, Content: services.BuildFieldSystemPrompt(fieldName, language)},
	}
	return append(messages, session.Messages...)
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
