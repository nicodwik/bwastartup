package campaign

import (
	"bwastartup/user"
	"time"

	"github.com/leekchan/accounting"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User             user.User
}

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  bool
}

func (c Campaign) GoalAmountIDRFormatter() string {
	ac := accounting.Accounting{Symbol: "Rp. ", Thousand: "."}
	return ac.FormatMoney(c.GoalAmount)
}
