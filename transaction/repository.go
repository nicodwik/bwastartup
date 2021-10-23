package transaction

import "gorm.io/gorm"

type Repository interface {
	GetCampaignByID(campaignID int) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetCampaignByID(campaignID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Where("campaign_id = ?", campaignID).Preload("User").Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
