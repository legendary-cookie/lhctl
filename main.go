package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/akamensky/argparse"
	"log"
	"os"
	"os/user"
)

//go:embed assets/default.toml
var defaultconf string

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

func Run(args []string) {
	parser := argparse.NewParser("lhctl", "Longhorn CLI to manage different aspects of longhorn")
	volumeCommand := parser.NewCommand("volume", "Manage and view volumes")
	volumeListCommand := volumeCommand.NewCommand("list", "List all volumes")
        listOutput := volumeListCommand.String("o", "output", &argparse.Options{Default: "table", Help: "Output format, valid are table, csv and json"})
	volumeCreateCommand := volumeCommand.NewCommand("create", "Create a new volume")
        volumeCreateName := volumeCreateCommand.String("n", "name", &argparse.Options{Required: true, Help: "The name of the volume to be created"})
        volumeCreateSize := volumeCreateCommand.String("s", "size", &argparse.Options{Required: true, Help: "The size of the volume to be created"})
	volumeDeleteCommand := volumeCommand.NewCommand("delete", "Delete a volume")
        volumeDeleteName := volumeDeleteCommand.String("n", "name", &argparse.Options{Required: true, Help: "The name of the volume to be deleted"})
	volumePvcCommand := volumeCommand.NewCommand("pvc", "Create a PV & PVC")
        volumePvcName := volumePvcCommand.String("v", "volume", &argparse.Options{Required: true, Help:"The volume for which the PV and PVC should be created"})
        volumePvcNamespace := volumePvcCommand.String("n", "namespace", &argparse.Options{Required: true, Help: "The namespace in which the PV and PVC should be created"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}
	if volumeCommand.Happened() {
		if volumeListCommand.Happened() {
			ListVolumesCmd(*listOutput)
		} else if volumeCreateCommand.Happened() {
			VolumeCreateCmd(*volumeCreateName, *volumeCreateSize)
		} else if volumeDeleteCommand.Happened() {
			VolumeDeleteCmd(*volumeDeleteName)
		} else if volumePvcCommand.Happened() {
			VolumePvcCmd(*volumePvcName, *volumePvcNamespace)
		}
	} else {
		// Should never get hit
		err := fmt.Errorf("bad arguments, please check usage")
		fmt.Print(parser.Usage(err))
	}
}
