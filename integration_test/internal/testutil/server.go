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

package testutil

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/pkg/httpserver"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-server/v2/config"
	"github.com/palantir/witchcraft-go-server/v2/witchcraft"
	"github.com/stretchr/testify/require"
)

func StartTestServer(t *testing.T, init witchcraft.InitFunc) (httpclient.Client, func()) {
	port, err := httpserver.AvailablePort()
	require.NoError(t, err)
	server := witchcraft.NewServer().
		WithInstallConfig(&config.Install{
			Server: config.Server{
				Address: "localhost",
				Port:    port,
			},
			UseConsoleLog: true,
		}).
		WithECVKeyProvider(witchcraft.ECVKeyNoOp()).
		WithRuntimeConfig(&config.Runtime{}).
		WithLoggerStdoutWriter(os.Stdout).
		WithDisableGoRuntimeMetrics().
		WithInitFunc(init).
		WithSelfSignedCertificate()

	serverChan := make(chan error)
	go func() {
		serverChan <- server.Start()
	}()
	client, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{fmt.Sprintf("https://localhost:%d", port)}),
		httpclient.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	)
	require.NoError(t, err)
	success := <-httpserver.Ready(func() (*http.Response, error) {
		resp, err := client.Get(context.Background(), httpclient.WithPath("/status/readiness"))
		return resp, err
	}, httpserver.WaitTimeoutParam(5*time.Second))
	if !success {
		errMsg := "timed out waiting for server to start"
		select {
		case err := <-serverChan:
			errMsg = fmt.Sprintf("%s: %+v", errMsg, err)
		default:
		}
		require.Fail(t, errMsg)
	}
	cleanup := func() {
		if err := server.Close(); err != nil {
			svc1log.FromContext(TestContext()).Error(err.Error(), svc1log.Stacktrace(err))
			panic(err)
		}
	}
	return client, cleanup
}
