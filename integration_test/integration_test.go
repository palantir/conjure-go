// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/godel/pkg/products/v2/products"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLI(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "TestCLIProject-")
	defer cleanup()
	require.NoError(t, err)

	cli, err := products.Bin("conjure-go")

	for currCaseNum, tc := range []struct {
		name   string
		irFile string
		want   map[string]string
	}{
		{
			name:   "base case",
			irFile: "testdata/example-service.json",
			want: map[string]string{
				"foundry/catalog/api/datasets/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package datasets

import (
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type BackingFileSystem struct {
	// The name by which this file system is identified.
	FileSystemId  string            ` + "`" + `json:"fileSystemId" conjure-docs:"The name by which this file system is identified."` + "`" + `
	BaseUri       string            ` + "`" + `json:"baseUri"` + "`" + `
	Configuration map[string]string ` + "`" + `json:"configuration"` + "`" + `
}

func (o BackingFileSystem) MarshalJSON() ([]byte, error) {
	if o.Configuration == nil {
		o.Configuration = make(map[string]string, 0)
	}
	type BackingFileSystemAlias BackingFileSystem
	return safejson.Marshal(BackingFileSystemAlias(o))
}

func (o *BackingFileSystem) UnmarshalJSON(data []byte) error {
	type BackingFileSystemAlias BackingFileSystem
	var rawBackingFileSystem BackingFileSystemAlias
	if err := safejson.Unmarshal(data, &rawBackingFileSystem); err != nil {
		return err
	}
	if rawBackingFileSystem.Configuration == nil {
		rawBackingFileSystem.Configuration = make(map[string]string, 0)
	}
	*o = BackingFileSystem(rawBackingFileSystem)
	return nil
}

func (o BackingFileSystem) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BackingFileSystem) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type Dataset struct {
	FileSystemId string ` + "`" + `json:"fileSystemId"` + "`" + `
	// Uniquely identifies this dataset.
	Rid rid.ResourceIdentifier ` + "`" + `json:"rid" conjure-docs:"Uniquely identifies this dataset."` + "`" + `
}

func (o Dataset) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Dataset) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
				"foundry/catalog/api/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type CreateDatasetRequest struct {
	FileSystemId string ` + "`" + `json:"fileSystemId"` + "`" + `
	Path         string ` + "`" + `json:"path"` + "`" + `
}

func (o CreateDatasetRequest) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *CreateDatasetRequest) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
				"test/api/services.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/palantir/conjure-go-runtime/conjure-go-client/httpclient"
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/rid"

	"github.com/palantir/conjure-go/integration_test/{{currCaseTmpDir}}/foundry/catalog/api"
	"github.com/palantir/conjure-go/integration_test/{{currCaseTmpDir}}/foundry/catalog/api/datasets"
)

// A Markdown description of the service.
type ExampleServiceClient interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context, authHeader bearertoken.Token) (map[string]datasets.BackingFileSystem, error)
	CreateDataset(ctx context.Context, cookieToken bearertoken.Token, requestArg api.CreateDatasetRequest) (datasets.Dataset, error)
	GetDataset(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) (*datasets.Dataset, error)
	GetBranches(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) ([]string, error)
	// Gets all branches of this dataset.
	GetBranchesDeprecated(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) ([]string, error)
	ResolveBranch(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier, branchArg string) (*string, error)
	TestParam(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*string, error)
	TestBoolean(ctx context.Context, authHeader bearertoken.Token) (bool, error)
	TestDouble(ctx context.Context, authHeader bearertoken.Token) (float64, error)
	TestInteger(ctx context.Context, authHeader bearertoken.Token) (int, error)
}

type exampleServiceClient struct {
	client httpclient.Client
}

func NewExampleServiceClient(client httpclient.Client) ExampleServiceClient {
	return &exampleServiceClient{client: client}
}

func (c *exampleServiceClient) GetFileSystems(ctx context.Context, authHeader bearertoken.Token) (map[string]datasets.BackingFileSystem, error) {
	var returnVal map[string]datasets.BackingFileSystem
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetFileSystems"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/fileSystems"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make(map[string]datasets.BackingFileSystem, 0)
	}
	return returnVal, nil
}

