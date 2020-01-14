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

package ping

import (
	"net/http"
	"testing"

	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"
	dtoBaseV2dot0 "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/base"
	dtoErrorV2dot0 "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/error"
	dtoV2dot0 "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/ping"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/infrastructure/test"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common/batchdto"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/batch"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// UseCaseTest verifies ping endpoint returns expected result; common implementation intended to be executed by
// each service that includes ping support.
func UseCaseTest(t *testing.T, router *mux.Router) {
	type versionVariation struct {
		name             string
		request          []byte
		expectedResponse interface{}
		expectedStatus   int
	}

	requestIDOne := test.FactoryRandomString()
	requestIDTwo := test.FactoryRandomString()
	invalidJSON := test.InvalidJSON()
	versionVariations := map[string][]versionVariation{
		application.Version2dot0: {
			{
				name:             test.Join(test.TypeValid, test.One),
				request:          test.Marshal(t, dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne))),
				expectedResponse: dtoV2dot0.NewResponse(dtoBaseV2dot0.NewResponseForSuccess(requestIDOne)),
				expectedStatus:   http.StatusOK,
			},
			{
				name: test.Join(test.TypeValid, test.Two),
				request: test.Marshal(
					t,
					[]interface{}{
						dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne)),
						dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDTwo)),
					},
				),
				expectedResponse: []interface{}{
					dtoV2dot0.NewResponse(dtoBaseV2dot0.NewResponseForSuccess(requestIDOne)),
					dtoV2dot0.NewResponse(dtoBaseV2dot0.NewResponseForSuccess(requestIDTwo)),
				},
				expectedStatus: http.StatusMultiStatus,
			},
			{
				name:             test.Join(test.TypeCannotUnmarshal, test.One),
				request:          test.Marshal(t, invalidJSON),
				expectedResponse: dtoErrorV2dot0.NewResponse(dtoBaseV2dot0.NewResponse("", string(test.Marshal(t, invalidJSON)), application.StatusUseCaseUnmarshalFailure)),
				expectedStatus:   http.StatusBadRequest,
			},
			{
				name: test.Join(test.TypeCannotUnmarshal, test.Two),
				request: test.Marshal(
					t,
					[]interface{}{
						dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne)),
						invalidJSON,
					},
				),
				expectedResponse: []interface{}{
					dtoV2dot0.NewResponse(dtoBaseV2dot0.NewResponseForSuccess(requestIDOne)),
					dtoErrorV2dot0.NewResponse(dtoBaseV2dot0.NewResponse("", string(test.Marshal(t, invalidJSON)), application.StatusUseCaseUnmarshalFailure)),
				},
				expectedStatus: http.StatusMultiStatus,
			},
			{
				name:             test.Join(test.TypeEmptyRequestId, test.One),
				request:          test.Marshal(t, dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(""))),
				expectedResponse: dtoErrorV2dot0.NewResponse(dtoBaseV2dot0.NewResponse("", dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest("")), application.StatusRequestIdEmptyFailure)),
				expectedStatus:   http.StatusBadRequest,
			},
			{
				name: test.Join(test.TypeEmptyRequestId, test.Two),
				request: test.Marshal(
					t,
					[]interface{}{
						dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne)),
						dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest("")),
					},
				),
				expectedResponse: []interface{}{
					dtoV2dot0.NewResponse(dtoBaseV2dot0.NewResponseForSuccess(requestIDOne)),
					dtoErrorV2dot0.NewResponse(dtoBaseV2dot0.NewResponse("", dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest("")), application.StatusRequestIdEmptyFailure)),
				},
				expectedStatus: http.StatusMultiStatus,
			},
		},
	}
	for version := range versionVariations {
		for _, variation := range versionVariations[version] {
			for _, m := range test.ValidMethods() {
				t.Run(
					test.Name(m, version, variation.name),
					func(t *testing.T) {
						w := test.SendRequest(t, router, m, endpoint, variation.request)

						switch m {
						case method:
							assert.Equal(t, variation.expectedStatus, w.Code)
							test.AssertContentTypeIsJSON(t, w.Header())
							test.AssertJSONBody(t, variation.expectedResponse, test.RecastDTOs(t, w.Body.Bytes(), dtoErrorV2dot0.NewEmptyResponse, dtoV2dot0.NewEmptyResponse))
						default:
							assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
						}
					},
				)
			}
		}
	}
}

// BatchTest verifies ping requests sent to batch endpoint return expected results; common implementation
// intended to be executed by each service that includes ping support.
func BatchTest(t *testing.T, router *mux.Router) {
	type versionVariation struct {
		name             string
		request          []byte
		expectedResponse interface{}
		expectedStatus   int
	}

	requestID := test.FactoryRandomString()
	versionVariations := map[string][]versionVariation{
		application.Version2dot0: {
			{
				name:             test.Join(test.One, test.TypeValid),
				request:          test.Marshal(t, []interface{}{batchdto.NewTestRequest(batchdto.NewCommon(application.Version2, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestID)))}),
				expectedResponse: []interface{}{batchdto.NewTestRequest(batchdto.NewCommon(application.Version2, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewResponse(dtoBaseV2dot0.NewResponseForSuccess(requestID)))},
				expectedStatus:   http.StatusMultiStatus,
			},
			{
				name: test.Join(test.Two, test.TypeValid),
				request: test.Marshal(
					t,
					[]interface{}{
						batchdto.NewTestRequest(batchdto.NewCommon(application.Version2, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestID))),
						batchdto.NewTestRequest(batchdto.NewCommon(application.Version2, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestID))),
					},
				),
				expectedResponse: []interface{}{
					batchdto.NewTestRequest(batchdto.NewCommon(application.Version2, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewResponse(dtoBaseV2dot0.NewResponseForSuccess(requestID))),
					batchdto.NewTestRequest(batchdto.NewCommon(application.Version2, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewResponse(dtoBaseV2dot0.NewResponseForSuccess(requestID))),
				},
				expectedStatus: http.StatusMultiStatus,
			},
		},
	}
	for version := range versionVariations {
		for _, variation := range versionVariations[version] {
			t.Run(
				test.Name(version, batch.Method, variation.name),
				func(t *testing.T) {
					w := test.SendRequest(t, router, batch.Method, batch.Endpoint, variation.request)

					assert.Equal(t, variation.expectedStatus, w.Code)
					test.AssertContentTypeIsJSON(t, w.Header())
					test.AssertJSONBody(t, variation.expectedResponse, test.RecastDTOs(t, w.Body.Bytes(), dtoErrorV2dot0.NewEmptyResponse, batchdto.NewEmptyResponse))
				},
			)
		}
	}
}
