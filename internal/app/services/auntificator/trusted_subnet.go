package auntificator

import (
	"errors"
	"net"

	"github.com/northmule/shorturl/config"
)

// CheckTrustedSubnet проверка запроса на принадлежность к сети
type CheckTrustedSubnet struct {
	configApp *config.Config
}

// NewTrustedSubnet конструктор
func NewTrustedSubnet(configApp *config.Config) *CheckTrustedSubnet {
	return &CheckTrustedSubnet{
		configApp: configApp,
	}
}

// GrantAccess предоставить доступ
func (c *CheckTrustedSubnet) GrantAccess(ip string) error {

	var err error
	if c.configApp.TrustedSubnet == "" {
		return errors.New("trusted network is not set, access is limited")
	}

	var expectedIP net.IP
	var actualIP net.IP
	var expectedNet *net.IPNet

	actualIP = net.ParseIP(ip)
	if actualIP == nil {
		return errors.New("no IP address has been transmitted, access is restricted")
	}

	expectedIP, expectedNet, err = net.ParseCIDR(c.configApp.TrustedSubnet)
	if err != nil {
		return errors.New("the configuration address is not recognized, access is limited")
	}

	if ok := expectedIP.Equal(actualIP); ok {
		return nil
	}

	if ok := expectedNet.Contains(actualIP); !ok {
		return errors.New("the address is not allowed, access is limited")
	}

	return nil

}
