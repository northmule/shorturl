package auntificator

import (
	"testing"

	"github.com/northmule/shorturl/config"
	"github.com/stretchr/testify/assert"
)

func TestGrantAccess_NoTrustedSubnet(t *testing.T) {
	configApp := &config.Config{TrustedSubnet: ""}
	checker := NewTrustedSubnet(configApp)
	err := checker.GrantAccess("192.168.1.1")
	assert.Error(t, err, "trusted network is not set, access is limited")
}

func TestGrantAccess_InvalidIP(t *testing.T) {
	configApp := &config.Config{TrustedSubnet: "192.168.1.0/24"}
	checker := NewTrustedSubnet(configApp)
	err := checker.GrantAccess("invalid_ip")
	assert.Error(t, err, "no IP address has been transmitted, access is restricted")
}

func TestGrantAccess_InvalidCIDR(t *testing.T) {
	configApp := &config.Config{TrustedSubnet: "invalid_cidr"}
	checker := NewTrustedSubnet(configApp)
	err := checker.GrantAccess("192.168.1.1")
	assert.Error(t, err, "the configuration address is not recognized, access is limited")
}

func TestGrantAccess_IPMatchesExactly(t *testing.T) {
	configApp := &config.Config{TrustedSubnet: "192.168.1.1/32"}
	checker := NewTrustedSubnet(configApp)
	err := checker.GrantAccess("192.168.1.1")
	assert.NoError(t, err)
}

func TestGrantAccess_IPInSubnet(t *testing.T) {
	configApp := &config.Config{TrustedSubnet: "192.168.1.0/24"}
	checker := NewTrustedSubnet(configApp)
	err := checker.GrantAccess("192.168.1.100")
	assert.NoError(t, err)
}

func TestGrantAccess_IPNotInSubnet(t *testing.T) {
	configApp := &config.Config{TrustedSubnet: "192.168.1.0/24"}
	checker := NewTrustedSubnet(configApp)
	err := checker.GrantAccess("192.168.2.100")
	assert.Error(t, err, "the address is not allowed, access is limited")
}
