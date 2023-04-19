package handler

import (
	"campaigns/campaign"
	"campaigns/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ambil semua data campaign atau ambil data sesuai dengan user yg login

// tangkap parameter di handler
// handler ke service
// service menentukan repository mana yang di call
// repository akses campaign : GetAll dan GetByUserID
// db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id")) // returnya pasti string, jd pake strconv buatkonfersi ke int
	campaign, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns",http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return 
	}
	response := helper.APIResponse("List ofcampaigns", http.StatusOK, "success", campaign)
	c.JSON(http.StatusOK, response)
}