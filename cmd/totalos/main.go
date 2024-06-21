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
	FormatStorageDisk bool         `json:"format_storage_disk"`
	Image             string       `json:"image"`
	Rebooting         bool         `json:"rebooting"`
	Config            string       `json:"config"`
	StorageDisk       totalos.Disk `json:"storage_disk"`
	SystemDisk        totalos.Disk `json:"system_disk"`
}

type Report struct {
	Installation Installation    `json:"installation"`
	Machine      totalos.Machine `json:"machine"`
}

type CallArgs struct {
	IP                string
	Port              uint16
	User              string
	Password          string
	KeyPath           string
	Image             string
	Webhook           string
	Config            string
	Reboot            bool
	FormatStorageDisk bool
}

func NewCallArgs() *CallArgs {
	ip := flag.String("ip", "", "IP of the server")
	port := flag.Uint("port", 22, "SSH port of the server")
	user := flag.String("user", "root", "name of the user")
	password := flag.String("password", "", "password of the user (optional)")
	keyPath := flag.String("key", "", "path to the private key (optional)")
	image := flag.String("image", "", "URL to raw.xz image (optional)")
	webhook := flag.String(
		"webhook",
		"",
		"Endpoint that should receive the report through HTTP POST (optional)",
	)
	config := flag.String("config", "", "URL at which the machine configuration data may be found (optional)")
	versionFlag := flag.Bool("version", false, "prints the version")
	rebootFlag := flag.Bool("reboot", false, "reboot the server")
	formatStorageDiskFlag := flag.Bool("format-storage-disk", false, "format storage disk (optional)")

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
	return &CallArgs{
		IP:                *ip,
		Port:              uint16(*port),
		User:              *user,
		Password:          *password,
		KeyPath:           *keyPath,
		Image:             *image,
		Webhook:           *webhook,
		Config:            *config,
		Reboot:            *rebootFlag,
		FormatStorageDisk: *formatStorageDiskFlag,
	}
}

func main() {
	// Parse and validate arguments and populate CallArgs. Might exit early (--version).
	args := NewCallArgs()

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
	var mach totalos.Machine
	var g errgroup.Group
	g.SetLimit(5)
	g.Go(func() error {
		arch, err := command.Arch(srv, cb)
		mach.Arch = arch
		return err
	})
	g.Go(func() error {
		ipv4, err := command.IPv4(srv, cb)
		mach.IPv4Network.IP = ipv4.String()
		mach.Hostname = fmt.Sprintf("talos-%s", strings.ReplaceAll(ipv4.String(), ".", "-"))
		return err
	})
	g.Go(func() error {
		nm, err := command.IPv4Netmask(srv, cb)
		mach.IPv4Network.Netmask = nm.String()
		return err
	})
	g.Go(func() error {
		gw, err := command.IPv4Gateway(srv, cb)
		mach.IPv4Network.Gateway = gw.String()
		return err
	})
	g.Go(func() error {
		mac, err := command.MAC(srv, cb)
		mach.MAC = mac
		return err
	})
	g.Go(func() error {
		uuid, err := command.SystemUUID(srv, cb)
		mach.UUID = uuid
		return err
	})
	g.Go(func() error {
		cpuName, err := command.CPUName(srv, cb)
		mach.CPU.Name = cpuName
		return err
	})
	g.Go(func() error {
		cpuCores, err := command.CPUCores(srv, cb)
		mach.CPU.Cores = cpuCores
		return err
	})
	g.Go(func() error {
		cpuThreads, err := command.CPUThreads(srv, cb)
		mach.CPU.Threads = cpuThreads
		return err
	})
	g.Go(func() error {
		memory, err := command.Memory(srv, cb)
		mach.Memory = memory
		return err
	})
	g.Go(func() error {
		disks, err := command.Disks(srv, cb)
		mach.Disks = disks
		return err
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	// Installation
	installation := Installation{
		Image:             args.Image,
		Rebooting:         args.Reboot,
		FormatStorageDisk: args.FormatStorageDisk,
		Config:            args.Config,
	}
	// If image is not given, query the latest one
	if installation.Image == "" {
		url, err := totalos.LatestImageURL(ctx, mach.Arch, client)
		if err != nil {
			log.Fatal(err)
		}
		installation.Image = url
	}
	// Reset disks
	if err := command.SoftwareRAIDNotExists(srv, cb); err != nil {
		log.Fatal(err)
	}
	if err := command.WipeFileSystemSignatures(srv, cb); err != nil {
		log.Fatal(err)
	}
	// Select system disk and write image data on it
	systemDisk, err := command.SelectSystemDisk(mach.Disks)
	if err != nil {
		log.Fatal(err)
	}
	installation.SystemDisk = systemDisk
	if err := command.InstallImage(srv, installation.Image, installation.SystemDisk.Device(), cb); err != nil {
		log.Fatal(err)
	}
	// If config is given, set it as talos.config option in grub.cfg
	if args.Config != "" {
		if err := command.SetConfigURL(srv, args.Config, installation.SystemDisk.Device(), cb); err != nil {
			log.Fatal(err)
		}
	}
	// Select storage disk and format it (if requested)
	storageDisk, err := command.SelectStorageDisk(mach.Disks, systemDisk)
	if err != nil {
		log.Fatal(err)
	}
	installation.StorageDisk = storageDisk
	if args.FormatStorageDisk {
		if err := command.FormatExt4(srv, "storage", "61291e61-291e-6129-1e61-291e61291e00", installation.StorageDisk.Device(), cb); err != nil {
			log.Fatal(err)
		}
	}
	// Create report and print it to stdout
	report := Report{
		Installation: installation,
		Machine:      mach,
	}
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
	// Send report to webhook (if given)
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
	// Reboot (if requested)
	if args.Reboot {
		command.Reboot(srv, cb)
	}
}
