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

// PentestListResponse for paginated pentest list
type PentestListResponse struct {
	Pentests   []*PentestData `json:"pentests"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// SchedulePentestResponse returned after scheduling a pentest
type SchedulePentestResponse struct {
	PentestID string `json:"pentest_id"`
}

// ReportFormat for downloading reports
type ReportFormat string

const (
	ReportFormatPDF      ReportFormat = "pdf"
	ReportFormatJSON     ReportFormat = "json"
	ReportFormatMarkdown ReportFormat = "md"
)

// DownloadReportRequest for requesting a pentest report
type DownloadReportRequest struct {
	PentestID string       `json:"pentest_id"`
	Format    ReportFormat `json:"format"`
	OutputDir string       `json:"output_dir"`
}
