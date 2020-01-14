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
	"testing"

	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/infrastructure/test"
	commonTest "github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/common/test"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/http/api/system/agent"
)

func TestUseCaseTest(t *testing.T) {
	cancel, wg, router := agent.NewSUT(t, []string{test.EnvNoSecurity}, []string{})
	commonTest.UseCaseTest(t, router)
	cancel()
	wg.Wait()
}

func TestBatchTest(t *testing.T) {
	cancel, wg, router := agent.NewSUT(t, []string{test.EnvNoSecurity}, []string{})
	commonTest.BatchTest(t, router)
	cancel()
	wg.Wait()
}
