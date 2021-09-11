package orm

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Device struct {
	Model
	DeviceKey   string `gorm:"unique,index" json:"device_key,omitempty"`
	DeviceToken string `gorm:"unique,index" json:"device_token,omitempty"`
	DeviceType  string `gorm:"default:ios" json:"device_type,omitempty"`
}

type User struct {
	Model
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Ticket   string `json:"ticket,omitempty"`
}

type PushMessage struct {
	Model
	DeviceToken string `form:"-" json:"-" xml:"-" query:"-"`
	DeviceKey   string `form:"device_key,omitempty" json:"device_key,omitempty" xml:"device_key,omitempty" query:"device_key,omitempty"`
	DeviceType  string `form:"device_type,omitempty" json:"device_type,omitempty" xml:"device_type,omitempty" query:"device_type,omitempty"`
	Category    string `form:"category,omitempty" json:"category,omitempty" xml:"category,omitempty" query:"category,omitempty"`
	Title       string `form:"title,omitempty" json:"title,omitempty" xml:"title,omitempty" query:"title,omitempty"`
	Body        string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty" query:"body,omitempty"`
	// ios notification sound(system sound please refer to http://iphonedevwiki.net/index.php/AudioServices)
	Sound     string                 `form:"sound,omitempty" json:"sound,omitempty" xml:"sound,omitempty" query:"sound,omitempty"`
	Group     string                 `form:"group,omitempty" json:"group,omitempty" xml:"group,omitempty" query:"group,omitempty"`
	ExtParams map[string]interface{} `gorm:"-" form:"ext_params,omitempty" json:"ext_params,omitempty" xml:"ext_params,omitempty" query:"ext_params,omitempty"`
}
