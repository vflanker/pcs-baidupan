package pcsweb

import (
	"github.com/iikira/BaiduPCS-Go/internal/pcsweb/webhelper"
	"github.com/nytimes/gziphandler"
	"net/http"
	"strings"
)

func (pw *PCSWeb) authMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions { // 忽略OPTIONS
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			return
		}

		// 验证token
		if !pw.checkToken(r.Header) {
			w.WriteHeader(http.StatusUnauthorized)
			webhelper.WriteJSONObject(w, &ErrInfo{
				Errno: http.StatusUnauthorized,
				Msg:   http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		next(w, r)
	}
}

func crossSiteMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		{
			header.Set("Access-Control-Allow-Origin", "http://192.168.50.41:4200")
			header.Add("Access-Control-Allow-Headers", strings.Join([]string{
				HeaderContentType,
				HeaderAuthorization,
			}, ","))
		}
		next(w, r)
	}
}

func jsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "application/json")
		next(w, r)
	}
}

func gzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gziphandler.GzipHandler(http.HandlerFunc(next)).ServeHTTP(w, r)
	}
}
