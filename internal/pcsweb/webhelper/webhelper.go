package webhelper

import (
	"github.com/json-iterator/go"
	"io"
)

// WriteJSONObject 写入json对象
func WriteJSONObject(w io.Writer, obj interface{}) error {
	e := jsoniter.NewEncoder(w)
	err := e.Encode(obj)
	if err != nil {
		return err
	}
	return nil
}

// MustWriteJSONObject 写入json对象, 出错则panic
func MustWriteJSONObject(w io.Writer, obj interface{}) {
	err := WriteJSONObject(w, obj)
	if err != nil {
		panic(err)
	}
}
