package random_test

import (
	"giftcard-engine/utils"
	"giftcard-engine/utils/random"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGiftCardPublicKey(t *testing.T) {
	t.Parallel()
	public := random.GiftCardPublicKey()
	assert.NotEmpty(t, public)
	assert.Equal(t, len(public), utils.GiftCardPublicKeyLength)
}

func TestGiftCardSecretKey(t *testing.T) {
	t.Parallel()
	secret := random.GiftCardSecretKey()
	assert.NotEmpty(t, secret)
	assert.Equal(t, len(secret), utils.GiftCardSecretKeyLength)
}
func BenchmarkTestGiftCardSecretKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		random.GiftCardSecretKey()
	}
}

func BenchmarkGiftCardPublicKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		random.GiftCardPublicKey()
	}
}
