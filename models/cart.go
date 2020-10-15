package models

//Cart ...
type Cart struct {
	ID          int64       `gorm:"primary_key" json:"id"`
	UserID      int64       `gorm:"not null" json:"user_id"`
	ISPurchased bool        `json:"is_purchased"`
	CartItems   []*CartItem `json:"cart_items"`
	CreatedAt   int64       `json:"created_at"`
}

//CartItem ...
type CartItem struct {
	CartID int64 `json:"cart_id"`
	ItemID int64 `json:"item_id"`
}
