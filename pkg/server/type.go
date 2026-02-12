package server

import "fmt"

type (
	GigaByte int
	Mbps     int
	MHz      int
)

type Disk struct {
	Model     string `json:"model"`
	Name      string `json:"name"`
	Serial    string `json:"serial"`
	Size      int    `json:"size"`
	Type      string `json:"type"`
	Transport string `json:"tran"`
	WWN       string `json:"wwn"`
}

func (d *Disk) Device() string {
	return fmt.Sprintf("/dev/%s", d.Name)
}

type Network struct {
	IP          string   `json:"ip"`
	Netmask     string   `json:"netmask"`
	CIDR        string   `json:"cidr"`
	Gateway     string   `json:"gateway"`
	ResolversV4 []string `json:"resolvers_v4"`
	ResolversV6 []string `json:"resolvers_v6"`
}

type CPU struct {
	Name        string `json:"name"`
	Cores       int    `json:"cores"`
	CoreFreqMin MHz    `json:"core_freq_min_mhz"`
	CoreFreqMax MHz    `json:"core_freq_max_mhz"`
	Threads     int    `json:"threads"`
}

type Ethernet struct {
	Device     string     `json:"device"`
	MAC        string     `json:"mac"`
	Speed      Mbps       `json:"speed_mbps"`
	IDNetNames IDNetNames `json:"id_net_names"`
}

// IDNetNames define the inputs used by systemd/udev to generate network
// interface names. Despite being called “predictable”, the resulting names
// are selected from a hard-coded priority list: some naming schemes are
// almost always available, while others apply only in specific situations.
// As a result, the chosen name depends on which identifiers are present
// on the system at the time.
// See https://github.com/systemd/systemd/blob/main/src/udev/net/link-config.c#L738
type IDNetNames struct {
	FromDatabase string `json:"from_database"`
	Onboard      string `json:"onboard"`
	Slot         string `json:"slot"`
	Path         string `json:"path"`
	MAC          string `json:"mac"`
}

// InterfaceName returns the interface name that Talos picks when predictable
// naming is not disabled.
// Might return empty string when there are no interface names to choose from,
// which must be checked for.
// See https://docs.siderolabs.com/talos/v1.12/networking/predictable-interface-names
func (idnn *IDNetNames) InterfaceName() string {
	for _, name := range []string{idnn.Onboard, idnn.Slot, idnn.Path, idnn.MAC} {
		if name != "" {
			return name
		}
	}
	return ""
}

type Memory struct {
	Size    GigaByte `json:"size_gb"`
	Modules []string `json:"modules"`
}

type System struct {
	Manufacturer string `json:"manufacturer"`
	ProductName  string `json:"product_name"`
	Version      string `json:"version"`
	Family       string `json:"family"`
	UUID         string `json:"uuid"`
	SerialNumber string `json:"serial_number"`
	SKUNumber    string `json:"sku_number"`
}

type Machine struct {
	Arch        string   `json:"arch"`
	IPv4Network Network  `json:"ipv4_network"`
	Hostname    string   `json:"hostname"`
	Disks       []Disk   `json:"disks"`
	CPU         CPU      `json:"cpu"`
	Memory      Memory   `json:"memory"`
	System      System   `json:"system"`
	Ethernet    Ethernet `json:"ethernet"`
}
