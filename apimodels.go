package main

type VolumeGetResponse struct {
	Data []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Ready       bool   `json:"ready"`
		Size        string `json:"size"`
		State       string `json:"state"`
		Controllers []struct {
			ActualSize string `json:"actualSize"`
		} `json:"controllers"`
	} `json:"data"`
}

type VolumeCreateRequest struct {
	Name                    string   `json:"name"`
	Size                    string   `json:"size"`
	NumberOfReplicas        int      `json:"numberOfReplicas"`
	Frontend                string   `json:"frontend"`
	DataLocality            string   `json:"dataLocality"`
	AccessMode              string   `json:"accessMode"`
	BackingImage            string   `json:"backingImage"`
	ReplicaAutoBalance      string   `json:"replicaAutoBalance"`
	RevisionCounterDisabled bool     `json:"revisionCounterDisabled"`
	Encrypted               bool     `json:"encrypted"`
	NodeSelector            []string `json:"nodeSelector"`
	DiskSelector            []string `json:"diskSelector"`
	StaleReplicaTimeout     int      `json:"staleReplicaTimeout"`
	FromBackup              string   `json:"fromBackup"`
}
