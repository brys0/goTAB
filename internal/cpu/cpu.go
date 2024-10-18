package cpu

import (
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"strings"
)

type CPU struct {
	Product      string  `json:"product"`
	Vendor       string  `json:"vendor"`
	Cores        int32   `json:"cores"`
	Architecture string  `json:"architecture"`
	HzAdvertised float64 `json:"hz_advertised"`
}

func GetCPUInfo() (*CPU, error) {
	processors, err := cpu.Info()

	if err != nil {
		return nil, err
	}

	// GET the first cpu, we don't support multiple CPU support yet.
	processor := processors[0]

	ghwProcessors, err := ghw.CPU()

	if err != nil {
		return nil, err
	}

	arch, err := host.Info()

	if err != nil {
		return nil, err
	}

	return &CPU{
		Product:      processor.ModelName,
		Vendor:       TrimCPUName(processor.VendorID),
		Cores:        int32(ghwProcessors.TotalCores),
		Architecture: arch.KernelArch,
		HzAdvertised: processor.Mhz,
	}, nil
}

func (c *CPU) GetCPUArch() string {
	if strings.Contains(c.Architecture, "64") {
		return "amd64"
	}

	if strings.Contains(c.Architecture, "aarch64") {
		return "arm64"
	}

	if strings.Contains(c.Architecture, "arm64") {
		return "arm"
	}

	return c.Architecture
}

func TrimCPUName(name string) string {
	if strings.Contains(strings.ToLower(name), "intel") {
		return "Intel"
	}

	if strings.Contains(strings.ToLower(name), "amd") {
		return "Amd"
	}

	return name
}
