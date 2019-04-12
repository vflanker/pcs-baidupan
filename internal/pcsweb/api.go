package pcsweb

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/internal/pcsconfig"
	"github.com/iikira/BaiduPCS-Go/internal/pcsweb/webhelper"
	"io"
	"net/http"
	"path"
)

func (pw *PCSWeb) pcsAPIMiddleware(pcsAPI http.HandlerFunc) http.HandlerFunc {
	return gzipMiddleware(crossSiteMiddleware(jsonMiddleware(pw.authMiddleWare(pcsAPI))))
}

func (pw *PCSWeb) pcsFileList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		pcspath = path.Clean(r.Form.Get("path"))
		orderBy = baidupcs.OrderBy(r.Form.Get("by"))
		order   = baidupcs.Order(r.Form.Get("order"))
	)

	if orderBy == "" {
		orderBy = baidupcs.OrderByName
	}
	if order == "" {
		order = baidupcs.OrderAsc
	}

	respBody, err := pcsconfig.Config.ActiveUserBaiduPCS().PrepareFilesDirectoriesList(pcspath, &baidupcs.OrderOptions{
		By:    orderBy,
		Order: order,
	})
	if err != nil {
		webhelper.MustWriteJSONObject(w, NewErrInfo(ErrnoPCSAPIError, err.Error()))
		return
	}

	defer respBody.Close()
	io.Copy(w, respBody)
}

func (pw *PCSWeb) pcsMkdir(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		pcspath = path.Clean(r.Form.Get("path"))
	)

	respBody, err := pcsconfig.Config.ActiveUserBaiduPCS().PrepareMkdir(pcspath)
	if err != nil {
		webhelper.MustWriteJSONObject(w, NewErrInfo(ErrnoPCSAPIError, err.Error()))
		return
	}

	defer respBody.Close()
	io.Copy(w, respBody)
}
