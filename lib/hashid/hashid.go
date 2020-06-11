package hashid

import (
	"errors"
	"sync"

	"github.com/speps/go-hashids"
)

var hashidEncodeEmptyErr = errors.New("id must gt 0")

const defaultLength = 12

var dataPool = &sync.Pool{
	New: func() interface{} {
		return hashids.NewData()
	},
}

// Encode 生成hashid
func Encode(salt string, id int, minLength ...int) (string, error) {
	if id <= 0 {
		return "", hashidEncodeEmptyErr
	}

	length := defaultLength
	if len(minLength) > 0 {
		length = minLength[0]
	}

	h, err := getHashID(salt, length)
	if err != nil {
		return "", err
	}

	hid, err := h.Encode([]int{id})
	if err != nil {
		return "", err
	}

	return hid, nil
}

func getHashID(salt string, minLength int) (*hashids.HashID, error) {
	hd := dataPool.Get().(*hashids.HashIDData)
	dataPool.Put(hd)

	hd.Salt = salt
	hd.MinLength = minLength

	return hashids.NewWithData(hd)
}
