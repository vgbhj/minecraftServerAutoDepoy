package pkg

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func DeployCommands(client *ssh.Client, commands []string) (string, error) {
	var fullOutput bytes.Buffer

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

		fullOutput.WriteString(cmdOutput.String() + "\n")

		if err != nil {
			output := fullOutput.String()
			// Берем последние 500 символов
			if len(output) > 500 {
				output = output[len(output)-500:]
			}
			return output, fmt.Errorf("command failed: %w", err)
		}
	}

	return fullOutput.String(), nil
}
