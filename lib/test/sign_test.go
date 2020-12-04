// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package test

import (
	"testing"

	"github.com/juetun/base-wrapper/lib/common/sign"
)

func TestSign(t *testing.T) {
	sign.Sign().
		SignTopRequest(map[string]string{
			"KK": "111",
			"AA": "mmmk",
		}, "aaa", sign.CHARSET_UTF_8)
}
