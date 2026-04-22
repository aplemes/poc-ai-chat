package models

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type Message struct {
	Role       Role       `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function FuncCall `json:"function"`
}

type FuncCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Session struct {
	ID              string        `json:"id"`
	Messages        []Message     `json:"messages"`
	Status          string        `json:"status"` // "collecting" | "complete"
	PendingFormData *FormFillData `json:"-"`
}

type FormFillData struct {
	Title               string   `json:"title"`
	BusinessLine        string   `json:"businessLine"`
	RequesterBU         string   `json:"requesterBU"`
	BusInterested       []string `json:"busInterested"`
	TimeSensitive       string   `json:"timeSensitive"`
	WhyDemand           string   `json:"whyDemand"`
	WhoIsImpacted       string   `json:"whoIsImpacted"`
	BenefitCategory     string   `json:"benefitCategory"`
	BenefitHypothesis   string   `json:"benefitHypothesis"`
	MeasureBenefits     string   `json:"measureBenefits"`
	LowConfidenceFields []string `json:"lowConfidenceFields,omitempty"`
}
