package transaction

import (
	"campaigns/campaign"
	"errors"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTansactionByCampaignID(input GetTransactionsCampaignInput) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTansactionByCampaignID(input GetTransactionsCampaignInput) ([]Transaction, error) {
	// get campaign, check userid yg buat campaign
	getCampaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if getCampaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the campaign")
	}
	
	transaction, err := s.repository.GetCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
