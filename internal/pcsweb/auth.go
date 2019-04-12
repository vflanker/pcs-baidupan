package pcsweb

import (
	"github.com/iikira/BaiduPCS-Go/internal/pcsweb/webhelper"
	"net/http"
	"strings"
)

type (
	jwtJSON struct {
		*ErrInfo
		Token string `json:"token"`
	}
)

func (pw *PCSWeb) issueToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		pwd = r.Form.Get("pwd")
	)

	// 验证密码
	if pw.Pwd != pwd {
		webhelper.MustWriteJSONObject(w, NewErrInfo(ErrnoAuthError, "pwd is invalid"))
		return
	}

	j := jwtJSON{
		ErrInfo: &ErrInfo{},
		Token:   pw.webjwt.SignStandardClaims(),
	}
	webhelper.MustWriteJSONObject(w, &j)
	return
}

func (pw *PCSWeb) checkToken(header http.Header) bool {
	tokenString := ParseToken(header)
	return pw.webjwt.Verify(tokenString)
}

// ParseToken 从header解析token
func ParseToken(header http.Header) string {
	a := header.Get(HeaderAuthorization)
	return strings.TrimPrefix(a, "Bearer ")
}
