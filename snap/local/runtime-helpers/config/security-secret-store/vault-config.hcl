# This file maintains a small number of changes vs. the vault configuration
# used for docker deployment, namely:
#
# - all hostnames are localhost instead of the docker name
# - tls is disabled (per ADR 0015)
# - tls listener parameters are all set to ""
#
listener "tcp" { 
  address = "localhost:8200" 
  tls_disable = "1"
  cluster_address = ""
  tls_min_version = ""
  tls_client_ca_file =""
  tls_cert_file =""
  tls_key_file = ""
}

backend "consul" {
  path = "vault/"
  address = "localhost:8500"
  scheme = "http"
  redirect_addr = "http://localhost:8200"
  cluster_addr = "http://localhost:8201"
}

default_lease_ttl = "168h"
max_lease_ttl = "720h"
