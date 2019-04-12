// Package pcsweb web前端包
package pcsweb

import (
	"github.com/GeertJohan/go.rice"
	"github.com/iikira/BaiduPCS-Go/internal/pcsconfig"
	"github.com/iikira/BaiduPCS-Go/internal/pcsweb/webjwt"
	"io"
	"net/http"
	"path"
	"strings"
)

const (
	// APIVersion api 版本
	APIVersion = "2.0"
	// HeaderContentType Content-Type
	HeaderContentType = "Content-Type"
	// HeaderAuthorization 用户鉴权
	HeaderAuthorization = "Authorization"
)

type (
	// PCSWeb web对象
	PCSWeb struct {
		*pcsconfig.WebConfig

		fileBox *rice.Box
		mux     *http.ServeMux
		webjwt  *webjwt.WebJWT

		indexContents []byte
	}
)

// New 初始化 PCSWeb
func New() *PCSWeb {
	return &PCSWeb{}
}

func (pw *PCSWeb) lazyInit() {
	if pw.WebConfig == nil {
		pw.WebConfig = pcsconfig.Config.WebConfig()
	}
	if pw.WebConfig.Addr == "" {
		pw.WebConfig.Addr = ":4203"
	}
	if pw.mux == nil {
		pw.mux = http.NewServeMux()
	}
	if pw.webjwt == nil {
		pw.webjwt = webjwt.NewWebJWT()
	}
}

func (pw *PCSWeb) initFileBox() (err error) {
	if pw.fileBox != nil {
		return
	}

	pw.fileBox, err = rice.FindBox("filebox") // 文件盒子
	if err != nil {
		return err
	}

	pw.indexContents, err = pw.fileBox.Bytes("index.html") // 主页
	if err != nil {
		return err
	}

	return nil
}

// SetAddr 设置addr
func (pw *PCSWeb) SetAddr(addr string) {
	pw.lazyInit()
	pw.Addr = addr
}

// Serve 启动服务
func (pw *PCSWeb) Serve() (err error) {
	err = pw.initFileBox()
	if err != nil {
		return err
	}

	pw.lazyInit()

	// 前端
	pw.mux.Handle("/", gzipMiddleware(pw.index))

	// auth
	pw.mux.Handle("/auth/register", crossSiteMiddleware(jsonMiddleware(pw.issueToken)))

	// pcs rest api
	pw.mux.Handle(pw.getPCSServePath("/file/mkdir"), pw.pcsAPIMiddleware(pw.pcsMkdir))
	pw.mux.Handle(pw.getPCSServePath("/file/list"), pw.pcsAPIMiddleware(pw.pcsFileList))

	return http.ListenAndServe(pw.Addr, pw.mux)
}

func (pw *PCSWeb) getPCSServePath(suffix string) string {
	return "/api/" + APIVersion + "/pcs" + suffix
}

// index 所有的请求都交给index.html处理
func (pw *PCSWeb) index(w http.ResponseWriter, r *http.Request) {
	bpath := strings.TrimLeft(r.URL.Path, "/")
	file, err := pw.fileBox.Open(bpath)
	if file != nil {
		defer file.Close()
		stat, err := file.Stat()
		if err != nil || stat.IsDir() {
			w.Write(pw.indexContents)
			return
		}
	}
	if err != nil {
		w.Write(pw.indexContents)
		return
	}

	// 获取资源
	var (
		ext    = path.Ext(bpath)
		header = w.Header()
	)

	switch ext {
	case ".html":
		header.Set(HeaderContentType, "text/html; charset=utf-8")
	case ".css":
		header.Set(HeaderContentType, "text/css")
	case ".js":
		header.Set(HeaderContentType, "application/javascript")
	case ".svg":
		header.Set(HeaderContentType, "image/svg+xml")
	}

	if ext != ".html" {
		header.Set("Cache-Control", "public, max-age=31536000")
	}

	io.Copy(w, file)
}
