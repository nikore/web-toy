package loadgen

import "fmt"
import "log"
import "sync"
import "time"

type MemLoadGenerator struct {
	gate    *sync.Once
	mem_mb  int64
	running bool
}

func NewMemLoadGenerator(mem_mb int64) *MemLoadGenerator {
	if mem_mb < 1 {
		log.Panicf("Invalid mem_mb passed to NewMemLoadGenerator() (%d)", mem_mb)
	}
	return &MemLoadGenerator{
		gate:    &sync.Once{},
		mem_mb:  mem_mb,
		running: false,
	}
}

func (m *MemLoadGenerator) Start() {
	m.gate.Do(m.genMemLoad)
	m.running = true
}

func (m *MemLoadGenerator) Status() string {
	if m.running {
		return fmt.Sprintf("Allocated %d MiB of memory\n", m.mem_mb)
	}
	return fmt.Sprintf("Memory not allocated (target will be %d MiB)\n", m.mem_mb)
}

func (m *MemLoadGenerator) genMemLoad() {
	mem := allocateMem(m.mem_mb)
	go func() {
		for {
			// Sweep the memory once a minute to keep it paged in
			for _, block := range mem {
				for idx := range block {
					if idx&0x1 == 1 {
						block[idx] ^= 0xAA
					} else {
						block[idx] ^= 0xFF
					}
				}
			}
			time.Sleep(time.Minute)
		}
	}()
}

func allocateMem(mem_mb int64) [][]byte {
	mem := make([][]byte, 0, 0)
	for i := int64(0); i < mem_mb; i++ {
		mem = append(mem, make([]byte, 1024*1024))
	}
	return mem
}
