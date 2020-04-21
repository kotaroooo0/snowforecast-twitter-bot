package responses

type GetTwitterWebhookCrcCheckResponse struct {
	Token string `json:"response_token"`
}

func NewGetTwitterWebhookCrcCheckResponse() GetTwitterWebhookCrcCheckResponse {
	return GetTwitterWebhookCrcCheckResponse{}
}
