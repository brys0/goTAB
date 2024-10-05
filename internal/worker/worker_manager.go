package worker

import (
	"github.com/brys0/goTAB/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Create functional testing loop as defined in client whitepaper,
// see: https://github.com/JPVenson/Jellyfin.HardwareVisualizer/blob/main/HWA-Client-Whitepaper.md#testing-loop
func CreateWorkManager(p *tea.Program, worker_args string, video string) {
	done_channel := make(chan int, 5)

	// TODO: Calculate workers based on ffmpeg speed
	workers_to_spawn := 1

	for range workers_to_spawn {
		go StartWorker("ffmpeg/bin/ffmpeg", worker_args, video, p, done_channel)
	}

	workers_done := 0
	for {
		v, ok := <-done_channel

		if ok == false {
			break
		}

		workers_done += v

		if workers_done == workers_to_spawn {
			p.Send(tui.WorkerReset(true))
			return
		}
	}

}
