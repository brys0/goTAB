package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/brys0/goTAB/internal/app"
)

func main() {
	var noCPU bool
	var server string // Offical jellyfin HWA server

	flag.StringVar(&server, "server", "https://hwa.jellyfin.org/", "Your https://github.com/JPVenson/Jellyfin.HardwareVisualizer server, leave blank to use offical")
	flag.BoolVar(&noCPU, "nocpu", false, "Would you like to do software based encoding/decoding. This can take longer.")

	defer timer("Application")()
	application := app.CreateNewTAB(noCPU, server)
	// application.DisplayBanner()
	application.GetHardware()
	application.ValidatePlatform()
	application.PromptGPU()

	application.FetchPlatformTests()
	application.PromptConfirmTests()

	// model := tui.CreateTUIDisplay()

	// p := tea.NewProgram(model)

	// go p.Send(tui.ProgressInfo{
	// 	Type:     "task",
	// 	Progress: .25,
	// })

	// go p.Send(tui.TaskInfo("ffmpeg 4k"))
	// p.Run()

	//ti.UpdateTaskProgress(tui.ProgressInfo{Type: "task", Progress: 25})

}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
