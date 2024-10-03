package memory

import (
	"bytes"
	"errors"
	"os/exec"

	"gopkg.in/xmlpath.v2"
)

type MemoryBank struct {
	Descr   string `json:"description"`
	Size    string `json:"size"`
	Vendor  string `json:"vendor"`
	Product string `json:"product"`
}

type Memory struct {
	Descr string       `json:"description"`
	Banks []MemoryBank `json:"banks"`
}

type CPUCache struct {
	Descr string `json:"description"`
	Size  string `json:"size"`
}

type CPU struct {
	Descr   string     `json:"description"`
	Version string     `json:"version"`
	Size    string     `json:"size"`
	Width   string     `json:"width"`
	Cache   []CPUCache `json:"cache"`
}

type DiskVolume struct {
	Descr       string `json:"description"`
	LogicalName string `json:"logical_name"`
	Size        string `json:"size"`
}

type Disk struct {
	Descr   string       `json:"description"`
	Product string       `json:"product"`
	Serial  string       `json:"serial"`
	Size    string       `json:"size"`
	Volumes []DiskVolume `json:"volumes"`
}

type Firmware struct {
	Descr    string `json:"description"`
	Vendor   string `json:"vendor"`
	Date     string `json:"date"`
	Size     string `json:"size"`
	Capacity string `json:"capacity"`
}

type Core struct {
	Descr    string     `json:"description"`
	Firmware []Firmware `json:"firmware"`
	CPU      []CPU      `json:"cpu"`
	Memory   []Memory   `json:"memory"`
	Disks    []Disk     `json:"disk"`
}

type Hardware struct {
	Descr   string `json:"description"`
	Product string `json:"product"`
	Serial  string `json:"serial"`
	Vendor  string `json:"vendor"`
	Core    Core   `json:"core"`
}

func GetHardware() (*Hardware, error) {
	cmdOut, err := exec.Command("lshw", "-quiet", "-xml").Output()
	if err != nil {
		return nil, err
	}
	return hardwareFromText(cmdOut)
}

