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
	"golang.org/x/sync/errgroup"

	"github.com/fabiant7t/totalos/internal/totalos"
	"github.com/fabiant7t/totalos/internal/totalos/command"
)

var (
	version = "dev" // default version, redacted when building
)

type Installation struct {
	Image      string       `json:"image"`
	SystemDisk totalos.Disk `json:"system_disk"`
	Rebooting  bool         `json:"rebooting"`
}

type Report struct {
	Installation Installation    `json:"installation"`
	Machine      totalos.Machine `json:"machine"`
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
	var m totalos.Machine
	var g errgroup.Group
	g.SetLimit(5)
	g.Go(func() error {
		arch, err := command.Arch(srv, cb)
		m.Arch = arch
		return err
	})
	g.Go(func() error {
		ipv4, err := command.IPv4(srv, cb)
		m.IPv4Network.IP = ipv4.String()
		m.Hostname = fmt.Sprintf("talos-%s", strings.ReplaceAll(ipv4.String(), ".", "-"))
		return err
	})
	g.Go(func() error {
		nm, err := command.IPv4Netmask(srv, cb)
		m.IPv4Network.Netmask = nm.String()
		return err
	})
	g.Go(func() error {
		gw, err := command.IPv4Gateway(srv, cb)
		m.IPv4Network.Gateway = gw.String()
		return err
	})
	g.Go(func() error {
		mac, err := command.MAC(srv, cb)
		m.MAC = mac
		return err
	})
	g.Go(func() error {
		uuid, err := command.SystemUUID(srv, cb)
		m.UUID = uuid
		return err
	})
	g.Go(func() error {
		cpuName, err := command.CPUName(srv, cb)
		m.CPU.Name = cpuName
		return err
	})
	g.Go(func() error {
		cpuCores, err := command.CPUCores(srv, cb)
		m.CPU.Cores = cpuCores
		return err
	})
	g.Go(func() error {
		cpuThreads, err := command.CPUThreads(srv, cb)
		m.CPU.Threads = cpuThreads
		return err
	})
	g.Go(func() error {
		memory, err := command.Memory(srv, cb)
		m.Memory = memory
		return err
	})
	g.Go(func() error {
		disks, err := command.Disks(srv, cb)
		m.Disks = disks
		return err
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	// Installation
	installation := Installation{
		Image:     args.Image,
		Rebooting: args.Reboot,
	}
	if installation.Image == "" {
		url, err := totalos.LatestImageURL(ctx, m.Arch, client)
		if err != nil {
			log.Fatal(err)
		}
		installation.Image = url
	}
	if err := command.SoftwareRAIDNotExists(srv, cb); err != nil {
		log.Fatal(err)
	}
	if err := command.WipeFileSystemSignatures(srv, cb); err != nil {
		log.Fatal(err)
	}
	disk, err := command.SelectSystemDisk(m.Disks)
	if err != nil {
		log.Fatal(err)
	}
	installation.SystemDisk = disk
	if err := command.InstallImage(srv, installation.Image, installation.SystemDisk.Device(), cb); err != nil {
		log.Fatal(err)
	}

	// Report
	report := Report{
		Installation: installation,
		Machine:      m,
	}
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	if args.Webhook != "" {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			args.Webhook,
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
