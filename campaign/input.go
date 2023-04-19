package campaign

// get from url lgsgs
type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"`
}
