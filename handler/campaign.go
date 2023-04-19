package handler

import (
	"campaigns/campaign"
	"campaigns/helper"
	"campaigns/user"
	"fmt"
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

// read campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id")) // returnya pasti string, jd pake strconv buatkonfersi ke int
	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetOneCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

// create campaigns
func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput
	// fmt.Println(input)
	err := c.ShouldBindJSON(&input)
	fmt.Println(err)
	if err != nil {

		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// get data user from jwt
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newcampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newcampaign))
	c.JSON(http.StatusOK, response)
}

// update campaign
func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	
	err := c.ShouldBindUri(&inputID)
	
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var inputData campaign.CreateCampaignInput
	cekInputData := c.ShouldBindJSON(&inputData)
	if cekInputData != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// yg bisa update cuma user yg login aja sesuai dgn post yg dia buat
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	response := helper.APIResponse("Update campaign success!", http.StatusOK, "success", updatedCampaign)
	c.JSON(http.StatusOK, response)

}

// upload image
func (h *campaignHandler) UploadImageCampaign(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err:= c.ShouldBind(&input) // harus kosong, kalau kosong brrti ga eror
	fmt.Println("ini error",err)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}
		response := helper.APIResponse("Failed to upload campaign images",http.StatusBadRequest,"error",errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	// baca inputan gambar dari FE, namanya hars file
	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign images",http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return 
	}
	// currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	path := fmt.Sprintf("images/campaigns/%d-%s",userID,file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded":false}
		response := helper.APIResponse("Failed to upload campaign images",http.StatusBadRequest, "error",data)
		c.JSON(http.StatusBadRequest, response)
		return 
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded":false}
		response := helper.APIResponse("Failed to upload campaign images",http.StatusBadRequest,"error",data)
		c.JSON(http.StatusBadRequest,response)
		return
	}

	data := gin.H{"is_uploaded":true}
	response := helper.APIResponse("Campaign images successfully uploaded",http.StatusOK,"success",data)
	c.JSON(http.StatusOK,response) 
}