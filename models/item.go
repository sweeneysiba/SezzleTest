package models

//Item ...
type Item struct {
	ID        int64  `gorm:"primary_key" json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}
