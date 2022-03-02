package main

import "log"
import "fmt"
import "github.com/jedib0t/go-pretty/v6/table"
import "os"
import "strconv"

func volumeHelp() {
	fmt.Println("Available subcommands for <volume> are:\n- list\n- create <name> <size>\n- delete <name>\n- snapshot <volume> <snapshotname>")
}

func VolumeCommand(args []string) {
	if len(args) == 0 {
		volumeHelp()
	} else {
		switch args[0] {
		case "list":
			volumes := GetVolumes().Data
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Name", "Size", "Actual Size", "Ready", "State"})
			for _, volume := range volumes {
				size, err := strconv.ParseInt(volume.Size, 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				actualSize, err := strconv.ParseInt(volume.Controllers[0].ActualSize, 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				t.AppendRow([]interface{}{volume.Name, ByteCountIEC(size), ByteCountIEC(actualSize), volume.Ready, volume.State})
			}
			t.Render()
		case "create":
			if len(args) < 3 {
				print(colorRed)
				fmt.Println("You have to specify name and size!")
				print(colorReset)
				break
			}
			name := args[1]
			size := args[2]
			CreateVolume(name, size)
		case "delete":
			if len(args) < 2 {
				print(colorRed)
				fmt.Println("You have to specify the name of the volume you want to delete!")
				print(colorReset)
				break
			}
			name := args[1]
			c := askForConfirmation("Do you really want to delete this volume?")
			if c {
				DeleteVolume(name)
			}
		case "snapshot":
			if len(args) < 3 {
				print(colorRed)
				fmt.Println("You have to specify the volume and the name of the snapshot!")
				print(colorReset)
				break
			}
			volume := args[1]
			name := args[2]
			CreateSnapshot(volume, name)
		default:
			volumeHelp()
		}
	}
}
