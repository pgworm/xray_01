package logstash

import (
	"errors"
	"github.com/mono83/udpwriter"
	"github.com/mono83/xray"
	"github.com/mono83/xray/out"
	"github.com/mono83/xray/text"
	"net"
	"time"
)

// Config holds information for filtered receiver
// Tags for JSON, YAML and TOML are configured
type Config struct {
	Address       string   `json:"address" yaml:"address" toml:"address"`
	MinLevel      string   `json:"level" yaml:"level" toml:"level"`
	Buffer        int      `json:"bufferMillis" yaml:"bufferMillis" toml:"bufferMillis"`
	ArgsWhiteList []string `json:"argsWhiteList" yaml:"argsWhiteList" toml:"argsWhiteList"`
	ArgsBlackList []string `json:"argsBlackList" yaml:"argsBlackList" toml:"argsBlackList"`
}

// Validate validates configuration contents
func (c Config) Validate() error {
	if len(c.Address) == 0 {
		return errors.New("empty Logstash binding address")
	}
	if c.Buffer < 100 {
		return errors.New("at least 100ms buffering must be configured")
	}

	return nil
}

// Build builds asynchronous buffered logstash receiver with log-level filtering
func (c Config) Build() (xray.Handler, error) {
	hld, err := c.Raw()
	if err != nil {
		return nil, err
	}

	level := text.ParseLevel(c.MinLevel)

	return out.Filter(
		out.Buffer(
			hld,
			time.Duration(c.Buffer)*time.Millisecond,
		),
		func(event xray.Event) bool {
			if event == nil {
				return false
			}
			l, ok := event.(xray.LogEvent)
			return ok && l.GetLevel() >= level
		},
	), nil
}

// Raw methods builds synchronous sender for logstash.
// It is not recommended to use this method, call Build() instead
func (c Config) Raw() (xray.Handler, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	netAddr, err := net.ResolveUDPAddr("udp", c.Address)
	if err != nil {
		return nil, err
	}

	// Building sender
	send := &sender{
		target:    udpwriter.New(netAddr),
		argFilter: xray.ArgFilterDoubleList(c.ArgsWhiteList, c.ArgsBlackList),
	}

	return send.handle, nil
}
