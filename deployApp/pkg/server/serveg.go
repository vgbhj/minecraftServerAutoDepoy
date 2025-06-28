package pkg

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func DeployCommands(client *ssh.Client, commands []string) (string, error) {
	for _, cmd := range commands {
		session, err := client.NewSession()
		if err != nil {
			return "", fmt.Errorf("failed to create SSH session: %w", err)
		}

		var cmdOutput bytes.Buffer
		session.Stdout = &cmdOutput
		session.Stderr = &cmdOutput

		err = session.Run(cmd)
		session.Close()

		if err != nil {
			fullOutput := cmdOutput.String()
			// Берем последние 500 символов
			truncated := ""
			if len(fullOutput) > 500 {
				truncated = fullOutput[len(fullOutput)-500:]
			} else {
				truncated = fullOutput
			}
			return truncated, fmt.Errorf("command failed: %w", err)
		}
	}

	return "All commands executed successfully", nil
}
