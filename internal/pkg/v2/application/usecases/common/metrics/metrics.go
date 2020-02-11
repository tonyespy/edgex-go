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

// metrics implements behavior for the metrics use case.
package metrics

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"
	dtoBase "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/base"
	dto "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/metrics"
	domain "github.com/edgexfoundry/edgex-go/internal/pkg/v2/domain/metrics"
)

// UseCase contains references to dependencies required by the corresponding Routable contract implementation.
type UseCase struct {
	service *domain.Service
}

// New is a factory method that returns an initialized UseCase receiver struct.
func New(service *domain.Service) *UseCase {
	return &UseCase{
		service: service,
	}
}

// Execute encapsulates the behavior for this use case.
func (uc *UseCase) Execute(r interface{}) (interface{}, application.Status) {
	request, _ := r.(*dto.Request)
	alloc, totalAlloc, sys, mallocs, frees, liveObjects, cpuBusyAvg := uc.service.Get()
	return dto.NewResponse(
			*dtoBase.NewResponseForSuccess(request.RequestID),
			alloc,
			totalAlloc,
			sys,
			mallocs,
			frees,
			liveObjects,
			cpuBusyAvg,
		),
		application.StatusSuccess
}

// EmptyRequest returns an empty request associated with this use case.
func (_ *UseCase) EmptyRequest() interface{} {
	return dto.NewEmptyRequest()
}