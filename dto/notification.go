package dto

type Notification struct {
	PushToken string `json:"push_token"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Os        string `json:"os"`
}
