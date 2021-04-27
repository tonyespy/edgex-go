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

const (
	// EnvConfig is the prefix used for configure hook keys used for
	// EdgeX configuration overrides.
	EnvConfig = "env"
	// ServiceConsul is the service key for Consul.
	ServiceConsul = "consul"
	// ServiceRedis is the service key for Redis.
	ServiceRedis = "redis"
	// ServiceData is the service key for EdgeX Core Data.
	ServiceData = "core-data"
	// ServiceMetadata is the service key for EdgeX Core MetaData.
	ServiceMetadata = "core-metadata"
	// ServiceCommand is the service key for EdgeX Core Command.
	ServiceCommand = "core-command"
	// ServiceNotify is the service key for EdgeX Support Notifications.
	ServiceNotify = "support-notifications"
	// ServiceSched is the service key for EdgeX Support Scheduler.
	ServiceSched = "support-scheduler"
	// ServiceAppCfg is the service key for EdgeX App Service Configurable.
	ServiceAppCfg = "app-service-configurable"
	// ServiceDevVirt is the service key for EdgeX Device Virtual.
	ServiceDevVirt = "device-virtual"
	// ServiceSecStore is the service key for EdgeX Security Secret Store (aka Vault).
	ServiceSecStore = "security-secret-store"
	// ServiceProxy is the service key for EdgeX API Gateway (aka Kong).
	ServiceProxy = "security-proxy"
	// ServiceSysMgmt is the service key for EdgeX SMA (sys-mgmt-agent).
	ServiceSysMgmt = "sys-mgmt-agent"
	// ServiceKuiper is the service key for the Kuiper rules engine.
	ServiceKuiper   = "kuiper"
	snapEnv         = "SNAP"
	snapCommonEnv   = "SNAP_COMMON"
	snapDataEnv     = "SNAP_DATA"
	snapInstNameEnv = "SNAP_INSTANCE_NAME"
	snapNameEnv     = "SNAP_NAME"
	snapRevEnv      = "SNAP_REVISION"
)

// ConfToEnv defines mappings from snap config keys to EdgeX environment variable
// names that are used to override individual service configuration values via a
// .env file read by the snap service wrapper.
//
// The syntax to set a configuration key is:
//
// env.<service name>.<section>.<keyname>
//
var ConfToEnv = map[string]string{
	// [Writable] - not yet supported
	// conf_to_env["writable.log-level"]="BootTimeout"
	// [Service]
	"service.boot-timeout":     "SERVICE_BOOTTIMEOUT",
	"service.check-interval":   "SERVICE_CHECKINTERVAL",
	"service.host":             "SERVICE_HOST",
	"service.server-bind-addr": "SERVICE_SERVERBINDADDR",
	"service.port":             "SERVICE_PORT",
	"service.protocol":         "SERVICE_PROTOCOL",
	"service.max-result-count": "SERVICE_MAXRESULTCOUNT",
	"service.read-max-limit":   "SERVICE_READMAXLIMIT",
	"service.startup-msg":      "SERVICE_STARTUPMSG",
	"service.timeout":          "SERVICE_TIMEOUT",

	// [Registry] -- not yet supported, would also require consul changes

	// [Clients.Command]
	"clients.command.port": "CLIENTS_COMMAND_PORT",

	// [Clients.CoreData]
	"clients.coredata.port": "CLIENTS_COREDATA_PORT",

	// [Clients.Data]
	// There are two client keys for CoreData because device-sdk-go uses
	// this key, and all the core services uses the previous key.
	"clients.data.port": "CLIENTS_DATA_PORT",

	// [Clients.Metadata]
	"clients.metadata.port": "CLIENTS_METADATA_PORT",

	// [Clients.Notifications]
	"clients.notifications.port": "CLIENTS_NOTIFICATIONS_PORT",

	// [Clients.Scheduler]
	"clients.scheduler.port": "CLIENTS_SCHEDULER_PORT",

	// [Databases] -- not supported

	// [MessageQueue] -- core-data only
	"messagequeue.topic": "MESSAGEQUEUE_TOPIC",

	// [MessageQueue.Optional] - not yet supported

	// [SecretStore]
	"secretstore.additional-retry-attempts": "SECRETSTORE_ADDITIONALRETRYATTEMPTS",
	"secretstore.retry-wait-period":         "SECRETSTORE_RETRYWAITPERIOD",

	// [SecretStore.Authentication] -- not supported

	// [SecretStoreExclusive] -- app-service-config only; not supported

	// Note the key used vs. the key used for the edgexfoundry snap (e.g. Data).
	// [Binding] -- app-service-config only; not supported

	// [MessageBus] -- not supported

	// support-notifications options
	// [Smtp]
	"smtp.host":                    "SMTP_HOST",
	"smtp.username":                "SMTP_USERNAME",
	"smtp.password":                "SMTP_PASSWORD",
	"smtp.port":                    "SMTP_PORT",
	"smtp.sender":                  "SMTP_SENDER",
	"smtp.enable-self-signed-cert": "SMTP_ENABLE_SELF_SIGNED_CERT",

	// security-proxy
	//
	// ADD_PROXY_ROUTE is a csv list of URLs to be added to the
	// API Gateway (aka Kong). For references:
	//
	// https://docs.edgexfoundry.org/1.3/microservices/security/Ch-APIGateway/
	//
	// NOTE - this setting is not a configuration override, it's a top-level
	// environment variable used by the security-proxy-setup.
	//
	// TODO: validation
	//
	"add-proxy-route": "ADD_PROXY_ROUTE",

	// [AuthType]
	"authtype.name": "AUTHTYPE_NAME",

	// security-secret-store
	//
	// ADD_SECRETSTORE_TOKENS is a csv list of service keys to be added to the
	// list of Vault tokens that security-file-token-provider (launched by
	// security-secretstore-setup) creates.
	//
	// NOTE - this setting is not a configuration override, it's a top-level
	// environment variable used by the security-secretstore-setup.
	//
	// TODO: validation
	//
	"add-secretstore-tokens": "ADD_SECRETSTORE_TOKENS",
}

// Services is a string array of all of the edgexfoundry snap services.
var Services = []string{
	// base services
	ServiceConsul,
	ServiceRedis,
	// core services
	ServiceData,
	ServiceMetadata,
	ServiceCommand,
	// support services
	ServiceNotify,
	ServiceSched,
	// app-services
	ServiceAppCfg,
	// device services
	ServiceDevVirt,
	// security services
	ServiceSecStore,
	ServiceProxy,
	// sys mgmt services
	ServiceSysMgmt,
	// rules-engine
	ServiceKuiper,
}
