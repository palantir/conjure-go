package dj

import (
	"io"
)

var CODEC = codec{}

type codec struct{}

func (codec) Accept() string {
	return "application/json"
}

func (codec) Decode(r io.Reader, v interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (codec) ContentType() string {
	return "application/json"
}

func (codec) Encode(w io.Writer, v interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (codec) Marshal(v interface{}) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
