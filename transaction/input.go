package transaction

import "campaigns/user"

type GetTransactionsCampaignInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}