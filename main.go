package main

import (
	_ "embed"
	"errors"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"os/user"
)

//go:embed assets/default.toml
var defaultconf string

//go:embed assets/help.txt
var help string

var conf Config

type Config struct {
	ApiUrl string
	Auth   bool
	Color  bool
	Debug  bool
}

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir := user.HomeDir
	confPath := homeDir + "/.config/longhorn.toml"
	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(confPath, []byte(defaultconf), 0666); err != nil {
			log.Fatal(err)
		} else {
			println("Wrote default config to " + confPath)
			println("Exiting!")
			os.Exit(0)
		}
	}
	_, err = toml.DecodeFile(confPath, &conf)
	if err != nil {
		log.Fatal(err)
	}
	Run(os.Args[1:])
}

func HelpCommand() {
	print(help)
}

func Run(args []string) {
	if len(args) != 0 {
		switch args[0] {
		case "volume":
			VolumeCommand(args[1:])
		default:
			HelpCommand()
		}
	} else {
		HelpCommand()
	}
}
