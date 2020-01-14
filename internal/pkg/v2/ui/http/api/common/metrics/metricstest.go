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

package metrics

import (
	"net/http"
	"testing"

	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"
	dtoBaseV2dot0 "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/base"
	dtoErrorV2dot0 "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/error"
	dtoV2dot0 "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/metrics"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/infrastructure/test"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common/batchdto"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/batch"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/metrics/v2dot0"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// UseCaseTest verifies metrics endpoint returns expected result; common implementation intended to be executed by
// each service that includes metrics support.
func UseCaseTest(t *testing.T, router *mux.Router) {
	type versionVariation struct {
		name                string
		request             []byte
		assertValidResponse func(t *testing.T, actual []byte, requestIDs []string)
		expectedStatus      int
	}

	requestIDOne := test.FactoryRandomString()
	requestIDTwo := test.FactoryRandomString()
	invalidJSON := test.InvalidJSON()
	versionVariations := map[string][]versionVariation{
		application.Version2dot0: {
			{
				name:                test.Join(test.TypeValid, test.One),
				request:             test.Marshal(t, dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne))),
				assertValidResponse: v2dot0.AssertValidV2dot0UseCaseMetricsResponse,
				expectedStatus:      http.StatusOK,
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
				assertValidResponse: v2dot0.AssertValidV2dot0UseCaseMetricsResponse,
				expectedStatus:      http.StatusMultiStatus,
			},
			{
				name:    test.Join(test.TypeCannotUnmarshal, test.One),
				request: test.Marshal(t, invalidJSON),
				assertValidResponse: func(t *testing.T, actual []byte, _ []string) {
					test.AssertJSONBody(t, dtoErrorV2dot0.NewResponse(dtoBaseV2dot0.NewResponse("", string(test.Marshal(t, invalidJSON)), application.StatusUseCaseUnmarshalFailure)), actual)
				},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name: test.Join(test.TypeValid, test.TypeCannotUnmarshal, test.Two),
				request: test.Marshal(
					t,
					[]interface{}{
						dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne)),
						invalidJSON,
					},
				),
				assertValidResponse: func(t *testing.T, actual []byte, _ []string) {
					v2dot0.AssertTwoWithOneValidAndOneError(t, actual, requestIDOne, application.StatusUseCaseUnmarshalFailure)
				},
				expectedStatus: http.StatusMultiStatus,
			},
			{
				name:    test.Join(test.TypeEmptyRequestId, test.One),
				request: test.Marshal(t, dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(""))),
				assertValidResponse: func(t *testing.T, actual []byte, _ []string) {
					test.AssertJSONBody(t, dtoErrorV2dot0.NewResponse(dtoBaseV2dot0.NewResponse("", dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest("")), application.StatusRequestIdEmptyFailure)), actual)
				},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name: test.Join(test.TypeValid, test.TypeEmptyRequestId, test.Two),
				request: test.Marshal(
					t,
					[]interface{}{
						dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne)),
						dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest("")),
					},
				),
				assertValidResponse: func(t *testing.T, actual []byte, _ []string) {
					v2dot0.AssertTwoWithOneValidAndOneError(t, actual, requestIDOne, application.StatusRequestIdEmptyFailure)
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
							variation.assertValidResponse(t, test.RecastDTOs(t, w.Body.Bytes(), dtoErrorV2dot0.NewEmptyResponse, dtoV2dot0.NewEmptyResponse), []string{requestIDOne, requestIDTwo})
						default:
							assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
						}
					},
				)
			}
		}
	}
}

// BatchTest verifies metrics requests sent to batch endpoint return expected results; common implementation
// intended to be executed by each service that includes metrics support.
func BatchTest(t *testing.T, router *mux.Router) {
	type versionVariation struct {
		name                string
		request             []byte
		assertValidResponse func(t *testing.T, actual []byte, version, kind, action string, requestIDs []string)
		expectedStatus      int
	}

	requestIDOne := test.FactoryRandomString()
	requestIDTwo := test.FactoryRandomString()
	versionVariations := map[string][]versionVariation{
		application.Version2dot0: {
			{
				name:                test.Join(test.One, test.TypeValid),
				request:             test.Marshal(t, []interface{}{batchdto.NewTestRequest(batchdto.NewCommon(version, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne)))}),
				assertValidResponse: v2dot0.AssertValidV2dot0BatchMetricsResponse,
				expectedStatus:      http.StatusMultiStatus,
			},
			{
				name: test.Join(test.Two, test.TypeValid),
				request: test.Marshal(
					t,
					[]interface{}{
						batchdto.NewTestRequest(batchdto.NewCommon(version, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDOne))),
						batchdto.NewTestRequest(batchdto.NewCommon(version, kind, action, batchdto.StrategySynchronous), dtoV2dot0.NewRequest(dtoBaseV2dot0.NewRequest(requestIDTwo))),
					},
				),
				assertValidResponse: v2dot0.AssertValidV2dot0BatchMetricsResponse,
				expectedStatus:      http.StatusMultiStatus,
			},
		},
	}
	for v := range versionVariations {
		for _, variation := range versionVariations[v] {
			t.Run(
				test.Name(v, batch.Method, variation.name),
				func(t *testing.T) {
					w := test.SendRequest(t, router, batch.Method, batch.Endpoint, variation.request)

					assert.Equal(t, variation.expectedStatus, w.Code)
					test.AssertContentTypeIsJSON(t, w.Header())
					variation.assertValidResponse(t, test.RecastDTOs(t, w.Body.Bytes(), dtoErrorV2dot0.NewEmptyResponse, batchdto.NewEmptyResponse), version, kind, action, []string{requestIDOne, requestIDTwo})
				},
			)
		}
	}
}
