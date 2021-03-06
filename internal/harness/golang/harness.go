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

// Package golang provides the Match Making Function service for Open Match golang harness.
package golang

import (
	"google.golang.org/grpc"
	"open-match.dev/open-match/internal/config"
	"open-match.dev/open-match/internal/pb"
	"open-match.dev/open-match/internal/rpc"

	"github.com/sirupsen/logrus"
)

var (
	harnessLogger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "matchfunction.golang.harness",
	})
)

// FunctionSettings is a collection of parameters used to customize matchfunction views.
type FunctionSettings struct {
	FunctionName string
	Func         matchFunction
}

// RunMatchFunction is a hook for the main() method in the main executable.
func RunMatchFunction(settings *FunctionSettings) {
	cfg, err := config.Read()
	if err != nil {
		harnessLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("cannot read configuration.")
	}
	p, err := rpc.NewServerParamsFromConfig(cfg, "api.functions")
	if err != nil {
		harnessLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("cannot construct server.")
	}

	if err := BindService(p, cfg, settings); err != nil {
		harnessLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("failed to bind functions service.")
	}

	rpc.MustServeForever(p)
}

// BindService creates the function service to the server Params.
func BindService(p *rpc.ServerParams, cfg config.View, fs *FunctionSettings) error {
	service, err := newMatchFunctionService(cfg, fs)
	if err != nil {
		return err
	}

	p.AddHandleFunc(func(s *grpc.Server) {
		pb.RegisterMatchFunctionServer(s, service)
	}, pb.RegisterMatchFunctionHandlerFromEndpoint)

	return nil
}
