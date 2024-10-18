package fio

import (
	"github.com/charmbracelet/log"
	"os"
	"os/exec"
)

func Unzip(src, dest string) error {
	os.MkdirAll(dest, 0755)
	unxz := exec.Command("tar", "-xf", src, "-C", dest)

	log.Debug("Exec command", "args", unxz.Args)
	err := unxz.Start()

	err = unxz.Wait()

	if err != nil {
		return err
	}

	return nil
}
