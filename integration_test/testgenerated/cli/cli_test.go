package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go/v6/integration_test/testgenerated/cli/api"
	api_mock "github.com/palantir/conjure-go/v6/internal/generated/mocks/github.com/palantir/conjure-go/v6/integration_test/testgenerated/cli/api"
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/uuid"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
	"github.com/palantir/witchcraft-go-server/v2/wrouter/whttprouter"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testBearerToken = "token"
)

func TestCommand_Echo(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"echo",
			"--bearer_token",
			testBearerToken,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("Echo", mock.Anything, bearertoken.Token(testBearerToken)).Return(nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("Echo", mock.Anything, bearertoken.Token(testBearerToken)).Return(fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "echo failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid param", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing required bearer token", func(t *testing.T) {
			args := []string{
				"",
				"echo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "bearer_token is a required argument")
		})
	})
}

func TestCommand_EchoStrings(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"echoStrings",
			"--body",
			`["string1","string2"]`,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoStrings", mock.Anything, []string{"string1", "string2"}).Return([]string{"string1", "string2"}, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "[\n    \"string1\",\n    \"string2\"\n]\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoStrings", mock.Anything, []string{"string1", "string2"}).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "echoStrings failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"echoStrings",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "body is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"echoStrings",
				"--body",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "invalid value for body argument")
		})
	})
}

func TestCommand_EchoCustomObject(t *testing.T) {
	customObject := api.CustomObject{
		Data: []byte("bytes"),
	}
	customObjectBytes, err := json.Marshal(customObject)
	require.NoError(t, err)
	tmpDir := t.TempDir()
	filepath := path.Join(tmpDir, "data")
	require.NoError(t, os.WriteFile(filepath, customObjectBytes, 0755))

	t.Run("valid input - json-encoded string", func(t *testing.T) {
		args := []string{
			"",
			"echoCustomObject",
			"--body",
			`{"data": "Ynl0ZXM="}`,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoCustomObject", mock.Anything, &customObject).Return(&customObject, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "{\n    \"data\": \"Ynl0ZXM=\"\n}\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoCustomObject", mock.Anything, &customObject).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "echoCustomObject failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("valid input - from file", func(t *testing.T) {
		args := []string{
			"",
			"echoCustomObject",
			"--body",
			"@" + filepath,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoCustomObject", mock.Anything, &customObject).Return(&customObject, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "{\n    \"data\": \"Ynl0ZXM=\"\n}\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoCustomObject", mock.Anything, &customObject).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "echoCustomObject failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("nil optional param and return value works", func(t *testing.T) {
		args := []string{
			"",
			"echoCustomObject",
		}
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		impl.On("EchoCustomObject", mock.Anything, (*api.CustomObject)(nil)).Return(nil, nil).Times(1)
		executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "null\n")
	})
	t.Run("invalid input", func(t *testing.T) {
		args := []string{
			"",
			"echoCustomObject",
			"--body",
			"foo",
		}
		t.Run("invalid body param value", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			executeAndAssertError(t, testServiceCommand, args, impl, "invalid value for body argument")
		})
	})
}

func TestCommand_EchoOptionalAlias(t *testing.T) {
	val := 123
	optionalIntAlias := api.OptionalIntegerAlias{
		Value: &val,
	}

	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"echoOptionalAlias",
			"--body",
			`123`,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoOptionalAlias", mock.Anything, optionalIntAlias).Return(optionalIntAlias, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "123\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoOptionalAlias", mock.Anything, optionalIntAlias).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "echoOptionalAlias failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("nil optional param and return value works", func(t *testing.T) {
		args := []string{
			"",
			"echoOptionalAlias",
		}
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		impl.On("EchoOptionalAlias", mock.Anything, api.OptionalIntegerAlias{}).Return(api.OptionalIntegerAlias{}, nil).Times(1)
		executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "null\n")
	})
	t.Run("invalid input", func(t *testing.T) {
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"echoOptionalAlias",
				"--body",
				"foo",
			}
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			executeAndAssertError(t, testServiceCommand, args, impl, `failed to parse "body" as integer`)
		})
	})
}

