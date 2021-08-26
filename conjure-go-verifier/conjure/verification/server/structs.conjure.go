// This file was generated by Conjure and should not be manually edited.

package server

import (
	safejson "github.com/palantir/pkg/safejson"
)

type ClientTestCases struct {
	AutoDeserialize         map[EndpointName]PositiveAndNegativeTestCases `json:"autoDeserialize"`
	SingleHeaderService     map[EndpointName][]string                     `json:"singleHeaderService"`
	SinglePathParamService  map[EndpointName][]string                     `json:"singlePathParamService"`
	SingleQueryParamService map[EndpointName][]string                     `json:"singleQueryParamService"`
}

func (o ClientTestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o ClientTestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"autoDeserialize\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.AutoDeserialize {
				var err error
				out, err = k.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				out = append(out, ':')
				var err error
				out, err = v.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				i++
				if i < len(o.AutoDeserialize) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singleHeaderService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SingleHeaderService {
				var err error
				out, err = k.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				out = append(out, ':')
				out = append(out, '[')
				{
					for i := range v {
						out = safejson.AppendQuotedString(out, v[i])
						if i < len(v)-1 {
							out = append(out, ',')
						}
					}
				}
				out = append(out, ']')
				i++
				if i < len(o.SingleHeaderService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singlePathParamService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SinglePathParamService {
				var err error
				out, err = k.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				out = append(out, ':')
				out = append(out, '[')
				{
					for i := range v {
						out = safejson.AppendQuotedString(out, v[i])
						if i < len(v)-1 {
							out = append(out, ',')
						}
					}
				}
				out = append(out, ']')
				i++
				if i < len(o.SinglePathParamService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singleQueryParamService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SingleQueryParamService {
				var err error
				out, err = k.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				out = append(out, ':')
				out = append(out, '[')
				{
					for i := range v {
						out = safejson.AppendQuotedString(out, v[i])
						if i < len(v)-1 {
							out = append(out, ',')
						}
					}
				}
				out = append(out, ']')
				i++
				if i < len(o.SingleQueryParamService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
}

type IgnoredClientTestCases struct {
	AutoDeserialize         map[EndpointName][]string `json:"autoDeserialize"`
	SingleHeaderService     map[EndpointName][]string `json:"singleHeaderService"`
	SinglePathParamService  map[EndpointName][]string `json:"singlePathParamService"`
	SingleQueryParamService map[EndpointName][]string `json:"singleQueryParamService"`
}

func (o IgnoredClientTestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o IgnoredClientTestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"autoDeserialize\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.AutoDeserialize {
				var err error
				out, err = k.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				out = append(out, ':')
				out = append(out, '[')
				{
					for i := range v {
						out = safejson.AppendQuotedString(out, v[i])
						if i < len(v)-1 {
							out = append(out, ',')
						}
					}
				}
				out = append(out, ']')
				i++
				if i < len(o.AutoDeserialize) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singleHeaderService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SingleHeaderService {
				var err error
				out, err = k.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				out = append(out, ':')
				out = append(out, '[')
				{
					for i := range v {
						out = safejson.AppendQuotedString(out, v[i])
						if i < len(v)-1 {
							out = append(out, ',')
						}
					}
				}
				out = append(out, ']')
				i++
				if i < len(o.SingleHeaderService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singlePathParamService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SinglePathParamService {
				var err error
				out, err = k.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				out = append(out, ':')
				out = append(out, '[')
				{
					for i := range v {
						out = safejson.AppendQuotedString(out, v[i])
						if i < len(v)-1 {
							out = append(out, ',')
						}
					}
				}
				out = append(out, ']')
				i++
				if i < len(o.SinglePathParamService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singleQueryParamService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SingleQueryParamService {
				var err error
				out, err = k.AppendJSON(out)
				if err != nil {
					return nil, err
				}
				out = append(out, ':')
				out = append(out, '[')
				{
					for i := range v {
						out = safejson.AppendQuotedString(out, v[i])
						if i < len(v)-1 {
							out = append(out, ',')
						}
					}
				}
				out = append(out, ']')
				i++
				if i < len(o.SingleQueryParamService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
}

type IgnoredTestCases struct {
	Client IgnoredClientTestCases `json:"client"`
}

func (o IgnoredTestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o IgnoredTestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"client\":"...)
		var err error
		out, err = o.Client.AppendJSON(out)
		if err != nil {
			return nil, err
		}
	}
	out = append(out, '}')
	return out, nil
}

type PositiveAndNegativeTestCases struct {
	Positive []string `json:"positive"`
	Negative []string `json:"negative"`
}

func (o PositiveAndNegativeTestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o PositiveAndNegativeTestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"positive\":"...)
		out = append(out, '[')
		{
			for i := range o.Positive {
				out = safejson.AppendQuotedString(out, o.Positive[i])
				if i < len(o.Positive)-1 {
					out = append(out, ',')
				}
			}
		}
		out = append(out, ']')
		out = append(out, ',')
	}
	{
		out = append(out, "\"negative\":"...)
		out = append(out, '[')
		{
			for i := range o.Negative {
				out = safejson.AppendQuotedString(out, o.Negative[i])
				if i < len(o.Negative)-1 {
					out = append(out, ',')
				}
			}
		}
		out = append(out, ']')
	}
	out = append(out, '}')
	return out, nil
}

type TestCases struct {
	Client ClientTestCases `json:"client"`
}

func (o TestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o TestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"client\":"...)
		var err error
		out, err = o.Client.AppendJSON(out)
		if err != nil {
			return nil, err
		}
	}
	out = append(out, '}')
	return out, nil
}
