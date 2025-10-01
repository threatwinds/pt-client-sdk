package twpt_client_sdk

import (
	"time"
)

type TaskType string

const (
	TaskTypeRecon          TaskType = "recon"
	TaskTypeExploitTool    TaskType = "exploit-tool"
	TaskTypeExploitScript  TaskType = "exploit-script"
	TaskTypeExploitLateral TaskType = "exploit-lateral"
	TaskTypeReport         TaskType = "report"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Pentest struct {
	ID         string     `json:"id"`
	Status     Status     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Result     *string    `json:"result"`
	Targets    []Target   `json:"targets,omitempty"`
}

type Target struct {
	ID         string     `json:"id"`
	PentestID  string     `json:"pentest_id"`
	Target     string     `json:"target"`
	Status     Status     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Result     *string    `json:"result"`
	Pentest    *Pentest   `json:"pentest,omitempty"`
	Tasks      []Task     `json:"tasks,omitempty"`
}

type Task struct {
	ID       string   `json:"id"`
	TargetID string   `json:"target_id"`
	Type     TaskType `json:"type"`
	Status   Status   `json:"status"`
	Result   *string  `json:"result"`
	Phase    int      `json:"phase"`
	Target   *Target  `json:"target,omitempty"`
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
	Targets []string `json:"targets"`
}

type SchedulePentestResponse struct {
	PentestID string `json:"pentest_id"`
}

type ReportFormat string

const (
	ReportFormatPDF      ReportFormat = "pdf"
	ReportFormatAll      ReportFormat = "all"
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
