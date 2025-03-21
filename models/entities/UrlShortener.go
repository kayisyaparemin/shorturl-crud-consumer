package entities

type UrlShortener struct {
	Id                int    `gorm:"primaryKey"`
	TelephoneNumber   string `gorm:"column:telephonenumber"`
	Key               string `gorm:"column:key"`
	LongUrl           string `gorm:"column:longurl"`
	PartialId         int    `gorm:"column:partialid"`
	ChannelCampaignId int    `gorm:"column:channelcampaignid"`
	IsActive          bool   `gorm:"column:isactive"`
	CreatedBy         int    `gorm:"column:createdby"`
}
