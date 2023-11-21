FROM golang:1.12-alpine

ENV name="goat-os" \
    user="goat-os"
ENV project="/go/src/github.com/goat-project/${name}/" \
    homeDir="/var/lib/${user}/" \
    logDir="/var/${name}/log/" \
    outputDir="/var/${name}/output/"

LABEL application=${name} \
    description="Extracting data from Openstack" \
    maintainer="mattim314@github" 

# Install tools required for project
RUN apk add --no-cache git shadow
RUN go get -u github.com/goat-project/goat-os

# List project dependencies with Gopkg.toml and Gopkg.lock
COPY go.modbgo.sum ${project}
# Install library dependencies
WORKDIR ${project}
RUN dep ensure -vendor-only

# Create user and directories
RUN useradd --system --shell /bin/false --home ${homeDir} --create-home --uid 1000 ${user} && \
    usermod -L ${user} && \
    mkdir -p ${logDir} ${outputDir} ${templatesDir} && \
    chown -R ${user}:${user} ${logDir} ${outputDir} ${templatesDir}

# Copy the entire project and build it
COPY . ${project}
RUN go build -o /bin/${name}

# Switch user
USER ${user}

# Expose port for goat
# Unnecessary since access to goat-os is not required during data gathering
# EXPOSE 9623

# Run main command with following options:
#   -d, --debug string                 debug
#   -e, --endpoint string              goat server [GOAT_SERVER_ENDPOINT] (required)
#   -h, --help                         help for goat-os
#   -i, --identifier string            goat identifier [IDENTIFIER] (required)
#       --log-path string              path to log file
#   -o, --openstack-endpoint string    Openstack endpoint [OPENSTACK_ENDPOINT] (required)
#   -s, --openstack-secret string      Openstack secret [OPENSTACK_SECRET] (required)
#   -p, --records-for-period string    records for period [TIME PERIOD]
#   -f, --records-from string          records from [TIME]
#   -t, --records-to string            records to [TIME]
#       --version                      version for goat-os

CMD /bin/goat-os