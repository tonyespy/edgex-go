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

package common

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/batch"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/metrics"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/ping"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/test"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/version"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/routing"
)

// V2Routes provides a single cross-service implementation of common V2 API route definitions.
func V2Routes(inAcceptanceTestMode bool, controllers []routing.Controller) []routing.Controller {
	commonRoutes := []routing.Controller{
		batch.New(),
		metrics.New(),
		ping.New(),
		version.New(),
	}

	if inAcceptanceTestMode {
		commonRoutes = append(
			commonRoutes,
			[]routing.Controller{
				test.New(),
			}...,
		)
	}

	return append(commonRoutes, controllers...)
}
