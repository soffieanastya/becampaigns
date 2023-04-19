package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserid(userID int) ([]Campaign, error)
}

// ini di akes di file ini aja
type repository struct {
	db *gorm.DB
}

// supaya bisa diakses di luar pckace
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindByUserid(userID int) ([]Campaign, error) {
	var campaigns []Campaign
	// preload campaign images, untuk load relasi dari campaigns, yaitu campaign_images
	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
