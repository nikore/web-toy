package loadgen

import (
	"math"
	"sync"
)

type CpuLoadGenerator struct {
	gate    *sync.Once
	running bool
}

func NewCpuLoadGenerator() *CpuLoadGenerator {
	return &CpuLoadGenerator{
		gate:    &sync.Once{},
		running: false,
	}
}

func (l *CpuLoadGenerator) Start() {
	l.gate.Do(genLoad)
	l.running = true
}

func (l *CpuLoadGenerator) Status() string {
	if l.running {
		return "CPU load generation 100%\n"
	}
	return "CPU load generation not started\n"
}

func genLoad() {
	go func() {
		var x uint64 = 0
		var _ float64
		for {
			x ^= 0xFFFFFFFFFFFFFFFF
			_ = math.Pow(32768, 4)
		}

	}()
}
