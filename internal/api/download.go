package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

type progressWriter struct {
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(float64)
}

func (pw *progressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw.file, io.TeeReader(pw.reader, pw))
	if err != nil {
		log.Fatal("Could not downlaod ffmpeg file")
	}
}
func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}
	return len(p), nil
}

func DownloadFile(url string, name string, onProgress func(ratio float64)) {

	resp, err := http.Get(url) // nolint:gosec

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatal("Could not download", "url", url, "name", name)
	}

	defer resp.Body.Close()

	file, err := create(name)

	if err != nil {
		log.Fatal("Could not create file", "file", name, "err", err)
	}

	defer file.Close() // nolint:errcheck

	pw := &progressWriter{
		total:      int(resp.ContentLength),
		file:       file,
		reader:     resp.Body,
		onProgress: onProgress,
	}

	pw.Start()
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
