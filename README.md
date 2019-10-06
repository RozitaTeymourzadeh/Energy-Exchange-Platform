# Energy Exchange Platform
To exchange energy between the householders 

## Dependancy

1. Goland IDE environment
2. Git 
3. EdgeX Foundary 


## Import Device-sdk-go

```linux
mkdir -p ~/go/src/github.com/edgexfoundry
cd ~/go/src/github.com/edgexfoundry
git clone https://github.com/edgexfoundry/device-sdk-go.git
mkdir device-simple
cp -rf ./device-sdk-go/example/* ./device-simple/
cp ./device-sdk-go/Makefile ./device-simple
cp ./device-sdk-go/Version ./device-simple/
cp ./device-sdk-go/version.go ./device-simple/

```

## Configure Device Service 

1.Main.go: 

import library

```linux
"github.com/edgexfoundary/device-simple/driver"
```
2. Configure Makefile:

```linux
MICROSERVICES=cmd/device-simple/device-simple
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-simple.Version=$(VERSION)"
cmd/device-simple/device-simple:
$(GO) build $(GOFLAGS) -o $@ ./cmd/device-simple
```
3. go.mod

```linux
GO111MODULE=on go mode init
```

## Build the Project

Make file is used to build the project. Build command in the build directory is:
```linux
make build
```
