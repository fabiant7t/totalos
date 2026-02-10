package kernel

import "fmt"

// IPOptionStaticV4 is used to build the ip= option of the kernel commandline
type IPOptionStaticV4 struct {
	ClientIP  string
	GatewayIP string
	Netmask   string
	Hostname  string
	Device    string
	DNS0IP    string
	DNS1IP    string
	NTP0IP    string
}

// String returns the ip= option value
func (c *IPOptionStaticV4) String() string {
	return fmt.Sprintf(
		"%s::%s:%s:%s:%s:off:%s:%s:%s",
		c.ClientIP,
		c.GatewayIP,
		c.Netmask,
		c.Hostname,
		c.Device,
		c.DNS0IP,
		c.DNS1IP,
		c.NTP0IP,
	)
}
