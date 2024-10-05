package fio

import (
	"os/exec"
)

func Unzip(src, dest string) error {
	unxz := exec.Command("tar", "-xf", src, "-C", dest)

	err := unxz.Start()

	unxz.Wait()

	if err != nil {
		return err
	}

	return nil
}
