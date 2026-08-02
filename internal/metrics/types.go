package metrics

type SystemMetrics struct {
	CPUPercent float64 `json:"cpu_percent"`
	RAMUsed    float64 `json:"ram_used"`
	RAMTotal   float64 `json:"ram_total"`

	DiskUsed  float64 `json:"disk_used"`
	DiskTotal float64 `json:"disk_total"`
}
