package hardware

import (
	"github.com/brys0/goTAB/internal/cpu"
	"github.com/brys0/goTAB/internal/gpu"
	"github.com/brys0/goTAB/internal/memory"
	"github.com/brys0/goTAB/internal/os"
)

type HardwareInfo struct {
	FFmpeg *FFmpegInfo     `json:"ffmpeg"`
	OS     *os.OS          `json:"os"`
	CPU    *cpu.CPU        `json:"cpu"`
	Memory []memory.Memory `json:"memory"`
	GPU    []gpu.GPU       `json:"gpu"`
}

type FFmpegInfo struct {
	FFmpegSourceURL string
	FFmpegVersion   string
	FFmpegHashes    []HashedFile
}

type HashedFile struct {
	Type string
	Hash string
}
