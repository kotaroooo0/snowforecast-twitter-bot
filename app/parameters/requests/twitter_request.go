package requests

// GET Request (CRC Check)
type GetTwitterWebhookRequest struct {
	CrcToken string `json:"crc_token" form:"crc_token" binding:"required"`
}

func NewGetTwitterWebhookRequest() GetTwitterWebhookRequest {
	return GetTwitterWebhookRequest{}
}

func (r *GetTwitterWebhookRequest) Validate() error {
	return nil
}
