package worker

import (
	"context"
	"github.com/charmbracelet/log"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/brys0/goTAB/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// TODO: This code is terrible, my brain hurts. go doesn't like arguments in exec.Command with just a string, so I have to manually recreate it.
// Brilliant.
func StartWorker(ffmpeg_path string, args string, video string, p *tea.Program, done chan int) {
	log.Infof("%v %v", ffmpeg_path, args)
	id := uuid.New()
	p.Send(tui.TuiWorkerInfo{
		ID:       id.String(),
		Speed:    -1,
		Frame:    0,
		Done:     false,
		Errored:  false,
		ErrorStr: nil,
	})

	no_quotes := strings.ReplaceAll(args, "\"", "")
	s := strings.Split(strings.ReplaceAll(no_quotes, "{video_file}", video), " ")

	var wg sync.WaitGroup
	cmd := exec.CommandContext(context.Background(), ffmpeg_path, s...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		copyLogs(stdout)
	}()

	go func() {
		defer wg.Done()
		copyLogs(stderr)
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {

		errStr := err.Error()
		p.Send(tui.TuiWorkerInfo{
			ID:       id.String(),
			Speed:    -1,
			Frame:    0,
			Done:     false,
			Errored:  true,
			ErrorStr: &errStr,
		})

		done <- 1
		// println("ffmpeg/bin/ffmpeg " + args)
		// println("Actual command: ffmpeg/bin/ffmpeg " + strings.Join(s, " "))
		// panic(err)
		return
	}

	p.Send(tui.TuiWorkerInfo{
		ID:       id.String(),
		Speed:    -1,
		Frame:    0,
		Done:     true,
		Errored:  false,
		ErrorStr: nil,
	})
	done <- 1
}

func splitCharsInclusive(s, chars string) (out []string) {
	for {
		m := strings.IndexAny(s, chars)
		if m < 0 {
			break
		}
		out = append(out, s[:m], s[m:m+1])
		s = s[m+1:]
	}
	out = append(out, s)
	return
}

func copyLogs(r io.Reader) {
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			decodeStatsFrames(string(buf[0:n]))
		}
		if err != nil {
			break
		}
	}
}

type FrameStat struct {
	Frame int
	FPS   float64
	Speed float64
}

func decodeStatsFrames(stat string) {
	//log.Error(stat)
	// println(stat)
	// if !strings.Contains(stat, "frame") {
	// 	return
	// }

	// splitStr := strings.Split(stat, "=")

	// log.Println(strings.Join(splitStr, ","))
}
