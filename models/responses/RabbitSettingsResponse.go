package responses

type RabbitSettings struct {
	Host   string `json:"host"`
	Password string `json:"password"`
	UserName string `json:"username"`
	VirtualHost string `json:"virtualhost"`
}