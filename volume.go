package main

import "log"
import "fmt"
import "github.com/jedib0t/go-pretty/v6/table"
import "os"
import "strconv"
import "github.com/akamensky/argparse"

func volumeHelp() {
	fmt.Println("Available subcommands for <volume> are:\n- list\n- create <name> <size>\n- delete <name>\n- snapshot <volume> <snapshotname>\n- createpvc <volume> <namespace>")
}

func VolumeCommand(args []string) {
	if len(args) == 0 {
		volumeHelp()
	} else {
		switch args[0] {
		case "list":
			parser := argparse.NewParser(os.Args[0]+" volume list", "")
			o := parser.String("o", "output", &argparse.Options{Help: "Specify output format", Default: "table"})
			err := parser.Parse(args)
			if err != nil {
				// In case of error print error and print usage
				// This can also be done by passing -h or --help flags
				fmt.Print(parser.Usage(err))
				break
			}
			volumes := GetVolumes().Data
			switch *o {
			case "table":
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
			case "csv":
				for _, volume := range volumes {
					size, err := strconv.ParseInt(volume.Size, 10, 64)
					if err != nil {
						log.Fatal(err)
					}
					actualSize, err := strconv.ParseInt(volume.Controllers[0].ActualSize, 10, 64)
					if err != nil {
						log.Fatal(err)
					}
					print(volume.Name)
					print(",")
					print(size)
					print(",")
					print(actualSize)
					print(",")
					print(volume.Ready)
					print(",")
					print(volume.State)
					print("\n")
				}
			default:
				fmt.Println("That output format doesn't exist!\nValid formats are: [table/json/csv]")
			}
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
		case "createpvc":
			if len(args) < 3 {
				print(colorRed)
				fmt.Println("You have to specify the volume and the namespace!")
				print(colorReset)
				break
			}
			volume := args[1]
			ns := args[2]
			CreatePv(volume, volume)
			CreatePvc(volume, volume, ns)
		default:
			volumeHelp()
		}
	}
}