func (c *exampleServiceClient) CreateDataset(ctx context.Context, cookieToken bearertoken.Token, requestArg api.CreateDatasetRequest) (datasets.Dataset, error) {
	var defaultReturnVal datasets.Dataset
	var returnVal *datasets.Dataset
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("CreateDataset"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithHeader("Cookie", fmt.Sprint("PALANTIR_TOKEN=", cookieToken)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(requestArg))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return defaultReturnVal, err
	}
	_ = resp
	if returnVal == nil {
		return defaultReturnVal, fmt.Errorf("returnVal cannot be nil")
	}
	return *returnVal, nil
}

func (c *exampleServiceClient) GetDataset(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) (*datasets.Dataset, error) {
	var returnVal *datasets.Dataset
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetDataset"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s", url.PathEscape(fmt.Sprint(datasetRidArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	return returnVal, nil
}

func (c *exampleServiceClient) GetBranches(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) ([]string, error) {
	var returnVal []string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetBranches"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s/branches", url.PathEscape(fmt.Sprint(datasetRidArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make([]string, 0)
	}
	return returnVal, nil
}

func (c *exampleServiceClient) GetBranchesDeprecated(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) ([]string, error) {
	var returnVal []string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetBranchesDeprecated"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s/branchesDeprecated", url.PathEscape(fmt.Sprint(datasetRidArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make([]string, 0)
	}
	return returnVal, nil
}

func (c *exampleServiceClient) ResolveBranch(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier, branchArg string) (*string, error) {
	var returnVal *string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("ResolveBranch"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s/branches/%s/resolve", url.PathEscape(fmt.Sprint(datasetRidArg)), url.PathEscape(fmt.Sprint(branchArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	return returnVal, nil
}

func (c *exampleServiceClient) TestParam(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*string, error) {
	var returnVal *string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("TestParam"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s/testParam", url.PathEscape(fmt.Sprint(datasetRidArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	return returnVal, nil
}

func (c *exampleServiceClient) TestBoolean(ctx context.Context, authHeader bearertoken.Token) (bool, error) {
	var defaultReturnVal bool
	var returnVal *bool
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("TestBoolean"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/boolean"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return defaultReturnVal, err
	}
	_ = resp
	if returnVal == nil {
		return defaultReturnVal, fmt.Errorf("returnVal cannot be nil")
	}
	return *returnVal, nil
}

func (c *exampleServiceClient) TestDouble(ctx context.Context, authHeader bearertoken.Token) (float64, error) {
	var defaultReturnVal float64
	var returnVal *float64
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("TestDouble"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/double"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return defaultReturnVal, err
	}
	_ = resp
	if returnVal == nil {
		return defaultReturnVal, fmt.Errorf("returnVal cannot be nil")
	}
	return *returnVal, nil
}

func (c *exampleServiceClient) TestInteger(ctx context.Context, authHeader bearertoken.Token) (int, error) {
	var defaultReturnVal int
	var returnVal *int
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("TestInteger"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/integer"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return defaultReturnVal, err
	}
	_ = resp
	if returnVal == nil {
		return defaultReturnVal, fmt.Errorf("returnVal cannot be nil")
	}
	return *returnVal, nil
}

// A Markdown description of the service.
type ExampleServiceClientWithAuth interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context) (map[string]datasets.BackingFileSystem, error)
	CreateDataset(ctx context.Context, requestArg api.CreateDatasetRequest) (datasets.Dataset, error)
	GetDataset(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*datasets.Dataset, error)
	GetBranches(ctx context.Context, datasetRidArg rid.ResourceIdentifier) ([]string, error)
	// Gets all branches of this dataset.
	GetBranchesDeprecated(ctx context.Context, datasetRidArg rid.ResourceIdentifier) ([]string, error)
	ResolveBranch(ctx context.Context, datasetRidArg rid.ResourceIdentifier, branchArg string) (*string, error)
	TestParam(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*string, error)
	TestBoolean(ctx context.Context) (bool, error)
	TestDouble(ctx context.Context) (float64, error)
	TestInteger(ctx context.Context) (int, error)
}

func NewExampleServiceClientWithAuth(client ExampleServiceClient, authHeader bearertoken.Token, cookieToken bearertoken.Token) ExampleServiceClientWithAuth {
	return &exampleServiceClientWithAuth{client: client, authHeader: authHeader, cookieToken: cookieToken}
}

type exampleServiceClientWithAuth struct {
	client      ExampleServiceClient
	authHeader  bearertoken.Token
	cookieToken bearertoken.Token
}

func (c *exampleServiceClientWithAuth) GetFileSystems(ctx context.Context) (map[string]datasets.BackingFileSystem, error) {
	return c.client.GetFileSystems(ctx, c.authHeader)
}

func (c *exampleServiceClientWithAuth) CreateDataset(ctx context.Context, requestArg api.CreateDatasetRequest) (datasets.Dataset, error) {
	return c.client.CreateDataset(ctx, c.cookieToken, requestArg)
}

func (c *exampleServiceClientWithAuth) GetDataset(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*datasets.Dataset, error) {
	return c.client.GetDataset(ctx, c.authHeader, datasetRidArg)
}

func (c *exampleServiceClientWithAuth) GetBranches(ctx context.Context, datasetRidArg rid.ResourceIdentifier) ([]string, error) {
	return c.client.GetBranches(ctx, c.authHeader, datasetRidArg)
}

func (c *exampleServiceClientWithAuth) GetBranchesDeprecated(ctx context.Context, datasetRidArg rid.ResourceIdentifier) ([]string, error) {
	return c.client.GetBranchesDeprecated(ctx, c.authHeader, datasetRidArg)
}

func (c *exampleServiceClientWithAuth) ResolveBranch(ctx context.Context, datasetRidArg rid.ResourceIdentifier, branchArg string) (*string, error) {
	return c.client.ResolveBranch(ctx, c.authHeader, datasetRidArg, branchArg)
}

func (c *exampleServiceClientWithAuth) TestParam(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*string, error) {
	return c.client.TestParam(ctx, datasetRidArg)
}

func (c *exampleServiceClientWithAuth) TestBoolean(ctx context.Context) (bool, error) {
	return c.client.TestBoolean(ctx, c.authHeader)
}

func (c *exampleServiceClientWithAuth) TestDouble(ctx context.Context) (float64, error) {
	return c.client.TestDouble(ctx, c.authHeader)
}

func (c *exampleServiceClientWithAuth) TestInteger(ctx context.Context) (int, error) {
	return c.client.TestInteger(ctx, c.authHeader)
}
`,
			},
		},
	} {
		currCaseTmpDir := path.Join(tmpDir, fmt.Sprintf("case-%d", currCaseNum))
		err := os.Mkdir(currCaseTmpDir, 0755)
		require.NoError(t, err, "Case %d: %s", currCaseNum, tc.name)

		outputDir := path.Join(currCaseTmpDir, "output")
		err = os.Mkdir(outputDir, 0755)
		require.NoError(t, err, "Case %d: %s", currCaseNum, tc.name)

		cmd := exec.Command(cli, "--output", outputDir, tc.irFile)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "command %v failed with output:\n%s", cmd.Args, string(output))
		for k, wantSrc := range tc.want {
			wantSrc = strings.Replace(wantSrc, "{{currCaseTmpDir}}", path.Join(currCaseTmpDir, "output"), -1)
			bytes, err := ioutil.ReadFile(path.Join(currCaseTmpDir, "output", k))
			require.NoError(t, err)
			gotSrc := string(bytes)
			assert.Equal(t, strings.Split(wantSrc, "\n"), strings.Split(gotSrc, "\n"), "Case %d: %s: Unexpected content for file %s:\n%s\nWanted:\n%s", currCaseNum, tc.name, k, gotSrc, wantSrc)
		}
	}
}

func TestCLIInModule(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "TestCLIProject-")
	defer cleanup()
	require.NoError(t, err)

	cli, err := products.Bin("conjure-go")

	for currCaseNum, tc := range []struct {
		name   string
		irFile string
		want   map[string]string
	}{
		{
			name:   "base case",
			irFile: "testdata/example-service.json",
			want: map[string]string{
				"foundry/catalog/api/datasets/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package datasets

import (
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type BackingFileSystem struct {
	// The name by which this file system is identified.
	FileSystemId  string            ` + "`" + `json:"fileSystemId" conjure-docs:"The name by which this file system is identified."` + "`" + `
	BaseUri       string            ` + "`" + `json:"baseUri"` + "`" + `
	Configuration map[string]string ` + "`" + `json:"configuration"` + "`" + `
}

func (o BackingFileSystem) MarshalJSON() ([]byte, error) {
	if o.Configuration == nil {
		o.Configuration = make(map[string]string, 0)
	}
	type BackingFileSystemAlias BackingFileSystem
	return safejson.Marshal(BackingFileSystemAlias(o))
}

func (o *BackingFileSystem) UnmarshalJSON(data []byte) error {
	type BackingFileSystemAlias BackingFileSystem
	var rawBackingFileSystem BackingFileSystemAlias
	if err := safejson.Unmarshal(data, &rawBackingFileSystem); err != nil {
		return err
	}
	if rawBackingFileSystem.Configuration == nil {
		rawBackingFileSystem.Configuration = make(map[string]string, 0)
	}
	*o = BackingFileSystem(rawBackingFileSystem)
	return nil
}

func (o BackingFileSystem) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BackingFileSystem) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type Dataset struct {
	FileSystemId string ` + "`" + `json:"fileSystemId"` + "`" + `
	// Uniquely identifies this dataset.
	Rid rid.ResourceIdentifier ` + "`" + `json:"rid" conjure-docs:"Uniquely identifies this dataset."` + "`" + `
}

func (o Dataset) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Dataset) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
				"foundry/catalog/api/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type CreateDatasetRequest struct {
	FileSystemId string ` + "`" + `json:"fileSystemId"` + "`" + `
	Path         string ` + "`" + `json:"path"` + "`" + `
}

func (o CreateDatasetRequest) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *CreateDatasetRequest) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
				"test/api/services.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/palantir/conjure-go-runtime/conjure-go-client/httpclient"
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/rid"
	"{{currModulePath}}/foundry/catalog/api"
	"{{currModulePath}}/foundry/catalog/api/datasets"
)

// A Markdown description of the service.
type ExampleServiceClient interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context, authHeader bearertoken.Token) (map[string]datasets.BackingFileSystem, error)
	CreateDataset(ctx context.Context, cookieToken bearertoken.Token, requestArg api.CreateDatasetRequest) (datasets.Dataset, error)
	GetDataset(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) (*datasets.Dataset, error)
	GetBranches(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) ([]string, error)
	// Gets all branches of this dataset.
	GetBranchesDeprecated(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) ([]string, error)
	ResolveBranch(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier, branchArg string) (*string, error)
	TestParam(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*string, error)
	TestBoolean(ctx context.Context, authHeader bearertoken.Token) (bool, error)
	TestDouble(ctx context.Context, authHeader bearertoken.Token) (float64, error)
	TestInteger(ctx context.Context, authHeader bearertoken.Token) (int, error)
}

type exampleServiceClient struct {
	client httpclient.Client
}

func NewExampleServiceClient(client httpclient.Client) ExampleServiceClient {
	return &exampleServiceClient{client: client}
}

func (c *exampleServiceClient) GetFileSystems(ctx context.Context, authHeader bearertoken.Token) (map[string]datasets.BackingFileSystem, error) {
	var returnVal map[string]datasets.BackingFileSystem
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetFileSystems"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/fileSystems"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make(map[string]datasets.BackingFileSystem, 0)
	}
	return returnVal, nil
}

func (c *exampleServiceClient) CreateDataset(ctx context.Context, cookieToken bearertoken.Token, requestArg api.CreateDatasetRequest) (datasets.Dataset, error) {
	var defaultReturnVal datasets.Dataset
	var returnVal *datasets.Dataset
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("CreateDataset"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithHeader("Cookie", fmt.Sprint("PALANTIR_TOKEN=", cookieToken)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(requestArg))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return defaultReturnVal, err
	}
	_ = resp
	if returnVal == nil {
		return defaultReturnVal, fmt.Errorf("returnVal cannot be nil")
	}
	return *returnVal, nil
}

func (c *exampleServiceClient) GetDataset(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) (*datasets.Dataset, error) {
	var returnVal *datasets.Dataset
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetDataset"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s", url.PathEscape(fmt.Sprint(datasetRidArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	return returnVal, nil
}

func (c *exampleServiceClient) GetBranches(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) ([]string, error) {
	var returnVal []string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetBranches"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s/branches", url.PathEscape(fmt.Sprint(datasetRidArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make([]string, 0)
	}
	return returnVal, nil
}

func (c *exampleServiceClient) GetBranchesDeprecated(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier) ([]string, error) {
	var returnVal []string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetBranchesDeprecated"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s/branchesDeprecated", url.PathEscape(fmt.Sprint(datasetRidArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make([]string, 0)
	}
	return returnVal, nil
}

func (c *exampleServiceClient) ResolveBranch(ctx context.Context, authHeader bearertoken.Token, datasetRidArg rid.ResourceIdentifier, branchArg string) (*string, error) {
	var returnVal *string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("ResolveBranch"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s/branches/%s/resolve", url.PathEscape(fmt.Sprint(datasetRidArg)), url.PathEscape(fmt.Sprint(branchArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	return returnVal, nil
}

func (c *exampleServiceClient) TestParam(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*string, error) {
	var returnVal *string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("TestParam"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets/%s/testParam", url.PathEscape(fmt.Sprint(datasetRidArg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	return returnVal, nil
}

func (c *exampleServiceClient) TestBoolean(ctx context.Context, authHeader bearertoken.Token) (bool, error) {
	var defaultReturnVal bool
	var returnVal *bool
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("TestBoolean"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/boolean"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return defaultReturnVal, err
	}
	_ = resp
	if returnVal == nil {
		return defaultReturnVal, fmt.Errorf("returnVal cannot be nil")
	}
	return *returnVal, nil
}

func (c *exampleServiceClient) TestDouble(ctx context.Context, authHeader bearertoken.Token) (float64, error) {
	var defaultReturnVal float64
	var returnVal *float64
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("TestDouble"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/double"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return defaultReturnVal, err
	}
	_ = resp
	if returnVal == nil {
		return defaultReturnVal, fmt.Errorf("returnVal cannot be nil")
	}
	return *returnVal, nil
}

func (c *exampleServiceClient) TestInteger(ctx context.Context, authHeader bearertoken.Token) (int, error) {
	var defaultReturnVal int
	var returnVal *int
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("TestInteger"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/integer"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return defaultReturnVal, err
	}
	_ = resp
	if returnVal == nil {
		return defaultReturnVal, fmt.Errorf("returnVal cannot be nil")
	}
	return *returnVal, nil
}

// A Markdown description of the service.
type ExampleServiceClientWithAuth interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context) (map[string]datasets.BackingFileSystem, error)
	CreateDataset(ctx context.Context, requestArg api.CreateDatasetRequest) (datasets.Dataset, error)
	GetDataset(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*datasets.Dataset, error)
	GetBranches(ctx context.Context, datasetRidArg rid.ResourceIdentifier) ([]string, error)
	// Gets all branches of this dataset.
	GetBranchesDeprecated(ctx context.Context, datasetRidArg rid.ResourceIdentifier) ([]string, error)
	ResolveBranch(ctx context.Context, datasetRidArg rid.ResourceIdentifier, branchArg string) (*string, error)
	TestParam(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*string, error)
	TestBoolean(ctx context.Context) (bool, error)
	TestDouble(ctx context.Context) (float64, error)
	TestInteger(ctx context.Context) (int, error)
}

func NewExampleServiceClientWithAuth(client ExampleServiceClient, authHeader bearertoken.Token, cookieToken bearertoken.Token) ExampleServiceClientWithAuth {
	return &exampleServiceClientWithAuth{client: client, authHeader: authHeader, cookieToken: cookieToken}
}

type exampleServiceClientWithAuth struct {
	client      ExampleServiceClient
	authHeader  bearertoken.Token
	cookieToken bearertoken.Token
}

func (c *exampleServiceClientWithAuth) GetFileSystems(ctx context.Context) (map[string]datasets.BackingFileSystem, error) {
	return c.client.GetFileSystems(ctx, c.authHeader)
}

func (c *exampleServiceClientWithAuth) CreateDataset(ctx context.Context, requestArg api.CreateDatasetRequest) (datasets.Dataset, error) {
	return c.client.CreateDataset(ctx, c.cookieToken, requestArg)
}

func (c *exampleServiceClientWithAuth) GetDataset(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*datasets.Dataset, error) {
	return c.client.GetDataset(ctx, c.authHeader, datasetRidArg)
}

func (c *exampleServiceClientWithAuth) GetBranches(ctx context.Context, datasetRidArg rid.ResourceIdentifier) ([]string, error) {
	return c.client.GetBranches(ctx, c.authHeader, datasetRidArg)
}

func (c *exampleServiceClientWithAuth) GetBranchesDeprecated(ctx context.Context, datasetRidArg rid.ResourceIdentifier) ([]string, error) {
	return c.client.GetBranchesDeprecated(ctx, c.authHeader, datasetRidArg)
}

func (c *exampleServiceClientWithAuth) ResolveBranch(ctx context.Context, datasetRidArg rid.ResourceIdentifier, branchArg string) (*string, error) {
	return c.client.ResolveBranch(ctx, c.authHeader, datasetRidArg, branchArg)
}

func (c *exampleServiceClientWithAuth) TestParam(ctx context.Context, datasetRidArg rid.ResourceIdentifier) (*string, error) {
	return c.client.TestParam(ctx, datasetRidArg)
}

func (c *exampleServiceClientWithAuth) TestBoolean(ctx context.Context) (bool, error) {
	return c.client.TestBoolean(ctx, c.authHeader)
}

func (c *exampleServiceClientWithAuth) TestDouble(ctx context.Context) (float64, error) {
	return c.client.TestDouble(ctx, c.authHeader)
}

func (c *exampleServiceClientWithAuth) TestInteger(ctx context.Context) (int, error) {
	return c.client.TestInteger(ctx, c.authHeader)
}
`,
			},
		},
	} {
		currCaseTmpDir := path.Join(tmpDir, fmt.Sprintf("case-%d", currCaseNum))
		err := os.Mkdir(currCaseTmpDir, 0755)
		require.NoError(t, err, "Case %d: %s", currCaseNum, tc.name)
		modulePath := fmt.Sprintf("github.com/test-org/test-repo-%d", currCaseNum)
		err = ioutil.WriteFile(path.Join(currCaseTmpDir, "go.mod"), []byte("module "+modulePath), 0644)
		require.NoError(t, err, "Case %d: %s", currCaseNum, tc.name)

		outputDir := path.Join(currCaseTmpDir, "output")
		err = os.Mkdir(outputDir, 0755)
		require.NoError(t, err, "Case %d: %s", currCaseNum, tc.name)

		cmd := exec.Command(cli, "--output", outputDir, tc.irFile)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "command %v failed with output:\n%s", cmd.Args, string(output))
		for k, wantSrc := range tc.want {
			wantSrc = strings.Replace(wantSrc, "{{currModulePath}}", path.Join(modulePath, "output"), -1)
			bytes, err := ioutil.ReadFile(path.Join(currCaseTmpDir, "output", k))
			require.NoError(t, err)
			gotSrc := string(bytes)
			assert.Equal(t, strings.Split(wantSrc, "\n"), strings.Split(gotSrc, "\n"), "Case %d: %s: Unexpected content for file %s:\n%s\nWanted:\n%s", currCaseNum, tc.name, k, gotSrc, wantSrc)
		}
	}
}
