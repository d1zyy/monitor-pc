package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"math"
	"time"
)

type WindowsCollector struct {
	diskPath string
}

func NewWindowsCollector() *WindowsCollector {
	return &WindowsCollector{
		diskPath: "C:",
	}
}

func (c *WindowsCollector) Collect() (*SystemMetrics, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	ram, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskUsage, err := disk.Usage(c.diskPath)
	if err != nil {
		return nil, err
	}

	return &SystemMetrics{
		CPUPercent: round_windows(cpuPercent[0]),
		RAMUsed:    round_windows(float64(ram.Used) / 1024 / 1024 / 1024),
		RAMTotal:   round_windows(float64(ram.Total) / 1024 / 1024 / 1024),
		DiskUsed:   round_windows(float64(diskUsage.Used) / 1024 / 1024 / 1024),
		DiskTotal:  round_windows(float64(diskUsage.Total) / 1024 / 1024 / 1024),
	}, nil
}

func round_windows(val float64) float64 {
	return math.Round(val*100) / 100
}
