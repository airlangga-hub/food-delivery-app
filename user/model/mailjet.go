package model

type MailjetRequest struct {
	Messages []MessageRequest `json:"Messages"`
}

type MessageRequest struct {
	From     Person   `json:"From"`
	To       []Person `json:"To"`
	Subject  string   `json:"Subject"`
	TextPart string   `json:"TextPart"`
	HTMLPart string   `json:"HTMLPart"`
}

type Person struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

type MailjetResponse struct {
	Messages     []MessageResponse `json:"Messages"`
	ErrorMessage string            `json:"ErrorMessage"`
	StatusCode   int               `json:"StatusCode"`
}

type MessageResponse struct {
	Status   string `json:"Status"`
	CustomID string `json:"CustomID"`
	To       []To   `json:"To"`
	Cc       []any  `json:"Cc"`
	Bcc      []any  `json:"Bcc"`
}

type To struct {
	Email       string `json:"Email"`
	MessageUUID string `json:"MessageUUID"`
	MessageID   int64  `json:"MessageID"`
	MessageHref string `json:"MessageHref"`
}