package models

//User ...
type User struct {
	ID        int64  `gorm:"primary_key;unique;not null"  json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	UserName  string `gorm:"unique;not null" json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Token     string `gorm:"type:varchar(1000)" json:"token,omitempty"`
	CartID    int64  `gorm:"unique;not null" json:"cart_id,omitempty"`
	CreatedAt int64  ` gorm:"autoUpdateTime:milli" json:"created_at,omitempty"`
}
