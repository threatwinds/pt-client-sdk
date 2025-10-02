package twpt_client_sdk

import (
	"time"
)

type Status string

const (
	StatusRecongnizing Status = "recognizing"
	StatusExploiting   Status = "exploiting"
	StatusReporting    Status = "reporting"

	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Scope string

const (
	ScopeHolistic Scope = "HOLISTIC"
	ScopeTargeted Scope = "TARGETED"
)

type Type string

const (
	TypeBlackBox Type = "BLACK_BOX"
	TypeWhiteBox Type = "WHITE_BOX"
)

type Style string

const (
	StyleAggressive Style = "AGGRESSIVE"
	StyleSafe       Style = "SAFE"
)

type Base struct {
	ID         string     `json:"id"`
	Status     Status     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

type Pentest struct {
	Base
	Style   Style          `json:"style"`
	Exploit bool           `json:"exploit"`
	Summary map[string]any `json:"summary,omitempty"`
	Targets []Target       `json:"targets,omitempty"`
}

type Target struct {
	Base
	PentestID string  `json:"pentest_id"`
	Target    string  `json:"target"`
	Scope     Scope   `json:"scope"`
	Type      Type    `json:"type"`
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty"`
}

type Credentials struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PentestListResponse struct {
	Pentests   []Pentest `json:"pentests"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

type SchedulePentestRequest struct {
	Pentest
}

type SchedulePentestResponse struct {
	PentestID string `json:"pentest_id"`
}

type ReportFormat string

const (
	ReportFormatPDF      ReportFormat = "pdf"
	ReportFormatJson     ReportFormat = "json"
	ReportFormatMarkdown ReportFormat = "md"
)

type DownloadReportRequest struct {
	PentestID string       `json:"pentest_id"`
	Format    ReportFormat `json:"format"`
	OutputDir string       `json:"output_dir"`
}

type PentestSubscription struct {
	Updates  <-chan Pentest
	Messages <-chan string
	Errors   <-chan error
}
