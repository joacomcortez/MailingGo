package rest

import (
	"mailinggo/mailoffer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Send emails
// @Description Send emails to all users with a list of articles
// @Tags emails
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully sent emails"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /send-emails [post]

func PostCartNotification(c *gin.Context) {
	err := mailoffer.SendCartNotification(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send cart notifications", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emails sent successfully"})
}
