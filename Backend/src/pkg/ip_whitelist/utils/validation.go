package utils

import (
	"fmt"
	"net"
)

// Validate validates the IP whitelist request
func ValidateIPAddress(ipAddress string) error {
	if _, _, err := net.ParseCIDR(ipAddress); err != nil {
		return fmt.Errorf("must be a valid CIDR (e.g. 192.168.1.1/24)")
	}
	return nil
}
