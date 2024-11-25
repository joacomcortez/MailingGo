package rest

import (
	"mailinggo/mailoffer"
	"mailinggo/rabbit/rschema"
)

// @Summary Send emails
// @Description Send emails to all users with a list of articles
// @Tags emails
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully sent emails"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /send-emails [post]

func PostPriceUpdate(updatedProduct rschema.PriceChangeNotification) error {
	// Pass the token to SendPriceNotification
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNjZlNzQ2MWY3MDhmZTdlMGIyMmRlZDMzIiwidXNlcklEIjoiNjZkYzUzY2M4ZGFmNDEwODc5NWY5MWZjIn0.6n7UvEykHQuw5rfJSg7_98CQfFOpMscirptpnc8qkaU" 
	err := mailoffer.SendPriceNotification(updatedProduct.ArticleId, token)
	if err != nil {
		return err
	}

	return nil
}
