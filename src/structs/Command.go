package structs

type Command struct {
	Contract string `json:"contract"`
	Args     string `json: "args"`
}
