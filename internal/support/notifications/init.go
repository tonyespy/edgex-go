/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package notifications

import (
	"context"
	"sync"

	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/delegate"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common/middleware/debugging"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/routing"

	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/di"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"

	"github.com/gorilla/mux"
)

// Bootstrap contains references to dependencies required by the BootstrapHandler.
type Bootstrap struct {
	muxRouter            *mux.Router
	inDebugMode          bool
	inAcceptanceTestMode bool
}

// NewBootstrap is a factory method that returns an initialized Bootstrap receiver struct.
func NewBootstrap(muxRouter *mux.Router, indDebugMode, inAcceptanceTestMode bool) *Bootstrap {
	return &Bootstrap{
		muxRouter:            muxRouter,
		inDebugMode:          indDebugMode,
		inAcceptanceTestMode: inAcceptanceTestMode,
	}
}

// BootstrapHandler fulfills the BootstrapHandler contract and performs initialization for the notifications service.
func (b *Bootstrap) BootstrapHandler(_ context.Context, _ *sync.WaitGroup, _ startup.Timer, dic *di.Container) bool {
	loadV1Routes(b.muxRouter, dic)
	b.loadV2Routes(dic, container.LoggingClientFrom(dic.Get))
	return true
}

// loadV2Routes creates a new command-query router and handles the related mux.Router initialization for API V2 routes.
func (b *Bootstrap) loadV2Routes(dic *di.Container, lc logger.LoggingClient) {
	handlers := []delegate.Handler{}
	if b.inDebugMode {
		handlers = append(handlers, debugging.New(lc).Handler)
	}

	routing.Initialize(
		dic,
		b.muxRouter,
		handlers,
		common.V2Routes(
			b.inAcceptanceTestMode,
			[]routing.Controller{},
		),
	)
}
