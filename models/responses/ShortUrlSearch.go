package responses

type ShortUrlSearch struct {
	Key string `json:"key"`
	LongUrl string `json:"longUrl"`
	TelephoneNumber string `json:"telephoneNumber"`
}