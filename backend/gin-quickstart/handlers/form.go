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

	// TakeAndClearPendingFormData checks session existence and atomically clears the pending data (BE-01/BE-05).
	formData, sessionExists := h.conversations.TakeAndClearPendingFormData(req.SessionID)
	if !sessionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	if formData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no pending form data to confirm"})
		return
	}

	setSSEHeaders(c, req.SessionID)

	h.conversations.AddMessage(req.SessionID, models.Message{
		Role:    models.RoleUser,
		Content: msgUserConfirmed,
	})

	argsBytes, err := json.Marshal(formData)
	if err != nil {
		writeEvent(c, services.ChatEvent{Type: "error", Content: "failed to serialize form data"})
		return
	}
	toolCallID := fmt.Sprintf("confirm_%s_%d", req.SessionID[:8], len(h.conversations.GetMessages(req.SessionID)))
	syntheticToolCall := &models.ToolCall{
		ID:   toolCallID,
		Type: models.ToolCallTypeFunction,
		Function: models.FuncCall{
			Name:      "fill_demand_form",
			Arguments: string(argsBytes),
		},
	}
	addToolCallMessages(h, req.SessionID, syntheticToolCall, msgFormFilled)

	writeEvent(c, services.ChatEvent{Type: "form_fill", Data: *formData})
	writeEvent(c, services.ChatEvent{Type: "done", Content: req.SessionID})
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

	setSSEHeaders(c, "")

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
