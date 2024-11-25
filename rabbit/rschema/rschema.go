package rschema

// PriceChangeNotification represents the message structure for price change notifications.
type PriceChangeNotification struct {
	ArticleId string  `json:"articleId"`
	Price     float32 `json:"price"`
}