// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package test

import (
	"testing"

	"github.com/juetun/base-wrapper/lib/common/signencrypt"
)

func TestSign(t *testing.T) {
	signencrypt.Sign().
		SignTopRequest(map[string]string{
			"KK": "111",
			"AA": "mmmk",
		}, "aaa", signencrypt.CharsetUtf8)
}
