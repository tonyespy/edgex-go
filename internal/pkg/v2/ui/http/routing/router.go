/*******************************************************************************
 * Copyright 2020 Dell Inc.
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

package routing

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/delegate"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common/routable"

	"github.com/edgexfoundry/go-mod-bootstrap/di"

	"github.com/gorilla/mux"
)

// Controller defines the contract fulfilled by ui-level controllers.
type Controller interface {
	// Add wires up zero or more routes in the provided mux.Router.
	Add(dic *di.Container, muxRouter *mux.Router, router *Router)

	// supported returns a slice of supported (a list of supported behaviors).
	Supported() []common.Supported
}

// supported is the common type definition for the map of supported behavior used by Router.
type supported map[string]common.Routable

// Router contains references to dependencies required by the commandQuery routing implementation.
type Router struct {
	supported supported
}

// Initialize takes a list of Controller implementations, adds zero or more corresponding routes to the supplied mux.Router,
// and updates the receiver's supported field with the code to execute.
func Initialize(dic *di.Container, muxRouter *mux.Router, handlers []delegate.Handler, controllers []Controller) {
	r := Router{
		supported: make(supported),
	}

	for i := range controllers {
		controllers[i].Add(dic, muxRouter, &r)
		for _, s := range controllers[i].Supported() {
			r.supported[r.envelopeToKey(s.Version, s.Kind, s.Action)] =
				routable.NewDelegate(
					s.Routable,
					delegate.Apply(
						application.NewBehavior(s.Version, s.Kind, s.Action),
						s.Routable.Execute,
						handlers,
					).Execute,
				)
		}
	}
}

// FindRoute returns whether or not a ui.Routable exists for a specific version, kind, and action (and the Routable if
// it does).
func (r *Router) FindRoute(version, kind, action string) (common.Routable, bool) {
	routableBehavior, exists := r.supported[r.envelopeToKey(version, kind, action)]
	return routableBehavior, exists
}

// envelopeToKey converts an action, kind, and version to the receiver's supported map's key.
func (r *Router) envelopeToKey(version, kind, action string) string {
	return version + "_" + kind + "_" + action
}
