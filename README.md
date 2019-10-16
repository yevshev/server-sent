# sse-project

## Information
The `client` directory contains the source for the [collector]() docker image
The `server` directory contains the source for the [sse-server]() docker image

## Deploying to Docker Swarm
First, generate a personal access token from your dockerhub [settings](https://hub.docker.com/settings/security) page.

Because our image registry is private, we'll need to log in to our dockerhub account in order to pull our images.

```sh
# log in to your docker account
$ docker login --username <username>

# At the password prompt, enter your personal access token.
```
Deploy 100 containers with each running our go binary, and name it 'test-cluster':

```sh
$ docker service create --replicas 100 --name test-cluster yevshev/emulator-push

# We can provide addtional parameters as needed
```
Additional docker service [commands](https://docs.docker.com/engine/reference/commandline/service/
):
```sh
# List running docker services
$ docker service ls

# Display detailed information on the service
$ docker service inspect test-cluster --pretty

# Scale our service to 1000 containers
$ docker service scale test-cluster=1000

# Remove the service
$ docker service rm test-cluster

# stop all running containers
$ docker stop $(docker ps -a -q)

# Delete all containers
$ docker rm $(docker ps -a -q)
```
