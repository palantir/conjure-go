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

package verifier_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go/v6/conjure-go-verifier/conjure/verification/server"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/pkg/httpserver"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const (
	serverURI = "http://localhost:8000"
)

var (
	testDefinitions struct {
		Client server.ClientTestCases `yaml:"client"`
	}

	ignoredTestCases struct {
		Client server.IgnoredClientTestCases `yaml:"client"`
	}

	behaviors = map[bool]string{
		true:  "succeeded",
		false: "failed",
	}
)

func TestMain(m *testing.M) {
	os.Exit(runTestMain(m))
}

func runTestMain(m *testing.M) int {
	// start verification server using Docker if it is not already running
	dockerContainerID, err := startDockerServerIfNotRunning()
	if dockerContainerID != "" {
		// if verification server was started using "docker run", terminate container using "docker kill" on teardown
		defer func() {
			cmd := exec.Command("docker", "kill", dockerContainerID)
			if output, err := cmd.CombinedOutput(); err != nil {
				fmt.Printf("command %v failed:\nOutput:\n%s\nError:\n%v\n", cmd.Args, string(output), err)
			}
		}()
	}
	if err != nil {
		panic(fmt.Sprintf("failed to start docker: %v", err))
	}

	// read test cases from verification-server-test-cases.json using conjure-go generated definitions
	bytes, err := ioutil.ReadFile("verification-server-test-cases.json")
	if err != nil {
		panic(fmt.Sprintf("failed to read verification-server-test-cases.json: %v", err))
	}
	if err := json.Unmarshal(bytes, &testDefinitions); err != nil {
		panic(fmt.Sprintf("failed to unmarshal verification-server-test-cases.json: %v", err))
	}

	// read ignored test cases from ignored-test-cases.yml.
	ignoredBytes, err := ioutil.ReadFile("ignored-test-cases.yml")
	if err != nil {
		panic(fmt.Sprintf("failed to read ignored-test-cases.json: %v", err))
	}
	if err := yaml.Unmarshal(ignoredBytes, &ignoredTestCases); err != nil {
		panic(fmt.Sprintf("failed to unmarshal ignored-test-cases.json: %v", err))
	}

	return m.Run()
}

func startDockerServerIfNotRunning() (string, error) {
	// verification server is already running
	if resp, err := http.Get(serverURI + "/body/receiveDoubleExample/0"); err == nil && resp.StatusCode == http.StatusOK {
		return "", nil
	}

	// run verification server in docker
	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"-p",
		"8000:8000",
		fmt.Sprintf("palantirtechnologies/conjure-verification-server:%s", verificationServerVersion),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.WithStack(err)
	}
	dockerContainerID := strings.TrimSpace(string(output))

	serverReady := <-httpserver.URLReady(serverURI + "/body/receiveDoubleExample/0")
	if !serverReady {
		return dockerContainerID, errors.Errorf("timed out waiting for verification server to become available")
	}
	return dockerContainerID, nil
}

