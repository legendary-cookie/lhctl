package main

import "log"
import "fmt"
import "github.com/jedib0t/go-pretty/v6/table"
import "os"
import "strconv"
import "encoding/json"

type VolumeListEntry struct {
	Name       string `json:"name"`
	Size       string `json:"size"`
	ActualSize string `json:"actualSize"`
	Ready      bool   `json:"ready"`
	State      string `json:"state"`
}


func ListVolumesCmd(output string) {
	volumes := GetVolumes().Data
	switch output {
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
        case "json":
          var entries []VolumeListEntry
          for _, volume := range volumes {
                        entry := VolumeListEntry{
                          Name: volume.Name,
                          Size: volume.Size,
                          ActualSize: volume.Controllers[0].ActualSize,
                          Ready: volume.Ready,
                          State: volume.State,
                        }
                        entries = append(entries, entry)
          }
          jsonThing, err := json.Marshal(entries)
          if err != nil {
            log.Fatal(err)
          }
          println(string(jsonThing))
	default:
		fmt.Println("That output format doesn't exist!\nValid formats are: [table/json/csv]")
	}
}
func VolumeCreateCmd(name, size string) {
	CreateVolume(name, size)
}

func VolumeDeleteCmd(name string) {
	c := askForConfirmation("Do you really want to delete this volume?")
	if c {
		DeleteVolume(name)
	}
}
func VolumePvcCmd(volume, ns string) {
	CreatePv(volume, volume)
	CreatePvc(volume, volume, ns)
}
