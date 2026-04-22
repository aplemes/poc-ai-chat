package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gin-quickstart/models"
)

const groqAPIURL = "https://api.groq.com/openai/v1/chat/completions"
const llamaModel = "llama-3.3-70b-versatile"

// sseMaxTokenBytes is the scanner buffer size for reading SSE lines (BE-M6).
// 1 MiB is generous enough to handle any realistic LLM response chunk.
const sseMaxTokenBytes = 1 << 20

type LlamaService struct {
	apiKey string
	client *http.Client
}

func NewLlamaService(apiKey string) *LlamaService {
	return &LlamaService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 90 * time.Second},
	}
}

type ChatEvent struct {
	Type    string      `json:"type"`
	Content string      `json:"content,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content   string           `json:"content"`
			ToolCalls []streamToolCall `json:"tool_calls"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type streamToolCall struct {
	Index    int    `json:"index"`
	ID       string `json:"id"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type toolCallAccum struct {
	ID   string
	Name string
	Args string
}

// StreamChat streams tokens to emit func and returns a ToolCall if the model invoked one.
// Pass nil tools to skip tool calling (e.g. for the analysis endpoint).
func (s *LlamaService) StreamChat(
	ctx context.Context,
	messages []models.Message,
	tools []map[string]interface{},
	emit func(ChatEvent),
) (*models.ToolCall, error) {
	requestBody := map[string]interface{}{
		"model":    llamaModel,
		"messages": messages,
		"stream":   true,
	}
	if len(tools) > 0 {
		requestBody["tools"] = tools
		requestBody["tool_choice"] = "auto"
	}
	reqBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", groqAPIURL, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("groq API error %d: %s", resp.StatusCode, errBody)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, sseMaxTokenBytes), sseMaxTokenBytes) // BE-M6
	toolCalls := map[int]*toolCallAccum{}

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			log.Printf("WARN: skipping malformed SSE chunk: %v", err) // BE-M7
			continue
		}
		if len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta

		if delta.Content != "" {
			emit(ChatEvent{Type: "token", Content: delta.Content})
		}

		for _, tc := range delta.ToolCalls {
			if _, ok := toolCalls[tc.Index]; !ok {
				toolCalls[tc.Index] = &toolCallAccum{}
			}
			acc := toolCalls[tc.Index]
			if tc.ID != "" {
				acc.ID = tc.ID
			}
			if tc.Function.Name != "" {
				acc.Name = tc.Function.Name
			}
			acc.Args += tc.Function.Arguments
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if acc, ok := toolCalls[0]; ok && acc.Name != "" {
		return &models.ToolCall{
			ID:   acc.ID,
			Type: models.ToolCallTypeFunction,
			Function: models.FuncCall{
				Name:      acc.Name,
				Arguments: acc.Args,
			},
		}, nil
	}

	return nil, nil
}
