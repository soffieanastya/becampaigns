package handler

import (
	"campaigns/helper"
	"campaigns/user"
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