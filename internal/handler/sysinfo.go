package handler

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	CPUCores   int     `json:"cpu_cores"`
	CPUModel   string  `json:"cpu_model"`
	MemTotalMB uint64  `json:"mem_total_mb"`
	MemUsedMB  uint64  `json:"mem_used_mb"`
	MemFreeMB  uint64  `json:"mem_free_mb"`
	MemPercent float64 `json:"mem_percent"`
	CPUPercent float64 `json:"cpu_percent"`
	GoVersion  string  `json:"go_version"`
	GOOS       string  `json:"goos"`
	GOARCH     string  `json:"goarch"`
	Hostname   string  `json:"hostname"`
}

func SysInfo(c echo.Context) error {
	model := cpuModel()
	cpuPct := currentCPUPercent()
	totalMB, usedMB, freeMB, memPct := memStats()
	hostname, _ := os.Hostname()

	return c.JSON(http.StatusOK, SystemInfo{
		CPUCores:   runtime.NumCPU(),
		CPUModel:   model,
		MemTotalMB: totalMB,
		MemUsedMB:  usedMB,
		MemFreeMB:  freeMB,
		MemPercent: memPct,
		CPUPercent: cpuPct,
		GoVersion:  runtime.Version(),
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
		Hostname:   hostname,
	})
}

func cpuModel() string {
	info, err := cpu.Info()
	if err != nil || len(info) == 0 {
		return "unknown"
	}
	return info[0].ModelName
}

func currentCPUPercent() float64 {
	percents, err := cpu.Percent(200*time.Millisecond, false)
	if err != nil || len(percents) == 0 {
		return 0
	}
	return percents[0]
}

func memStats() (totalMB, usedMB, freeMB uint64, pct float64) {
	vmStat, err := mem.VirtualMemory()
	if err != nil || vmStat == nil {
		return
	}
	totalMB = vmStat.Total / 1024 / 1024
	usedMB = vmStat.Used / 1024 / 1024
	freeMB = vmStat.Free / 1024 / 1024
	pct = vmStat.UsedPercent
	return
}
