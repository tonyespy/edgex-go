/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
 * Copyright (c) 2020 Canonical
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

package command

import (
	"fmt"
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/delegate"
	useCase "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/usecases/core/command/awe"
	validator "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/validator/awe"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common/routable"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/handle"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/routing"

	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/gorilla/mux"
)

const (
	version  = application.Version2
	kind     = "awe"
	action   = application.ActionCommand
	method   = api.ActionCommand
	aweEndpoint = api.BaseURL + "/awe"
)

// controller contains references to dependencies required by the corresponding Controller contract implementation.
type aweController struct{}

// New is a factory function that returns an initialized controller receiver struct.
func NewAWE() *aweController {
	return &aweController{}
}

// Add wires up zero or more routes in the provided mux.Router.
func (c *aweController) Add(_ *di.Container, muxRouter *mux.Router, router *routing.Router) {
	fmt.Printf("awe controller.Add: %s %s %s\n", version, kind, action)

	muxRouter.HandleFunc(
		aweEndpoint,
		func(w http.ResponseWriter, r *http.Request) {
			handle.UseCaseRequest(w, r, version, kind, action, router)
		}).Methods(method)
}

// Supported returns a slice of Supported (a list of supported behaviors).
func (c *aweController) Supported() []common.Supported {
	behavior := application.NewBehavior(version, kind, action)

	fmt.Printf("awe controller.Supported: new behavior: %s %s %s\n", version, kind, action)

	uc := useCase.New()
	return []common.Supported{
		common.NewSupported(
			behavior,
			routable.NewDelegate(
				uc,
				delegate.Apply(
					behavior,
					uc.Execute,
					[]delegate.Handler{
						validator.Validate,
					},
				).Execute,
			),
		),
	}
}
