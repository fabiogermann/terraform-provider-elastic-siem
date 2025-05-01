package transferobjects

import (
	"terraform-provider-elastic-siem/tools"
	"time"
)

type ExceptionItemResponse struct {
	ExceptionItemBase
	HTags        []interface{}               `json:"_tags,omitempty"`
	Comments     []ExceptionCommentsResponse `json:"comments,omitempty"`
	CreatedAt    time.Time                   `json:"created_at,omitempty"`
	CreatedBy    string                      `json:"created_by,omitempty"`
	TieBreakerID string                      `json:"tie_breaker_id,omitempty"`
	UpdatedAt    time.Time                   `json:"updated_a,omitemptyt"`
	UpdatedBy    string                      `json:"updated_by,omitempty"`
}

type ExceptionComments struct {
	Comment string `json:"comment,omitempty"`
}

type ExceptionCommentsResponse struct {
	ExceptionComments
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy string    `json:"created_by,omitempty"`
	ID        string    `json:"id,omitempty"`
}

type ExceptionItemEntry struct {
	Field    string            `json:"field,omitempty"`
	Operator string            `json:"operator,omitempty"`
	Type     string            `json:"type,omitempty"`
	Value    tools.StringSlice `json:"value,omitempty"`
}

type ExceptionItemBase struct {
	Description   string               `json:"description,omitempty"`
	Entries       []ExceptionItemEntry `json:"entries,omitempty"`
	ID            string               `json:"id,omitempty"`
	ListID        string               `json:"list_id,omitempty"`
	ItemID        string               `json:"item_id,omitempty"`
	Name          string               `json:"name,omitempty"`
	NamespaceType string               `json:"namespace_type,omitempty"`
	Tags          []string             `json:"tags,omitempty"`
	Type          string               `json:"type,omitempty"`
}

type ExceptionItem struct {
	ExceptionItemBase
	Comments []ExceptionComments `json:"comments,omitempty"`
}
