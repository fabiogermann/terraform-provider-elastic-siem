package transferobjects

import "time"

type ThreatItem struct {
	Framework string `json:"framework,omitempty"`
	Tactic    struct {
		ID        string `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Reference string `json:"reference,omitempty"`
	} `json:"tactic,omitempty"`
	Technique []struct {
		ID        string `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Reference string `json:"reference,omitempty"`
	} `json:"technique,omitempty"`
}

type ExecutionHistoryItem struct {
	LastExecution struct {
		Date        time.Time `json:"date,omitempty"`
		Status      string    `json:"status,omitempty"`
		StatusOrder int       `json:"status_order,omitempty"`
		Message     string    `json:"message,omitempty"`
		Metrics     struct {
			TotalSearchDurationMs     int `json:"total_search_duration_ms,omitempty"`
			TotalIndexingDurationMs   int `json:"total_indexing_duration_ms,omitempty"`
			TotalEnrichmentDurationMs int `json:"total_enrichment_duration_ms,omitempty"`
		} `json:"metrics,omitempty"`
	} `json:"last_execution,omitempty"`
}

type RiskScoreMapping struct {
	Field    string `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

type SeverityMapping struct {
	Field    string `json:"field,omitempty"`
	Value    string `json:"value,omitempty"`
	Operator string `json:"operator,omitempty"`
	Severity string `json:"severity,omitempty"`
}

type ActionItem struct {
	Group  string `json:"group,omitempty"`
	ID     string `json:"id,omitempty"`
	Params struct {
		Body string `json:"body,omitempty"`
	} `json:"params,omitempty"`
	ActionTypeID string `json:"action_type_id,omitempty"`
}

type MetaItem struct {
	From             string `json:"from,omitempty"`
	KibanaSiemAppURL string `json:"kibana_siem_app_url,omitempty"`
}

type ExceptionListItem struct {
	ID            string `json:"id,omitempty"`
	ListID        string `json:"list_id,omitempty"`
	Type          string `json:"type,omitempty"`
	NamespaceType string `json:"namespace_type,omitempty"`
}

type DetectionRuleResponse struct {
	DetectionRule
	CreatedAt        time.Time            `json:"created_at,omitempty"`
	CreatedBy        string               `json:"created_by,omitempty"`
	ExceptionsList   []ExceptionListItem  `json:"exceptions_list"`
	ExecutionSummary ExecutionHistoryItem `json:"execution_summary,omitempty"`
	Meta             MetaItem             `json:"meta,omitempty"`
	UpdatedAt        time.Time            `json:"updated_at,omitempty"`
}

type DetectionRule struct {
	Actions             []ActionItem       `json:"actions,omitempty"`
	Author              []string           `json:"author,omitempty"`
	Description         string             `json:"description,omitempty"`
	Enabled             bool               `json:"enabled,omitempty"`
	FalsePositives      []interface{}      `json:"false_positives,omitempty"`
	Filters             []interface{}      `json:"filters,omitempty"`
	From                string             `json:"from,omitempty"`
	ID                  string             `json:"id,omitempty"`
	Immutable           bool               `json:"immutable,omitempty"`
	Interval            string             `json:"interval,omitempty"`
	Language            string             `json:"language,omitempty"`
	MaxSignals          int                `json:"max_signals,omitempty"`
	Name                string             `json:"name,omitempty"`
	OutputIndex         string             `json:"output_index,omitempty"`
	Query               string             `json:"query,omitempty"`
	References          []interface{}      `json:"references,omitempty"`
	RelatedIntegrations []interface{}      `json:"related_integrations,omitempty"`
	RequiredFields      []interface{}      `json:"required_fields,omitempty"`
	RiskScore           int                `json:"risk_score,omitempty"`
	RiskScoreMapping    []RiskScoreMapping `json:"risk_score_mapping,omitempty"`
	RuleID              string             `json:"rule_id,omitempty"`
	Setup               string             `json:"setup,omitempty"`
	Severity            string             `json:"severity,omitempty"`
	SeverityMapping     []SeverityMapping  `json:"severity_mapping,omitempty"`
	Tags                []string           `json:"tags,omitempty"`
	Threat              []ThreatItem       `json:"threat,omitempty"`
	Throttle            string             `json:"throttle,omitempty"`
	To                  string             `json:"to,omitempty"`
	Type                string             `json:"type,omitempty"`
	UpdatedBy           string             `json:"updated_by,omitempty"`
	Version             int                `json:"version,omitempty"`
}
