package github

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func CreateRepoInTemp(giturl string) (string, error) {
	tempDir, err := os.MkdirTemp("", "cloudhub-*")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(
		"git",
		"clone",
		giturl,
		tempDir,
	)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf(
			"git clone failed: %v\n%s",
			err,
			string(output),
		)
	}
	log.Println("workspace:", tempDir)

	return tempDir, nil
}
