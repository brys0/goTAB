package gpu

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
)

func GetGPUInfo() ([]GPU, error) {
	// Consider not using lshw for linux, for now it's the best option.
	_, err := os.OpenFile("/usr/bin/lshw", os.O_RDONLY, 0444)

	if err != nil {
		return nil, errors.New("lshw command not found in /usr/bin/lshw, please make sure lshw is installed")
	}

	cmdOut, err := exec.Command("lshw", "-C", "display", "-quiet", "-json").Output()
	if err != nil {
		return nil, err
	}

	gpus := make([]GPU, 5)

	err = json.Unmarshal(cmdOut, &gpus)

	if err != nil {
		return nil, err
	}

	return gpus, nil
}

func (gpu *GPU) GetDeviceName() string {
	return strings.Replace(gpu.BusInfo, "@", "-", 1)
}
