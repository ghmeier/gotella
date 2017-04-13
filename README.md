# gotella
A gnutella node in go.

# To Run
1. Install Go ([getting started](https://golang.org/doc/install))
2. Run redis locally ([quickstart](https://redis.io/topics/quickstart))
3. `go get github.com/ghmeier/gotella`
4. `go install github.com/ghmeier/gotella`
5. Run:

```
gotella listen_port redis_port discovery_host
```

`listen_port`: the port where this gotella node will listen for descriptors (ex: `9000`)

`redis_port`: the port where your local redis server is listening (ex: `6379`).

`discovery_host`: (optional) the ip:port of the first node to connect to the network (ex: `10.27.252.211:8000`)

### To Request a File:

Type the desired filename into a running instance of gotella and it'll search the network for the desired file. All files will be loaded into the `./public` folder of your running gotella node. All files in the folder will be available to be streamed to connected nodes.
