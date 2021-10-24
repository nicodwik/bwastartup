package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"errors"
)

type Service interface {
	GetTransactionsByCampaignID(input GetTransactionsDetailInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, payment payment.Service) *service {
	return &service{repository, campaignRepository, payment}
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
func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	var transaction Transaction

	transaction.Amount = input.Amount
	transaction.CampaignID = input.CampaignID
	transaction.UserID = input.User.Id
	transaction.Status = "pending"

	savedTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return savedTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     savedTransaction.ID,
		Amount: savedTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return savedTransaction, err
	}

	savedTransaction.PaymentUrl = paymentURL
	savedTransaction, err = s.repository.Update(savedTransaction)
	if err != nil {
		return savedTransaction, err
	}

	return savedTransaction, nil
}
