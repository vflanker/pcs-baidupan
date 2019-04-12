package pcsweb

import (
	"fmt"
	"github.com/json-iterator/go"
)

const (
	// ErrnoSuccess 成功
	ErrnoSuccess int = iota
	// ErrnoAuthError 验证错误
	ErrnoAuthError
	// ErrnoPCSAPIError pcs api 错误
	ErrnoPCSAPIError
)

// ErrInfo web 错误详情
type ErrInfo struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

// NewErrInfo 初始化错误信息
func NewErrInfo(errno int, msg string) *ErrInfo {
	return &ErrInfo{
		Errno: errno,
		Msg:   msg,
	}
}

func (ei *ErrInfo) Error() string {
	return fmt.Sprintf("errno: %d, msg: %s", ei.Errno, ei.Msg)
}

// JSON 将错误信息打包成 json
func (ei *ErrInfo) JSON() (data []byte) {
	var err error
	data, err = jsoniter.MarshalIndent(ei, "", " ")
	checkErr(err)

	return
}

// checkErr 遇到错误就退出
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
