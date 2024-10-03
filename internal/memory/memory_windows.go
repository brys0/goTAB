package memory

func Get_memory_info() ([]*MemoryInfo, error) {
	mem, err := ghw.Memory()

	if err != nil {
		return nil, err
	}

	slots := make([]*MemoryInfo, len(mem.Modules))

	for dimId := 0; dimId < len(mem.Modules); dimId++ {
		ffff

		slots = append(slots, &MemoryInfo{
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
