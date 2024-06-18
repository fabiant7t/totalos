package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fabiant7t/totalos/internal/totalos"
	"github.com/fabiant7t/totalos/internal/totalos/services"
)

var (
	version = "dev" // default version, redacted when building
)

type Report struct {
	Image              string `json:"image"`
	IP                 string `json:"ip"`
	Netmask            string `json:"netmask"`
	Gateway            string `json:"gateway"`
	IPKernelParameters string `json:"ip_kernel_parameters"`
	Disk               string `json:"disk"`
	Hostname           string `json:"hostname"`
	Status             string `json:"status"`
}

func main() {
	ip := flag.String("ip", "", "IP of the server")
	port := flag.Uint("port", 22, "SSH port of the server")
	user := flag.String("user", "root", "name of the user")
	password := flag.String("password", "", "password of the user (optional)")
	keyPath := flag.String("key", "", "path to the private key (optional)")
	image := flag.String("image", "", "URL to ISO image (optional)")
	versionFlag := flag.Bool("version", false, "prints the version")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\n", version)
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

	srv := totalos.NewServer(*ip, *user, uint16(*port), *password, nil)
	if *keyPath != "" {
		if err := srv.SetKeyFromFile(*keyPath); err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()

	if *image == "" {
		arch, err := services.Arch(srv, nil)
		if err != nil {
			log.Fatal(err)
		}
		url, err := totalos.LatestImageURL(ctx, arch)
		if err != nil {
			log.Fatal(err)
		}
		*image = url
	}
	ipv4, err := services.IPv4(srv, nil)
	if err != nil {
		log.Fatal(err)
	}
	ipv4nm, err := services.IPv4Netmask(srv, nil)
	if err != nil {
		log.Fatal(err)
	}
	ipv4gw, err := services.IPv4Gateway(srv, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := services.SoftwareRAIDNotExists(srv, nil); err != nil {
		log.Fatal(err)
	}
	if err := services.WipeFileSystemSignatures(srv, nil); err != nil {
		log.Fatal(err)
	}
	disk, err := services.NominateInstallDisk(srv, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := services.InstallImage(srv, *image, disk, nil); err != nil {
		log.Fatal(err)
	}
	hostname := fmt.Sprintf("talos-%s", strings.ReplaceAll(ipv4.String(), ".", "-"))
	report := Report{
		Image:              *image,
		IP:                 ipv4.String(),
		Netmask:            ipv4nm.String(),
		Gateway:            ipv4gw.String(),
		IPKernelParameters: fmt.Sprintf("ip=%s::%s:%s:%s:eth0:off:8.8.8.8:8.8.4.4:216.239.35.8", ipv4, ipv4gw, ipv4nm, hostname),
		Disk:               disk,
		Hostname:           hostname,
		Status:             "Rebooting to Talos Linux in maintenance mode",
	}
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	services.Reboot(srv, nil)
}
