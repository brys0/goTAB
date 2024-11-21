package worker

import (
	"github.com/brys0/goTAB/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

// TODO: Create functional testing loop as defined in client whitepaper,
// see: https://github.com/JPVenson/Jellyfin.HardwareVisualizer/blob/main/HWA-Client-Whitepaper.md#testing-loop
func CreateWorkManager(p *tea.Program, worker_args string, video string, total_workers int) {
	done_channel := make(chan []*FrameStat, total_workers)

	log.Info("Spawning workers", "workers", total_workers)
	for range total_workers {
		// log.Info("Spawning worker")
		go StartWorker("ffmpeg/bin/ffmpeg", worker_args, video, p, done_channel)
	}

	workers_done := 0
	for {
		v, ok := <-done_channel

		if ok == false {
			break
		}

		workers_done += 1

		log.Info("A worker just finished with stats", "workers_done", workers_done, "v", v)
		if len(v) >= 1 && v[len(v)-1].Speed < 1 {
			log.Error("Reached limited")
			return
		}

		if workers_done == total_workers {

			log.Info("Last speed was")
			log.Info("Spawning ", "workers", total_workers+1)
			CreateWorkManager(p, worker_args, video, total_workers+1)
			p.Send(tui.WorkerReset(true))
		}
	}

}
