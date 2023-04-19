package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// tamilkan semua campaign berdasarkan usernya dan kalau gaada usernya
func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.FindByUserid(userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	// kalau uer ga login, panggil semua campaign
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

// get one campaign
func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

// create campaign
func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}

	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID

	// create slug
	// you must install slug
	// go get -u github.com/gosimple/slug

	strSlug := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(strSlug) // tambahid user agar menghindari nama sug yg sama tiap user beda bikin jududlnya


	// panggil repo
	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

// update campaigns
// dapetin campaignnya dulu, tangkap parameter
func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	dataCampaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return dataCampaign, err
	}
	// bisa di update kalau campaign itu buatan user tsb
	if dataCampaign.UserID != inputData.User.ID {
		return dataCampaign, errors.New("Not an owner of the campaign")
	}

	// masukin sesuai inputan
	dataCampaign.Name = inputData.Name
	dataCampaign.ShortDescription = inputData.ShortDescription
	dataCampaign.Description = inputData.Description
	dataCampaign.Perks = inputData.Perks
	dataCampaign.GoalAmount = inputData.GoalAmount

	// simpan ke db
	update, err := s.repository.Update(dataCampaign)
	if err != nil {
		return update, err
	}
	return update, nil
}

// save image
func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) { 
	// cek user yg ngedit sama ga kaya user yang upload/buat post?
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}
	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("Not an owner of the campaign")
	}

	isprimary := 0
	if input.IsPrimary {
		isprimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	} 

	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isprimary
	campaignImage.FileName = fileLocation

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}