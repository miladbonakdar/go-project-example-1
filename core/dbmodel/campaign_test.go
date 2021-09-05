package dbmodel_test

import (
	"giftcard-engine/core/dbmodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCampaign(t *testing.T) {
	t.Parallel()
	camp := dbmodel.NewCampaign("test")

	assert.NotEmpty(t, camp)
	assert.NotEmpty(t, camp.Title)
	assert.Empty(t, camp.ID)
	assert.Empty(t, camp.CreatedAt)
	assert.Empty(t, camp.UpdatedAt)
	assert.Empty(t, camp.DeletedAt)
	assert.Equal(t, "test", camp.Title)
}

func TestEmptyCampaign(t *testing.T) {
	t.Parallel()
	camp := dbmodel.EmptyCampaign()
	assert.Empty(t, camp.Title)
	assert.Empty(t, camp.ID)
	assert.Empty(t, camp.CreatedAt)
	assert.Empty(t, camp.UpdatedAt)
	assert.Empty(t, camp.DeletedAt)

}

func TestUpdateCampaign(t *testing.T) {
	t.Parallel()
	camp := dbmodel.NewCampaign("test")
	camp.Update("dastan")
	assert.NotEqual(t, "test", camp.Title)
	assert.Equal(t, "dastan", camp.Title)
}
