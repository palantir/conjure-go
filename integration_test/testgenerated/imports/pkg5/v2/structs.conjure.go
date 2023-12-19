// This file was generated by Conjure and should not be manually edited.

package v2

import (
	"io"

	dj "github.com/palantir/conjure-go/v6/dj"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
)

type DifferentPackageEndingInVersion struct {
	Name string `json:"name"`
}

func (o DifferentPackageEndingInVersion) MarshalJSON() ([]byte, error) {
	out := make([]byte, 0)
	if _, err := o.WriteJSON(dj.NewAppender(&out)); err != nil {
		return nil, err
	}
	return out, dj.Valid(out)
}

func (o DifferentPackageEndingInVersion) WriteJSON(w io.Writer) (int, error) {
	var out int
	if n, err := dj.WriteOpenObject(w); err != nil {
		return 0, err
	} else {
		out += n
	}
	{
		if n, err := dj.WriteLiteral(w, "\"name\":"); err != nil {
			return 0, err
		} else {
			out += n
		}
		if n, err := dj.WriteString(w, o.Name); err != nil {
			return 0, err
		} else {
			out += n
		}
	}
	if n, err := dj.WriteCloseObject(w); err != nil {
		return 0, err
	} else {
		out += n
	}
	return out, nil
}

func (o *DifferentPackageEndingInVersion) UnmarshalJSON(data []byte) error {
	value, err := dj.Parse(data)
	if err != nil {
		return err
	}
	return o.UnmarshalJSONResult(value, false)
}

func (o *DifferentPackageEndingInVersion) UnmarshalJSONStrict(data []byte) error {
	value, err := dj.Parse(data)
	if err != nil {
		return err
	}
	return o.UnmarshalJSONResult(value, true)
}

func (o *DifferentPackageEndingInVersion) UnmarshalJSONString(data string) error {
	value, err := dj.Parse(data)
	if err != nil {
		return err
	}
	return o.UnmarshalJSONResult(value, false)
}

func (o *DifferentPackageEndingInVersion) UnmarshalJSONStringStrict(data string) error {
	value, err := dj.Parse(data)
	if err != nil {
		return err
	}
	return o.UnmarshalJSONResult(value, true)
}

func (o *DifferentPackageEndingInVersion) UnmarshalJSONResult(value dj.Result, disallowUnknownFields bool) error {
	var seenName bool
	var unknownFields []string
	iter, idx, err := value.ObjectIterator(0)
	if err != nil {
		return err
	}
	for iter.HasNext(value, idx) {
		var fieldKey, fieldValue dj.Result
		fieldKey, fieldValue, idx, err = iter.Next(value, idx)
		if err != nil {
			return err
		}
		switch fieldKey.Str {
		case "name":
			if seenName {
				return dj.UnmarshalDuplicateFieldError{Index: fieldKey.Index, Type: "DifferentPackageEndingInVersion", Field: "name"}
			}
			seenName = true
			o.Name, err = fieldValue.String()
			if err != nil {
				return werror.Convert(dj.UnmarshalFieldError{Index: fieldValue.Index, Type: "DifferentPackageEndingInVersion", Field: "name", Err: err})
			}
		default:
			if disallowUnknownFields {
				unknownFields = append(unknownFields, fieldKey.Str)
			}
		}
	}
	var missingFields []string
	if !seenName {
		missingFields = append(missingFields, "name")
	}
	if len(missingFields) > 0 {
		return werror.Convert(dj.UnmarshalMissingFieldsError{Index: value.Index, Type: "DifferentPackageEndingInVersion", Fields: missingFields})
	}
	if disallowUnknownFields && len(unknownFields) > 0 {
		return werror.Convert(dj.UnmarshalUnknownFieldsError{Index: value.Index, Type: "DifferentPackageEndingInVersion", Fields: unknownFields})
	}
	return nil
}

func (o DifferentPackageEndingInVersion) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *DifferentPackageEndingInVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
