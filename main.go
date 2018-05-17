package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

func enableProxySettings() {
	fmt.Println("Configuring proxy in Network settings")
	cmd := exec.Command("networksetup", "-setsocksfirewallproxy", "Wi-Fi", "127.0.0.1", "3333")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Setting proxy settings failed with %s\n", err)
	}
	fmt.Println("Proxy now configured")
}

func disableProxySettings() {
	fmt.Println("\nRemoving proxy in Network settings")
	cmd := exec.Command("networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "off")
	println("Done. ✨")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Removing proxy settings failed with %s\n", err)
	}
}

func openSSHTunnel() {
	fmt.Println("Opening SOCKs tunnel on port: 3333")
	cmd := exec.Command("ssh", "-N", "socks.vpn")
	cmd.CombinedOutput()
}

func startVPN() {
	enableProxySettings()
	openSSHTunnel()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			disableProxySettings()
		}
	}()

	startVPN()
}
