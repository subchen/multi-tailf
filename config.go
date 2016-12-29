package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/subchen/gstack/stringutil"
	"github.com/ttacon/chalk"
)

type (
	Config struct {
		SshPass        string   `json:"ssh-pass"`
		SshKey         string   `json:"ssh-key"`
		SshFile        string   `json:"ssh-file"`
		UrlList        []string `json:"files"`
		serverFileList []*ServerFile
	}

	ServerFile struct {
		host        string
		port        string
		username    string
		password    string
		key         string
		file        string
		promptsHost string
		promptsFile string
	}
)

func (c *Config) init() {
	// create server file list
	for _, url := range c.UrlList {
		f := c.newServerFile(url)
		c.serverFileList = append(c.serverFileList, f)
	}

	// compute prompts max width
	maxHostLen := 0
	maxFileLen := 0
	for _, f := range c.serverFileList {
		f.promptsHost = f.host
		if f.port != "" && f.port != "22" {
			f.promptsHost = f.promptsHost + ":" + f.port
		}
		if len(f.promptsHost) > maxHostLen {
			maxHostLen = len(f.promptsHost)
		}

		f.promptsFile = filepath.Base(f.file)
		if len(f.promptsFile) > maxFileLen {
			maxFileLen = len(f.promptsFile)
		}
	}

	// set prompts for each serverfile
	for _, f := range c.serverFileList {
		f.promptsHost = stringutil.RightPad(f.promptsHost, " ", maxHostLen)
		f.promptsFile = stringutil.RightPad(f.promptsFile, " ", maxFileLen)
	}
}

func (c *Config) newServerFile(url string) *ServerFile {
	if !strings.ContainsAny(url, "@:") {
		return &ServerFile{
			host: "localhost",
			file: url,
		}
	}

	RE_USER_PASS := "(\\w+)(:(\\w+))?"
	RE_HOST_PORT := "@([^:@]+)(:([0-9]+))?"
	RE_FILE := "(:([^:]+))?"

	re := regexp.MustCompile("^" + RE_USER_PASS + RE_HOST_PORT + RE_FILE + "$")
	matched := re.FindAllStringSubmatch(url, -1)
	if len(matched) != 1 {
		panic("invalid format: " + url)
	}

	f := ServerFile{
		host:     matched[0][4],
		port:     matched[0][6],
		username: matched[0][1],
		password: matched[0][3],
		file:     matched[0][8],
	}

	if f.password == "" {
		f.password = c.SshPass
	}
	if f.password == "" {
		f.key = c.SshKey
	}
	if f.file == "" {
		f.file = c.SshFile
	}

	return &f
}

func (f *ServerFile) MakeCommand() *Command {
	prompts := fmt.Sprintf("%s %s -> ", chalk.Green.Color(f.promptsHost), chalk.Cyan.Color(f.promptsFile))

	if f.host == "localhost" {
		return &Command{
			prompts: prompts,
			cmd:     exec.Command("tail", "-f", f.file),
		}
	}

	command := "ssh"
	parameters := []string{}
	if f.password != "" {
		command = "sshpass"
		parameters = append(parameters, "-p", f.password, "ssh")
	} else if f.key != "" {
		parameters = append(parameters, "-i", f.key)
	}
	if f.port != "" && f.port != "22" {
		parameters = append(parameters, "-p", f.port)
	}
	parameters = append(parameters, "-o", "LogLevel=quiet")
	parameters = append(parameters, "-o", "UserKnownHostsFile=/dev/null")
	parameters = append(parameters, "-o", "StrictHostKeyChecking=no")
	parameters = append(parameters, f.username+"@"+f.host)
	parameters = append(parameters, "tail -f "+f.file)

	return &Command{
		prompts: prompts,
		cmd:     exec.Command(command, parameters...),
	}
}
