# Server-Sent Events Simulation
*HTTP Polling repository can be found [here](https://github.com/yevshev/server-sent)*

## Information
The `client` directory contains the source for the [sse-client](https://hub.docker.com/repository/docker/yevshev/sse-client) docker image 

The `server` directory contains the source for the [sse-server](https://hub.docker.com/repository/docker/yevshev/sse-server) docker image

## Deploying to Docker Swarm
Deploy the sse server containers defined in [servers.yml](https://github.com/yevshev/server-sent/blob/master/servers.yml), each running our Go sse server binary, and name it 'sse':

```sh
$ docker stack deploy -c servers.yml sse
```
Deploy the sse client container defined in [client.yml](https://github.com/yevshev/server-sent/blob/master/client.yml), runnin our Go sse client binary, and name it 'collector':

```sh
$ docker stack deploy -c client.yml events
```

## Client logs
View a real-time feed of all the data the client receives from the servers:
```sh
$ docker service logs events_client -f
```

## Performance Metrics
View a real-time feed of resource utilization and performance for each runing container:
```sh
$ docker stats
```