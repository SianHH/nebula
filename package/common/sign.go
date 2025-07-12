package common

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/slackhq/nebula/cert"
	"golang.org/x/crypto/curve25519"
)

type SignConfig struct {
	Name       string
	Ip         string
	CaKeyFile  string
	CaCertFile string
	Groups     []string
	Subnet     []string
}

func Sign(cfg SignConfig) (outCert string, outKey string, err error) {
	var curve cert.Curve
	var caKey []byte

	// naively attempt to decode the private key as though it is not encrypted
	caKey, _, curve, err = cert.UnmarshalSigningPrivateKey([]byte(cfg.CaKeyFile))
	if err != nil {
		fmt.Println("1")
		return "", "", err
	}

	caCert, _, err := cert.UnmarshalNebulaCertificateFromPEM([]byte(cfg.CaCertFile))
	if err != nil {
		fmt.Println("2")
		return "", "", err
	}

	if err := caCert.VerifyPrivateKey(curve, caKey); err != nil {
		fmt.Println("3")
		return "", "", err
	}

	issuer, err := caCert.Sha256Sum()
	if err != nil {
		fmt.Println("4")
		return "", "", err
	}

	if caCert.Expired(time.Now()) {
		fmt.Println("5")
		return "", "", err
	}

	duration := time.Until(caCert.Details.NotAfter) - time.Second*1

	ip, ipNet, err := net.ParseCIDR(cfg.Ip)
	if err != nil {
		fmt.Println("6")
		return "", "", err
	}
	if ip.To4() == nil {
		fmt.Println("7")
		return "", "", err
	}
	ipNet.IP = ip

	subnets := []*net.IPNet{}
	for _, rs := range cfg.Subnet {
		rs := strings.Trim(rs, " ")
		if rs != "" {
			_, s, err := net.ParseCIDR(rs)
			if err != nil {
				fmt.Println("8")
				return "", "", err
			}
			if s.IP.To4() == nil {
				fmt.Println("9")
				return "", "", err
			}
			subnets = append(subnets, s)
		}
	}

	pub, rawPriv := newKeypair(curve)

	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      cfg.Name,
			Ips:       []*net.IPNet{ipNet},
			Groups:    cfg.Groups,
			Subnets:   subnets,
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(duration),
			PublicKey: pub,
			IsCA:      false,
			Issuer:    issuer,
			Curve:     curve,
		},
	}

	if err := nc.CheckRootConstrains(caCert); err != nil {
		fmt.Println("10")
		return "", "", err
	}

	if err = nc.Sign(curve, caKey); err != nil {
		fmt.Println("11")
		return "", "", err
	}

	b, err := nc.MarshalToPEM()
	if err != nil {
		fmt.Println("12")
		return "", "", err
	}

	return string(cert.MarshalPrivateKey(curve, rawPriv)), string(b), nil
}

func newKeypair(curve cert.Curve) ([]byte, []byte) {
	switch curve {
	case cert.Curve_CURVE25519:
		return x25519Keypair()
	case cert.Curve_P256:
		return p256Keypair()
	default:
		return nil, nil
	}
}

func x25519Keypair() ([]byte, []byte) {
	privkey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, privkey); err != nil {
		panic(err)
	}

	pubkey, err := curve25519.X25519(privkey, curve25519.Basepoint)
	if err != nil {
		panic(err)
	}

	return pubkey, privkey
}

func p256Keypair() ([]byte, []byte) {
	privkey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey := privkey.PublicKey()
	return pubkey.Bytes(), privkey.Bytes()
}
