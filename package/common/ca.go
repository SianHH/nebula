package common

import (
	"crypto/rand"
	"net"
	"strings"
	"time"

	"github.com/slackhq/nebula/cert"
	"golang.org/x/crypto/ed25519"
)

type CaConfig struct {
	Name     string
	Duration time.Duration // 有效期 24 * time.Hour
	Groups   []string      // CA 用于限制签发的配置
	Ips      []string      // CA 用于限制签发的配置
	Subnet   []string      // CA 用于限制签发的配置
}

func Ca(cfg CaConfig) (outCert string, outKey string, err error) {
	var ips []*net.IPNet
	for _, rs := range cfg.Ips {
		rs := strings.Trim(rs, " ")
		if rs != "" {
			ip, ipNet, err := net.ParseCIDR(rs)
			if err != nil {
				return "", "", err
			}
			if ip.To4() == nil {
				return "", "", err
			}

			ipNet.IP = ip
			ips = append(ips, ipNet)
		}
	}

	var subnets []*net.IPNet
	for _, rs := range cfg.Subnet {
		rs := strings.Trim(rs, " ")
		if rs != "" {
			_, s, err := net.ParseCIDR(rs)
			if err != nil {
				return "", "", err
			}
			if s.IP.To4() == nil {
				return "", "", err
			}
			subnets = append(subnets, s)
		}
	}

	var curve cert.Curve
	var pub, rawPriv []byte
	curve = cert.Curve_CURVE25519
	pub, rawPriv, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      cfg.Name,
			Groups:    cfg.Groups,
			Ips:       ips,
			Subnets:   subnets,
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(cfg.Duration),
			PublicKey: pub,
			IsCA:      true,
			Curve:     curve,
		},
	}

	err = nc.Sign(curve, rawPriv)
	if err != nil {
		return "", "", err
	}

	b, err := nc.MarshalToPEM()
	if err != nil {
		return "", "", err
	}
	return string(b), string(cert.MarshalSigningPrivateKey(curve, rawPriv)), nil
}
