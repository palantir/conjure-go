package cj

import (
	"io"

	"github.com/go-json-experiment/json"
)

type CodecJSONV2 struct{}

func (CodecJSONV2) Accept() string {
	return "application/json"
}

func (CodecJSONV2) Decode(r io.Reader, v interface{}) error {
	return json.UnmarshalRead(r, *&v)
}

func (CodecJSONV2) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, *&v)
}

func (CodecJSONV2) ContentType() string {
	return "application/json"
}

func (CodecJSONV2) Encode(w io.Writer, v interface{}) error {
	return json.MarshalWrite(w, v)
}

func (CodecJSONV2) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
