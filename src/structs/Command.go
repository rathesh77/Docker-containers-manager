package structs

type Command struct {
	Contract string `json:"contract"`
	Deployment
	Service
}
