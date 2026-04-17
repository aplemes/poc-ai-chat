package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gin-quickstart/models"
)

const groqAPIURL = "https://api.groq.com/openai/v1/chat/completions"
const llamaModel = "llama-3.3-70b-versatile"

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
			Content   string              `json:"content"`
			ToolCalls []streamToolCall    `json:"tool_calls"`
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
func (s *LlamaService) StreamChat(
	ctx context.Context,
	messages []models.Message,
	emit func(ChatEvent),
) (*models.ToolCall, error) {
	body, err := json.Marshal(map[string]interface{}{
		"model":       llamaModel,
		"messages":    messages,
		"stream":      true,
		"tools":       buildTools(),
		"tool_choice": "auto",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", groqAPIURL, bytes.NewReader(body))
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("groq API error %d: %s", resp.StatusCode, body)
	}

	scanner := bufio.NewScanner(resp.Body)
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
			Type: "function",
			Function: models.FuncCall{
				Name:      acc.Name,
				Arguments: acc.Args,
			},
		}, nil
	}

	return nil, nil
}

func buildTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "fill_demand_form",
				"description": "Preenche o formulário de demanda com as informações coletadas do usuário. Chame esta função apenas quando tiver todos os campos obrigatórios.",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Título da demanda: começa com verbo no infinitivo + escopo. Ex: 'Add the new Payment Method XXX on the website only'",
						},
						"businessLine": map[string]interface{}{
							"type":        "string",
							"description": "ID da linha de negócio/organização",
							"enum":        []string{"18518", "18519", "18520", "18521", "18522", "18523", "18524", "19033", "18525"},
						},
						"requesterBU": map[string]interface{}{
							"type":        "string",
							"description": "ID da BU solicitante (formato ADEO-XXXX)",
							"enum": []string{
								"ADEO-8052", "ADEO-35430", "ADEO-35424", "ADEO-8078",
								"ADEO-35426", "ADEO-35427", "ADEO-8062", "ADEO-8087",
								"ADEO-35431", "ADEO-36214", "ADEO-8089", "ADEO-8054",
								"ADEO-8095", "ADEO-8064", "ADEO-8075", "ADEO-8100",
								"ADEO-8067", "ADEO-8057", "ADEO-40345", "ADEO-8091",
								"ADEO-8061", "ADEO-8084", "ADEO-8071", "ADEO-8051",
								"ADEO-8053", "ADEO-8092", "ADEO-8070", "ADEO-23566",
								"ADEO-8074", "ADEO-8055", "ADEO-8056", "ADEO-8060",
							},
						},
						"busInterested": map[string]interface{}{
							"type":        "string",
							"description": "ID da BU interessada/impactada",
							"enum": []string{
								"20047", "20048", "20049", "20050", "20051", "20052",
								"20053", "20054", "20055", "20056", "20057", "20058",
								"20059", "20060", "20061", "20062", "20063", "20064",
								"20065", "20066", "20067", "20068", "20069", "20070",
								"20071", "20072", "20073", "20074", "20075", "20076",
								"20077", "20078",
							},
						},
						"timeSensitive": map[string]interface{}{
							"type":        "string",
							"description": "Se a demanda é urgente por questão legal, de segurança ou não há urgência",
							"enum":        []string{"No", "Legal", "Security"},
						},
						"whyDemand": map[string]interface{}{
							"type":        "string",
							"description": "Descreva a situação atual, os pontos de dor, comparação com concorrentes e o contexto que motiva a demanda",
						},
						"whoIsImpacted": map[string]interface{}{
							"type":        "string",
							"description": "Quem é impactado: tipos de usuários/personas e estimativa de quantidade. Ex: 'Customers online: 5% of total = 50,000. Coworkers: not concerned.'",
						},
						"benefitCategory": map[string]interface{}{
							"type":        "string",
							"description": "Categoria do benefício esperado da demanda",
							"enum":        []string{"Cost efficiency", "Environmental & social sustainability", "Service quality & security risk", "Customer satisfaction & revenue", "Innovation", "Other"},
						},
						"benefitHypothesis": map[string]interface{}{
							"type":        "string",
							"description": "Hipóteses do usuário para atingir os benefícios esperados",
						},
						"measureBenefits": map[string]interface{}{
							"type":        "string",
							"description": "KPIs e timing para medir os benefícios. Ex: 'GMV per payment method in first 3 months after activation.'",
						},
					},
					"required": []string{"title", "businessLine", "requesterBU", "busInterested", "timeSensitive", "whyDemand", "whoIsImpacted", "benefitCategory", "benefitHypothesis", "measureBenefits"},
				},
			},
		},
	}
}

var languageNames = map[string]string{
	"pt": "Portuguese",
	"en": "English",
	"es": "Spanish",
	"fr": "French",
}

