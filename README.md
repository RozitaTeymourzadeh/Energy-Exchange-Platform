# Energy Exchange Platform
To exchange energy between the householders 

## Setup Docker on Raspberry 

```linux
curl -sSL https://get.docker.com | sh
docker run --rm hello-world
```

If permission error is observed, run the following options:

**Option #1**
```linux
sudo docker run --rm hello-world
```
**Option #2**
```linux
usermod -aG docker <your_user>
```

**Install Python modules:**
```linux
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py && sudo python3 get-pip.py
sudo apt-get install libssl-dev
sudo apt-get install libffi-dev
```
**Install vim:**
```linux
sudo apt-get update && sudo apt-get install vim -y
```

**Install docker-compose:**
```linux
sudo pip3 install docker-compose
vim docker-compose.yml
```
You can add your device in docke-compose.yml file.

```linux
version: '3'
services:
  webapp:
    ports:
      - 5000:8000
    image: python:3.7-alpine
    command: "python -m http.server 8000"
```
**Run the docker compose:**

```linux
sudo docker-compose up -d
sudo docker-compose down
```

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
## Run the Project

Run the binery generated in the follwing directory, in this example simple-device

```linux
cd cmd
cd simple-device
./simple-device
```
## Kill or delete the generated devices

For re-run the application, it is required to kill or delete the devices that is running in docker background to avoid panic error in docker while is running to do so, do the following:

***Delete using Postman  

```linux
DELETE clocalhost:48081/api/v1/device/id/DEVICE_ID
```
```linux
docker-compose up -d --remove-orphans
```
