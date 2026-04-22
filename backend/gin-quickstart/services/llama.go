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
	tools []map[string]interface{},
	emit func(ChatEvent),
) (*models.ToolCall, error) {
	body, err := json.Marshal(map[string]interface{}{
		"model":       llamaModel,
		"messages":    messages,
		"stream":      true,
		"tools":       tools,
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

func formProperties() map[string]interface{} {
	return map[string]interface{}{
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
			"type":        "array",
			"description": "IDs das BUs já alinhadas/interessadas nesta demanda",
			"items": map[string]interface{}{
				"type": "string",
				"enum": []string{
					"20047", "20048", "20049", "20050", "20051", "20052",
					"20053", "20054", "20055", "20056", "20057", "20058",
					"20059", "20060", "20061", "20062", "20063", "20064",
					"20065", "20066", "20067", "20068", "20069", "20070",
					"20071", "20072", "20073", "20074", "20075", "20076",
					"20077", "20078",
				},
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
		"lowConfidenceFields": map[string]interface{}{
			"type":        "array",
			"description": "Field names where you inferred the value from context rather than explicit user input. Include only fields you are genuinely unsure about.",
			"items":       map[string]interface{}{"type": "string"},
		},
	}
}

var formRequiredFields = []string{
	"title", "businessLine", "requesterBU", "busInterested", "timeSensitive",
	"whyDemand", "whoIsImpacted", "benefitCategory", "benefitHypothesis", "measureBenefits",
}

func BuildTools() []map[string]interface{} {
	props := formProperties()
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "propose_form_data",
				"description": "Propõe ao usuário um resumo de todos os campos coletados para revisão e confirmação ANTES de preencher o formulário. Chame esta função quando tiver todos os campos obrigatórios prontos.",
				"parameters": map[string]interface{}{
					"type":       "object",
					"properties": props,
					"required":   formRequiredFields,
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "fill_demand_form",
				"description": "Preenche efetivamente o formulário. Use APENAS após o usuário confirmar a proposta, ou para aplicar correções solicitadas após uma confirmação anterior.",
				"parameters": map[string]interface{}{
					"type":       "object",
					"properties": props,
					"required":   formRequiredFields,
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
3. Never invent IDs — use only the values listed above

## Two-phase submission flow

**Phase 1 — Collect:** Ask questions until you have ALL required fields.

**Phase 2 — Propose:** When you have all fields, call "propose_form_data" (NOT "fill_demand_form"). This shows the user a summary to review.
- If the user confirms (says yes, looks good, confirm, etc.) → call "fill_demand_form" with all fields
- If the user requests changes BEFORE confirming → apply corrections and call "propose_form_data" again with the updated values

**Phase 3 — Corrections:** After "fill_demand_form" has been called and the form is filled, if the user asks for further changes, call "fill_demand_form" directly with the corrected values. No need to propose again.

## Uncertainty marking

When you infer a field value from context (not directly stated by the user), add that field name to the "lowConfidenceFields" array. This flags it for the user to double-check.

Examples of uncertain fields:
- You inferred businessLine from the project type described
- You estimated whoIsImpacted from usage context
- You guessed requesterBU from country or team name

Do NOT add a field to "lowConfidenceFields" if the user explicitly stated the value.

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

// BuildFieldTool returns the fill_field tool scoped to a single form field.
func BuildFieldTool(fieldName string) []map[string]interface{} {
	props := formProperties()
	fieldProp, ok := props[fieldName]
	if !ok {
		fieldProp = map[string]interface{}{"type": "string", "description": "Value for the field"}
	}
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "fill_field",
				"description": "Fill the field with the collected value. Call this when you have enough context to provide a confident answer.",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"value": fieldProp,
					},
					"required": []string{"value"},
				},
			},
		},
	}
}

// BuildFieldSystemPrompt returns a focused system prompt for assisting with a single form field.
func BuildFieldSystemPrompt(fieldName, language string) string {
	lang := languageNames[language]
	if lang == "" {
		lang = "English"
	}
	prompt, ok := fieldSystemPrompts[fieldName]
	if !ok {
		prompt = "Help the user fill in this field with appropriate content. When you have the answer, call fill_field."
	}
	return "Converse with the user in " + lang + ". Field values must always be written in English.\n\nYou are a focused assistant for ONE specific field of a Demand form. Do NOT collect other fields.\nAsk at most ONE clarifying question. When you have enough context, call fill_field immediately.\n\n" + prompt
}

var fieldSystemPrompts = map[string]string{
	"title": `You are helping fill the **Title** field.

Format: infinitive verb + what + where/scope.
Examples:
- "Add the new Payment Method X on the website"
- "Improve checkout performance on mobile app"
- "Enable loyalty points program for B2B customers"

Rules: start with an infinitive verb, be specific about scope, keep it concise.
Ask what the user wants to achieve and for whom, then call fill_field with a well-formed title.`,

	"businessLine": `You are helping fill the **Business Line** field.

This identifies the organisational area responsible for the demand. Valid options:
1. Omnicommerce Experience (ID: 18518)
2. Services & Renovation (ID: 18519)
3. Supply Chain & Delivery (ID: 18520)
4. Offer & Industry (ID: 18521)
5. Finance (ID: 18522)
6. Positive Impacts (ID: 18523)
7. Human & Sharing (ID: 18524)
8. Executive Succession Plan (ID: 19033)
9. Digital Data Tech (ID: 18525)

Ask which area best represents the user's team, then call fill_field with the matching ID.`,

	"requesterBU": `You are helping fill the **Requester BU** field.

This is the Business Unit making the request. Valid options:
Adeo Marketplace Services (ADEO-8052), Adeo Production (ADEO-35430), Adeo Services Chine (ADEO-35424),
Adeo Services France (ADEO-8078), Adeo Services Poland (ADEO-35426), Adeo Services Vietnam (ADEO-35427),
Bricocenter Italy (ADEO-8062), Bricoman Poland (ADEO-8087), Enki Home (ADEO-35431), GO XL (ADEO-36214),
Golilla (ADEO-8089), Kbane France (ADEO-8054), Leroy Merlin Brazil (ADEO-8095),
Leroy Merlin France (ADEO-8064), Leroy Merlin Greece/Cyprus (ADEO-8075), Leroy Merlin Italy (ADEO-8100),
Leroy Merlin Poland (ADEO-8067), Leroy Merlin Portugal (ADEO-8057), Leroy Merlin Renovation (ADEO-40345),
Leroy Merlin Romania (ADEO-8091), Leroy Merlin South Africa (ADEO-8061), Leroy Merlin Spain (ADEO-8084),
Leroy Merlin Ukraine (ADEO-8071), Obramat Portugal (ADEO-8051), Obramat Spain (ADEO-8053),
Obramax Brazil (ADEO-8092), Quotatis (ADEO-8070), Saint Maclou France (ADEO-23566),
Tecnomat France (ADEO-8074), Tecnomat Italy (ADEO-8055), Terra Incognita (ADEO-8056), Weldom France (ADEO-8060).

Ask which BU the user belongs to, then call fill_field with the matching ADEO-XXXX ID.`,

	"busInterested": `You are helping fill the **BUs Interested** field.

This is a Business Unit already aligned with and interested in this demand. Valid options (use the numeric ID):
Adeo Marketplace Services (20047), Adeo Production (20048), Adeo Services Chine (20049),
Adeo Services France (20050), Adeo Services Poland (20051), Adeo Services Vietnam (20052),
Bricocenter Italy (20053), Bricoman Poland (20054), Enki Home (20055), GO XL (20056),
Golilla (20057), Kbane France (20058), Leroy Merlin Brazil (20059), Leroy Merlin France (20060),
Leroy Merlin Greece/Cyprus (20061), Leroy Merlin Italy (20062), Leroy Merlin Poland (20063),
Leroy Merlin Portugal (20064), Leroy Merlin Renovation (20065), Leroy Merlin Romania (20066),
Leroy Merlin South Africa (20067), Leroy Merlin Spain (20068), Leroy Merlin Ukraine (20069),
Obramat Portugal (20070), Obramat Spain (20071), Obramax Brazil (20072), Quotatis (20073),
Saint Maclou France (20074), Tecnomat France (20075), Tecnomat Italy (20076),
Terra Incognita (20077), Weldom France (20078).

Ask which BU has already shown interest, then call fill_field with the numeric ID.`,

	"timeSensitive": `You are helping fill the **Time Sensitive** field.

Options:
- "No" — no urgency driven by compliance or security
- "Legal" — there is a legal or regulatory deadline
- "Security" — there is a security risk requiring urgent action

Ask if the demand has a hard deadline due to legal or security reasons. Call fill_field with "No", "Legal", or "Security".`,

	"whyDemand": `You are helping fill the **Why Demand** field.

This should cover in one comprehensive text:
1. The context or event that motivates the demand
2. The current situation and tools in use today
3. The pain points, inefficiencies, or problems to solve

Example: "Customers cannot complete purchases using Apple Pay, which is the #1 requested payment method in our surveys. Currently we support only Visa/Mastercard. This causes a 12% cart abandonment rate on iOS devices."

Ask the user to describe the situation and problem. Then synthesise their answer into a clear paragraph and call fill_field.`,

	"whoIsImpacted": `You are helping fill the **Who Is Impacted** field.

This should list the types of users affected and estimate quantities.
Example: "Online customers: ~50,000 per month (5% of total). Store coworkers: not concerned. B2B partners: not concerned."

Ask who will be affected by this demand and how many, then call fill_field.`,

	"benefitCategory": `You are helping fill the **Benefit Category** field.

Choose the category that best describes the primary expected benefit:
1. Cost efficiency — reduces costs or improves operational efficiency
2. Environmental & social sustainability — positive impact on environment or society
3. Service quality & security risk — improves quality, reliability, or reduces security risks
4. Customer satisfaction & revenue — improves customer experience or grows revenue
5. Innovation — introduces new capabilities or disruptive approaches
6. Other — does not fit the above

Ask the user to describe the expected benefit, then map it to one category and call fill_field.`,

	"benefitHypothesis": `You are helping fill the **Benefit Hypothesis** field.

This is the user's hypothesis on HOW the demand will achieve the expected benefit.
Format: "We believe that by doing X, we will achieve Y because Z."
Example: "We believe that by adding Apple Pay, checkout conversion will increase by 8% because iOS users currently abandon at payment step due to lack of their preferred method."

Ask the user to explain their reasoning, then synthesise it into a hypothesis statement and call fill_field.`,

	"measureBenefits": `You are helping fill the **Measure Benefits** field.

This defines KPIs and the timeframe to verify the benefits.
Example: "Conversion rate on iOS devices during first 3 months after launch. Cart abandonment rate on payment step. GMV per payment method."

Ask the user how they will know the demand succeeded and by when. Then call fill_field with the KPIs and timing.`,
}
