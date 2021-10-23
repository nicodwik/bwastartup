package transaction

import (
	"bwastartup/campaign"
	"errors"
)

type Service interface {
	GetTransactionsByCampaignID(input GetTransactionsDetailInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionsByCampaignID(input GetTransactionsDetailInput) ([]Transaction, error) {
	transactions := []Transaction{}

	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return transactions, err
	}

	if campaign.UserID != input.User.Id {
		return transactions, errors.New("Unauthorized")
	}

	transactions, err = s.repository.GetCampaignByID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionsByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
