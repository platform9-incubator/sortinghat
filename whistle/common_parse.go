package whistle

type DataBody struct {
	Message           string `json:"Message"`
	Host              string `json:"Host"`
	Category          string `json:"Category"`
	Time              int    `json:"Time"`
	Name              string `json:"Name"`
}
