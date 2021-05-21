package model

type PushMessage struct {
	Source      string `form:"source,omitempty" json:"source,omitempty" xml:"source,omitempty" query:"source,omitempty"`
	DeviceToken string `form:"-" json:"-" xml:"-" query:"-"`
	DeviceKey   string `form:"device_key,omitempty" json:"device_key,omitempty" xml:"device_key,omitempty" query:"device_key,omitempty"`
	Category    string `form:"category,omitempty" json:"category,omitempty" xml:"category,omitempty" query:"category,omitempty"`
	Title       string `form:"title,omitempty" json:"title,omitempty" xml:"title,omitempty" query:"title,omitempty"`
	Body        string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty" query:"body,omitempty"`
	// ios notification sound(system sound please refer to http://iphonedevwiki.net/index.php/AudioServices)
	Sound     string                 `form:"sound,omitempty" json:"sound,omitempty" xml:"sound,omitempty" query:"sound,omitempty"`
	ExtParams map[string]interface{} `form:"ext_params,omitempty" json:"ext_params,omitempty" xml:"ext_params,omitempty" query:"ext_params,omitempty"`
}

type Page struct {
	PageNo    int `form:"pageNo,1" json:"pageNo,1" xml:"pageNo,1" query:"pageNo,1"`
	PageSize  int `form:"pageSize,10" json:"pageSize,10" xml:"pageSize,10" query:"pageSize,10"`
	Total     int
	TotalPage int
	Data      interface{}
}
