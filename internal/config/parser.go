package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"path/filepath"

	"github.com/byawitz/ggh/internal/theme"
	"github.com/charmbracelet/bubbles/table"
)

type SSHConfig struct {
	Name                  string   `json:"name"`
	Host                  string   `json:"host" ssh:"HostName"`
	Port                  string   `json:"port" ssh:"Port"`
	User                  string   `json:"user" ssh:"User"`
	Key                   string   `json:"key" ssh:"IdentityFile"`
	UserKnownHostsFile    string   `json:"userknownhostsfile" ssh:"UserKnownHostsFile"`
	StrictHostKeyChecking string   `json:"stricthostkeychecking" ssh:"StrictHostKeyChecking"`
	LogLevel              string   `json:"loglevel" ssh:"LogLevel"`
	SetEnv                []string `json:"setenv" ssh:"SetEnv" repeatable:"true"`
	ConnectTimeout        string   `json:"connecttimeout" ssh:"ConnectTimeout"`
}

type optionSetter func(*SSHConfig, string)

var supportedOptions = map[string]optionSetter{}

func init() {
	t := reflect.TypeOf(SSHConfig{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("ssh")
		if tag == "" {
			continue
		}

		if field.Tag.Get("repeatable") == "true" {
			supportedOptions[tag] = makeSliceFieldSetter(i)
		} else {
			supportedOptions[tag] = makeFieldSetter(i)
		}
	}
}

func makeFieldSetter(fieldIndex int) optionSetter {
	return func(c *SSHConfig, v string) {
		reflect.ValueOf(c).Elem().Field(fieldIndex).SetString(v)
	}
}

func makeSliceFieldSetter(fieldIndex int) optionSetter {
	return func(c *SSHConfig, v string) {
		field := reflect.ValueOf(c).Elem().Field(fieldIndex)
		field.Set(reflect.Append(field, reflect.ValueOf(v)))
	}
}

func SetOption(c *SSHConfig, option string, value string) {
	if setter, ok := supportedOptions[option]; ok {
		setter(c, value)
	}
}

const (
	DirectSSH     = "──────────────"
	MissingConfig = "❗"
)

func (c *SSHConfig) IsDirectSSH() bool {
	return c.Name == "" || c.Name == DirectSSH
}

func (c *SSHConfig) UniqueKey() string {
	if !c.IsDirectSSH() {
		return c.Name
	}
	return fmt.Sprintf("%s%s%s", c.Host, c.Port, c.User)
}

func (c *SSHConfig) CleanName() {
	if c.Name == DirectSSH {
		c.Name = ""
	}
	c.Name = strings.TrimPrefix(c.Name, MissingConfig)
}

func Parse(configFile string) ([]SSHConfig, error) {
	return ParseWithSearch("", configFile)
}

func ParseWithSearch(search string, configFile string) ([]SSHConfig, error) {
	configsStrings := strings.Split(strings.ReplaceAll(configFile, "\r\n", "\n"), "Host ")
	var configs = make([]SSHConfig, 0)

	for _, config := range configsStrings {
		lines := strings.Split(config, "\n")

		if strings.Trim(lines[0], " ") == "" {
			continue
		}

		sshConfig := SSHConfig{
			Name: lines[0],
			Port: "",
			User: "",
		}

		for _, line := range lines {
			if len(line) == 0 || line[0] == '#' {
				continue
			}

			line = strings.ReplaceAll(strings.TrimLeft(line, " \t"), "\t", " ")
			m1 := regexp.MustCompile(` ( *)`)
			line = m1.ReplaceAllString(line, " ")

			lineData := strings.Split(line, " ")
			option := lineData[0]
			value := ""
			if len(lineData) > 1 {
				value = lineData[1]
			}
			if strings.Contains(line, "Include") {
				result, err := ParseInclude(search, value)
				if err != nil {
					panic(err)
				}
				configs = append(configs, result...)
				continue
			}

			if setter, ok := supportedOptions[option]; ok {
				setter(&sshConfig, value)
			}
		}

		if sshConfig.Host == "" || !strings.Contains(sshConfig.Name, search) {
			continue
		}

		configs = append(configs, sshConfig)

	}

	return configs, nil
}

func ParseInclude(search string, path string) ([]SSHConfig, error) {
	var results = make([]SSHConfig, 0)

	var isAbsolute = path[0] == '/' || path[0] == '~'

	var paths []string
	var err error

	if isAbsolute {
		if path[0] == '~' {
			path = filepath.Join(HomeDir(), path[2:])
		}
	} else {
		path = filepath.Join(GetSshDir(), path)
	}

	paths, err = filepath.Glob(path)

	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		info, err := os.Stat(path)

		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			continue
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		items, err := ParseWithSearch(search, string(fileContent))
		if err != nil {
			return nil, err
		}
		results = append(results, items...)
	}

	return results, nil
}

func Print() {
	list, err := Parse(GetConfigFile())

	if err != nil {
		log.Fatal(err)
	}

	if len(list) == 0 {
		fmt.Println("No configs found in ~/.ssh/config.")
		return
	}

	var rows []table.Row
	for _, history := range list {
		rows = append(rows, table.Row{history.Name, history.Host, history.Port, history.User, history.Key})
	}
	fmt.Println(theme.PrintTable(rows, theme.PrintConfig))

}
