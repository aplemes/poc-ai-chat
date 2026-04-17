package models

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role      Role       `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
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
	ID       string    `json:"id"`
	Messages []Message `json:"messages"`
	Status   string    `json:"status"` // "collecting" | "complete"
}

type FormFillData struct {
	Title            string   `json:"title"`
	DemandScope      string   `json:"demandScope"`
	BusinessLine     string   `json:"businessLine"`
	RequesterBU      string   `json:"requesterBU"`
	BusInterested    []string `json:"busInterested"`
	DemandContext    string   `json:"demandContext"`
	CurrentSituation string   `json:"currentSituation"`
	ProblemsToSolve  string   `json:"problemsToSolve"`
	WhoIsImpacted    string   `json:"whoIsImpacted"`
	MeasureBenefits  string   `json:"measureBenefits"`
}
