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

package v2dot0

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"
	dtoErrorV2dot0 "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/error"
	dtoV2dot0 "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/metrics"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common/batchdto"

	"github.com/stretchr/testify/assert"
)

// assertValid validates metrics response; we can't predict what will be returned so we consider valid
// response to have all non-zero fields.
func assertValid(t *testing.T, response *dtoV2dot0.Response, requestIDs []string) {
	inList := func() bool {
		for index := range requestIDs {
			if response.RequestID == requestIDs[index] {
				return true
			}
		}
		return false
	}

	if !inList() {
		assert.Fail(t, fmt.Sprintf("requestID %s not in list %v", response.RequestID, requestIDs))
	}
	assert.NotEqual(t, 0, response.Alloc)
	assert.NotEqual(t, 0, response.TotalAlloc)
	assert.NotEqual(t, 0, response.Sys)
	assert.NotEqual(t, 0, response.Mallocs)
	assert.NotEqual(t, 0, response.Frees)
	assert.NotEqual(t, 0, response.LiveObjects)
	assert.NotEqual(t, 0, response.CpuBusyAvg)
}

// AssertValidV2dot0UseCaseMetricsResponse validates metrics response; we can't predict what will be returned so we consider valid
// response to have all non-zero fields.
func AssertValidV2dot0UseCaseMetricsResponse(t *testing.T, actual []byte, requestIDs []string) {
	var responseDTO dtoV2dot0.Response

	// single response?
	err := json.Unmarshal(actual, &responseDTO)
	if err == nil {
		assertValid(t, &responseDTO, requestIDs)
		return
	}

	// multiple  responses?
	var responseDTOs []dtoV2dot0.Response
	err = json.Unmarshal(actual, &responseDTOs)
	if err == nil {
		for i := range responseDTOs {
			assertValid(t, &responseDTOs[i], requestIDs)
		}
		return
	}

	assert.Fail(t, "unable to validate metrics response: %s", err.Error())
}

// AssertTwoWithOneValidAndOneError validates one successful result and one error result of a specific type.
func AssertTwoWithOneValidAndOneError(t *testing.T, actual []byte, validRequestID string, status application.Status) {
	var responseDTOs []*json.RawMessage
	if err := json.Unmarshal(actual, &responseDTOs); err != nil {
		assert.Fail(t, "unable to unmarshal: %s", err.Error())
		return
	}

	assert.Equal(t, 2, len(responseDTOs))

	var validExists = false
	for i := range responseDTOs {
		var responseDTO dtoV2dot0.Response
		if err := json.Unmarshal(*responseDTOs[i], &responseDTO); err != nil {
			continue
		}
		if responseDTO.StatusCode != application.StatusSuccess {
			continue
		}

		validExists = true
		assertValid(t, &responseDTO, []string{validRequestID})
	}
	assert.True(t, validExists)

	var invalidExists = false
	for i := range responseDTOs {
		var responseDTO dtoErrorV2dot0.Response
		if err := json.Unmarshal(*responseDTOs[i], &responseDTO); err != nil {
			continue
		}
		if responseDTO.StatusCode == application.StatusSuccess {
			continue
		}

		invalidExists = true
		assert.Equal(t, status, responseDTO.StatusCode)
	}
	assert.True(t, invalidExists)

}

// AssertValidV2dot0BatchMetricsResponse validates metrics response; we can't predict what will be returned so we consider valid
// response to have all non-zero fields.
func AssertValidV2dot0BatchMetricsResponse(
	t *testing.T,
	actual []byte,
	version string,
	kind string,
	action string,
	_ []string) {

	responseDTO := batchdto.TestResponse{}
	assertValidBatch := func(response *batchdto.TestResponse) error {
		assert.Equal(t, version, response.Version)
		assert.Equal(t, kind, response.Kind)
		assert.Equal(t, action, response.Action)
		err := json.Unmarshal(*response.Content, &responseDTO)
		if err == nil {
			//assertValid(t, &responseDTO.Content, requestIDs)
		}
		return err
	}

	// single response?
	err := json.Unmarshal(actual, &responseDTO)
	if err == nil {
		//assertValid(&responseDTO)
		return
	}

	// multiple batch responses?
	responseDTOs := batchdto.EmptyTestResponseSlice()
	err = json.Unmarshal(actual, &responseDTOs)
	if err == nil {
		for i := range responseDTOs {
			if err := assertValidBatch(&responseDTOs[i]); err != nil {
				assert.Fail(t, "unable to unmarshal metrics responseDTO: %s", err.Error())
			}
		}
		return
	}

	assert.Fail(t, "unable to validate metrics response: %s", err.Error())
}
