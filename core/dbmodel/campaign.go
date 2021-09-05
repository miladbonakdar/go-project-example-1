package dbmodel

type Campaign struct {
	AbstractModel
	Title string `gorm:"column:Title;unique_index;not null"`
}

func NewCampaign(title string) *Campaign {
	return &Campaign{
		Title: title,
	}
}

func EmptyCampaign() Campaign {
	return Campaign{}
}

func (c *Campaign) Update(title string) {
	c.Title = title
}

func (*Campaign) TableName() string {
	return "Campaign"
}
