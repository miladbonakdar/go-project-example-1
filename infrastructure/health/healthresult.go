package health

type HealthResultDto struct {
	Status      Health            `json:"status"`
	Duration    string            `json:"duration"`
	Exception   string            `json:"exception,omitempty"`
	Description string            `json:"description,omitempty"`
	Data        map[string]string `json:"data"`
}
