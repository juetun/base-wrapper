package hashid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	a := assert.New(t)

	hid, err := Encode("tableName", 1)
	a.Nil(err)
	a.Len(hid, defaultLength)

}
