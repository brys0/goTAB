package app

import (
	"github.com/brys0/goTAB/internal/cpu"
	"github.com/brys0/goTAB/internal/gpu"
	"github.com/brys0/goTAB/internal/hardware"
	"github.com/brys0/goTAB/internal/memory"
	"github.com/brys0/goTAB/internal/os"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

type App struct {
	SelectedGraphics []int   `json:"-"`
	CPUTest          bool    `json:"-"`
	Errors           []error `json:"-"`

	Token  *string                `json:"token"`
	HwInfo *hardware.HardwareInfo `json:"hwinfo"`
	Tests  []int                  `json:"-"`
}

func CreateNewTAB(testCPU bool) *App {
	return &App{
		SelectedGraphics: []int{},
		CPUTest:          testCPU,
		Errors:           []error{},

		// json data
		Token:  nil,
		HwInfo: nil,
		Tests:  []int{},
	}
}

func (app *App) GetHardware() {
	memoryInfo, err := memory.GetMemoryInfo()

	if err != nil {
		log.Fatal("Error occurred trying to get system memoryInfo info", err)
	}

	cpuInfo, err := cpu.GetCPUInfo()

	if err != nil {
		log.Fatal("Error occurred trying to get system CPU info", err)
	}

	gpuInfo, err := gpu.GetGPUInfo()

	if err != nil {
		log.Fatal("Error occurred trying to get system GPU info", err)
	}

	osInfo, err := os.GetOSInfo()

	if err != nil {
		log.Fatal("Error occurred trying to get system info", err)
	}

	app.HwInfo = &hardware.HardwareInfo{
		FFmpeg: nil,
		OS:     osInfo,
		CPU:    cpuInfo,
		Memory: memoryInfo,
		GPU:    gpuInfo,
	}
}

func (app *App) PromptGPU() error {
	if app.HwInfo == nil {
		log.Fatal("No hardware info")
	}

	gpuSelectOptions := make([]huh.Option[int], len(app.HwInfo.GPU))

	for id := range app.HwInfo.GPU {
		card := app.HwInfo.GPU[id]
		gpuSelectOptions[id] = huh.Option[int]{Key: card.Description, Value: id}
	}

	gpuMultiSelect := huh.NewMultiSelect[int]().
		Title("Select a GPU").
		Options(gpuSelectOptions...).
		Description("Press [Space] to select and [Enter] to submit").
		Value(&app.SelectedGraphics)

	gpuSelectForm := huh.NewForm(huh.NewGroup(gpuMultiSelect))

	err := gpuSelectForm.Run()

	if err != nil {
		log.Fatal("Error occurred trying to select a GPU", err)
	}

	for i := range app.SelectedGraphics {
		id := app.SelectedGraphics[i]
		log.Info("Running test for", "gpu", app.HwInfo.GPU[id].Description)
	}

	return nil
}
