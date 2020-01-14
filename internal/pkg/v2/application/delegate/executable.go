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

package delegate

import "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"

// Handler defines the contract for executable handlers.
type Handler func(
	request interface{},
	behavior *application.Behavior,
	execute application.Executable) (response interface{}, status application.Status)

// executable contains references to dependencies required by the corresponding Executable contract.
type executable struct {
	execute application.Executable
}

// Apply is a factory function that returns an Executable; it wraps a provided Executable with behavior.
func Apply(behavior *application.Behavior, execute application.Executable, handlers []Handler) *executable {
	return &executable{
		execute: apply(behavior, execute, handlers),
	}
}

// apply creates a closure chain that wraps the provided application.Executable with the provided handlers' behavior.
func apply(behavior *application.Behavior, execute application.Executable, handlers []Handler) application.Executable {
	for i := len(handlers) - 1; i >= 0; i-- {
		execute = func(handler Handler, execute application.Executable) application.Executable {
			return func(request interface{}) (response interface{}, status application.Status) {
				return handler(request, behavior, execute)
			}
		}(handlers[i], execute)
	}
	return execute
}

// Execute delegates execution to curried executable.
func (r *executable) Execute(content interface{}) (interface{}, application.Status) {
	return r.execute(content)
}
