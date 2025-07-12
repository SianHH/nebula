package main

import (
	"github.com/slackhq/nebula/package/common"
	"os"
	"time"
)

func main() {
	//"example-ca", "./example-ca.key", "./example-ca.crt"
	caCert, caKey, err := common.Ca(common.CaConfig{
		Name:     "example-ca",
		Duration: time.Hour * 24 * 365,
		Groups:   nil,
		Ips:      nil,
		Subnet:   nil,
	})
	if err != nil {
		panic(err)
	}
	os.WriteFile("example-ca.crt", []byte(caCert), 0644)
	os.WriteFile("example-ca.key", []byte(caKey), 0644)
	signCert, signKey, err := common.Sign(common.SignConfig{
		Name:       "device1",
		Ip:         "10.100.100.1/24",
		CaKeyFile:  caKey,
		CaCertFile: caCert,
		Groups:     nil,
		Subnet:     nil,
	})
	if err != nil {
		panic(err)
	}
	os.WriteFile("example-device1.crt", []byte(signCert), 0644)
	os.WriteFile("example-device1.key", []byte(signKey), 0644)
}
