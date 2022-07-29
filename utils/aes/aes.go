package aes

import (
	"github.com/Artexus/api-matthew-backend/constant"
	"github.com/speps/go-hashids"
)

var hashData *hashids.HashIDData
var hashID *hashids.HashID

func initialize() {
	if hashData != nil {
		return
	}
	hashData = hashids.NewData()
	hashData.MinLength = constant.AESMinLength
	hashData.Salt = constant.AESKey

	hashID, _ = hashids.NewWithData(hashData)
}

func EncryptID(id int) string {
	initialize()

	res, _ := hashID.Encode([]int{id})
	return res
}

func DecryptID(id string) (res int) {
	initialize()

	res = hashID.Decode(id)[0]
	return
}
