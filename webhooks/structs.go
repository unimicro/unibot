package webhooks

const (
	PhaseQueued    = "QUEUED"
	PhaseStarted   = "STARTED"
	PhaseCompleted = "COMPLETED"
	PhaseFinalized = "FINALIZED"
	StatusFailure  = "FAILURE"
	StatusSuccess  = "SUCCESS"
	StatusAborted  = "ABORTED"
)

type JenkinsNotification struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	URL         string `json:"url"`
	Build       struct {
		FullURL   string `json:"full_url"`
		Number    int    `json:"number"`
		QueueID   int    `json:"queue_id"`
		Timestamp int64  `json:"timestamp"`
		Phase     string `json:"phase"`
		Status    string `json:"status"`
		URL       string `json:"url"`
		Scm       struct {
			URL    string `json:"url"`
			Branch string `json:"branch"`
			Commit string `json:"commit"`
		} `json:"scm"`
		Parameters struct {
			Branch string `json:"branch"`
		} `json:"parameters"`
		Log       string `json:"log"`
		Artifacts struct {
		} `json:"artifacts"`
	} `json:"build"`
}
