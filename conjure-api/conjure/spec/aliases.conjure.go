// This file was generated by Conjure and should not be manually edited.

package spec

import (
	"context"

	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

// Must be in lowerCamelCase. Numbers are permitted, but not at the beginning of a word. Allowed argument names: "fooBar", "build2Request". Disallowed names: "FooBar", "2BuildRequest".
type ArgumentName string

func (a ArgumentName) String() string {
	return string(a)
}

func (a *ArgumentName) UnmarshalString(data string) error {
	rawArgumentName := data
	*a = ArgumentName(rawArgumentName)
	return nil
}

func (a ArgumentName) MarshalJSON() ([]byte, error) {
	size, err := a.JSONSize()
	if err != nil {
		return nil, err
	}
	return a.AppendJSON(make([]byte, 0, size))
}

func (a ArgumentName) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, string(a))
	return out, nil
}

func (a ArgumentName) JSONSize() (int, error) {
	var out int
	out += safejson.QuotedStringLength(string(a))
	return out, nil
}

func (a *ArgumentName) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ArgumentName")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *ArgumentName) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ArgumentName")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *ArgumentName) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawArgumentName string
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "ArgumentName expected JSON string")
		return err
	}
	rawArgumentName = value.Str
	*a = ArgumentName(rawArgumentName)
	return nil
}

func (a ArgumentName) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ArgumentName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type Documentation string

func (a Documentation) String() string {
	return string(a)
}

func (a *Documentation) UnmarshalString(data string) error {
	rawDocumentation := data
	*a = Documentation(rawDocumentation)
	return nil
}

func (a Documentation) MarshalJSON() ([]byte, error) {
	size, err := a.JSONSize()
	if err != nil {
		return nil, err
	}
	return a.AppendJSON(make([]byte, 0, size))
}

func (a Documentation) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, string(a))
	return out, nil
}

func (a Documentation) JSONSize() (int, error) {
	var out int
	out += safejson.QuotedStringLength(string(a))
	return out, nil
}

func (a *Documentation) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Documentation")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *Documentation) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Documentation")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *Documentation) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawDocumentation string
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "Documentation expected JSON string")
		return err
	}
	rawDocumentation = value.Str
	*a = Documentation(rawDocumentation)
	return nil
}

func (a Documentation) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *Documentation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

// Should be in lowerCamelCase.
type EndpointName string

func (a EndpointName) String() string {
	return string(a)
}

func (a *EndpointName) UnmarshalString(data string) error {
	rawEndpointName := data
	*a = EndpointName(rawEndpointName)
	return nil
}

func (a EndpointName) MarshalJSON() ([]byte, error) {
	size, err := a.JSONSize()
	if err != nil {
		return nil, err
	}
	return a.AppendJSON(make([]byte, 0, size))
}

func (a EndpointName) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, string(a))
	return out, nil
}

func (a EndpointName) JSONSize() (int, error) {
	var out int
	out += safejson.QuotedStringLength(string(a))
	return out, nil
}

func (a *EndpointName) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for EndpointName")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *EndpointName) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for EndpointName")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *EndpointName) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawEndpointName string
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "EndpointName expected JSON string")
		return err
	}
	rawEndpointName = value.Str
	*a = EndpointName(rawEndpointName)
	return nil
}

func (a EndpointName) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *EndpointName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type ErrorNamespace string

func (a ErrorNamespace) String() string {
	return string(a)
}

func (a *ErrorNamespace) UnmarshalString(data string) error {
	rawErrorNamespace := data
	*a = ErrorNamespace(rawErrorNamespace)
	return nil
}

func (a ErrorNamespace) MarshalJSON() ([]byte, error) {
	size, err := a.JSONSize()
	if err != nil {
		return nil, err
	}
	return a.AppendJSON(make([]byte, 0, size))
}

func (a ErrorNamespace) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, string(a))
	return out, nil
}

func (a ErrorNamespace) JSONSize() (int, error) {
	var out int
	out += safejson.QuotedStringLength(string(a))
	return out, nil
}

func (a *ErrorNamespace) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ErrorNamespace")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *ErrorNamespace) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ErrorNamespace")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *ErrorNamespace) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawErrorNamespace string
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "ErrorNamespace expected JSON string")
		return err
	}
	rawErrorNamespace = value.Str
	*a = ErrorNamespace(rawErrorNamespace)
	return nil
}

func (a ErrorNamespace) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ErrorNamespace) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

// Should be in lowerCamelCase, but kebab-case and snake_case are also permitted.
type FieldName string

func (a FieldName) String() string {
	return string(a)
}

func (a *FieldName) UnmarshalString(data string) error {
	rawFieldName := data
	*a = FieldName(rawFieldName)
	return nil
}

func (a FieldName) MarshalJSON() ([]byte, error) {
	size, err := a.JSONSize()
	if err != nil {
		return nil, err
	}
	return a.AppendJSON(make([]byte, 0, size))
}

func (a FieldName) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, string(a))
	return out, nil
}

func (a FieldName) JSONSize() (int, error) {
	var out int
	out += safejson.QuotedStringLength(string(a))
	return out, nil
}

func (a *FieldName) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for FieldName")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *FieldName) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for FieldName")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *FieldName) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawFieldName string
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "FieldName expected JSON string")
		return err
	}
	rawFieldName = value.Str
	*a = FieldName(rawFieldName)
	return nil
}

func (a FieldName) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *FieldName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type HttpPath string

func (a HttpPath) String() string {
	return string(a)
}

func (a *HttpPath) UnmarshalString(data string) error {
	rawHttpPath := data
	*a = HttpPath(rawHttpPath)
	return nil
}

func (a HttpPath) MarshalJSON() ([]byte, error) {
	size, err := a.JSONSize()
	if err != nil {
		return nil, err
	}
	return a.AppendJSON(make([]byte, 0, size))
}

func (a HttpPath) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, string(a))
	return out, nil
}

func (a HttpPath) JSONSize() (int, error) {
	var out int
	out += safejson.QuotedStringLength(string(a))
	return out, nil
}

func (a *HttpPath) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for HttpPath")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *HttpPath) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for HttpPath")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *HttpPath) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawHttpPath string
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "HttpPath expected JSON string")
		return err
	}
	rawHttpPath = value.Str
	*a = HttpPath(rawHttpPath)
	return nil
}

func (a HttpPath) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *HttpPath) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

// For header parameters, the parameter id must be in Upper-Kebab-Case. For query parameters, the parameter id must be in lowerCamelCase. Numbers are permitted, but not at the beginning of a word.
type ParameterId string

func (a ParameterId) String() string {
	return string(a)
}

func (a *ParameterId) UnmarshalString(data string) error {
	rawParameterId := data
	*a = ParameterId(rawParameterId)
	return nil
}

func (a ParameterId) MarshalJSON() ([]byte, error) {
	size, err := a.JSONSize()
	if err != nil {
		return nil, err
	}
	return a.AppendJSON(make([]byte, 0, size))
}

func (a ParameterId) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, string(a))
	return out, nil
}

func (a ParameterId) JSONSize() (int, error) {
	var out int
	out += safejson.QuotedStringLength(string(a))
	return out, nil
}

func (a *ParameterId) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ParameterId")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *ParameterId) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ParameterId")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *ParameterId) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawParameterId string
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "ParameterId expected JSON string")
		return err
	}
	rawParameterId = value.Str
	*a = ParameterId(rawParameterId)
	return nil
}

func (a ParameterId) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ParameterId) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
