package utilities

import(
	"urlshortener-crud-consumer/models/entities"
)
func UpdateUrlShortenerModel(urlShortener entities.UrlShortener, channelCampaignId int, partialId int, longUrl string, telephoneNumber string, userId int) entities.UrlShortener {
	urlShortener.ChannelCampaignId = channelCampaignId
	urlShortener.PartialId = partialId
	urlShortener.LongUrl = longUrl
	urlShortener.TelephoneNumber = telephoneNumber
	urlShortener.CreatedBy = userId
	urlShortener.IsActive = true

	return urlShortener
}