# one
One ctl to rule them all

![](https://serhiy.s3.eu-central-1.amazonaws.com/Github_repo/one/cover.png)

## Disclaimer
It's a personal project that helps to manage a lot of microservices in different environments and with different configurations. The project is in the early development stage.
The main purpose of the project is to simplify the management of microservices like building the docker image, showing service logs, run/stop/restart the service, etc.

## Features
* build the docker container with predefined params
* easy install/uninstall service
* update the service with zero-downtime
* specifying the context per app/project
* easy logs and better list

## Installation
### Binary release
You can manually download a binary release [here](https://github.com/exelban/one/releases).

### curl
```bash
curl -sfL https://api.serhiy.io/one/install.sh | sh
```

### go
```bash
go install github.com/exelban/one@latest
```

## Commands
There are a few commands that are available in the app:

- `one build` / `one b` - build the docker image.
- `one start` / `one install` - start the service. Use `-c` to copy docker-compose to the remote server.
- `one stop` / `one uninstall` - stop the service. Use `-c` to copy docker-compose to the remote server.
- `one restart` / `one r` - restart the service. Use `-c` to copy docker-compose to the remote server.
- `one logs` / `one l` - show the logs of the service. Use `-f` to follow the logs.
- `one list` - list all available services.
- `one context` - managing the app context.

## Configuration
Project configuration is stored in the `.one` file. It contains the information about the project and environment.  
If the name or image is not provided in the configuration cli will try to detect the docker-compose file. If the file is found the name or image will be used from the docker-compose file.

Parameters:

- `name` - the name of the service
- `context` - the name or id of the context
- `build` - the build configuration
- `ssh` - the ssh configuration

### build
- `file` - the path to the Dockerfile
- `image` - the docker image name
- `push` - push the docker image to the registry after the build
- `platforms` - the list of platforms for the docker image. Dockerx will be used to build the image for different platforms instead of the default docker build.
- `args` - the list of build arguments
- `force` - force restart the service when upgrade

### ssh
- `host` - the host of the remote server
- `port` - the port of the remote server (22 by default)
- `username` - the username of the remote server
- `password` - the password of the remote server
- `privateKey` - the path to the private key
- `swarmMode` - the docker swarm mode

#### Config with context:
```yaml
name: "ping"
context: "prod"
docker:
  image: "exelban/ping:latest"
```

#### Config without context:
```yaml
name: "ping"
docker:
  image: "exelban/ping:latest"
  platforms: ["linux/amd64", "linux/arm64"]
  force: true
ssh:
  host: "staging"
  username: "root"
  privateKey: "/path/to/private/key"
```

## Context
The context keeps the information about the remote server. It allows to have multiple environments and switch between them easily. The context could be active and `one` will use that context to execute the commands.

### Commands
- `one context` - show the active context
- `one context list` - list all available contexts
- `one context add` - add a new context
- `one context delete [name/id]` - remove the context
- `one context activate [name/id]` - activate the context
- `one context deactivate` - deactivate active context

### Creating the context
Parameters for the context:

- `name` - the name of the context
- `host` - the host of the remote server
- `username` - the username of the remote server
- `password` - the password of the remote server
- `private-key` - the path to the private key
- `swarm` - the docker swarm mode

#### Example creating the context
```bash
one context add --name=prod --host=prod.example.com --username=root --private-key=/path/to/private/key
```

## License
[MIT License](https://github.com/exelban/one/blob/master/LICENSE)