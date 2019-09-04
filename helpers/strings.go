package helpers

import (
	"crypto/rand"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz"
const Maxlen = 32
const ml = 8

// Copy of auth.GenerateSecureKey to prevent cyclic import with auth library
func GenerateSecureKey() string {
	k := make([]byte, 32)
	io.ReadFull(rand.Reader, k)
	return fmt.Sprintf("%x", k)
}

func CheckToBool(str string) bool {
	if str == "true" {
		return true
	}
	return false
}

func CapFirstLetter(s string) string {
	a := []rune(s)
	a[0] = unicode.ToUpper(a[0])
	s = string(a)
	return s
}

func PrettifyString(s string) string {
	//let's make pretty urls from title
	reg, _ := regexp.Compile("[^A-Za-z0-9]+")
	s = reg.ReplaceAllString(s, "-")
	s = strings.ToLower(strings.Trim(s, "-"))
	return s
}

/*

func EncodeHashIds(i64 int64, salt string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = ml
	h := hashids.NewWithData(hd)
	s, err := h.EncodeInt64([]int64{i64})
	if err != nil {
		return "", err
	}
	return s, nil
}

func DecodeHashIds(s string, salt string) (int64, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = ml
	h := hashids.NewWithData(hd)
	i64, err := h.DecodeInt64WithError(s)
	if err != nil {
		return 0, err
	}
	return i64[0], err
}
*/
