package campaign

import "campaigns/user"

// get from url lgsgs
type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateCampaignInput struct {
	Name             string    `json:"name" binding:"required"`
	ShortDescription string    `json:"short_description" binding:"required"`
	Description      string    `json:"description" binding:"required"`
	GoalAmount       int       `json:"goal_amount" binding:"required"`
	Perks            string    `json:"perks" binding:"required"`
	User             user.User // get data user dr jwt
}

type CreateCampaignImageInput struct {
	CampaignID int  `form:"campaign_id" binding:"required"`
	IsPrimary  bool `form:"is_primary" binding:"required"`
	User       user.User
}
