package main

import (
	"fmt"
	"github.com/slackhq/nebula/package/core"
)

func main() {
	ctrl, err := core.NewCtrl(core.Config{
		Pki: core.ConfigPki{
			CA:   "./example-ca.crt",
			Cert: "./example-a.crt",
			Key:  "./example-a.key",
		},
		StaticHostMap: make(map[string][]string),
		Lighthouse: core.ConfigLighthouse{
			AmLighthouse: true,
			ServeDNS:     false,
			DNS:          core.ConfigLighthouseDNS{},
			Interval:     0,
			Hosts: []string{
				"111.111.111.111",
			},
		},
		Listen: core.ConfigListen{
			Host: "0.0.0.0",
			Port: 4242,
		},
		Punchy: core.ConfigPunchy{
			Punch:   true,
			Respond: true,
		},
		Relay: core.ConfigRelay{
			Relays:    nil,
			AmRelay:   false,
			UseRelays: false,
		},
		Firewall: core.ConfigFirewall{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	ctrl.Start()
	//ctrl.ShutdownBlock()
}
