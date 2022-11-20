package config

// 微信小程序配置
type WxConf struct {
	AppId      string `json:"AppId"`  //微信appId
	Secret     string `json:"Secret"` //微信secret
	Grant_type string `json:"Grant_type"`
}
