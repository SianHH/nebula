package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula"
	"github.com/slackhq/nebula/config"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	Pki           ConfigPki           `yaml:"pki,omitempty"`
	StaticHostMap map[string][]string `yaml:"static_host_map,omitempty"`
	Lighthouse    ConfigLighthouse    `yaml:"lighthouse,omitempty"`
	Listen        ConfigListen        `yaml:"listen,omitempty"`
	Punchy        ConfigPunchy        `yaml:"punchy,omitempty"`
	Relay         ConfigRelay         `yaml:"relay,omitempty"`
	Firewall      ConfigFirewall      `yaml:"firewall,omitempty"`
}

type ConfigPki struct {
	CA   string `yaml:"ca,omitempty"`
	Cert string `yaml:"cert,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

type ConfigLighthouse struct {
	AmLighthouse bool                `yaml:"am_lighthouse,omitempty"`
	ServeDNS     bool                `yaml:"serve_dns,omitempty"`
	DNS          ConfigLighthouseDNS `yaml:"dns,omitempty"`
	Interval     int                 `yaml:"interval,omitempty"`
	Hosts        []string            `yaml:"hosts,omitempty"`
}

type ConfigLighthouseDNS struct {
	Host string `yaml:"host,omitempty"`
	Port bool   `yaml:"port,omitempty"`
}

type ConfigListen struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

type ConfigPunchy struct {
	Punch        bool          `yaml:"punch,omitempty"`
	Respond      bool          `yaml:"respond,omitempty"`
	Delay        time.Duration `yaml:"delay,omitempty"`
	RespondDelay time.Duration `yaml:"respondDelay,omitempty"`
}

type ConfigRelay struct {
	Relays    []string `yaml:"relays,omitempty"`
	AmRelay   bool     `yaml:"am_relay,omitempty"`
	UseRelays bool     `yaml:"use_relays,omitempty"`
}

type ConfigFirewall struct {
	OutboundAction string                `json:"outbound_action,omitempty"`
	InboundAction  string                `json:"inbound_action,omitempty"`
	Outbound       []ConfigFirewallBound `json:"outbound,omitempty"`
	Inbound        []ConfigFirewallBound `json:"inbound,omitempty"`
}

type ConfigFirewallBound struct {
	Port      interface{} `json:"port,omitempty"`
	Proto     string      `json:"proto,omitempty"`
	Host      string      `json:"host,omitempty,omitempty"`
	Groups    []string    `json:"groups,omitempty,omitempty"`
	Group     string      `json:"group,omitempty,omitempty"`
	LocalCidr string      `json:"local_cidr,omitempty,omitempty"`
}

func NewCtrl(cfg Config) (*nebula.Control, error) {
	l := logrus.New()
	l.Out = os.Stdout
	marshal, _ := yaml.Marshal(cfg)
	fmt.Println(string(marshal))
	c := config.NewC(l)
	err := c.LoadString(string(marshal))
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %s", err)
	}
	return nebula.Main(c, false, "v0.0.1", l, nil)
}
