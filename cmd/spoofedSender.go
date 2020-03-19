package main

import (
	"flag"
	"fmt"
	"github.com/skysoft-atm/gorillaz"
	"github.com/skysoft-atm/gorillaz/stream"
	"github.com/skysoft-atm/supercaster/network"
	"time"
)

func init() {
	flag.String(network.ConfigNetworkInterface, "Ethernet", "Network interface to use.")
}

// Sends test values on a UDP raw socket

// to run it first give the permission for raw sockets
// sudo setcap cap_net_raw=ep ./spoofedSender

// go build -o spoofedSender && sudo setcap cap_net_raw=ep ./spoofedSender && ./spoofedSender --network.interface=wlp59s0

const hostPort = "224.0.0.23:9999"

func main() {
	g := gorillaz.New(gorillaz.WithServiceName("SpoofedUdpSender"))
	<-g.Run()

	udpPub := network.UdpPub{
		HostPort:      hostPort,
		InterfaceName: g.Viper.GetString(network.ConfigNetworkInterface),
		Type:          network.Multicast,
	}

	err := network.StreamToUdpSpoofSourceAddr(MockedSource{}, udpPub)
	if err != nil {
		panic(err)
	}
	select {}
}

type MockedSource struct {
}

func (MockedSource) EvtChan() chan *stream.Event {
	res := make(chan *stream.Event, 1)
	go sendOnChan(res)
	return res

}

func sendOnChan(sink chan *stream.Event) {
	tick := time.NewTicker(1 * time.Second)
	i := 0
	for range tick.C {
		sink <- &stream.Event{
			Key:   []byte("10.0.0.5:6789"), // the source IP to mock on the publication
			Value: []byte(fmt.Sprintf("mocked value %d", i)),
		}
		i++
	}
}
