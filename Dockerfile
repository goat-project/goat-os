FROM golang:1.17-alpine

ARG branch=master
ARG version

ENV name="goat-os" \
    user="goat"
ENV project="/go/src/github.com/goat-project/${name}/" \
    homeDir="/var/lib/${user}/" \
    logDir="/var/${name}/log/"

LABEL application=${name} \
      description="Exporting Openstack accounting data" \
      maintainer="svetlovska@cesnet.cz" \
      version=${version} \
      branch=${branch}

# Install tools required for project
RUN apk add --no-cache git shadow

WORKDIR ${project}

# Create user and log directory
RUN useradd --system --shell /bin/false --home ${homeDir} --create-home --uid 1000 ${user} && \
    usermod -L ${user} && \
    mkdir -p ${logDir} && \
    chown -R ${user}:${user} ${logDir}

# Copy the entire project and build it, dependencies are installed with build
COPY . ${project}
RUN go build -o /bin/${name}

# Switch user
USER ${user}

# Run main command with subcommands and options:
# No subcommand runs goat-os, configures it from config file (config/goat-os.yml)
# and extracts virtual machine, network and storage data in the same time.
# To extract only specific data, use subcommand:
#   network     Extract network data
#   storage     Extract storage data
#   vm          Extract virtual machine data
#   help        Help about any command
#
# The configuration from file should be rewrite using the following flags:
#      --allow-reauth string                    Openstack authentication allow reauth. [OS_ALLOW_REAUTH]
#      --application-credential-id string       Openstack application credential TenantID [OS_APPCREDENTIAL_ID]
#      --application-credential-name string     Openstack application credential name [OS_APPCREDENTIAL_NAME]
#      --application-credential-secret string   Openstack application credential secret [OS_APPCREDENTIAL_SECRET]
#      --availability string                    Openstack endpoint availability (public, internal, admin)
#  -d, --debug string                           debug
#      --domain-id string                       Openstack authentication domain TenantID [OS_DOMAIN_ID]
#      --domain-name string                     Openstack authentication domain name [OS_DOMAIN_NAME]
#  -e, --endpoint string                        goat server [GOAT_SERVER_ENDPOINT] (required)
#  -h, --help                                   help for goat-os
#  -i, --identifier string                      goat identifier [IDENTIFIER] (required)
#      --log-path string                        path to log file
#      --name string                            Openstack endpoint name [OS_ENDPOINT_NAME]
#  -o, --openstack-identity-endpoint string     Openstack identity endpoint [OS_IDENTITY_ENDPOINT] (required)
#      --passcode string                        Openstack authentication passcode [OS_PASSCODE]
#      --password string                        Openstack authentication password [OS_PASSWORD]
#  -p, --records-for-period string              records for period [TIME PERIOD]
#  -f, --records-from string                    records from [TIME]
#  -t, --records-to string                      records to [TIME]
#      --region string                          Openstack endpoint region [OS_ENDPOINT_REGION]
#      --scope-domain-id string                 Openstack scope domain TenantID [OS_SCOPE_DOMAIN_ID]
#      --scope-domain-name string               Openstack scope domain name [OS_SCOPE_DOMAIN_NAME]
#      --scope-project-id string                Openstack scope project TenantID [OS_SCOPE_PROJECT_ID]
#      --scope-project-name string              Openstack scope project name [OS_SCOPE_PROJECT_NAME]
#      --scope-system string                    Openstack scope system [OS_SCOPE_SYSTEM]
#      --tenant-id string                       Openstack authentication tenant TenantID [OS_TENANT_ID]
#      --tenant-name string                     Openstack authentication tenant name [OS_TENANT_NAME]
#      --token-id string                        Openstack authentication token TenantID [OS_TOKEN_ID]
#      --type string                            Openstack endpoint type [OS_ENDPOINT_TYPE]
#      --user-id string                         Openstack authentication user TenantID [OS_USER_ID]
#      --username string                        Openstack authentication username [OS_USERNAME]
#  -v, --version                                version for goat-os
#
# Example:
# - extract virtual machine data from the last 5 years and save it with idetifier 'goat-vm'
# CMD /bin/goat-os vm --log-path=${logDir}${name}.log -p=5y -i=goat-vm
#
# - extract storage data
# CMD /bin/goat-os storage
#
# Extract all data about servers (vm), storage, ip (network)
CMD /bin/goat-os --log-path=${logDir}${name}.log