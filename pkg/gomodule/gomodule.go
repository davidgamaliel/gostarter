package gomodule

import (
	"bytes"
	"log/slog"
	"os/exec"
)

func GenerateGoMod(tempFolderPath, name string) error {
	err := ExecuteCmd("go", []string{"mod", "init", name}, tempFolderPath)
	if err != nil {
		slog.Error("Failed to generate go.mod file", err)
		return err

	}
	return nil
}

func GoGetPackage(appDir string, packages []string) error {
	for _, packageName := range packages {
		if err := ExecuteCmd("go",
			[]string{"get", "-u", packageName},
			appDir); err != nil {
			return err
		}
	}

	return nil
}

func ExecuteCmd(name string, args []string, dir string) error {
	command := exec.Command(name, args...)
	command.Dir = dir
	var out bytes.Buffer
	command.Stdout = &out
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}
