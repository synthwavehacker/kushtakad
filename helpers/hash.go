package helpers

import (
	hashids "github.com/speps/go-hashids"
)

const HashMin = 16

func HashToId(hash, salt string) (int64, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 8
	h, err := hashids.NewWithData(hd)
	id, err := h.DecodeInt64WithError(hash)
	return id[0], err
}

func IdToHash(id int64, salt string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 8
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	hashid, err := h.EncodeInt64([]int64{id})
	if err != nil {
		return "", err
	}
	return hashid, nil
}
