package repositories

import (
	"log"
	"urlshortener-crud-consumer/models/entities"

	"gorm.io/gorm"
)
type IUrlShortenerRepository interface {
	BulkUpdateAvailableToUsedKeys(entityModels []entities.UrlShortener) error
	GetAvailableKeys(size int) ([]entities.UrlShortener, error)
}
type UrlShortenerRepository struct {
	db *gorm.DB
}
func NewUrlShortenerRepository(db *gorm.DB) *UrlShortenerRepository {
	return &UrlShortenerRepository{db: db}
}
func (usr *UrlShortenerRepository) BulkUpdateAvailableToUsedKeys(entityModels []entities.UrlShortener) error {
	if len(entityModels) == 0 {
		return nil
	}

	tx := usr.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back due to panic: %v", r)
		}
	}()

	updateValues := make([]map[string]interface{}, len(entityModels))

	for i, entity := range entityModels {
		updateValues[i] = map[string]interface{}{
			"id":                entity.Id,
			"telephonenumber":   entity.TelephoneNumber,
			"longurl":           entity.LongUrl,
			"partialid":         entity.PartialId,
			"channelcampaignid": entity.ChannelCampaignId,
			"isactive":          entity.IsActive,
			"createdby":         entity.CreatedBy,
		}
	}
	if err := tx.Model(&entities.UrlShortener{}).
		Where("id IN (?)", getEntityIds(entityModels)).
		Updates(updateValues).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
func (usr *UrlShortenerRepository) GetAvailableKeys(size int) ([]entities.UrlShortener, error) {
	var keys []entities.UrlShortener

	err := usr.db.
		Select("id, key").
		Where("telephonenumber IS NULL").
		Limit(size).
		Find(&keys).Error

	return keys, err
}
func getEntityIds(entities []entities.UrlShortener) []int {
	ids := make([]int, len(entities))
	for i, entity := range entities {
		ids[i] = entity.Id
	}
	return ids
}
