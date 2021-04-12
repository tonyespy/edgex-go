// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	debug bool = false
	log   *syslog.Writer
	// Snap contains the value of the SNAP environment variable.
	Snap string
	// SnapConf contains the expanded path '$SNAP/config'.
	SnapConf string
	// SnapCommon contains the value of the SNAP_COMMON environment variable.
	SnapCommon string
	// SnapData contains the value of the SNAP_DATA environment variable.
	SnapData string
	// SnapDataConf contains the expanded path '$SNAP_DATA/config'.
	SnapDataConf string
	// SnapInst contains the value of the SNAP_INSTANCE_NAME environment variable.
	SnapInst string
	// SnapName contains the value of the SNAP_NAME environment variable.
	SnapName string
	// SnapRev contains the value of the SNAP_REVISION environment variable.
	SnapRev string
)

// CtlCli is the test obj for overridding functions
type CtlCli struct{}

// SnapCtl interface provides abstration for unit testing
type SnapCtl interface {
	Config(key string) (string, error)
	SetConfig(key string, val string) error
	Stop(svc string, disable bool) error
}

// CopyFile copies a file within the snap
func CopyFile(srcPath, destPath string) error {

	inFile, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// TODO: check file perm
	err = ioutil.WriteFile(destPath, inFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CopyFileReplace copies a file within the snap and replaces strings using
// the string/replace values in the rStrings parameter.
func CopyFileReplace(srcPath, destPath string, rStrings map[string]string) error {

	inFile, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	rStr := string(inFile)
	for k, v := range rStrings {
		rStr = strings.Replace(rStr, k, v, 1)
	}

	// TODO: check file perm
	outBytes := []byte(rStr)
	err = ioutil.WriteFile(destPath, outBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Debug writes the given msg to sylog (sev=LOG_DEBUG) if the associated
// global snap 'debug' configuration flag is set to 'true'.
func Debug(msg string) {
	if debug {
		log.Debug(msg)
	}
}

// Error writes the given msg to sylog (sev=LOG_ERROR).
func Error(msg string) {
	log.Err(msg)
}

// Info writes the given msg to sylog (sev=LOG_INFO).
func Info(msg string) {
	log.Info(msg)
}

// Warn writes the given msg to sylog (sev=LOG_WARNING).
func Warn(msg string) {
	log.Err(msg)
}

// GetEnvVars populates global variables for each of the SNAP*
// variables defined in the snap's environment
func GetEnvVars() error {
	Snap = os.Getenv(snapEnv)
	if Snap == "" {
		return errors.New("SNAP is not set")
	}

	SnapCommon = os.Getenv(snapCommonEnv)
	if SnapCommon == "" {
		return errors.New("SNAP_COMMON is not set")
	}

	SnapData = os.Getenv(snapDataEnv)
	if SnapData == "" {
		return errors.New("SNAP_DATA is not set")
	}

	SnapInst = os.Getenv(snapInstNameEnv)
	if SnapInst == "" {
		return errors.New("SNAP_INSTANCE_NAME is not set")
	}

	SnapName = os.Getenv(snapNameEnv)
	if SnapName == "" {
		return errors.New("SNAP_NAME is not set")
	}

	SnapRev = os.Getenv(snapRevEnv)
	if SnapRev == "" {
		return errors.New("SNAP_REVISION_NAME is not set")
	}

	SnapConf = Snap + "/config"
	SnapDataConf = SnapData + "/config"

	return nil
}

// InitLog create a new syslog instance for the hook and sets the
// global debug flag based on the value of the setDebug parameter.
func InitLog(setDebug bool) error {
	var err error

	debug = setDebug

	log, err = syslog.New(syslog.LOG_INFO, "edgexfoundry:configure")
	if err != nil {
		return err
	}

	return nil
}

// NewSnapCtl returns a normal runtime client
func NewSnapCtl() *CtlCli {
	return &CtlCli{}
}

// Config uses snapctl to get a value from a key, or returns error.
func (cc *CtlCli) Config(key string) (string, error) {
	out, err := exec.Command("snapctl", "get", key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// SetConfig uses snapctl to set a config value from a key, or returns error.
func (cc *CtlCli) SetConfig(key string, val string) error {

	err := exec.Command("snapctl", "set", fmt.Sprintf("%s=%s", key, val)).Run()
	if err != nil {
		return fmt.Errorf("snapctl SET failed for %s - %v", key, err)
	}
	return nil
}

// Start uses snapctrl to start a service and optionally enable it
func (cc *CtlCli) Start(svc string, enable bool) error {
	var cmd *exec.Cmd

	name := SnapName + "." + svc
	if enable {
		cmd = exec.Command("snapctl", "start", "--enable", name)
	} else {
		cmd = exec.Command("snapctl", "start", name)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("snapctl start %s failed - %v", name, err)
	}

	return nil
}

// Stop uses snapctrl to stop a service and optionally disable it
func (cc *CtlCli) Stop(svc string, disable bool) error {
	var cmd *exec.Cmd

	name := SnapName + "." + svc
	if disable {
		cmd = exec.Command("snapctl", "stop", "--disable", name)
	} else {
		cmd = exec.Command("snapctl", "stop", name)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("snapctl stop %s failed - %v", name, err)
	}

	return nil
}

// p is the current prefix of the config key being processed (e.g. "service", "security.auth")
// k is the key name of the current JSON object being processed
// vJSON is the current object
// flatConf is a map containing the configuration keys/values processed thus far
func flattenConfigJSON(p string, k string, vJSON interface{}, flatConf map[string]string) {
	var mk string

	// top level keys don't include "env", so no separator needed
	if p == "" {
		mk = k
	} else {
		mk = fmt.Sprintf("%s.%s", p, k)
	}

	switch t := vJSON.(type) {
	case string:
		flatConf[mk] = t
	case bool:
		flatConf[mk] = strconv.FormatBool(t)
	case float64:
		flatConf[mk] = strconv.FormatFloat(t, 'f', -1, 64)
	case map[string]interface{}:

		for k, v := range t {
			flattenConfigJSON(mk, k, v, flatConf)
		}
	default:
		panic(fmt.Sprintf("internal error: invalid JSON configuration from snapd - prefix: %s key: %s obj: %v", p, k, t))
	}
}

// HandleEdgeXConfig processes snap configuration which can be used to override
// edgexfoundry configuration via environment variables sourced by the snap
// service wrapper script.
func HandleEdgeXConfig(service, envJSON string) error {

	if envJSON == "" {
		return nil
	}

	var m map[string]interface{}
	var flatConf = make(map[string]string)

	err := json.Unmarshal([]byte(envJSON), &m)
	if err != nil {
		return fmt.Errorf("failed to unmarshall EdgeX config - %v", err)
	}

	for k, v := range m {
		flattenConfigJSON("", k, v, flatConf)
	}

	b := bytes.Buffer{}
	for k, v := range flatConf {
		env, ok := ConfToEnv[k]
		if !ok {
			return errors.New("invalid EdgeX config option - " + k)
		}

		// TODO: add logic to check that v is allowable for service
		// e.g. service.read-max-limit is valid for app-service-cfg only

		_, err := fmt.Fprintf(&b, "export %s=%s\n", env, v)
		if err != nil {
			return err
		}
	}

	// Handle security-* service naming. The service names in this
	// hook historically do not align with the actual binary commands.
	// As such, when handling configuration settings for them, we need
	// to translate the hook name to the actual binary name.
	if service == "security-proxy" {
		service = "security-proxy-setup"
	} else if service == "security-secret-store" {
		service = "security-secretstore-setup"
	}

	path := fmt.Sprintf("%s/%s/res/%s.env", SnapDataConf, service, service)
	tmp := path + ".tmp"

	err = ioutil.WriteFile(tmp, b.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %.env file - %v", service, err)
	}

	err = os.Rename(tmp, path)
	if err != nil {
		return fmt.Errorf("failed to rename %s.env.tmp file - %v", service, err)
	}

	return nil
}
