package main

import (
	"flag"
	"fmt"
	"github.com/skysoft-atm/consul-util/consul"
	"github.com/skysoft-atm/gorillaz"
	"github.com/skysoft-atm/supercaster/multicast"
	"github.com/skysoft-atm/supercaster/network"
	"strings"
)

const configMaxDatagramSize = "udp.maxDatagramSize"

func init() {
	flag.String(network.ConfigNetworkInterface, "", "Network interface to use.")
	flag.Int(configMaxDatagramSize, 8192, "Max datagram size.")
}

func main() {
	g := gorillaz.New(consul.ActivateServiceDiscovery(), consul.ActivateViperRemoteConfig())
	g.Run()

	streamToUdp := g.Viper.GetString("streamToUdp")
	udpToStream := g.Viper.GetString("udpToStream")

	netInterface := network.GetNetworkInterface(g)
	maxDatagramSize := g.Viper.GetInt(configMaxDatagramSize)

	for _, udpPub := range strings.Split(streamToUdp, "|") {
		p := strings.Split(udpPub, ">")
		if len(p) != 2 {
			panic("Error parsing udp publication " + udpPub)
		}
		serviceStream := p[0]
		addr := p[1]
		ss := strings.Split(serviceStream, "/")
		if len(p) != 2 {
			panic("Error parsing service stream " + serviceStream)
		}
		service := ss[0]
		stream := ss[1]

		gorillaz.Sugar.Infof("Publishing %s to %s", serviceStream, addr)
		go panicIf(func() error {
			err := network.ServiceStreamToUdp(service, stream, addr, g)
			return fmt.Errorf("error publishing %s to %s : %w", serviceStream, addr, err)
		})
	}

	for _, udpSub := range strings.Split(udpToStream, "|") {
		p := strings.Split(udpSub, ">")
		if len(p) != 2 {
			panic("Error parsing udp subscription " + udpSub)
		}
		addr := p[0]
		stream := p[1]

		source := network.UdpSource{
			NetInterface:    netInterface,
			HostPort:        addr,
			MaxDatagramSize: maxDatagramSize,
		}
		gorillaz.Sugar.Infof("Publishing %s to %s", addr, stream)
		go panicIf(func() error {
			err := multicast.UdpToStream(g, source, stream)
			return fmt.Errorf("error publishing %s to %s : %w", addr, stream, err)
		})
	}

	g.SetReady(true)
	select {}
}

func panicIf(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}
