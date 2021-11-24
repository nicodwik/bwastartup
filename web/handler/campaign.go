package handler

import (
	"bwastartup/campaign"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	for _, campaign := range campaigns {
		fmt.Println(campaign)
		fmt.Println("--------------")
	}
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}
