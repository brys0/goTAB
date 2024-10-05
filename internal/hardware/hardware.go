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
	FFmpegSourceURL string       `json:"ffmpeg_source_url"`
	FFmpegVersion   string       `json:"ffmpeg_version"`
	FFmpegHashes    []HashedFile `json:"ffmpeg_hashs"`
}

type HashedFile struct {
	Type string `json:"type"`
	Hash string `json:"hash"`
}
