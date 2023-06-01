package structs

type Deployment struct {
	PodLabel    string `json:"pod-label"`
	DockerImage string `json:"docker-image"`
	Args        string `json:"args"`
}
