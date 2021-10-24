package payment

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type service struct {
	campaignRepository campaign.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService(campaignRepository campaign.Repository) *service {
	return &service{campaignRepository}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = ""
	midclient.ClientKey = ""
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustDetail{
			FName: user.Name,
			Email: user.Email,
		},
	}

	snapToken, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return snapToken.RedirectURL, err
	}

	return snapToken.RedirectURL, nil
}
