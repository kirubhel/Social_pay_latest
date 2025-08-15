package gin

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sqlc-dev/pqtype"

	ipWhitelistUsecase "github.com/socialpay/socialpay/src/pkg/ip_whitelist/usecase"
)

type IPCheckerMiddleware struct {
	ipWhitelistUsecase ipWhitelistUsecase.IPWhitelistUseCase
}

func NewIPCheckerMiddleware(ipWhitelistUsecase ipWhitelistUsecase.IPWhitelistUseCase) *IPCheckerMiddleware {
	return &IPCheckerMiddleware{
		ipWhitelistUsecase: ipWhitelistUsecase,
	}
}

// ConvertToIPv4CIDR converts IPv6 or IPv4 CIDR to IPv4 CIDR if possible.
// If the IP is IPv6 and not convertible, it returns the original.
func ConvertToIPv4CIDR(cidr string) (string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	// Check if the IP can be converted to IPv4
	if v4 := ip.To4(); v4 != nil {
		// If we can convert, recalculate the mask length for IPv4
		ones, _ := ipNet.Mask.Size()

		// For IPv4-mapped IPv6 (::ffff:x.x.x.x) or loopback ::1
		return fmt.Sprintf("%s/%d", v4.String(), ones-96), nil
	}

	// If it's a pure IPv6 address, return as-is
	return cidr, nil
}

// IPCheckerMiddleware creates a middleware that validates if the client ip is whitelisted or not
func (i *IPCheckerMiddleware) CheckIP() gin.HandlerFunc {
	fmt.Println("[IP Checker] Initializing IP Checker Middleware")
	return func(c *gin.Context) {
		// Get merchant ID from header
		merchantID, exists := GetMerchantIDFromContext(c)
		if !exists {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INVALID_REQUEST",
					Message: "Failed to get merchant id from context",
				},
			})
		}

		clientIP := c.ClientIP()

		// Parse the IP
		ip := net.ParseIP(clientIP)
		if ip == nil {
			c.String(400, "invalid ip")
			return
		}

		// Convert IPv4-mapped IPv6 (::ffff:192.168.1.10) to pure IPv4
		if ip.To4() != nil {
			ip = ip.To4()
		}

		convertedIpv4, err := ConvertToIPv4CIDR(ip.String() + "/32")
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INTERNAL_SERVER",
					Message: "Failed to convert ip to ipv4",
				},
			})
			c.Abort()
			return
		}
		fmt.Println("parsed client id -> ", convertedIpv4)

		// Parse the IP string to ensure it's valid
		_, ipNet, err := net.ParseCIDR(convertedIpv4)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INVALID_REQUEST",
					Message: "Invalid client IP",
				},
			})
			c.Abort()
			return
		}

		// Convert net.IPNet to pqtype.CIDR
		cidr := pqtype.CIDR{
			IPNet: *ipNet,
			Valid: true,
		}

		ipWhitelists, err := i.ipWhitelistUsecase.GetIPWhitelist(c, merchantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INTERNAL_SERVER",
					Message: "Failed to fetch whitelisted IPs from merchant id",
				},
			})
			c.Abort()
			return
		}

		for _, ipWhitelist := range ipWhitelists {
			if ipWhitelist.IpAddress == cidr.IPNet.String() && ipWhitelist.IsActive {
				c.Next()
			}
		}

		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER",
				Message: "IP address not whitelisted",
			},
		})
		c.Abort()
		return
	}
}
