package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

// RandomString isn't crypto secure or anything
// we just need a psuedo random key to be generated for our canary token
func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func GenerateLink(baseurl, namespace string, size int) (string, string) {
	key := RandomString(size)
	url := fmt.Sprintf("%s/%s/%s/i.png", baseurl, namespace, key)
	log.Debugf("Url %s : Key %s", url, key)
	return key, url
}
