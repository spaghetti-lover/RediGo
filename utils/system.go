package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type CPUMonitor struct{}

func (m *CPUMonitor) Check(ctx context.Context) (string, error) {
	percent, err := cpu.PercentWithContext(ctx, 1*time.Second, false)
	if err != nil {
		return "N/A", err
	}

	return fmt.Sprintf("%.2f", percent[0]) + "%", nil
}

type MemoryMonitor struct{}

func (m *MemoryMonitor) Check(ctx context.Context) (map[string]string, error) {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return map[string]string{
			"used_percent": "N/A",
			"used_mb":      "N/A",
			"total_mb":     "N/A",
		}, err
	}

	return map[string]string{
		"used_percent": fmt.Sprintf("%.2f%%", vmStat.UsedPercent),
		"used_mb":      fmt.Sprintf("%.2f", float64(vmStat.Used)/1024/1024),
		"total_mb":     fmt.Sprintf("%.2f", float64(vmStat.Total)/1024/1024),
	}, nil
}

// SystemStats provides combined system statistics
type SystemStats struct {
	CPU    *CPUMonitor
	Memory *MemoryMonitor
}

func NewSystemStats() *SystemStats {
	return &SystemStats{
		CPU:    &CPUMonitor{},
		Memory: &MemoryMonitor{},
	}
}

func (s *SystemStats) GetAllStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get CPU stats
	cpuStat, err := s.CPU.Check(ctx)
	if err != nil {
		stats["cpu_percent"] = "N/A"
	} else {
		stats["cpu_percent"] = cpuStat
	}

	// Get Memory stats
	memStats, err := s.Memory.Check(ctx)
	if err != nil {
		stats["memory_used_percent"] = "N/A"
		stats["memory_used_mb"] = "N/A"
		stats["memory_total_mb"] = "N/A"
	} else {
		stats["memory_used_percent"] = memStats["used_percent"]
		stats["memory_used_mb"] = memStats["used_mb"]
		stats["memory_total_mb"] = memStats["total_mb"]
	}

	return stats, nil
}
