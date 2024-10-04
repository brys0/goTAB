package cpu

import (
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/cpu"
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

	topology, err := ghw.Topology()

	if err != nil {
		return nil, err
	}

	return &CPU{
		Product:      processor.ModelName,
		Vendor:       processor.VendorID,
		Cores:        processor.Cores,
		Architecture: topology.Architecture.String(),
		HzAdvertised: processor.Mhz,
	}, nil
}
