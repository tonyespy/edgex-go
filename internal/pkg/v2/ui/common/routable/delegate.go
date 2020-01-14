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

package routable

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/ui/common"
)

// delegate contains references to Routable.
type delegate struct {
	routable common.Routable
	execute  application.Executable
}

// NewDelegate is a factory function that returns a routable; fulfills the ui.Routable contract.
func NewDelegate(routable common.Routable, execute application.Executable) *delegate {
	return &delegate{
		routable: routable,
		execute:  execute,
	}
}

// Execute delegates to the application.Execute instance provided to the delegate's factory function.
func (d *delegate) Execute(request interface{}) (response interface{}, status application.Status) {
	return d.execute(request)
}

// EmptyRequest delegates to the EmptyRequest method of the common.Routable provided to the delegate's factory function.
func (d *delegate) EmptyRequest() interface{} {
	return d.routable.EmptyRequest()
}
