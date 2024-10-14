package server

import "fmt"

type GigaByte int
type Mbps int
type MHz int

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
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

type CPU struct {
	Name        string `json:"name"`
	Cores       int    `json:"cores"`
	CoreFreqMin MHz    `json:"core_freq_min_mhz"`
	CoreFreqMax MHz    `json:"core_freq_max_mhz"`
	Threads     int    `json:"threads"`
}

type Ethernet struct {
	Device string `json:"device"`
	MAC    string `json:"mac"`
	Speed  Mbps   `json:"speed_mbps"`
}

type Memory struct {
	Size    GigaByte `json:"size_gb"`
	Modules []string `json:"modules"`
}

type Machine struct {
	Arch        string   `json:"arch"`
	IPv4Network Network  `json:"ipv4_network"`
	Hostname    string   `json:"hostname"`
	Disks       []Disk   `json:"disks"`
	CPU         CPU      `json:"cpu"`
	Memory      Memory   `json:"memory"`
	UUID        string   `json:"uuid"`
	Ethernet    Ethernet `json:"ethernet"`
}
