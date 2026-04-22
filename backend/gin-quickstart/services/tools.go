package services

// formProperties returns the JSON Schema properties shared by propose_form_data and fill_demand_form.
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

// BuildTools returns the two-phase form fill tool set:
//   - propose_form_data: show a summary for user review (Phase 1→2)
//   - fill_demand_form: apply fields to the form (Phase 2→3, or direct correction in Phase 3)
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

// BuildFieldTool returns the fill_field tool scoped to a single form field.
// The value schema is taken from formProperties so enum constraints are preserved.
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
