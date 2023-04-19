package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserid(userID int) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign) (Campaign, error)
	CreateImage(campaignImage CampaignImage)(CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
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

// ambil sebuah campaign berdsar id campaign
func (r *repository) FindByID(ID int) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}


// create campaign by user
func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

// update campaign by user
func (r *repository) Update(campaign Campaign) (Campaign, error) {
	// updat pake save juga
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

// uploadimage campaign byuser, masuk data ke campaign_images
func (r *repository) CreateImage(campaignImage CampaignImage)(CampaignImage, error) {
	// add photo
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}


// ubah primary dr true jadi false kalau ada dua true di photonya, yg terbaru yg true yg sblmnya false
func (r *repository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	// uodate campaign_images set is_primary = false(0) where campaign_id = 1
	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?",campaignID).Update("is_primary",false).Error
	if err != nil {
		return false, err
	}
	return true, nil
}