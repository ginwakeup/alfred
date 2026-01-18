package network

import (
	"os/exec"
)

func Ensure(name string) error {
	cmd := exec.Command("docker", "network", "inspect", name)
	if err := cmd.Run(); err == nil {
		return nil
	}

	create := exec.Command("docker", "network", "create", name)
	create.Stdout = nil
	create.Stderr = nil
	return create.Run()
}
