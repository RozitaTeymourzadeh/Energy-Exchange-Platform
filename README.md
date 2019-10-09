# Energy Exchange Platform
To exchange energy between the householders 

## Setup Docker on Raspberry 

```linux
curl -sSL https://get.docker.com | sh
docker run --rm hello-world
```

If permission error is observed, run the following options:

**Option #1
```linux
sudo docker run --rm hello-world
```
**Option #2
```linux
usermod -aG docker <your_user>
```

**Install Python modules:
```linux
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py && sudo python3 get-pip.py
sudo apt-get install libssl-dev
sudo apt-get install libffi-dev
```
**Install vim:
```linux
sudo apt-get update && sudo apt-get install vim -y
```

**Install docker-compose:**
```linux
sudo pip3 install docker-compose
vim docker-compose.yml
```
You can add your device in docke-compose.yml file

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
GO111MODULE=on go mod init
```
Alternatively, the go.mod from existing SDK device can be used in simple device example, directory root on go.mod file should be updated.

## Build the Project

Make file is used to build the project. Build command in the build directory is:
```linux
make build
```
