package pkg

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

func DeployCommands(client *ssh.Client, commands []string) (string, error) {
	var output strings.Builder

	for _, cmd := range commands {
		session, err := client.NewSession()
		if err != nil {
			return output.String(), fmt.Errorf("failed to create SSH session: %w", err)
		}

		var cmdOutput bytes.Buffer
		session.Stdout = &cmdOutput
		session.Stderr = &cmdOutput

		err = session.Run(cmd)
		session.Close()

		output.WriteString(fmt.Sprintf("$ %s\n%s\n", cmd, cmdOutput.String()))

		if err != nil {
			output.WriteString(fmt.Sprintf("Error: %v\n", err))
			return output.String(), err
		}
	}

	return output.String(), nil
}
