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
	flag.String(network.ConfigNetworkInterface, "Ethernet", "Network interface to use.")
	flag.Int(configMaxDatagramSize, 8192, "Max datagram size.")
}

// go build -o mesh && sudo setcap cap_net_raw=ep ./mesh && ./mesh --network.interface=wlp59s0

func main() {
	g := gorillaz.New(consul.ActivateServiceDiscovery(), consul.ActivateViperRemoteConfig())
	g.Run()

	multicastStream := g.Viper.GetString("multicastStream")
	broadcastStream := g.Viper.GetString("broadcastStream")
	udpToStream := g.Viper.GetString("udpToStream")
	interfaceName := g.Viper.GetString(network.ConfigNetworkInterface)

	netInterface := network.GetNetworkInterface(interfaceName)

	maxDatagramSize := g.Viper.GetInt(configMaxDatagramSize)

	if multicastStream != "" {
		createPublication(multicastStream, interfaceName, g, network.Multicast)
	} else {
		gorillaz.Log.Info("No UDP multicast configured")
	}

	if broadcastStream != "" {
		createPublication(broadcastStream, interfaceName, g, network.Broadcast)
	} else {
		gorillaz.Log.Info("No UDP broadcast configured")
	}

	if udpToStream != "" {
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
	} else {
		gorillaz.Log.Info("No UDP to stream configured")
	}

	g.SetReady(true)
	select {}
}

func createPublication(streamDef string, interfaceName string, g *gorillaz.Gaz, pubType network.UdpPubType) {
	for _, udpPub := range strings.Split(streamDef, "|") {
		p := strings.Split(udpPub, ">")
		if len(p) != 2 {
			panic("Error parsing udp publication " + udpPub)
		}
		serviceStream := p[0]
		hostPort := p[1]
		ss := strings.Split(serviceStream, "/")
		if len(ss) != 2 {
			panic("Error parsing service stream " + serviceStream)
		}
		service := ss[0]
		stream := ss[1]

		gorillaz.Sugar.Infof("Publishing %s to %s", serviceStream, hostPort)
		go panicIf(func() error {
			err := network.ServiceStreamToUdpSpoofSourceAddr(service, stream, interfaceName, hostPort, g, pubType)
			return fmt.Errorf("error publishing %s to %s : %w", serviceStream, hostPort, err)
		})
	}
}

func panicIf(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}
