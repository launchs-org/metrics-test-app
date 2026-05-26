package stress

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type State struct {
	mu          sync.Mutex
	cpuCancel   context.CancelFunc
	cpuRunning  bool
	memRunning  bool
	cpuPercent  float64
	memBytes    int64
	allocBlocks [][]byte
}

type Status struct {
	CPURunning bool    `json:"cpu_running"`
	MemRunning bool    `json:"mem_running"`
	CPUPercent float64 `json:"cpu_percent"`
	MemMB      int64   `json:"mem_mb"`
}

func (s *State) GetStatus() Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Status{
		CPURunning: s.cpuRunning,
		MemRunning: s.memRunning,
		CPUPercent: s.cpuPercent,
		MemMB:      s.memBytes / 1024 / 1024,
	}
}

func (s *State) StartCPU(percent float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cpuCancel != nil {
		s.cpuCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cpuCancel = cancel
	s.cpuRunning = true
	s.cpuPercent = percent

	for range runtime.NumCPU() {
		go runCPUWorker(ctx, percent)
	}
}

func (s *State) StopCPU() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cpuCancel != nil {
		s.cpuCancel()
		s.cpuCancel = nil
	}
	s.cpuRunning = false
	s.cpuPercent = 0
}

func (s *State) StartMem(mb int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 前回の確保を先に解放
	s.allocBlocks = nil

	blocks := make([][]byte, mb)
	for i := range mb {
		block := make([]byte, 1024*1024)
		for j := range block {
			block[j] = byte(j)
		}
		blocks[i] = block
	}
	s.allocBlocks = blocks
	s.memRunning = true
	s.memBytes = mb * 1024 * 1024
}

func (s *State) StopMem() {
	s.mu.Lock()
	s.allocBlocks = nil
	s.memRunning = false
	s.memBytes = 0
	s.mu.Unlock()

	// ポインタを nil にした後、GC を即時実行して OS にメモリを返す
	runtime.GC()
}

func runCPUWorker(ctx context.Context, percent float64) {
	workMs := int64(percent * 10)
	sleepMs := 1000 - workMs
	var counter int64

	for {
		deadline := time.Now().Add(time.Duration(workMs) * time.Millisecond)
		for time.Now().Before(deadline) {
			atomic.AddInt64(&counter, 1)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(sleepMs) * time.Millisecond):
		}
	}
}
