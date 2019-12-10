# Server-Sent Events Simulation
[![asciicast](https://asciinema.org/a/QC7TGAAnCR7l1t68GBuyKrMXs.png)](https://asciinema.org/a/QC7TGAAnCR7l1t68GBuyKrMXs)
*HTTP Polling repository can be found [here](https://github.com/yevshev/polling)*

## Information
The `client` directory contains the source for the [sse-client](https://hub.docker.com/repository/docker/yevshev/sse-client) docker image 

The `server` directory contains the source for the [sse-server](https://hub.docker.com/repository/docker/yevshev/sse-server) docker image

## Create an overlay network
Create a virtual network for all of our containers to communicate with eachother, and name it 'ssenet':
```sh
$ docker network create -d overlay --attachable ssenet
```

## Deploying to Docker Swarm
Deploy the sse server containers defined in [servers.yml](https://github.com/yevshev/server-sent/blob/master/servers.yml), each running our Go sse server binary, and name it 'sse':

```sh
$ docker stack deploy -c servers.yml sse
```
Deploy the sse client container defined in [client.yml](https://github.com/yevshev/server-sent/blob/master/client.yml), running our Go sse client binary, and name it 'collector':

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

## Cleanup
Stop and delete the 'sse' server stack and 'events' client stack:
```sh
$ docker stack rm sse events
```
