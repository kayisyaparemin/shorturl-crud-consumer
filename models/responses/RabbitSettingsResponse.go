package responses

import(
	"strings"
)
type RabbitSettings struct {
	Host   string `json:"host"`
	Password string `json:"password"`
	UserName string `json:"username"`
	VirtualHost string `json:"virtualhost"`
}


func GetRabbitSettings(s *[]Settings) *RabbitSettings {

	settings := new(RabbitSettings)
	for _, setting := range *s{
		keySplit := strings.Split(setting.Key, ".")
		if(len(keySplit) == 2){
			key := keySplit[1]
			switch key {
			case "host":
				settings.Host = setting.Value
			case "password":
				settings.Password = setting.Value
			case "username":
				settings.UserName = setting.Value
			case "virtualhost":
				settings.VirtualHost = setting.Value
			}
		}
		
	}
	return settings
}