package worker

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

// TODO: This code is terrible, my brain hurts. go doesn't like you exec.Command with just a string, so I have to manually recreate it.
// Brilliant.
func StartWorker(ffmpeg_path string, args string, video string) {
	s := strings.Split(args, "-")

	newStr := make([]string, 0)
	for l := range s {
		str := s[l]

		newStr = append(newStr, fmt.Sprintf("-%s", strings.ReplaceAll(str, "{video_file}", video)))
	}
	cmd := exec.Command(ffmpeg_path, newStr...)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	log.Print(cmd.String())

	// Execute the command
	if err := cmd.Run(); err != nil {
		panic(cmd.String())
	}

	log.Print(stdBuffer.String())
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
