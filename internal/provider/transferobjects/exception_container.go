package transferobjects

import "time"

type ExceptionContainer struct {
	Description   string   `json:"description,omitempty"`
	Name          string   `json:"name,omitempty"`
	ListID        string   `json:"list_id,omitempty"`
	Type          string   `json:"type,omitempty"`
	NamespaceType string   `json:"namespace_type,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	ID            string   `json:"id,omitempty"`
}

type ExceptionContainerResponse struct {
	HVersion      string        `json:"_version,omitempty"`
	HTags         []interface{} `json:"_tags,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	CreatedBy     string        `json:"created_by,omitempty"`
	Description   string        `json:"description,omitempty"`
	ID            string        `json:"id,omitempty"`
	Immutable     bool          `json:"immutable,omitempty"`
	ListID        string        `json:"list_id,omitempty"`
	Name          string        `json:"name,omitempty"`
	NamespaceType string        `json:"namespace_type,omitempty"`
	OSTypes       []string      `json:"os_types,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
	TieBreakerID  string        `json:"tie_breaker_id,omitempty"`
	Type          string        `json:"type,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
	UpdatedBy     string        `json:"updated_by,omitempty"`
	Version       int           `json:"version,omitempty"`
}
