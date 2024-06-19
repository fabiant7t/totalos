package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/fabiant7t/totalos/internal/totalos"
	"github.com/fabiant7t/totalos/internal/totalos/command"
)

var (
	version = "dev" // default version, redacted when building
)

type Disk struct {
	Device       string `json:"device"`
	SerialNumber string `json:"serial_number"`
}

type Network struct {
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

type CPU struct {
	Name    string `json:"name"`
	Cores   int    `json:"cores"`
	Threads int    `json:"threads"`
}

type Machine struct {
	Arch        string             `json:"arch"`
	IPv4Network Network            `json:"ipv4_network"`
	Hostname    string             `json:"hostname"`
	Storage     []totalos.GigaByte `json:"storage_gb"`
	CPU         CPU                `json:"cpu"`
	Memory      totalos.GigaByte   `json:"memory_gb"`
	MAC         string             `json:"mac"`
	UUID        string             `json:"uuid"`
}

type Installation struct {
	Image     string `json:"image"`
	Disk      Disk   `json:"disk"`
	Rebooting bool   `json:"rebooting"`
}

type Report struct {
	Installation Installation `json:"installation"`
	Machine      Machine      `json:"machine"`
}

type CallArgs struct {
	IP       string
	Port     uint16
	User     string
	Password string
	KeyPath  string
	Image    string
	Webhook  string
	Reboot   bool
}

func main() {
	// Parse and validate arguments and populate CallArgs. Might exit early (--version).
	ip := flag.String("ip", "", "IP of the server")
	port := flag.Uint("port", 22, "SSH port of the server")
	user := flag.String("user", "root", "name of the user")
	password := flag.String("password", "", "password of the user (optional)")
	keyPath := flag.String("key", "", "path to the private key (optional)")
	image := flag.String("image", "", "URL to ISO image (optional)")
	webhook := flag.String(
		"webhook",
		"",
		"Endpoint that should receive the report through HTTP POST (optional)",
	)
	versionFlag := flag.Bool("version", false, "prints the version")
	rebootFlag := flag.Bool("reboot", false, "reboot the server")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("totalos v%s\n", version)
		os.Exit(0)
	}
	if *ip == "" {
		fmt.Println("Error: --ip flag is required")
		flag.Usage()
		os.Exit(1)
	}
	if *password == "" && *keyPath == "" {
		fmt.Println("Error: --password or --key required")
		flag.Usage()
		os.Exit(1)
	}
	args := CallArgs{
		IP:       *ip,
		Port:     uint16(*port),
		User:     *user,
		Password: *password,
		KeyPath:  *keyPath,
		Image:    *image,
		Webhook:  *webhook,
		Reboot:   *rebootFlag,
	}

	// Context with deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	// HTTP client
	client := &http.Client{}
	// SSH host key callback, key is not being checked.
	cb := ssh.InsecureIgnoreHostKey()
	// Define target server
	srv := totalos.NewServer(args.IP, args.User, args.Port, args.Password, nil)
	if args.KeyPath != "" {
		if err := srv.SetKeyFromFile(args.KeyPath); err != nil {
			log.Fatal(err)
		}
	}

	// Machine
	arch, err := command.Arch(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	ipv4, err := command.IPv4(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	ipv4nm, err := command.IPv4Netmask(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	ipv4gw, err := command.IPv4Gateway(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	mac, err := command.MAC(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	uuid, err := command.SystemUUID(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	cpuName, err := command.CPUName(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	cpuCores, err := command.CPUCores(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	cpuThreads, err := command.CPUThreads(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	memory, err := command.Memory(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	hostname := fmt.Sprintf("talos-%s", strings.ReplaceAll(ipv4.String(), ".", "-"))
	storage, err := command.Storage(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	machine := Machine{
		Arch: arch,
		IPv4Network: Network{
			IP:      ipv4.String(),
			Netmask: ipv4nm.String(),
			Gateway: ipv4gw.String(),
		},
		Hostname: hostname,
		CPU: CPU{
			Name:    cpuName,
			Cores:   cpuCores,
			Threads: cpuThreads,
		},
		Memory:  memory,
		Storage: storage,
		MAC:     mac,
		UUID:    uuid,
	}

	// Installation
	if args.Image == "" {
		url, err := totalos.LatestImageURL(ctx, arch, client)
		if err != nil {
			log.Fatal(err)
		}
		args.Image = url
	}
	if err := command.SoftwareRAIDNotExists(srv, cb); err != nil {
		log.Fatal(err)
	}
	if err := command.WipeFileSystemSignatures(srv, cb); err != nil {
		log.Fatal(err)
	}
	device, sn, err := command.NominateInstallDisk(srv, cb)
	if err != nil {
		log.Fatal(err)
	}
	if err := command.InstallImage(srv, args.Image, device, cb); err != nil {
		log.Fatal(err)
	}
	installation := Installation{
		Image: *image,
		Disk: Disk{
			Device:       device,
			SerialNumber: sn,
		},
		Rebooting: *rebootFlag,
	}

	// Report
	report := Report{
		Installation: installation,
		Machine:      machine,
	}
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	if *webhook != "" {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			*webhook,
			bytes.NewReader(jsonData),
		)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", fmt.Sprintf("totalos/%s", version))
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
	}

	if args.Reboot {
		command.Reboot(srv, cb)
	}
}
