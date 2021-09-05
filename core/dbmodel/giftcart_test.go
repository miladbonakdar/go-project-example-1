package dbmodel_test

import (
	"giftcard-engine/core/common"
	"giftcard-engine/core/dbmodel"
	"giftcard-engine/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mssql"
)

func TestNewGiftCard(t *testing.T) {
	t.Parallel()
	date := time.Now().UTC()
	amount := int32(2000)
	card := dbmodel.NewGiftCard(amount, date)

	assert.NotEmpty(t, card)
	assert.NotEmpty(t, card.SecretCode)
	assert.NotEmpty(t, card.PublicCode)
	assert.Empty(t, card.UUN)
	assert.Equal(t, dbmodel.Empty, card.Status)
	assert.Equal(t, date, card.ExpireDate)
	assert.Equal(t, amount, card.Amount)
}

func TestIsValid(t *testing.T) {
	t.Parallel()
	date := time.Now().UTC()
	card1 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "", ExpireDate: date, Status: dbmodel.Empty}
	card2 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "", ExpireDate: date, Status: dbmodel.Approved}
	card3 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "", ExpireDate: time.Now().Add(-time.Hour * 25).UTC(), Status: dbmodel.Empty}
	card4 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "milawd", ExpireDate: time.Now().UTC(), Status: dbmodel.Empty}

	assert.Equal(t, true, card1.IsValid())
	assert.Equal(t, false, card2.IsValid())
	assert.Equal(t, false, card3.IsValid())
	assert.Equal(t, false, card4.IsValid())
}

func TestIsDateValid(t *testing.T) {
	t.Parallel()
	date := time.Now().Add(time.Hour * 25).UTC()
	card1 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "", ExpireDate: date, Status: dbmodel.Empty}
	card2 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "", ExpireDate: time.Now().Add(-time.Hour * 25).UTC(), Status: dbmodel.Empty}

	assert.Equal(t, true, card1.IsDateValid())
	assert.Equal(t, false, card2.IsDateValid())
}

func TestSetUUN(t *testing.T) {
	t.Parallel()
	date := time.Now().Add(time.Hour * 25).UTC()
	card1 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "", ExpireDate: date, Status: dbmodel.Empty}
	card2 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "milawd", ExpireDate: date, Status: dbmodel.Empty}

	err1 := card1.SetUUN("some_one")
	err2 := card2.SetUUN("some_one")

	assert.Empty(t, err1)
	assert.NotEmpty(t, err2)
	assert.Equal(t, common.GiftCardIsTaken, err2)
}

func TestSetCampaign(t *testing.T) {
	t.Parallel()
	date := time.Now().Add(time.Hour * 25).UTC()
	card1 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "", ExpireDate: date, Status: dbmodel.Empty}

	card1.SetCampaign(1)

	assert.Equal(t, uint(1), card1.CampaignId)
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	date := time.Now().Add(time.Hour * 25).UTC()
	dateForUpdate := date.Add(time.Hour)
	amountForUpdate := int32(4000)
	card1 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "", ExpireDate: date, Status: dbmodel.Empty}
	card2 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "milawd", ExpireDate: date, Status: dbmodel.Empty}

	err1 := card1.Update(amountForUpdate, dateForUpdate)
	err2 := card2.Update(amountForUpdate, dateForUpdate)

	assert.Empty(t, err1)
	assert.NotEmpty(t, err2)
	assert.Equal(t, amountForUpdate, card1.Amount)
	assert.Equal(t, dateForUpdate, card1.ExpireDate)
	assert.Equal(t, common.GiftCardIsNotValid, err2)
}

func TestRollBack(t *testing.T) {
	t.Parallel()
	date := time.Now().Add(time.Hour * 25).UTC()
	card1 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "milawd", ExpireDate: date, Status: dbmodel.Approved}

	card1.RollBack()

	assert.Equal(t, "", card1.UUN)
	assert.Equal(t, dbmodel.Empty, card1.Status)
}

func TestGenerateKey(t *testing.T) {
	t.Parallel()
	date := time.Now().Add(time.Hour * 25).UTC()
	card1 := dbmodel.GiftCard{Amount: int32(2000), PublicCode: "public",
		SecretCode: "secret", UUN: "milawd", ExpireDate: date, Status: dbmodel.Approved}

	card1.GenerateKey()

	assert.NotEqual(t, "public", card1.PublicCode)
	assert.NotEqual(t, "secret", card1.SecretCode)
	assert.Equal(t, utils.GiftCardPublicKeyLength, len(card1.PublicCode))
	assert.Equal(t, utils.GiftCardSecretKeyLength, len(card1.SecretCode))
}
