package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB

}

type Repository interface {
	GetCampaignID(campaignID int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// get transaction per-campaign
func (r *repository) GetCampaignID(campaignID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

// get transactions per-user yg login
func (r *repository) GetByUserID(userID int) ([]Transaction, error) {
	var transactions []Transaction
	// join 3 tabel/join banyak tabel
	err := r.db.Preload("Campaign.CampaignImages","campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}