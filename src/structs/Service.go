package structs

type Service struct {
	Name        string   `json:"name"`
	PodSelector string   `json:"pod-selector"`
	Port        string   `json:"port"`
	Pods        []string `json:"pods"`
}
