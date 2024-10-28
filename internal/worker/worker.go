package worker

import (
	"context"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/brys0/goTAB/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

// TODO: This code is terrible, my brain hurts. go doesn't like arguments in exec.Command with just a string, so I have to manually recreate it.
// Brilliant.
func StartWorker(ffmpeg_path string, args string, video string, p *tea.Program, done chan []*FrameStat) {
	id := uuid.New()
	frame_stats := make([]*FrameStat, 0)

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
		frame_stats = readLogs(stdout, frame_stats)
	}()

	go func() {
		defer wg.Done()
		frame_stats = readLogs(stderr, frame_stats)
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

		done <- frame_stats
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
	done <- frame_stats
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

func readLogs(r io.Reader, frame_stats []*FrameStat) []*FrameStat {
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			frame := decodeStatsFrames(string(buf[0:n]))

			if frame == nil {
				continue
			}

			frame_stats = append(frame_stats, frame)
		}
		if err != nil {
			break
		}
	}

	return frame_stats
}

type FrameStat struct {
	Frame int64
	FPS   float64
	Speed float64
}

func decodeStatsFrames(stat string) *FrameStat {
	if !strings.Contains(stat, "frame") {
		return nil
	}

	benchmarkArray := parseKeyValueString(stat)

	frame, err := strconv.ParseInt(benchmarkArray["frame"], 0, 32)
	if err != nil {
		log.Fatal("Could not parse frame number from ffmpeg")
	}
	speed, err := strconv.ParseFloat(benchmarkArray["speed"], 0)

	if err != nil {
		log.Fatal("Could not parse speed float from ffmpeg")
	}

	fps, err := strconv.ParseFloat(benchmarkArray["fps"], 0)

	return &FrameStat{
		Frame: frame,
		Speed: speed,
		FPS:   fps,
	}

}

func parseKeyValueString(input string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Fields(input)

	for i, pair := range pairs {
		keyValuePair := strings.Split(pair, "=")

		// TODO: Frame is weird and doesnt parse right. oh well
		if len(keyValuePair) >= 1 && isInt(keyValuePair[0]) {
			key := strings.TrimSpace(keyValuePair[0])
			value := strings.ReplaceAll(pairs[i-1], "=", "")

			result[value] = key
			continue
		}

		if len(keyValuePair) == 2 && keyValuePair[0] != "frame" {
			key := strings.TrimSpace(keyValuePair[0])
			value := strings.ReplaceAll(strings.TrimSpace(keyValuePair[1]), "x", "")

			if value == "" {
				value = "0"
			}
			result[key] = value
		}
	}

	return result
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
