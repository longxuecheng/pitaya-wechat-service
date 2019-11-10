package utils

import (
	"github.com/speps/go-hashids"
)

var encryptor *HashEncryptor

func InitEncryptor() {
	encryptor = NewHashEncryptor("nG7nnhzEsDkiYadK", 16)
}

func NewHashEncryptor(salt string, minLength int) *HashEncryptor {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	hasher, _ := hashids.NewWithData(hd)
	return &HashEncryptor{
		hashData: hd,
		hasher:   hasher,
	}
}

type HashEncryptor struct {
	hashData *hashids.HashIDData
	hasher   *hashids.HashID
}

func (e *HashEncryptor) EncodeIn64(id int64) string {
	s, _ := e.hasher.EncodeInt64([]int64{id})
	return s
}

func (e *HashEncryptor) DecodeIn64(str string) int64 {
	int64s := e.hasher.DecodeInt64(str)
	return int64s[0]
}

func EncodeIn64(id int64) string {
	return encryptor.EncodeIn64(id)
}

func DecodeIn64(str string) int64 {
	return encryptor.DecodeIn64(str)
}
