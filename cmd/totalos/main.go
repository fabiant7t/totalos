package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fabiant7t/totalos/internal/totalos"
	"github.com/fabiant7t/totalos/internal/totalos/services"
)

func main() {
	ip := flag.String("ip", "", "IP of the server")
	port := flag.Uint("port", 22, "SSH port of the server")
	user := flag.String("user", "root", "name of the user")
	password := flag.String("password", "", "password of the user (optional)")
	keyPath := flag.String("key", "", "path to the private key (optional)")
	image := flag.String("image", "", "URL to ISO image (optional)")

	flag.Parse()

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
	services.Reboot(srv, nil)
	fmt.Printf("Installed %s on disk %s!\n", *image, disk)
	fmt.Printf("Your server %s is now rebooting to Talos maintenance mode!\n", *ip)
}
