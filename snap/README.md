# EdgeX Foundry Core Snap
[![snap store badge](https://raw.githubusercontent.com/snapcore/snap-store-badges/master/EN/%5BEN%5D-snap-store-black-uneditable.png)](https://snapcraft.io/edgexfoundry)

This folder contains snap packaging for the EdgeX Foundry reference implementation.

The snap contains all of the EdgeX Go-based micro services from this repository, device-virtual, app-service-configurable (for Kuiper
integration), as well as all the necessary runtime components (e.g. Consul, Kong, Redis, ...) required to run an EdgeX instance.

The project maintains a rolling release of the snap on the `edge` channel that is rebuilt and published at least once daily.

The snap currently supports running on both `amd64` and `arm64` platforms.

## Installation

### Installing snapd
The snap can be installed on any system that supports snaps. You can see how to install 
snaps on your system [here](https://snapcraft.io/docs/installing-snapd).

However for full security confinement, the snap should be installed on an 
Ubuntu 16.04 LTS or later Desktop or Server, or a system running Ubuntu Core 16 or later.

### Installing EdgeX Foundry as a snap
The snap is published in the snap store at https://snapcraft.io/edgexfoundry.
You can see the current revisions available for your machine's architecture by running the command:

```bash
$ snap info edgexfoundry
```

The snap can be installed using `snap install`. To install the latest stable version:

```bash
$ sudo snap install edgexfoundry
```

To install the snap from the edge channel:

```bash
$ sudo snap install edgexfoundry --edge
```

**Note** - in general, installing from the edge channel is only recommended for development purposes. Depending on the state of the current development release, your mileage may vary.


You can also specify specific releases using the `--channel` option. For example to install the Hanoi release of the snap:

```bash
$ sudo snap install edgexfoundry --channel=hanoi

```

Lastly, on a system supporting it, the snap may be installed using GNOME (or Ubuntu) Software Center by searching for `edgexfoundry`.

**Note** - the snap has only been tested on Ubuntu Desktop/Server LTS releases (16.04 or later), as well as Ubuntu Core (16 or later).

**WARNING** - don't install the EdgeX snap on a system which is already running one of the included services (e.g. Consul, Redis, Vault, ...), as this
may result in resource conflicts (i.e. ports) which could cause the snap install to fail.

## Using the EdgeX snap

Upon installation, the following EdgeX services are automatically and immediately started:

* consul
* redis
* core-data
* core-command
* core-metadata
* security-services (see [note below](https://github.com/edgexfoundry/edgex-go/tree/master/snap#security-services))

The following services are disabled by default:

* app-service-configurable (required for Kuiper)
* device-virtual
* kuiper
* support-notifications
* support-scheduler
* sys-mgmt-agent

Any disabled services can be enabled and started up using `snap set`:

```bash
$ sudo snap set edgexfoundry support-notifications=on
```

To turn a service off (thereby disabling and immediately stopping it) set the service to off:

```bash
$ sudo snap set edgexfoundry support-notifications=off
```

All services which are installed on the system as systemd units, which if enabled will automatically start running when the system boots or reboots.

### Configuring individual services

All default configuration files are shipped with the snap inside `$SNAP/config`, however because `$SNAP` isn't writable, all of the config files are copied during snap installation (specifically during the install hook, see `snap/hooks/install` in this repository) to `$SNAP_DATA/config`.

**Note** - `$SNAP` resolves to the path `/snap/edgexfoundry/current/` and `$SNAP_DATA` resolves to `/var/snap/edgexfoundry/current`.

In the Geneva release of EdgeX, services were changed such that each became responsible for "self-seeding" its own configuration to Consul.
Currently the only way to effect configuration changes for services that are auto-started (e.g. Core Data, Core Metadata) is to change
configuration directly via Consul's UI or [kv REST API](https://www.consul.io/api/kv.html). Changes made to configuration in Consul
require services to be restarted in order for the changes to take effect; the one exception are changes made to configuration items in
a service's ```[Writable]``` section. Services that aren't started by default (see above) *will* pickup any changes made to their config
files when started.

Also it should be noted that use of Consul is enabled by default in the snap. It is not possible at this time to run the EdgeX services in
the snap with Consul disabled.

### Viewing logs
To view the logs for all services in the edgexfoundry snap use:

```bash
$ sudo snap logs edgexfoundry
```

Individual service logs may be viewed by specifying the service name:

```bash
$ sudo snap logs edgexfoundry.consul
```

Or by using the systemd unit name and `journalctl`:

```bash
$ journalctl -u snap.edgexfoundry.consul
```

### Configuration Overrides
The EdgeX snap supports configuration overrides via its configure and install hooks which generate service-specific .env files
which are used to provide a custom environment to the service, overriding the default configuration provided by the service's
```configuration.toml``` file. If a configuration override is made after a service has already started, then the service must
be **restarted** via command-line (e.g. ```snap restart edgexfoundry.<service>```), snapd's REST API, or the SMA (sys-mgmt-agent).
If the overrides are provided via the snap configuration defaults capability of a gadget snap, the overrides will be picked
up when the services are first started.

The following syntax is used to specify service-specific configuration overrides:

```env.<service>.<stanza>.<config option>```

For instance, to setup an override of Core Data's Port use:

```$ sudo snap set edgexfoundry env.core-data.service.port=2112```

And restart the service:

```$ sudo snap restart edgexfoundry.core-data```

**Note** - at this time changes to configuration values in the [Writable] section are not supported.

For details on the mapping of configuration options to Config options, please refer to "Service Environment Configuration Overrides".

### Security services

Currently, the security services are enabled by default. The security services consitute the following components:

 * Kong
 * PostgreSQL
 * Vault
 * security-secrets-setup
 * security-secretstore-setup
 * security-proxy-setup

#### Secret Store
Vault is used by EdgeX for secret management (e.g. certificates, keys, passwords, ...) and is referred to as the Secret Store.

Use of Secret Store by all services can be disabled globally, but doing so will also disable the API Gateway, as it depends on the Secret Store.
Thus the following command will disable both:

```bash
$ sudo snap set edgexfoundry security-secret-store=off
```

### API Gateway 
Kong is used for access control to the EdgeX services from external systems and is referred to as the API Gateway. 

For more details please refer to the EdgeX API Gateway [documentation](https://docs.edgexfoundry.org/1.3/microservices/security/Ch-APIGateway/).


The API Gateway can be disabled by using the following command:

```bash
$ sudo snap set edgexfoundry security-proxy=off
```

**Note** - by default all services in the snap except for the API Gateway are restricted to listening on 'localhost' (i.e. the services are
not addressable from another system). In order to make a service accessible remotely, the appropriate configuration override of the
'Service.ServerBindAddr' needs to be made (e.g. ```sudo snap set edgexfoundry env.core-data.service.server-bind-addr=0.0.0.0```).


#### API Gateway User Setup

##### JWT Tokens

Before the API Gateway can be used, a user and group must be created, and a JWT access token generated. This can be accomplised via the
`secrets-config` command, or by using `snap set` commands.

The first step is to add a user. You need to create a public/private keypair, which can be done with

```bash
# Create private key:
$ openssl ecparam -genkey -name prime256v1 -noout -out private.pem

# Create public key:
$ openssl ec -in private.pem -pubout -out public.pem
```

If you then create the user using the secrets-config command, then you need to provide
- The username
- The public key
- (optionally) ID. This is a unique string identifying the credential. It will be required in the next step to create the JWT token. If you don't specify it, 
then an autogenerated one will be output by the secrets-config command

```bash
$ edgexfoundry.secrets-config proxy adduser --token-type jwt --user user01 --algorithm ES256 --public_key public.pem --id USER_ID
```

Alternatively, to do this using `snap set`:

```bash
# set user=username,user id,algorithm (ES256 or RS256)
sudo snap set edgexfoundry env.security-proxy.user=user01,USER_ID,ES256

# set public-key to the contents of a PEM-encoded public key file
sudo snap set edgexfoundry env.security-proxy.public-key="$(cat public.pem)"
```

The second step is then to generate a token using the user ID which you specified:

```bash
$ edgexfoundry.secrets-config proxy jwt --algorithm ES256 --private_key private.pem --id USER_ID --expiration=1h
TOKEN="copy manually from output of secrets-config command"
```

Alternatively , you can generate the token on a different device using a bash script:

```bash
header='{
    "alg": "ES256",
    "typ": "JWT"
}'

TTL=$((EPOCHSECONDS+3600)) 

payload='{
    "iss":"USER_ID",
    "iat":'$EPOCHSECONDS', 
    "nbf":'$EPOCHSECONDS',
    "exp":'$TTL' 
}'

JWT_HEADER=`echo -n $header | openssl base64 -e -A | sed s/\+/-/ | sed -E s/=+$//`
JWT_PAYLOAD=`echo -n $payload | openssl base64 -e -A | sed s/\+/-/ | sed -E s/=+$//`
JWT_SIGNATURE=`echo -n "$JWT_HEADER.$JWT_PAYLOAD" | openssl dgst -sha256 -binary -sign private.pem  | openssl asn1parse -inform DER  -offset 2 | grep -o "[0-9A-F]\+$" | tr -d '\n' | xxd -r -p | base64 -w0 | tr -d '=' | tr '+/' '-_'`
TOKEN=$JWT_HEADER.$JWT_PAYLOAD.$JWT_SIGNATURE
```

The resulting JWT token must be included
via an HTTP `Authorization: Bearer <access-token>` header on any REST calls used to access EdgeX services via the API Gateway. 

Example:

```bash
$ curl -k -X GET https://localhost:8443/coredata/api/v1/ping? -H "Authorization: Bearer $TOKEN"
```

Additional users can be added by repeatedly calling the secrets-config command as above. Only one user can however be set at any time when using the snap configure hook, so
the current user must first be removed by setting the snap configuration settings to an empty string, before setting the values again:

```bash
$ sudo snap set edgexfoundry env.security-proxy.user=""
$ sudo snap set edgexfoundry env.security-proxy.public-key=""
$ sudo snap set edgexfoundry env.security-proxy.user=user02,USER_ID2,ES256
$ sudo snap set edgexfoundry env.security-proxy.public-key="$(cat public.pem)"
```

##### OAuth 2.0

OAuth2 is implemented in Kong using the [OAuth 2.0 plugin](https://docs.konghq.com/hub/kong-inc/oauth2/)

The OAuth2 implementation uses the [Client Credentials flow](https://tools.ietf.org/html/rfc6749#section-4.4). In that, the token is generated in two steps

1. The user is created with associated client_id and client_secret values, which can either be specified or generated by the plugin.

2. The client_id and client_secret values are then used to authorize with the authorization server (the Kong plugin) and request an access token from the token endpoint.

The secrets-config command supports this OAuth2 workflow. However, there are some limitations in the EdgeX implementation. 

- The authentication method (jwt or oauth2) needs to be set before EdgeX is installed. During initialization it will set up Kong with the correct authentication plugins and that can only be done once.
- The token can only be created on the system running EdgeX (i.e. running the secrets-config command)
- The generated tokens have infinite TTL.

Therefore, to change the authentication to oauth2,  you would need to use a custom gadget snap, which would set the following value to be used during the snap initialization:

```bash
$ sudo snap set edgexfoundry env.security-proxy.kongauth.name="oauth2"
```

Once you have this set up, then there are two steps to create the token 

To create the user and set up the OAuth2 authorization server, do:

```bash
$ edgexfoundry.secrets-config proxy adduser --token-type oauth2 --user user123
```

This will print out generated `client_secret` and `client_id` values.

To then get the token use the following:
```bash
$ edgexfoundry.secrets-config proxy oauth2 --client_id <clientid> --client_secret <clientsecret>
```

The token can then be used as before:
```bash 
curl -k -X GET https://localhost:8443/coredata/api/v1/ping? -H "Authorization: Bearer $TOKEN"
```

#### API Gateway TLS Certificate Setup

By default Kong is configured with an EdgeX signed TLS certificate. Client validation of this certificate requires the root CA certificate from the EdgeX instance. This file
(`ca.pem`) can be copied from directory `$SNAP_DATA/secrets/ca`.

It is also possible to install your own TLS certificate to be used by the gateway. The steps to do so are as follows:

Start by provisioning a TLS certificate to use. One way to do so for testing purposes is to use the edgeca snap:

```bash
$ sudo snap install edgeca
$ edgeca gencsr --cn localhost --csr csrfile --key csrkeyfile
$ edgeca gencert -o localhost.cert -i csrfile -k localhost.key
```

Then install the certificate:
```bash
$ sudo snap set edgexfoundry env.security-proxy.tls-certificate="$(cat localhost.cert)"
$ sudo snap set edgexfoundry env.security-proxy.tls-private-key="$(cat localhost.key)"
```

This sample certificate is signed by the EdgeCA root CA, so by specifying the Root CA certificate for validation then a connection can now be made using your new certificate:

```bash
$ curl -v --cacert /var/snap/edgeca/current/CA.pem -X GET https://localhost:8443/coredata/api/v1/ping? -H "Authorization: Bearer $TOKEN"
```

To set the certificate again, you first need to clear the current setting by setting the values to an empty string:

```bash
$ sudo snap set edgexfoundry env.security-proxy.tls-certificate=""
$ sudo snap set edgexfoundry env.security-proxy.tls-private-key=""
```

If you are using a different server name than localhost, then it needs to be specified first using the optional tls-sni configuration setting. Example:

```bash
$ edgeca gencsr --cn server01 --csr csrfile --key csrkeyfile
$ edgeca gencert -o server.cert -i csrfile -k server.key
$ sudo snap set edgexfoundry env.security-proxy.tls-certificate=""
$ sudo snap set edgexfoundry env.security-proxy.tls-private-key=""
$ sudo snap set edgexfoundry env.security-proxy.tls-sni="server01"
$ sudo snap set edgexfoundry env.security-proxy.tls-certificate="$(cat server.cert)"
$ sudo snap set edgexfoundry env.security-proxy.tls-private-key="$(cat server.key)"

$ curl -v --cacert /var/snap/edgeca/current/CA.pem -X GET https://server01:8443/coredata/api/v1/ping? -H "Authorization: Bearer $TOKEN"
```





## Limitations

[See the GitHub issues with label snap for current issues.](https://github.com/edgexfoundry/edgex-go/issues?q=is%3Aopen+is%3Aissue+label%3Asnap)

## Service Environment Configuration Overrides
**Note** - all of the configuration options below must be specified with the prefix: 'env.<service>.' where '<service>' is one of the following:

  - core-command
  - core-data
  - core-metadata
  - support-notifications
  - support-scheduler
  - sys-mgmt-agent
  - security-proxy
  - security-secret-store
  - device-virtual
  - app-service-config

```
[Service]
service.boot-timeout            // Service.BootTimeout
service.check-interval          // Service.CheckInterval
service.host                    // Service.Host
service.server-bind-addr        // Service.ServerBindAddr
service.port                    // Service.Port
service.protocol                // Service.Protocol
service.max-result-count        // Service.MaxResultCount
service.read-max-limit          // Service_ReadMaxLimit         // app-service-configurable only
service.startup-msg             // Service.StartupMsg
service.timeout                 // Service.Timeout

[Clients.CoreData]
clients.data.port                 // Clients.CoreData.Port

[Clients.MetaData]
clients.metadata.port             // Clients.MedataData.Port

[Clients.Notifications]
clients.notifications.port        // Clients.Notifications.Port

[Clients.Scheduler]
clients.scheduler.port            // Clients.Scheduler.Port   // sys-mgmt-only

[MessageQueue]
messagequeue.topic                // MessageQueue.Topic       // core-data only

[SecretStore]
secretstore.additional-retry-attempts    // SecretStore.AdditionalRetryAttempts
secretstore.retry-wait-period            // SecretStore.RetryWaitPeriod
```

### API Gateway Settings (prefix: env.security-proxy.)

```add-proxy-route```

The add-proxy-route setting is a csv list of URLs to be added to the
API Gateway (aka Kong). For references:

https://docs.edgexfoundry.org/1.3/microservices/security/Ch-APIGateway/

NOTE - this setting is not a configuration override, it's a top-level
environment variable used by the security-proxy-setup.


```
[KongAuth]
kongauth.name                  // KongAuth.Name [ 'jwt' (default) or 'oauth2'
```

### Secret Store Settings (prefix: env.security-secret-store.)
```add-secretstore-tokens```

The add-secretstore-tokens setting is a csv list of service keys to be added
to the list of Vault tokens that security-file-token-provider (launched by
security-secretstore-setup) creates.

NOTE - this setting is not a configuration override, it's a top-level
environment variable used by the security-secretstore-setup.

### Support Notifications Settings (prefix: env.support-notifications.)
```
[Smtp]
smtp.host                      // Smtp.Host
smtp.username                  // Smtp.Username
smtp.password                  // Smtp.Password
smtp.port                      // Smtp.Port
smtp.sender                    // Smtp.Sender
smtp.enable-self-signed-cert   // Smtp.EnableSelfSignedCert

```

## Building

The snap is built with [snapcraft](https://snapcraft.io), and the snapcraft.yaml recipe is located within `edgex-go`, so the first step for all build methods involves cloning this repository:

```bash
$ git clone https://github.com/edgexfoundry/edgex-go
$ cd edgex-go
```

### Installing snapcraft

There are a few different ways to install snapcraft and use it, depending on what OS you are building on. However after building, the snap can only be run on a Linux machine (either a VM or natively). To install snapcraft on a Linux distro, first [install support for snaps](https://snapcraft.io/docs/installing-snapd), then install snapcraft as a snap with:

```bash
$ sudo snap install snapcraft
```

(note you will be promted to acknowledge you are installing a classic snap - use the `--classic` flag to acknowledge this)

**Note** - currently the snap doesn't support cross-compilation, and must be built natively on the target architecture. Specifically, to support cross-compilation the kong/lua parts must be modified to support cross-compilation. The openresty part uses non-standard flags for handling cross-compiling so all the flags would have to manually passed to build that part. Also luarocks doesn't seem to easily support cross-compilation, so that would need to be figured out as well.

#### Running snapcraft on MacOS

To install snapcraft on MacOS, see [this link](https://snapcraft.io/docs/install-snapcraft-on-macos). After doing so, follow in the below build instructions for "Building with multipass"

#### Running snapcraft on Windows

To install snapcraft on Windows, you will need to run a Linux VM and follow the above instructions to install snapcraft as a snap. Note that if you are using WSL, only WSL2 with full Linux kernel support will work - you cannot use WSL with snapcraft and snaps. If you like, you can install multipass to launch a Linux VM if your Windows machine has Windows 10 Pro or Enterprise with Hyper-V support. See this [forum post](https://discourse.ubuntu.com/t/installing-multipass-for-windows/9547) for more details.

### Building with multipass

The easiest way to build the snap is using the multipass VM tool that snapcraft knows to use directly. After [installing multipass](https://multipass.run), just run 

```bash
$ snapcraft
```

### Building with LXD containers

Alternatively, you can instruct snapcraft to use LXD containers instead of multipass VM's. This requires installing LXD as documented [here](https://snapcraft.io/docs/build-on-lxd).

```bash
$ snapcraft --use-lxd
```

Note that if you are building on non-amd64 hardware, snapcraft won't be able to use it's default LXD container image, so you can follow the next section to create an LXD container to run snapcraft in destructive-mode natively in the container.

### Building inside external container/VM using native snapcraft

Finally, snapcraft can be run inside a VM, container or other similar build environment to build the snap without having snapcraft manage the environment (such as in a docker container where snaps are not available, or inside a VM launched from a build-farm without using nested VM's). 

This requires creating an Ubuntu 18.04 environment and running snapcraft (from the snap) inside the environment with `--destructive-mode`. 

#### LXD

Snaps run inside LXD containers just like they do outside the container, so all you need to do is launch an Ubuntu 18.04 container, install snapcraft and run snapcraft like follows:

```bash
$ lxc launch ubuntu:18.04 edgex
Creating edgex
Starting edgex
$ lxc exec edgex /bin/bash
root@edgex:~# sudo apt update && sudo apt install snapd squashfuse git -y
root@edgex:~# sudo snap install snapcraft --classic
root@edgex:~# git clone https://github.com/edgexfoundry/edgex-go
root@edgex:~# cd edgex-go && snapcraft --destructive-mode
```

#### Docker

Snapcraft is smart enough to detect when it is running inside a docker container specifically, to the point where no additional arguments are need to snapcraft when it is run inside the container. For example, the upstream snapcraft docker image can be used (only on x86_64 architectures unfortunately) like so:

```bash
$ docker run -it -v"$PWD":/build snapcore/snapcraft:stable bash -c "apt update && cd /build && snapcraft"
```

Note that if you are building your own docker image, you can't run snapd inside the container, and so to install snapcraft, the docker image must download the snapcraft snap and extract it as if it was installed normally inside `/snap` (same goes for the `core` and `core18` snaps). This is done by the Linux Foundation Jenkins server for the project's CI and you can see an example of that [here](https://github.com/edgexfoundry/ci-management/blob/master/shell/edgexfoundry-snapcraft.sh). The upstream docker image also does this, but only for x86_64 architectures.

#### Multipass / generic VM

To use multipass to create an Ubuntu 18.04 environment suitable for building the snap (i.e. when running natively on windows):

```bash
$ multipass launch bionic -n edgex-snap-build
$ multipass shell edgex-snap-build
multipass@ubuntu:~$ git clone https://github.com/edgexfoundry/edgex-go
multipass@ubuntu:~$ cd edgex-go
multipass~ubuntu:~$ sudo snap install snapcraft --classic
multipass~ubuntu:~$ snapcraft --destructive-mode
```

The process should be similar for other VM's such as kvm, VirtualBox, etc. where you create the VM, clone the git repository, then install snapcraft as a snap and run with `--destructive-mode`. 

### Developing the snap

After building the snap from one of the above methods, you will have a binary snap package called `edgexfoundry_<latest version>_<arch>.snap`, which can be installed locally with the `--devmode` flag:

```bash
$ sudo snap install --devmode edgexfoundry*.snap
```

In addition, if you are using snapcraft with multipass VM's, you can speedup development by not creating a *.snap file and instead running in "try" mode . This is done by running `snapcraft try` which results in a `prime` folder placed in the root project directory that can then be "installed" using `snap try`. For example:

```bash
$ snapcraft try # produces prime dir instead of *.snap file
...
You can now run `snap try /home/ubuntu/go/src/github.com/edgexfoundry/edgex-go/prime`.
$ sudo snap try --devmode prime # snap try works the same as snap install, but expects a directory
edgexfoundry 1.0.0-20190513+0620a8d1 mounted from /home/ubuntu/go/src/github.com/edgexfoundry/edgex-go/prime
$
```

#### Interfaces

After installing the snap, you will need to connect interfaces and restart the snap. The snap needs the `hardware-observe`, `mount-observe`, and `system-observe` interfaces connected. These are automatically connected using snap store assertions when installing from the store, but when developing the snap and installing a revision locally, use the following commands to connect the interfaces:

```bash
$ sudo snap connect edgexfoundry:hardware-observe
$ sudo snap connect edgexfoundry:system-observe
$ sudo snap connect edgexfoundry:mount-observe
```

After connecting these restart the services in the snap with:

```bash
$ sudo snap restart edgexfoundry
```

