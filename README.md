# MYOB test

This repo contains the source and to build a simple HTTP based API in Go with
the following endpoints:

- a simple root endpoint which responds in a simple manner
- a health endpoint which returns an appropriate response code
- a metadata endpoint which returns basic information about your application

It also contains a number of configurations to test, build and optionally
deploy the resulting binary using Travis, drone.io or Gitlab.

The goal is to demonstrate the building of a single deployable artifact and
describe the steps involved.

## Building

To build the binary there is a make file provided with the 'build' target to
produce a binary for the current host architecture. There is also the
'build-all' target that will build a suitable binary for Windows, Darwin, Linux
and FreeBSD:

```shell
# Single binary
make build

# Cross-compiled
make build-all
```

Alternatively the following commands can be used if make is not available:

```shell
# Single binary
go build -ldflags "-w -s -X main.version=$(git describe --always)" -o server ./cmd/api

# Cross-compiled for Darwin
GOOS=darwin go build -ldflags "-w -s -X main.version=$(git describe --always)" -o server ./cmd/api
```

The above command embed the current git version into the binary for use by one
of the API endpoints.

## Testing

The native Go testing framework is used to test the handlers. They also produce
a coverage report. They can be run via the following commands:

```shell
# Via make
make test

# or
go test -short -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -html=coverage.txt -o coverage.html
```

Current coverage is only about 63%.

## Deployment

The build procedure produces a single binary as the artifact. This can be
copied to the destination and run directly. For deployment to other
architectures from the build host the cross-compiled binaries will need to be
used.

A Dockerfile is provided for deployment via docker. This is a two stage
dockerfile producing a small image based on Alpine linux with a single binary
as the entrypoint. This can then be published to a registry for later
deployment.

Deployment to an AWS lambda function is done via the make target 'handler.zip'.
While this may not have been a requirement for the task (only a single handler
is published) it was just for fun.

## Usage

The executable accepts the following command line options for configuring the
listening address and port:

```shell
-h HOST
	Listen on HOST IP address. Default is 0.0.0.0.
-p PORT
	Listen at PORT. Default is 8080.
```

These can also be specified via environment variables as HOST and PORT
respectively. This enables easily customisable deployment via Docker etc.

## Assumptions, possible improvements or issues

- Building and deployment via Gitlab is based on experience but is untested.
- Building and deployment via drone.io is based on experience but is untested.
- Building and deployment via TravisCI is based on experience but is untested.
- Tests only cover the handlers and not the main process and thread
  creation/teardown.
- While logging is to STDOUT/STDERR for use by Docker, it could be structured
  for better consumption.
- Makefile assumes GNU make
