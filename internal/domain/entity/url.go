package entity

type URL struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Feeds []Feed `json:"feeds"`
}