func TestCommand_EchoOptionalListAlias(t *testing.T) {
	val := []string{"string1", "string2"}
	optionalListAlias := api.OptionalListAlias{
		Value: &val,
	}

	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"echoOptionalListAlias",
			"--body",
			`["string1", "string2"]`,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoOptionalListAlias", mock.Anything, optionalListAlias).Return(optionalListAlias, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "[\n    \"string1\",\n    \"string2\"\n]\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("EchoOptionalListAlias", mock.Anything, optionalListAlias).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "echoOptionalListAlias failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("nil optional param and return value works", func(t *testing.T) {
		args := []string{
			"",
			"echoOptionalListAlias",
		}
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		impl.On("EchoOptionalListAlias", mock.Anything, api.OptionalListAlias{}).Return(api.OptionalListAlias{}, nil).Times(1)
		executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "null\n")
	})
	t.Run("invalid input", func(t *testing.T) {
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"echoOptionalListAlias",
				"--body",
				"foo",
			}
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			executeAndAssertError(t, testServiceCommand, args, impl, `invalid value for body argument`)
		})
	})
}

func TestCommand_GetPathParam(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"getPathParam",
			"--myPathParam",
			`value`,
			"--bearer_token",
			testBearerToken,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetPathParam", mock.Anything, bearertoken.Token(testBearerToken), "value").Return(nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetPathParam", mock.Anything, bearertoken.Token(testBearerToken), "value").Return(fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getPathParam failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"getPathParam",
				"--bearer_token",
				testBearerToken,
			}
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			executeAndAssertError(t, testServiceCommand, args, impl, "myPathParam is a required argument")
		})
		t.Run("missing bearer token", func(t *testing.T) {
			args := []string{
				"",
				"getPathParam",
				"--myPathParam",
				`value`,
			}
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			executeAndAssertError(t, testServiceCommand, args, impl, "bearer_token is a required argument")
		})
	})
}

func TestCommand_GetListBoolean(t *testing.T) {
	jsonVal := `[true, false]`
	tmpDir := t.TempDir()
	filepath := path.Join(tmpDir, "data")
	require.NoError(t, os.WriteFile(filepath, []byte(jsonVal), 0755))
	t.Run("valid input - json-encoded string", func(t *testing.T) {
		args := []string{
			"",
			"getListBoolean",
			"--myQueryParam1",
			jsonVal,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetListBoolean", mock.Anything, []bool{true, false}).Return([]bool{true, false}, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "[\n    true,\n    false\n]\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetListBoolean", mock.Anything, []bool{true, false}).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getListBoolean failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("valid input - from file", func(t *testing.T) {
		args := []string{
			"",
			"getListBoolean",
			"--myQueryParam1",
			"@" + filepath,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetListBoolean", mock.Anything, []bool{true, false}).Return([]bool{true, false}, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "[\n    true,\n    false\n]\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetListBoolean", mock.Anything, []bool{true, false}).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getListBoolean failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"getListBoolean",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myQueryParam1 is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"getListBoolean",
				"--myQueryParam1",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "invalid value for myQueryParam1 argument")
		})
	})
}

func TestCommand_PutMapStringString(t *testing.T) {
	testMap := map[string]string{
		"key": "value",
		"foo": "bar",
	}
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"putMapStringString",
			"--myParam",
			`{"key": "value", "foo": "bar"}`,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutMapStringString", mock.Anything, testMap).Return(testMap, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "{\n    \"foo\": \"bar\",\n    \"key\": \"value\"\n}\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutMapStringString", mock.Anything, testMap).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "putMapStringString failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"putMapStringString",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"putMapStringString",
				"--myParam",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "invalid value for myParam argument")
		})
	})
}

func TestCommand_PutMapStringAny(t *testing.T) {
	testMap := map[string]any{
		"key": "value",
		"foo": json.Number("123"),
		"bar": true,
	}
	jsonStringArg := `{"key": "value", "foo": 123, "bar": true}`
	tmpDir := t.TempDir()
	filepath := path.Join(tmpDir, "data")
	require.NoError(t, os.WriteFile(filepath, []byte(jsonStringArg), 0755))

	t.Run("valid input - json-encoded string", func(t *testing.T) {
		args := []string{
			"",
			"putMapStringAny",
			"--myParam",
			jsonStringArg,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutMapStringAny", mock.Anything, testMap).Return(testMap, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "{\n    \"bar\": true,\n    \"foo\": 123,\n    \"key\": \"value\"\n}\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutMapStringAny", mock.Anything, testMap).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "putMapStringAny failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("valid input - from file", func(t *testing.T) {
		args := []string{
			"",
			"putMapStringAny",
			"--myParam",
			"@" + filepath,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutMapStringAny", mock.Anything, testMap).Return(testMap, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "{\n    \"bar\": true,\n    \"foo\": 123,\n    \"key\": \"value\"\n}\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutMapStringAny", mock.Anything, testMap).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "putMapStringAny failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("valid input - from stdin", func(t *testing.T) {
		args := []string{
			"",
			"putMapStringAny",
			"--myParam",
			"@-",
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutMapStringAny", mock.Anything, testMap).Return(testMap, nil).Times(1)
			executeAndAssertSuccessAndOutputWithStdin(t, testServiceCommand, args, impl, "{\n    \"bar\": true,\n    \"foo\": 123,\n    \"key\": \"value\"\n}\n", strings.NewReader(jsonStringArg))
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutMapStringAny", mock.Anything, testMap).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertErrorWithStdin(t, testServiceCommand, args, impl, "putMapStringAny failed: httpclient request failed: 500 Internal Server Error", strings.NewReader(jsonStringArg))
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"putMapStringAny",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"putMapStringAny",
				"--myParam",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "invalid value for myParam argument")
		})
	})
}

