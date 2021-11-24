package handler

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService     user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{campaignService, userService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) Create(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	c.HTML(http.StatusOK, "campaign_create.html", gin.H{"users": users})
}

func (h *campaignHandler) Store(c *gin.Context) {
	var input campaign.FormCreateCampaignInput
	users, _ := h.userService.GetAllUsers()
	input.Users = users

	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusInternalServerError, "campaign_create.html", input)
		return
	}

	user, err := h.userService.GetUserByID(input.UserID)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusInternalServerError, "campaign_create.html", input)
		return
	}

	createCampaign := campaign.CreateCampaignInput{}
	createCampaign.Name = input.Name
	createCampaign.ShortDescription = input.ShortDescription
	createCampaign.Description = input.Description
	createCampaign.GoalAmount = input.GoalAmount
	createCampaign.Perks = input.Perks
	createCampaign.User = user

	_, err = h.campaignService.CreateCampaign(createCampaign)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusInternalServerError, "campaign_create.html", input)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	id := c.Param("id")
	idParam, _ := strconv.Atoi(id)

	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"id": idParam})
}

func (h *campaignHandler) StoreImage(c *gin.Context) {
	id := c.Param("id")
	idParam, _ := strconv.Atoi(id)

	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	path := fmt.Sprintf("images/%d-%s", idParam, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	campaignData, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: idParam})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	CreateCampaignImageInput := campaign.CreateCampaignImageInput{
		CampaignID: idParam,
		IsPrimary:  true,
		User:       campaignData.User,
	}

	_, err = h.campaignService.SaveCampaignImage(CreateCampaignImageInput, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Edit(c *gin.Context) {
	id := c.Param("id")
	idParam, _ := strconv.Atoi(id)

	campaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: idParam})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	c.HTML(http.StatusOK, "campaign_edit.html", gin.H{"campaign": campaign})
}

func (h *campaignHandler) Update(c *gin.Context) {
	id := c.Param("id")
	idParam, _ := strconv.Atoi(id)

	var input campaign.FormUpdateCampaignInput

	err := c.ShouldBind(&input)
	if err != nil {
		input.ID = idParam
		c.HTML(http.StatusInternalServerError, "campaign_edit.html", gin.H{"error": err, "campaign": input})
	}

	campaignData, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: idParam})
	if err != nil {
		input.ID = idParam
		c.HTML(http.StatusInternalServerError, "campaign_edit.html", gin.H{"error": err, "campaign": input})
	}

	createCampaignInput := campaign.CreateCampaignInput{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		User:             campaignData.User,
	}

	_, err = h.campaignService.UpdateCampaign(campaign.GetCampaignDetailInput{ID: idParam}, createCampaignInput)
	if err != nil {
		input.ID = idParam
		c.HTML(http.StatusInternalServerError, "campaign_edit.html", gin.H{"error": err, "campaign": input})
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Show(c *gin.Context) {
	id := c.Param("id")
	idParam, _ := strconv.Atoi(id)

	campaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: idParam})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}

	c.HTML(http.StatusOK, "campaign_show.html", campaign)
}
