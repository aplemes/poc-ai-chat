package handlers

import (
	"encoding/json"
	"net/http"

	"gin-quickstart/models"
	"gin-quickstart/services"

	"github.com/gin-gonic/gin"
)

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

	setSSEHeaders(c, session.ID)

	msgs := h.conversations.GetMessages(session.ID)
	messages := buildMessagesWithSystem(msgs, services.BuildFieldSystemPrompt(req.FieldName, req.Language))
	tools := services.BuildFieldTool(req.FieldName)

	var assistantText string

	toolCall, err := h.llama.StreamChat(c.Request.Context(), messages, tools, func(event services.ChatEvent) {
		if event.Type == "token" {
			assistantText += event.Content
		}
		writeEvent(c, event)
	})

	if err != nil {
		writeEvent(c, services.ChatEvent{Type: "error", Content: err.Error()})
		return
	}

	if toolCall != nil && toolCall.Function.Name == "fill_field" {
		fieldValue, err := parseFieldValue(req.FieldName, toolCall.Function.Arguments)
		if err != nil {
			writeEvent(c, services.ChatEvent{Type: "error", Content: "failed to parse field value: " + err.Error()})
			return
		}
		addToolCallMessages(h, session.ID, toolCall, msgFieldFilled)
		writeEvent(c, services.ChatEvent{
			Type: "field_fill",
			Data: map[string]interface{}{"fieldName": req.FieldName, "value": fieldValue},
		})
	} else if assistantText != "" {
		h.conversations.AddMessage(session.ID, models.Message{
			Role:    models.RoleAssistant,
			Content: assistantText,
		})
	}

	writeEvent(c, services.ChatEvent{Type: "done", Content: session.ID})
}

// parseFieldValue deserialises the fill_field tool arguments.
// busInterested returns []string; all other fields return string.
func parseFieldValue(fieldName, arguments string) (interface{}, error) {
	var raw struct {
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal([]byte(arguments), &raw); err != nil {
		return nil, err
	}

	if fieldName == "busInterested" {
		var vals []string
		if err := json.Unmarshal(raw.Value, &vals); err != nil {
			return nil, err
		}
		return vals, nil
	}

	var s string
	if err := json.Unmarshal(raw.Value, &s); err != nil {
		return nil, err
	}
	return s, nil
}
