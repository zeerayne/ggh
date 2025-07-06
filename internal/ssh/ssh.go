package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/byawitz/ggh/internal/config"
)

func GenerateCommandArgs(c config.SSHConfig) []string {
	key, port := "", ""
	user := "root"

	if c.User != "" {
		user = c.User
	}

	if c.Key != "" {
		key = "-i " + c.Key
	}

	if c.Port != "" {
		port = "-p " + c.Port
	}

	if c.StrictHostKeyChecking != "" {
		port = "-o StrictHostKeyChecking=" + c.StrictHostKeyChecking
	}

	if c.UserKnownHostsFile != "" {
		port = "-o UserKnownHostsFile" + c.UserKnownHostsFile
	}

	if c.LogLevel != "" {
		port = "-o LogLevel" + c.LogLevel
	}

	return strings.Split(fmt.Sprintf("%s@%s %s %s", user, c.Host, key, port), " ")
}

func Run(args []string) {
	args = slices.DeleteFunc(args, func(s string) bool { return s == "" })

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
