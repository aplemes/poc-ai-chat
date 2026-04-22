package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gin-quickstart/models"
	"gin-quickstart/services"

	"github.com/gin-gonic/gin"
)

type confirmRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
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

type analyzeRequest struct {
	FormData models.FormFillData `json:"formData"`
	Language string              `json:"language"`
}

func (h *ChatHandler) AnalyzeForm(c *gin.Context) {
	var req analyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	formJSON, err := json.MarshalIndent(req.FormData, "", "  ")
	if err != nil {
		writeEvent(c, services.ChatEvent{Type: "error", Content: "failed to serialize form data"})
		return
	}
	messages := []models.Message{
		{Role: models.RoleSystem, Content: services.BuildAnalysisSystemPrompt(req.Language)},
		{Role: models.RoleUser, Content: string(formJSON)},
	}

	if _, err := h.llama.StreamChat(c.Request.Context(), messages, nil, func(event services.ChatEvent) {
		writeEvent(c, event)
	}); err != nil {
		writeEvent(c, services.ChatEvent{Type: "error", Content: err.Error()})
		return
	}

	writeEvent(c, services.ChatEvent{Type: "done"})
}
