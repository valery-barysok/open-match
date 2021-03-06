// Copyright 2019 Google LLC
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

package rpc

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"open-match.dev/open-match/internal/pb"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	shellTesting "open-match.dev/open-match/internal/testing"
	netlistenerTesting "open-match.dev/open-match/internal/util/netlistener/testing"
)

func TestInsecureStartStop(t *testing.T) {
	assert := assert.New(t)
	grpcLh := netlistenerTesting.MustListen()
	httpLh := netlistenerTesting.MustListen()
	ff := &shellTesting.FakeFrontend{}

	params := NewServerParamsFromListeners(grpcLh, httpLh)
	params.AddHandleFunc(func(s *grpc.Server) {
		pb.RegisterFrontendServer(s, ff)
	}, pb.RegisterFrontendHandlerFromEndpoint)
	s := newInsecureServer(grpcLh, httpLh)
	defer s.stop()
	waitForStart, err := s.start(params)
	assert.Nil(err)
	waitForStart()

	conn, err := grpc.Dial(fmt.Sprintf(":%d", grpcLh.Number()), grpc.WithInsecure())
	assert.Nil(err)
	defer conn.Close()

	endpoint := fmt.Sprintf("http://localhost:%d", httpLh.Number())
	httpClient := &http.Client{
		Timeout: time.Second,
	}
	runGrpcWithProxyTests(assert, s, conn, httpClient, endpoint)
}