func TestCommand_GetDateTime(t *testing.T) {
	dtArg := "2017-01-02T04:04:04Z"
	dt, err := datetime.ParseDateTime(dtArg)
	require.NoError(t, err)
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"getDateTime",
			"--myParam",
			dtArg,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetDateTime", mock.Anything, dt).Return(dt, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "2017-01-02T04:04:04Z\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetDateTime", mock.Anything, dt).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getDateTime failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"getDateTime",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"getDateTime",
				"--myParam",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "failed to parse \"myParam\" as datetime")
		})
	})
}

func TestCommand_GetDouble(t *testing.T) {
	doubleArg := "123456.789012"
	floatVal, err := strconv.ParseFloat(doubleArg, 64)
	require.NoError(t, err)
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"getDouble",
			"--myParam",
			doubleArg,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetDouble", mock.Anything, floatVal).Return(floatVal, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "123456.789012\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetDouble", mock.Anything, floatVal).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getDouble failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"getDouble",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"getDouble",
				"--myParam",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "failed to parse \"myParam\" as double")
		})
	})
}

func TestCommand_GetRid(t *testing.T) {
	ridArg := "ri.service.instance.resource.8bbb03fa-f3d8-423c-bbe4-c072b939a8ba"
	ridVal, err := rid.ParseRID(ridArg)
	require.NoError(t, err)
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"getRid",
			"--myParam",
			ridArg,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetRid", mock.Anything, ridVal).Return(ridVal, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "ri.service.instance.resource.8bbb03fa-f3d8-423c-bbe4-c072b939a8ba\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetRid", mock.Anything, ridVal).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getRid failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"getRid",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"getRid",
				"--myParam",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "failed to parse \"myParam\" as rid")
		})
	})
}

func TestCommand_GetSafeLong(t *testing.T) {
	safeLongArg := "9007199254740991"
	safelongVal, err := safelong.ParseSafeLong(safeLongArg)
	require.NoError(t, err)
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"getSafeLong",
			"--myParam",
			safeLongArg,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetSafeLong", mock.Anything, safelongVal).Return(safelongVal, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "9007199254740991\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetSafeLong", mock.Anything, safelongVal).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getSafeLong failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"getSafeLong",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"getSafeLong",
				"--myParam",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "failed to parse \"myParam\" as safelong")
		})
	})
}

func TestCommand_GetUuid(t *testing.T) {
	uuidArg := "8bbb03fa-f3d8-423c-bbe4-c072b939a8ba"
	uuidVal, err := uuid.ParseUUID(uuidArg)
	require.NoError(t, err)
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"getUuid",
			"--myParam",
			uuidArg,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetUuid", mock.Anything, uuidVal).Return(uuidVal, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "8bbb03fa-f3d8-423c-bbe4-c072b939a8ba\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetUuid", mock.Anything, uuidVal).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getUuid failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"getUuid",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
		t.Run("invalid body param value", func(t *testing.T) {
			args := []string{
				"",
				"getUuid",
				"--myParam",
				"foo",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "failed to parse \"myParam\" as uuid")
		})
	})
}

