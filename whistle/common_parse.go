package whistle

type DataBody struct {
	Message           string `json:"Message"`
        Time              string `json:"Time"`
	Host              string `json:"Host"`
	Category          string `json:"Category"`
}
