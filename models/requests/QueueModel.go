package requests


type QueueModel struct{
	UserId int `json:"userId"`
	ChannelCampaignId int `json:"channelCampaignId"`
	PartialId int `json:"partialId"`
	TelephoneNumbers []string `json:"telephoneNumbers"`
	Url string `json:"url"`
	Size int `json:"size"`
	Email string `json:"email"`
}