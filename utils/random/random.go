package random

import (
	"giftcard-engine/utils"
	"math/rand"
	"sync"
	"time"
)

var mu sync.Mutex

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GiftCardSecretKey() string {
	return complexString(utils.GiftCardSecretKeyLength)
}

func GiftCardPublicKey() string {
	return stringWithCharset(utils.GiftCardPublicKeyLength, utils.Numbers)
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	numberRange := len(charset)
	for i := range b {
		b[i] = charset[atomicIntN(numberRange)]
	}
	return string(b)
}

func complexString(length int) string {
	return stringWithCharset(length, utils.EnglishCharacters+utils.Numbers)
}

func atomicIntN(numberRange int) int {
	mu.Lock()
	defer mu.Unlock()
	return seededRand.Intn(numberRange)
}
