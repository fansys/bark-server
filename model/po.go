package model

type Device struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	DeviceKey   string `gorm:"index,unique" json:"device_key"`
	DeviceToken string `gorm:"index,unique" json:"device_token"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	BaseModel
}

type Message struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Source      string `gorm:"index" json:"source"`
	DeviceKey   string `gorm:"index,unique" json:"device_key"`
	DeviceToken string `gorm:"index,unique" json:"device_token"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Sound       string `json:"sound"`
	ExtParams   string `json:"ext_params"`
	Pushed      bool   `gorm:"default:false" json:"-"`
	ErrorMsg    string `json:"-"`
	BaseModel
}

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"index,unique" json:"username"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	OpenId    string `gorm:"index,unique" json:"open_id"`
	UniqueKey string `gorm:"index" json:"unique_key"`
	Enabled   bool   `gorm:"default:true" json:"enabled"`
	BaseModel
}

type UserBind struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UniqueKey string `gorm:"index" json:"unique_key"`
	BindType  string `json:"bind_type"`
	Key       string `json:"key"`
	Enabled   bool   `gorm:"default:true" json:"enabled"`
	Blocked   bool   `gorm:"default:false" json:"blocked"`
	BaseModel
}
