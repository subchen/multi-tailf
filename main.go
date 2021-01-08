package main

import (
	"fmt"
	"github.com/subchen/go-cli"
	"github.com/ttacon/chalk"
	"os"
)

const (
	VERSION = "1.0.0"
)

var (
	BuildVersion     string
	BuildDate        string
	BuildGitBranch   string
	BuildGitCommit   string
	BuildGitRevCount string

	AppLogo = chalk.Magenta.Color(`
    ###    ### ########  #####  ## ##      #######
    ####  ####    ##    ##   ## ## ##      ##
    ## #### ##    ##    ####### ## ##      #####
    ##  ##  ##    ##    ##   ## ## ##      ##
    ##      ##    ##    ##   ## ## ####### ##
	`)
)

func main() {
	app := cli.NewApp()
	app.Name = "mtailf"
	app.Version = VERSION
	app.Usage = "mtailf is equivalent of tail -f on multiple local or remote files at once"
	app.UsageText = "Usage: mtailf [OPTIONS] [ /path/file | user:pass@host:file ] ...\n" +
		"   or: mtailf [ --version | --help ]"
	app.Flags = []*cli.Flag{
		{
			Name:  "--ssh-pass",
			Usage: "default password for ssh",
		},
		{
			Name:  "--ssh-key",
			Usage: "default key file for ssh",
		},
		{
			Name:        "--ssh-file",
			Usage:       "default file for ssh tail",
			Placeholder: "~/.ssh/id_rsa",
		},
	}
	app.Examples = "* Local files\n" +
		"    mtailf /var/log/messages-1 /var/log/messages-2\n" +
		"* Multiple files on servers\n" +
		"    mtailf root@10.0.0.1 root@10.0.0.2 --ssh-pass=password --ssh-file=/var/log/messages\n" +
		"* Use SSH private key\n" +
		"    mtailf root@10.0.0.1:/var/log/messages --ssh-key=/tmp/ssh.key\n" +
		"* Use SSH passwords\n" +
		"    mtailf root:p1@10.0.0.1:/var/log/messages root:p2@10.0.0.2:/var/log/messages\n" +
		"* Use SSH port\n" +
		"    mtailf root@10.0.0.1:8022:/var/log/messages"

	if BuildVersion != "" {
		app.BuildInfo = &cli.BuildInfo{
			Timestamp:   BuildDate,
			GitBranch:   BuildGitBranch,
			GitCommit:   BuildGitCommit,
			GitRevCount: BuildGitRevCount,
		}
	}

	app.Action = executeApp

	app.Run(os.Args)
}

func executeApp(ctx *cli.Context) {
	fmt.Println(AppLogo)

	cfg := Config{
		SshPass: ctx.GetString("--ssh-pass"),
		SshKey:  ctx.GetString("--ssh-key"),
		SshFile: ctx.GetString("--ssh-file"),
		UrlList: ctx.Args(),
	}

	cfg.init()

	commandList := &CommandList{}
	for _, f := range cfg.serverFileList {
		commandList.AppendCommand(f.MakeCommand())
	}
	commandList.Start()
}
