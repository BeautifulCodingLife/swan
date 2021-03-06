package config

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

type AgentConfig struct {
	Listen    string   `json:"listen"` // only for ping -> pong service
	LogLevel  string   `json:"logLevel"`
	JoinAddrs []string `json:"joinAddrs"`
	DNS       *DNS     `json:"dns"`
	Janitor   *Janitor `json:"janitor"`
}

type DNS struct {
	Domain          string        `json:"domain"`
	RecurseOn       bool          `json:"recurseOn"`
	ListenAddr      string        `json:"listenAddr"`
	TTL             int           `json:"ttl"`
	Resolvers       []string      `json:"resolvers"`
	ExchangeTimeout time.Duration `json:"exchangeTimeout"`

	SOARname   string `json:"soarname"`
	SOAMname   string `json:"soamname"`
	SOASerial  uint32 `json:"soaserial"`
	SOARefresh uint32 `json:"soarefresh"`
	SOARetry   uint32 `json:"soaretry"`
	SOAExpire  uint32 `json:"soaexpire"`
}

type Janitor struct {
	ListenAddr    string `json:"listenAddr"`
	TLSListenAddr string `json:"tlsListenAddr"`
	TLSCertFile   string `json:"tlsCertFile"`
	TLSKeyFile    string `json:"tlsKeyFile"`
	Domain        string `json:"domain"`
	AdvertiseIP   string `json:"advertiseIP"`
}

func NewAgentConfig(c *cli.Context) (*AgentConfig, error) {
	cfg := &AgentConfig{
		Listen:    "0.0.0.0:9999",
		LogLevel:  "info",
		JoinAddrs: []string{"0.0.0.0:9999"},
		DNS: &DNS{
			Domain:          "swan.com",
			ListenAddr:      "0.0.0.0:53",
			RecurseOn:       true,
			TTL:             0,
			Resolvers:       []string{"114.114.114.114"},
			ExchangeTimeout: time.Second * 3,
		},
		Janitor: &Janitor{
			ListenAddr: "0.0.0.0:80",
			Domain:     "swan.com",
		},
	}

	if c.String("listen") != "" {
		cfg.Listen = c.String("listen")
	}

	if c.String("log-level") != "" {
		cfg.LogLevel = c.String("log-level")
	}

	if c.String("join-addrs") != "" {
		cfg.JoinAddrs = strings.Split(c.String("join-addrs"), ",")
	}

	if c.String("domain") != "" {
		cfg.DNS.Domain = c.String("domain")
		cfg.Janitor.Domain = c.String("domain")
	}

	if c.String("gateway-advertise-ip") != "" {
		cfg.Janitor.AdvertiseIP = c.String("gateway-advertise-ip")
	}

	if c.String("gateway-listen-addr") != "" {
		cfg.Janitor.ListenAddr = c.String("gateway-listen-addr")
		if cfg.Janitor.AdvertiseIP == "" {
			cfg.Janitor.AdvertiseIP, _, _ = net.SplitHostPort(cfg.Janitor.ListenAddr)
		}
	}

	if c.String("gateway-tls-listen-addr") != "" {
		cfg.Janitor.TLSListenAddr = c.String("gateway-tls-listen-addr")
	}

	if c.String("gateway-tls-cert-file") != "" {
		cfg.Janitor.TLSCertFile = c.String("gateway-tls-cert-file")
	}

	if c.String("gateway-tls-key-file") != "" {
		cfg.Janitor.TLSKeyFile = c.String("gateway-tls-key-file")
	}

	if c.String("dns-listen-addr") != "" {
		cfg.DNS.ListenAddr = c.String("dns-listen-addr")
	}

	if c.String("dns-resolvers") != "" {
		cfg.DNS.Resolvers = strings.Split(c.String("dns-resolvers"), ",")
	}

	if ttl := c.Int("dns-ttl"); ttl > 0 {
		cfg.DNS.TTL = ttl
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *AgentConfig) validate() error {
	// verify Janitor.AdvertiseIP is valid ip addr
	if ip := c.Janitor.AdvertiseIP; ip != "" {
		if net.ParseIP(ip) == nil {
			return fmt.Errorf("invalid janitor advertise ip: %v", ip)
		}
	}

	// verify Janitor.TLS cert/key files exist if gateway tls enabled
	if c.Janitor.TLSListenAddr != "" {
		if _, err := os.Stat(c.Janitor.TLSCertFile); err != nil {
			return fmt.Errorf("tsl cert file: %v", err)
		}
		if _, err := os.Stat(c.Janitor.TLSKeyFile); err != nil {
			return fmt.Errorf("tsl key file: %v", err)
		}
	}

	return nil
}
