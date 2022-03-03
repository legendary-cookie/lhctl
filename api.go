package main

import (
	"github.com/monaco-io/request"
	"github.com/monaco-io/request/response"
	"log"
)

func sendGet(path string) *response.Sugar {
	client := request.Client{
		URL:    conf.ApiUrl + path,
		Method: "GET",
	}
	resp := client.Send()
	if !resp.OK() {
		log.Fatal(resp.Error())
	}
	return resp
}

func sendPost(path string, body interface{}) {
	client := request.Client{
		URL:    conf.ApiUrl + path,
		Method: "POST",
		JSON:   body,
	}
	resp := client.Send()
	if !resp.OK() {
		log.Fatal(resp.Error())
	}
}

func sendDelete(path string) {
	client := request.Client{
		URL:    conf.ApiUrl + path,
		Method: "DELETE",
	}
	resp := client.Send()
	if !resp.OK() {
		log.Fatal(resp.Error())
	}
}

func GetVolumes() VolumeGetResponse {
	var result VolumeGetResponse
	sendGet("volumes").Scan(&result)
	return result
}

func CreateVolume(name, size string) {
	body := VolumeCreateRequest{
		Name:                    name,
		Size:                    size,
		NumberOfReplicas:        2,
		Frontend:                "blockdev",
		ReplicaAutoBalance:      "ignored",
		RevisionCounterDisabled: false,
		StaleReplicaTimeout:     20,
		Encrypted:               false,
		DataLocality:            "disabled",
		BackingImage:            "",
		AccessMode:              "rwo",
	}
	sendPost("volumes", body)
}

func DeleteVolume(name string) {
	sendDelete("volumes/" + name)
}

func CreateSnapshot(volume, name string) {
	var body = struct {
		Name string `json:"name"`
	}{Name: name}
	sendPost("volumes/"+volume, body)
}

func CreatePv(volume, pvname string) {
	var body = struct {
		PvName string `json:"pvName"`
		FsType string `json:"fsType"`
	}{PvName: pvname, FsType: "ext4"}
	sendPost("volumes/"+volume+"?action=pvCreate", body)
}

func CreatePvc(volume, pvcname, namespace string) {
	var body = struct {
		PvcName   string `json:"pvcName"`
		Namespace string `json:"namespace"`
	}{PvcName: pvcname, Namespace: namespace}
	sendPost("volumes/"+volume+"?action=pvcCreate", body)
}
