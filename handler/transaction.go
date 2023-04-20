package handler

import (
	"campaigns/helper"
	"campaigns/transaction"
	"campaigns/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

// di campaign ini, siapa aja yg ngirim transaksi
func (h *transactionHandler) GetCampaignTransaction(c *gin.Context) {
	var input transaction.GetTransactionsCampaignInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	transactions, err := h.service.GetTansactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Transaction on campaign", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

// create by s o f f i e  a p u t r i
// transaksi user sesuai user yang login
func (h *transactionHandler) GetUserTransaction(c *gin.Context){
	// get id user from middleware yg di set
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionsByUserID(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get user transactions",http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest,response)
		return 
	}
	response := helper.APIResponse("User transactions",http.StatusOK,"success",transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK,response)
}