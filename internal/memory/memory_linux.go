package memory

import (
	"github.com/jaypipes/ghw"
)

func GetMemoryInfo() ([]Memory, error) {
	mem, err := ghw.Memory()

	if err != nil {
		return nil, err
	}

	slots := make([]Memory, 0)

	for dimId := 0; dimId < len(mem.Modules); dimId++ {
		dim := mem.Modules[dimId]

		// Empty slot
		if dim == nil {
			continue
		}

		slots = append(slots, Memory{
			Id:         dim.Label,
			Class:      "memory",
			PhysId:     dim.SerialNumber,
			Units:      "gigabytes",
			Size:       dim.SizeBytes / GIGABYTES,
			Vendor:     dim.Vendor,
			Speed:      0,
			FormFactor: dim.Location,
		})
	}

	return slots, nil

}
