package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"slices"

	"github.com/byawitz/ggh/internal/config"
)

func GenerateCommandArgs(c config.SSHConfig) []string {
	var args []string

	// Handle User and HostName explicitly first
	if len(c.User) > 0 {
		args = append(args, fmt.Sprintf("%s@%s", c.User, c.Host))
	} else {
		args = append(args, fmt.Sprintf("root@%s", c.Host))
	}

	if len(c.Port) > 0 {
		args = append(args, "-p", c.Port)
	}

	if len(c.Key) > 0 {
		args = append(args, "-i", c.Key)
	}

	val := reflect.ValueOf(c)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("ssh")

		if tag == "" {
			continue
		}

		// Skip standard flags handled manually above
		if tag == "Port" || tag == "IdentityFile" || tag == "User" || tag == "HostName" {
			continue
		}

		fieldVal := val.Field(i)

		if field.Type.Kind() == reflect.Slice {
			for j := 0; j < fieldVal.Len(); j++ {
				args = append(args, "-o", fmt.Sprintf("%s %s", tag, fieldVal.Index(j).String()))
			}
		} else {
			value := fieldVal.String()
			if value == "" {
				continue
			}
			args = append(args, "-o", fmt.Sprintf("%s=%s", tag, value))
		}
	}

	return args
}

func Run(args []string) {
	args = slices.DeleteFunc(args, func(s string) bool { return s == "" })

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
