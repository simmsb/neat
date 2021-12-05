package script

import "os/exec"

func Run(args ...string) error {
	cmd := exec.Command("/bin/sh", args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
