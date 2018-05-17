package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strings"
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

func help() {
	help := `usage: socks

You must have an a host named "socks.vpn" in your
ssh config with a dynamic forward of 3333:

# ~/.ssh/config
Host socks.vpn
  Hostname my.ip.address
  DynamicForward 3333
`

	println(help)
	os.Exit(0)
}

func missingSSHHost() bool {
	usr, err := user.Current()
	if err != nil {
		return true
	}

	config := fmt.Sprintf("%s/.ssh/config", usr.HomeDir)

	b, err := ioutil.ReadFile(config)
	if err != nil {
		return true
	}

	return !strings.Contains(string(b), "socks.vpn")
}

func main() {

	if len(os.Args) > 1 || missingSSHHost() {
		help()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			disableProxySettings()
		}
	}()

	startVPN()
}
