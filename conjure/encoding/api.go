package encoding

// (Marshal|Write|Append)(JSON|Text)[Strict]

type JSONObject interface {
	MarshalJSON() ([]byte, error)
	AppendJSON([]byte) []byte

	UnmarshalJSON([]byte) error
	UnmarshalJSONStrict([]byte) error
}