func TestAutoDeserialize(t *testing.T) {
	ctx := context.Background()
	client := server.NewAutoDeserializeServiceClient(newHTTPClient(t, serverURI))
	confirmClient := server.NewAutoDeserializeConfirmServiceClient(newHTTPClient(t, serverURI))

	for endpointName, posAndNegTestCases := range testDefinitions.Client.AutoDeserialize {
		//	we explicitly use the conjure-go lib function call to do this to keep
		//	this test consistent with its behavior.
		methodName := transforms.Export(string(endpointName))
		method := reflect.ValueOf(client).MethodByName(methodName)
		// in the positive case, the index should be the case's index in the
		// positive test case list, and in the negative case, it should be the
		// number of positive test cases plus the index in the negative test case
		// list.
		i := 0
		for _, casesAndType := range []struct {
			cases    []string
			positive bool
		}{
			{posAndNegTestCases.Positive, true},
			{posAndNegTestCases.Negative, false},
		} {
			for _, val := range casesAndType.cases {
				t.Run(fmt.Sprintf("%s %d", endpointName, i), func(t *testing.T) {
					response := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(i)})
					result, ok := response[0].Interface(), response[1].IsNil()
					got := behaviors[ok]
					want := behaviors[casesAndType.positive]
					isIgnored := false
					for _, ignoredVal := range ignoredTestCases.Client.AutoDeserialize[endpointName] {
						if val == ignoredVal {
							isIgnored = true
						}
					}

					if isIgnored {
						wronglyIgnored := func() {
							t.Helper()
							t.Errorf("%v %d was ignored so is expected to misbehave, however it %v as it should in a correct implementation. If this test case was fixed, please remove this test from ignored-test-cases.yml", endpointName, i, got)
						}
						// if this test case is ignored, we error if got and want are the *same*
						if got == want {
							// if this is a positive case and deserialize succeeded, send it to the confirmation endpoint to verify round-tripping fails as it should
							if ok && casesAndType.positive {
								// if this is a positive case and deserialize succeeded, send it to the confirmation endpoint to verify round-tripping fails as it should
								if err := confirmClient.Confirm(ctx, endpointName, i, result); err == nil {
									wronglyIgnored()
								}
							} else {
								wronglyIgnored()
							}
						}
					} else {
						// in the usual case, we error if the got and want are different
						if got != want {
							t.Errorf("%v %d incorrectly %s: input=%v result=%v err=%v", endpointName, i, got, val, result, response[1].Interface())
						}
						// if this is a positive case, send it to the confirmation endpoint to verify round-tripping
						if ok && casesAndType.positive {
							if err := confirmClient.Confirm(ctx, endpointName, i, result); err != nil {
								t.Errorf("%v %d confirmation failed: input=%v result=%v err=%v", endpointName, i, val, result, err.Error())
							}
						}
					}
					i++
				})
			}
		}
	}
}

func TestSingleHeader(t *testing.T) {
	testSingleArg(t, server.NewSingleHeaderServiceClient(newHTTPClient(t, serverURI)), testDefinitions.Client.SingleHeaderService, ignoredTestCases.Client.SingleHeaderService)
}

func TestSinglePathParam(t *testing.T) {
	testSingleArg(t, server.NewSinglePathParamServiceClient(newHTTPClient(t, serverURI)), testDefinitions.Client.SinglePathParamService, ignoredTestCases.Client.SinglePathParamService)
}

func TestSingleQueryParam(t *testing.T) {
	testSingleArg(t, server.NewSingleQueryParamServiceClient(newHTTPClient(t, serverURI)), testDefinitions.Client.SingleQueryParamService, ignoredTestCases.Client.SingleQueryParamService)
}

func testSingleArg(t *testing.T, service interface{}, tests map[server.EndpointName][]string, ignored map[server.EndpointName][]string) {
	for endpointName, vals := range tests {
		methodName := transforms.Export(string(endpointName))
		method := reflect.ValueOf(service).MethodByName(methodName)
		for i, val := range vals {
			// These all have the context as the first argument and the reflected value as the 2nd
			argType := method.Type().In(2)
			arg := reflect.New(argType).Interface()
			err := json.Unmarshal([]byte(val), *&arg)
			require.NoError(t, err, "%v %d failed to unmarshal %v", endpointName, i, val)
			if err != nil {
				continue
			}
			in := []reflect.Value{
				reflect.ValueOf(context.Background()), reflect.ValueOf(i), reflect.ValueOf(arg).Elem(),
			}
			response := method.Call(in)
			isIgnored := false
			for _, ignoredVal := range ignored[endpointName] {
				if val == ignoredVal {
					isIgnored = true
				}
			}
			responseErr := response[0].Interface()
			got := behaviors[responseErr == nil]
			want := behaviors[true]
			if isIgnored {
				if got == want {
					t.Errorf("%v %d was ignored so is expected to misbehave, however it %v as it should in a correct implementation. If this test case was fixed, please remove this test from ignored-test-cases.yml", endpointName, i, got)
				}
			} else {
				if got != want {
					t.Errorf("%v %d failed call with %v: %v", endpointName, i, val, responseErr)
				}
			}
		}
	}
}

func newHTTPClient(t *testing.T, url string) httpclient.Client {
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{url}),
	)
	require.NoError(t, err)
	return httpClient
}
