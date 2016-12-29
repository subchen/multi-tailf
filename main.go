package main

import (
	"fmt"
	"runtime"

	"github.com/subchen/gstack/cli"
	"github.com/ttacon/chalk"
)

const (
	VERSION = "1.0.0"
)

var (
	BuildVersion   string
	BuildGitCommit string
	BuildDate      string

	APP_LOGO = chalk.Magenta.Color(`
    ###    ### ########  #####  ## ##      #######
    ####  ####    ##    ##   ## ## ##      ##
    ## #### ##    ##    ####### ## ##      #####
    ##  ##  ##    ##    ##   ## ## ##      ##
    ##      ##    ##    ##   ## ## ####### ##
	`)
)

func main() {
	app := cli.NewApp("mtailf", "mtailf is equivalent of tail -f on multiple local or remote files at once")

	app.Flag("--ssh-pass", "default password for ssh")
	app.Flag("--ssh-key", "default key file for ssh").Default("~/.ssh/id_rsa")
	app.Flag("--ssh-file", "default file for ssh tail")
	app.AllowArgumentCount(1, -1)

	if BuildVersion == "" {
		app.Version = VERSION
	} else {
		app.Version = func() {
			fmt.Printf("Version: %s-%s\n", VERSION, BuildVersion)
			fmt.Printf("Go version: %s\n", runtime.Version())
			fmt.Printf("Git commit: %s\n", BuildGitCommit)
			fmt.Printf("Built: %s\n", BuildDate)
			fmt.Printf("OS/Arch: %s-%s\n", runtime.GOOS, runtime.GOARCH)
		}
	}

	app.Usage = func() {
		fmt.Println("Usage: mtailf [OPTIONS] [ /path/file | user:pass@host:file ] ...")
		fmt.Println("   or: mtailf [ --version | --help ]")
	}

	app.MoreHelp = func() {
		fmt.Println("Examples:")
		fmt.Println("* Local files")
		fmt.Println("    mtailf /var/log/messages-1 /var/log/messages-2")
		fmt.Println("* Multiple files on servers")
		fmt.Println("    mtailf root@10.0.0.1 root@10.0.0.2 --ssh-pass=password --ssh-file=/var/log/messages")
		fmt.Println("* Use SSH private key")
		fmt.Println("    mtailf root@10.0.0.1:/var/log/messages --ssh-key=/tmp/ssh.key")
		fmt.Println("* Use SSH passwords")
		fmt.Println("    mtailf root:p1@10.0.0.1:/var/log/messages root:p2@10.0.0.2:/var/log/messages")
		fmt.Println("* Use SSH port")
		fmt.Println("    mtailf root@10.0.0.1:8022:/var/log/messages")
	}

	app.Execute = executeApp

	app.Run()
}

func executeApp(ctx *cli.Context) {
	fmt.Println(APP_LOGO)

	cfg := Config{
		SshPass: ctx.String("--ssh-pass"),
		SshKey:  ctx.String("--ssh-key"),
		SshFile: ctx.String("--ssh-file"),
		UrlList: ctx.Args(),
	}

	cfg.init()

	commandList := &CommandList{}
	for _, f := range cfg.serverFileList {
		commandList.AppendCommand(f.MakeCommand())
	}
	commandList.Start()
}