func TestCommand_GetCustomEnum(t *testing.T) {
	customEnum := api.New_CustomEnum(api.CustomEnum_STATE1)
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"getEnum",
			"--myParam",
			"STATE1",
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetEnum", mock.Anything, customEnum).Return(customEnum, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "STATE1\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetEnum", mock.Anything, customEnum).Return(api.CustomEnum{}, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getEnum failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		impl, testServiceCommand := getMockClientAndTestCommand(t)
		t.Run("missing body param", func(t *testing.T) {
			args := []string{
				"",
				"getEnum",
			}
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
		t.Run("invalid body param value is just an unknown value in the enum", func(t *testing.T) {
			args := []string{
				"",
				"getEnum",
				"--myParam",
				"foo",
			}
			invalidValue := api.New_CustomEnum("FOO")
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetEnum", mock.Anything, invalidValue).Return(invalidValue, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "FOO\n")
		})
	})
}

func TestCommand_PutBinary(t *testing.T) {
	bytesVal := []byte("somebytes")
	base64Val := base64.StdEncoding.EncodeToString(bytesVal)
	tmpDir := t.TempDir()
	filepath := path.Join(tmpDir, "data")
	require.NoError(t, os.WriteFile(filepath, bytesVal, 0755))

	t.Run("valid input - base64", func(t *testing.T) {
		args := []string{
			"",
			"putBinary",
			"--myParam",
			base64Val,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutBinary", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				body, err := io.ReadAll(args[1].(io.ReadCloser))
				require.NoError(t, err)
				require.Equal(t, bytesVal, body)
			}).Return(io.NopCloser(bytes.NewReader(bytesVal)), nil).Once()
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "somebytes")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutBinary", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				body, err := io.ReadAll(args[1].(io.ReadCloser))
				require.NoError(t, err)
				require.Equal(t, bytesVal, body)
			}).Return(nil, fmt.Errorf("bad")).Once()
			executeAndAssertError(t, testServiceCommand, args, impl, "putBinary failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("valid input - file", func(t *testing.T) {
		args := []string{
			"",
			"putBinary",
			"--myParam",
			"@" + filepath,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			readCloser := io.NopCloser(bytes.NewReader(bytesVal))
			impl.On("PutBinary", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				body, err := io.ReadAll(args[1].(io.ReadCloser))
				require.NoError(t, err)
				require.Equal(t, bytesVal, body)
			}).Return(readCloser, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "somebytes")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutBinary", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				body, err := io.ReadAll(args[1].(io.ReadCloser))
				require.NoError(t, err)
				require.Equal(t, bytesVal, body)
			}).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "putBinary failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("valid input - stdin", func(t *testing.T) {
		args := []string{
			"",
			"putBinary",
			"--myParam",
			"@-",
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			readCloser := io.NopCloser(bytes.NewReader(bytesVal))
			impl.On("PutBinary", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				body, err := io.ReadAll(args[1].(io.ReadCloser))
				require.NoError(t, err)
				require.Equal(t, bytesVal, body)
			}).Return(readCloser, nil).Times(1)
			executeAndAssertSuccessAndOutputWithStdin(t, testServiceCommand, args, impl, "somebytes", bytes.NewReader(bytesVal))
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			readCloser := io.NopCloser(bytes.NewReader(bytesVal))
			impl.On("PutBinary", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				body, err := io.ReadAll(args[1].(io.ReadCloser))
				require.NoError(t, err)
				require.Equal(t, bytesVal, body)
			}).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertErrorWithStdin(t, testServiceCommand, args, impl, "putBinary failed: httpclient request failed: 500 Internal Server Error", readCloser)
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		t.Run("unknown file", func(t *testing.T) {
			args := []string{
				"",
				"putBinary",
				"--myParam",
				"@badpath",
			}
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			executeAndAssertError(t, testServiceCommand, args, impl, "failed to open file for argument myParam: open badpath: no such file or directory")
		})
		t.Run("missing param", func(t *testing.T) {
			args := []string{
				"",
				"putBinary",
			}
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			executeAndAssertError(t, testServiceCommand, args, impl, "myParam is a required argument")
		})
	})
}

func TestCommand_GetOptionalBinary(t *testing.T) {
	bytesVal := []byte("somebytes")
	readCloser := io.NopCloser(bytes.NewReader(bytesVal))
	t.Run("valid input", func(t *testing.T) {
		args := []string{
			"",
			"getOptionalBinary",
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetOptionalBinary", mock.Anything).Return(&readCloser, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "somebytes")
		})
		t.Run("success empty result", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetOptionalBinary", mock.Anything).Return(nil, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("GetOptionalBinary", mock.Anything).Return(nil, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "getOptionalBinary failed: httpclient request failed: 500 Internal Server Error")
		})
	})
}

