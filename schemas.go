package pt_client_sdk

// Credentials for API authentication
type Credentials struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// PaginationParams for listing pentests
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// HTTPScope represents the scope of a pentest as a string
type HTTPScope string

const (
	HTTPScopeHolistic HTTPScope = "HOLISTIC"
	HTTPScopeTargeted HTTPScope = "TARGETED"
)

// HTTPType represents the type of pentest as a string
type HTTPType string

const (
	HTTPTypeBlackBox HTTPType = "BLACK_BOX"
	HTTPTypeWhiteBox HTTPType = "WHITE_BOX"
)

// HTTPStyle represents the testing style as a string
type HTTPStyle string

const (
	HTTPStyleAggressive HTTPStyle = "AGGRESSIVE"
	HTTPStyleSafe       HTTPStyle = "SAFE"
)

// HTTPStatus represents the status of a pentest or target as a string
type HTTPStatus string

const (
	HTTPStatusPending    HTTPStatus = "PENDING"
	HTTPStatusInProgress HTTPStatus = "IN_PROGRESS"
	HTTPStatusCompleted  HTTPStatus = "COMPLETED"
	HTTPStatusFailed     HTTPStatus = "FAILED"
)

// HTTPPhase represents the current phase of a pentest as a string
type HTTPPhase string

const (
	HTTPPhaseRecon           HTTPPhase = "RECON"
	HTTPPhaseInitialExploit  HTTPPhase = "INITIAL_EXPLOIT"
	HTTPPhaseDeepExploit     HTTPPhase = "DEEP_EXPLOIT"
	HTTPPhaseLateralMovement HTTPPhase = "LATERAL_MOVEMENT"
	HTTPPhaseReport          HTTPPhase = "REPORT"
	HTTPPhaseFinished        HTTPPhase = "FINISHED"
)

// HTTPSeverity represents the severity level as a string
type HTTPSeverity string

const (
	HTTPSeverityNone     HTTPSeverity = "NONE"
	HTTPSeverityLow      HTTPSeverity = "LOW"
	HTTPSeverityMedium   HTTPSeverity = "MEDIUM"
	HTTPSeverityHigh     HTTPSeverity = "HIGH"
	HTTPSeverityCritical HTTPSeverity = "CRITICAL"
)

// HTTPTargetRequest represents a target request for HTTP API
type HTTPTargetRequest struct {
	ID          *string   `json:"id,omitempty"`
	Target      string    `json:"target"`
	Scope       HTTPScope `json:"scope"`
	Type        HTTPType  `json:"type"`
	Credentials *string   `json:"credentials,omitempty"`
}

// HTTPSchedulePentestRequest represents a request to schedule a pentest via HTTP
type HTTPSchedulePentestRequest struct {
	ID      *string              `json:"id,omitempty"`
	Style   HTTPStyle            `json:"style"`
	Exploit bool                 `json:"exploit"`
	Targets []*HTTPTargetRequest `json:"targets"`
}

// HTTPTargetData represents target information for HTTP API
type HTTPTargetData struct {
	ID          string        `json:"id"`
	PentestID   string        `json:"pentest_id"`
	Target      string        `json:"target"`
	Scope       HTTPScope     `json:"scope"`
	Type        HTTPType      `json:"type"`
	Status      HTTPStatus    `json:"status"`
	Phase       *HTTPPhase    `json:"phase,omitempty"`
	CreatedAt   *string       `json:"created_at,omitempty"`
	StartedAt   *string       `json:"started_at,omitempty"`
	FinishedAt  *string       `json:"finished_at,omitempty"`
	Credentials *string       `json:"credentials,omitempty"`
	Severity    *HTTPSeverity `json:"severity,omitempty"`
	Findings    *int32        `json:"findings,omitempty"`
	Summary     *string       `json:"summary,omitempty"`
}

// HTTPPentestData represents complete pentest information for HTTP API
type HTTPPentestData struct {
	ID         string            `json:"id"`
	Status     HTTPStatus        `json:"status"`
	CreatedAt  *string           `json:"created_at,omitempty"`
	StartedAt  *string           `json:"started_at,omitempty"`
	FinishedAt *string           `json:"finished_at,omitempty"`
	Style      HTTPStyle         `json:"style"`
	Exploit    bool              `json:"exploit"`
	Summary    *string           `json:"summary,omitempty"`
	Targets    []*HTTPTargetData `json:"targets"`
	Severity   *HTTPSeverity     `json:"severity,omitempty"`
	Findings   *int32            `json:"findings,omitempty"`
}

// HTTPSchedulePentestResponse represents the response after scheduling a pentest
type HTTPSchedulePentestResponse struct {
	PentestID string `json:"pentest_id"`
}

// HTTPPentestListResponse represents a paginated list of pentests
type HTTPPentestListResponse struct {
	Pentests   []*HTTPPentestData `json:"pentests"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}
