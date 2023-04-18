package handler

import (
	"campaigns/helper"
	"campaigns/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
)

type userhandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userhandler {
	return &userhandler{userService}
}

func (h *userhandler) RegisterUser(c *gin.Context) {
	//tangkap input dari user
	// map inpu dari user ke struct registerUserInput
	// tsruct di atas passing sbg parameter service
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors} // gin.h buap map aja

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "tokentokentoken")

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}


func (h *userhandler) Login(c *gin.Context) {
	// masukkan inputan dari user (email dan password)
	var input user.LoginInput

	// input ditangkap handler/binding
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors} // gin.h buap map aja

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()} // gin.h buap map aja

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, "tokentoken")
	response := helper.APIResponse("Login success!", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
	// mapping dari input user ke input struct
	// inpit struct passing service
	// di service mencari dgn bantuan repository  user dengan email x
	// cocokkan password

	// jadi pertama-tama buat repositorynya dulu
	// handler terakhir

}

func (h *userhandler) CheckEmailAVailability(c *gin.Context) {
	// input email dari user
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors} // gin.h buap map aja

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"} // gin.h buap map aja

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	var metaMessage string

	if isEmailAvailable {
		metaMessage = "Email is available"
	} else {
		metaMessage = "Email has been registered"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
		
	// input email di mapping ke struct - handler
	//  struct input di passing ke service
	// service pnggil repo, email dah ada belum
	// repository - query ke db

}

// upload avatar user (pae update)
func (h *userhandler) UploadAvatar(c *gin.Context) {
	// masukkan gambar
	// simpan gambar ke folder "images/"
	// service panggil repo untuk 1. menentukan user mana yang akses. pake jwt, 2. simpna lokasi file
	// repo: ambil id user, update data user 
	 
	// ambil foto dari post yang key-nya avatar
	file, err := c.FormFile("avatar")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// harusnya dapet dari JWT, inin manual dulu
	userID := 4

	// lokasi simpen foto
	// path := "images/" + file.Filename
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded" : false}
		response := helper.APIResponse("Failed to upload profile picture", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}


	_, err = h.userService.SaveAvatar(userID, path)
	
	if err != nil {
		data := gin.H{"is_uploaded" : false}
		response := helper.APIResponse("Failed to upload profile picture", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded":true}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}