func BuildSystemPrompt(language string) string {
	lang := languageNames[language]
	if lang == "" {
		lang = "English"
	}
	return "Converse with the user exclusively in " + lang + ". CRITICAL: All form field values passed to fill_demand_form must ALWAYS be written in English, regardless of the conversation language.\n\n" + systemPromptBase
}

const systemPromptBase = `You are an expert assistant helping Adeo collaborators fill in a Demand form. You collect all required information before submitting.

## Form fields and how to fill them

**title** (required)
Start with an infinitive verb and complete with the scope of application.
Example: "Add the new Payment Method 'XXX' on the website only"

**businessLine** (required) — use only the IDs below:
  18518 → Omnicommerce Experience
  18519 → Services & Renovation
  18520 → Supply Chain & Delivery
  18521 → Offer & Industry
  18522 → Finance
  18523 → Positive Impacts
  18524 → Human & Sharing
  19033 → Executive Succession Plan
  18525 → Digital Data Tech

**requesterBU** (required) — the BU making the request, use only the IDs below:
  ADEO-8052 → Adeo Marketplace Services  | ADEO-35430 → Adeo Production
  ADEO-35424 → Adeo Services Chine       | ADEO-8078 → Adeo Services France
  ADEO-35426 → Adeo Services Poland      | ADEO-35427 → Adeo Services Vietnam
  ADEO-8062 → Bricocenter Italy           | ADEO-8087 → Bricoman Poland
  ADEO-35431 → Enki Home                 | ADEO-36214 → GO XL
  ADEO-8089 → Golilla                    | ADEO-8054 → Kbane France
  ADEO-8095 → Leroy Merlin Brazil        | ADEO-8064 → Leroy Merlin France
  ADEO-8075 → Leroy Merlin Greece/Cyprus | ADEO-8100 → Leroy Merlin Italy
  ADEO-8067 → Leroy Merlin Poland        | ADEO-8057 → Leroy Merlin Portugal
  ADEO-40345 → Leroy Merlin Renovation   | ADEO-8091 → Leroy Merlin Romania
  ADEO-8061 → Leroy Merlin South Africa  | ADEO-8084 → Leroy Merlin Spain
  ADEO-8071 → Leroy Merlin Ukraine       | ADEO-8051 → Obramat Portugal
  ADEO-8053 → Obramat Spain              | ADEO-8092 → Obramax Brazil
  ADEO-8070 → Quotatis                   | ADEO-23566 → Saint Maclou France
  ADEO-8074 → Tecnomat France            | ADEO-8055 → Tecnomat Italy
  ADEO-8056 → Terra Incognita            | ADEO-8060 → Weldom France

**busInterested** (required) — one BU already aligned on this demand.
Use numeric IDs 20047–20078 in the same order as the BU list above (20047 = Adeo Marketplace Services, 20078 = Weldom France).

**timeSensitive** (required) — "No", "Legal", or "Security". Ask if the demand has urgency due to legal or security reasons; if not, use "No".

**whyDemand** (required)
A single comprehensive description covering: the context/event that motivates the demand, the current situation and tools in use, and the pain points or inefficiencies to solve.

**whoIsImpacted** (required)
The different types of users (personas) impacted and estimated numbers.
Example: "Customers online: 5% of total = 50,000. Coworkers: not concerned. Partners: not concerned."

**benefitCategory** (required) — one of:
  Cost efficiency | Environmental & social sustainability | Service quality & security risk
  Customer satisfaction & revenue | Innovation | Other

**benefitHypothesis** (required)
The user's hypothesis on how the demand will achieve the expected benefits.

**measureBenefits** (required)
KPIs and timing to verify the benefits.
Example: "GMV per payment method during the first 3 months after activation."

## Instructions

1. Ask exactly ONE question per message — no preamble, no extra commentary
2. Infer what you can from context; only ask what you truly cannot infer
3. You MUST collect ALL fields before calling fill_demand_form — every single field is required
4. Never invent IDs — use only the values listed above
5. Only call fill_demand_form when ALL fields have been answered

## Formatting rules for option fields (select / enum)

When asking about a field that has a predefined set of options (businessLine, requesterBU, busInterested, benefitCategory):
- Start with one short sentence explaining what the field means and why it matters
- Then present **only the labels** (never the IDs) as a numbered or bulleted list, one option per line
- End with a short question asking the user to choose
- Keep the list clean — no IDs, no extra descriptions beside each item

Example format:
"A **Linha de Negócio** identifica a área organizacional responsável pela demanda. Escolha uma das opções abaixo:

1. Omnicommerce Experience
2. Services & Renovation
3. Supply Chain & Delivery
...

Qual delas representa melhor a sua área?"`
