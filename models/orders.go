package models

//Orders ...
type Orders struct {
	ID        int64 `gorm:"primary_key" json:"id"`
	UserID    int64 `json:"user_id"`
	CartID    int64 `json:"cart_id"`
	CreatedAt int64 `json:"created_at"`
}
