package events

type Events struct {
	EventName   string `json:"event_name"`
	Description string `json:"description"`
	Incharge    string `json:"name"`
	EventDate   string `json:"event_date"`
}