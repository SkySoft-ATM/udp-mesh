# udp-mesh

Use this project to encapsulate UDP broadcast and multicast into [gorillaz streams](https://github.com/SkySoft-ATM/gorillaz) over TCP in order to forward it from machine to machine when UDP broadcast and multicast is not supported by the network layer (VMs on the cloud for example). 

You need to pass a service name on startup with the flag ```service.name``` to define your app name.

You can optionally define an environment with the flag ```env``` that can be used in a multi-tenant environment to override default configs (see below)

Upon startup, it will look for a configuration in a [Consul](https://www.consul.io/) key-value service.

The consul endpoint can be configured with the flag ```consul.endpoint``` (defaults to localhost:8500).

The configuration will be loaded from the key corresponding to the service name. If a configuration exists with the key pattern "serviceName.env" then this configuration will override the one defined for the service name.

#### The configuration can contain any of these keys:
```broadcastToStream``` : list of host:port>streamName separated by ```|```, it is not necessary to provide a host for the broadcast

```multicastToStream``` : list of host:port>streamName separated by ```|```
 
```broadcastStream``` : list of streamName>host:port separated by ```|```, if no host:port is provided, the broadcast IP for the network interface and the original broadcast port will be used. You can override the port by simply providing the port.
 
```multicastStream``` : list of streamName>host:port separated by ```|``` , if no host:port is provided, the original multicast IP & port will be used.


#### Here is a configuration sample with several applications:

service name: trackgen

```
{
    "broadcastToStream": ":25910>tracks|:25909>status"
}
```

server01
```
{
    "broadcastStream": "trackgen/tracks|trackgen/status",
    "multicastToStream": "239.193.68.7:24011>setup|239.192.68.5:24005>distrib"
}
```

visu01
```
{
    "multicastStream": "server01/setup|server01/distrib",
    "broadcastStream": "trackgen/tracks|trackgen/status"
}
```
 

#### Necessary dependencies for running on Linux and Windows

This project is using raw sockets in order to spoof the source IP of the UDP packets that are republished.

Some reading:
https://css.bz/2016/12/08/go-raw-sockets.html

On windows, the pcap lib needs to be installed:
https://nmap.org/npcap/windows-10.html

On Linux:
sudo apt-get install libpcap-dev