func TestCommand_GetCustomUnion(t *testing.T) {
	customUnion := api.NewCustomUnionFromAsString("teststring")
	customUnionBytes, err := json.Marshal(customUnion)
	require.NoError(t, err)
	tmpDir := t.TempDir()
	filepath := path.Join(tmpDir, "data")
	require.NoError(t, os.WriteFile(filepath, customUnionBytes, 0755))

	t.Run("valid input - json-encoded string", func(t *testing.T) {
		args := []string{
			"",
			"putCustomUnion",
			"--myParam",
			`{"type": "asString", "asString": "teststring"}`,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutCustomUnion", mock.Anything, customUnion).Return(customUnion, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "{\n    \"type\": \"asString\",\n    \"asString\": \"teststring\"\n}\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutCustomUnion", mock.Anything, customUnion).Return(api.CustomUnion{}, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "putCustomUnion failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("valid input - json-encoded string", func(t *testing.T) {
		args := []string{
			"",
			"putCustomUnion",
			"--myParam",
			"@" + filepath,
		}
		t.Run("success", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutCustomUnion", mock.Anything, customUnion).Return(customUnion, nil).Times(1)
			executeAndAssertSuccessAndOutput(t, testServiceCommand, args, impl, "{\n    \"type\": \"asString\",\n    \"asString\": \"teststring\"\n}\n")
		})
		t.Run("client error", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			impl.On("PutCustomUnion", mock.Anything, customUnion).Return(api.CustomUnion{}, fmt.Errorf("bad")).Times(1)
			executeAndAssertError(t, testServiceCommand, args, impl, "putCustomUnion failed: httpclient request failed: 500 Internal Server Error")
		})
	})
	t.Run("invalid input", func(t *testing.T) {
		args := []string{
			"",
			"putCustomUnion",
			"--myParam",
			"foo",
		}
		t.Run("invalid body param value", func(t *testing.T) {
			impl, testServiceCommand := getMockClientAndTestCommand(t)
			executeAndAssertError(t, testServiceCommand, args, impl, "invalid value for myParam argument")
		})
	})
}

func executeAndAssertErrorWithStdin(t *testing.T, cmd *cobra.Command, args []string, client *api_mock.TestService, expectedErr string, stdin io.Reader) {
	_, err := executeCmd(t, cmd, args, stdin)
	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr)
	mock.AssertExpectationsForObjects(t, client)
}

func executeAndAssertError(t *testing.T, cmd *cobra.Command, args []string, client *api_mock.TestService, expectedErr string) {
	executeAndAssertErrorWithStdin(t, cmd, args, client, expectedErr, nil)
}

func executeAndAssertSuccessAndOutputWithStdin(t *testing.T, cmd *cobra.Command, args []string, impl *api_mock.TestService, expectedOutput string, stdin io.Reader) {
	buf, err := executeCmd(t, cmd, args, stdin)
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, buf.String())
	mock.AssertExpectationsForObjects(t, impl)
}

func executeAndAssertSuccessAndOutput(t *testing.T, cmd *cobra.Command, args []string, impl *api_mock.TestService, expectedOutput string) {
	executeAndAssertSuccessAndOutputWithStdin(t, cmd, args, impl, expectedOutput, nil)
}

func executeCmd(t *testing.T, cmd *cobra.Command, args []string, stdin io.Reader) (*bytes.Buffer, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	if stdin != nil {
		cmd.SetIn(stdin)
	}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf, err
}

func getMockClientAndTestCommand(t *testing.T) (*api_mock.TestService, *cobra.Command) {
	serverImpl := &api_mock.TestService{}
	router := wrouter.New(whttprouter.New())
	require.NoError(t, api.RegisterRoutesTestService(router, serverImpl))
	ts := httptest.NewServer(router)
	httpc, err := httpclient.NewClient(httpclient.WithBaseURLs([]string{ts.URL}), httpclient.WithMaxRetries(0))
	require.NoError(t, err)
	client := api.NewTestServiceClient(httpc)

	provider := testClientProvider{
		client: client,
	}
	t.Cleanup(ts.Close)
	return serverImpl, api.NewTestServiceCLICommandWithClientProvider(provider)
}

type testClientProvider struct {
	client api.TestServiceClient
}

func (t testClientProvider) Get(_ context.Context, _ *pflag.FlagSet) (api.TestServiceClient, error) {
	return t.client, nil
}