func hardwareFromText(cmdOut []byte) (*Hardware, error) {

	var hw Hardware

	root, err := xmlpath.Parse(bytes.NewReader(cmdOut))
	if err != nil {
		return nil, err
	}

	path := xmlpath.MustCompile("/list/node/description")
	if value, ok := path.String(root); ok {
		hw.Descr = value
	}

	path = xmlpath.MustCompile("/list/node/product")
	if value, ok := path.String(root); ok {
		hw.Product = value
	}

	path = xmlpath.MustCompile("/list/node/serial")
	if value, ok := path.String(root); ok {
		hw.Serial = value
	}

	path = xmlpath.MustCompile("/list/node/vendor")
	if value, ok := path.String(root); ok {
		hw.Vendor = value
	}

	core := Core{}

	path = xmlpath.MustCompile("/list/node/node[@id='core']/description")
	if value, ok := path.String(root); ok {
		core.Descr = value
	}

	path = xmlpath.MustCompile("/list/node/node[@id='core']/node[contains(@id, 'firmware')]")
	firmwareRoots := path.Iter(root)
	firmwares := []Firmware{}
	for firmwareRoots.Next() {
		firmwareRoot := firmwareRoots.Node()
		firmware := Firmware{}
		path = xmlpath.MustCompile("description")
		if value, ok := path.String(firmwareRoot); ok {
			firmware.Descr = value
		}
		path = xmlpath.MustCompile("vendor")
		if value, ok := path.String(firmwareRoot); ok {
			firmware.Vendor = value
		}
		path = xmlpath.MustCompile("date")
		if value, ok := path.String(firmwareRoot); ok {
			firmware.Date = value
		}
		path = xmlpath.MustCompile("size")
		if value, ok := path.String(firmwareRoot); ok {
			firmware.Size = value
		}
		path = xmlpath.MustCompile("capacity")
		if value, ok := path.String(firmwareRoot); ok {
			firmware.Capacity = value
		}
		firmwares = append(firmwares, firmware)
	}

	path = xmlpath.MustCompile("/list/node/node[@id='core']/node[contains(@id, 'cpu')]")
	cpuRoots := path.Iter(root)
	cpus := []CPU{}
	for cpuRoots.Next() {
		cpuRoot := cpuRoots.Node()
		cpu := CPU{}
		path = xmlpath.MustCompile("description")
		if value, ok := path.String(cpuRoot); ok {
			cpu.Descr = value
		}
		path = xmlpath.MustCompile("version")
		if value, ok := path.String(cpuRoot); ok {
			cpu.Version = value
		}
		path = xmlpath.MustCompile("size")
		if value, ok := path.String(cpuRoot); ok {
			cpu.Size = value
		}
		path = xmlpath.MustCompile("width")
		if value, ok := path.String(cpuRoot); ok {
			cpu.Width = value
		}

		path = xmlpath.MustCompile("node[contains(@id, 'cache')]")

		caches := []CPUCache{}

		cacheRoots := path.Iter(cpuRoot)
		for cacheRoots.Next() {
			cacheRoot := cacheRoots.Node()
			cache := CPUCache{}
			path = xmlpath.MustCompile("description")
			if value, ok := path.String(cacheRoot); ok {
				cache.Descr = value
			}
			path = xmlpath.MustCompile("size")
			if value, ok := path.String(cacheRoot); ok {
				cache.Size = value
			}
			caches = append(caches, cache)
		}
		cpu.Cache = caches

		cpus = append(cpus, cpu)
	}

	path = xmlpath.MustCompile("/list/node/node[@id='core']/node[contains(@id, 'memory')]")
	memoryRoots := path.Iter(root)
	memories := []Memory{}
	for memoryRoots.Next() {
		memoryRoot := memoryRoots.Node()

		memory := Memory{}

		path = xmlpath.MustCompile("description")
		if value, ok := path.String(memoryRoot); ok {
			memory.Descr = value
		}

		path = xmlpath.MustCompile("node[contains(@id, 'bank')]")

		banks := []MemoryBank{}

		bankRoots := path.Iter(memoryRoot)
		for bankRoots.Next() {
			bankRoot := bankRoots.Node()
			bank := MemoryBank{}
			path = xmlpath.MustCompile("description")
			if value, ok := path.String(bankRoot); ok {
				bank.Descr = value
			}
			path = xmlpath.MustCompile("size")
			if value, ok := path.String(bankRoot); ok {
				bank.Size = value
			}
			path = xmlpath.MustCompile("vendor")
			if value, ok := path.String(bankRoot); ok {
				bank.Vendor = value
			}
			path = xmlpath.MustCompile("product")
			if value, ok := path.String(bankRoot); ok {
				bank.Product = value
			}
			banks = append(banks, bank)
		}
		memory.Banks = banks

		memories = append(memories, memory)
	}

	path = xmlpath.MustCompile("//*/node[contains(@id, 'disk')]")
	diskRoots := path.Iter(root)
	disks := []Disk{}
	for diskRoots.Next() {
		diskRoot := diskRoots.Node()

		disk := Disk{}

		path = xmlpath.MustCompile("description")
		if value, ok := path.String(diskRoot); ok {
			disk.Descr = value
		}
		path = xmlpath.MustCompile("product")
		if value, ok := path.String(diskRoot); ok {
			disk.Product = value
		}
		path = xmlpath.MustCompile("serial")
		if value, ok := path.String(diskRoot); ok {
			disk.Serial = value
		}
		path = xmlpath.MustCompile("size")
		if value, ok := path.String(diskRoot); ok {
			disk.Size = value
		}

		path = xmlpath.MustCompile("node[contains(@id, 'volume')]")

		volumes := []DiskVolume{}

		volumeRoots := path.Iter(diskRoot)
		for volumeRoots.Next() {
			volumeRoot := volumeRoots.Node()
			volume := DiskVolume{}
			path = xmlpath.MustCompile("description")
			if value, ok := path.String(volumeRoot); ok {
				volume.Descr = value
			}
			path = xmlpath.MustCompile("logicalname")
			if value, ok := path.String(volumeRoot); ok {
				volume.LogicalName = value
			}
			path = xmlpath.MustCompile("size")
			if value, ok := path.String(volumeRoot); ok {
				volume.Size = value
			}
			volumes = append(volumes, volume)
		}
		disk.Volumes = volumes

		disks = append(disks, disk)
	}

	core.CPU = cpus
	core.Memory = memories
	core.Firmware = firmwares
	core.Disks = disks

	hw.Core = core

	return &hw, nil
}

func Get_memory_info() ([]*MemoryInfo, error) {
	lshwCmd, err := GetHardware()

	if err != nil {
		return nil, errors.New("An error occurred when creating a new lshw2 instance: " + err.Error())
	}

	memory := lshwCmd.Core.Memory[0].Banks

	slots := make([]*MemoryInfo, len(memory))

	for dimId := 0; dimId < len(memory); dimId++ {
		dim := memory[dimId]

		append(slots, &MemoryInfo{
			Id:     dim.Descr,
			Class:  "memory",
			PhysId: dim.Product,
			Units:  "gigabytes",
			Size:   int64(dim.Size),
		})
	}

	if err != nil {
		return nil, err
	}

}
