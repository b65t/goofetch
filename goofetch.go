package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCommand(cmdName string, args ...string) string {
	cmd := exec.Command(cmdName, args...)

	output, err := cmd.Output()
	if err != nil {
		return "Error: " + err.Error()
	}

	return strings.TrimSpace(string(output))
}

func getDistro() string {
	output := runCommand("cat", "/etc/os-release")
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.Split(line, "=")[1], "\"")
		}
	}
	return "Unknown Distro"
}

func getPackageManager(distro string) (string, string) {
	switch distro {
	case "ubuntu", "debian":
		return "dpkg", "dpkg -l | wc -l"
	case "fedora", "centos":
		return "rpm", "rpm -qa | wc -l"
	case "arch", "artix", "cachyos":
		return "pacman", "pacman -Q | wc -l"
	default:
		return "unknown", ""
	}
}

func countPackages(command string) string {
	if command == "" {
		return "Package manager not supported."
	}
	return runCommand("sh", "-c", command)
}

var ansiColors = []int{
	1, 2, 3, 4, 5, 6, 7,
}

func printTerminalColors() {
	fmt.Printf(" ")

	for _, color := range ansiColors {
		fmt.Printf("\033[48;5;%dm  \033[0m", color)
	}
	fmt.Println()
}

func getOSLogo(distro string) string {
	filename := "logo/" + distro + ".txt"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = "logo/unknown.txt"
	}
	return readLogoFile(filename)
}

func readLogoFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "Error reading ASCII file"
	}
	return string(data)
}

func printTerminalColorsInline() string {
	colorStf := "\033[1;34m\033[0m "
	for _, color := range ansiColors {
		colorStf += fmt.Sprintf("\033[48;5;%dm  \033[0m", color)
	}
	return colorStf
}

func main() {
	hostname := runCommand("hostname")
	username := runCommand("whoami")
	kernel := runCommand("uname", "-r")
	distro := getDistro()
	uptime := runCommand("uptime", "-p")
	shell := runCommand("basename", os.Getenv("SHELL"))
	//	cpu := runCommand("sh", "-c", "lscpu | awk -F': ' '/Model name/ {print $2}' | sed 's/^[ \t]*//'")
	//	memory := runCommand("sh", "-c", `free -m | awk 'NR==2{printf "%.2f GiB / %.2f GiB (%.0f%%)\n", $3/1024, $2/1024, $3/$2*100}'`)
	wm := runCommand("sh", "-c", `ps -e | grep -oE 'i3|kwin|mutter|openbox|awesome|fluxbox|xmonad|sway|bspwm|qtile|dwm|hyprland' | head -n 1`)

	packageManager, packageCommand := getPackageManager(distro)
	packageCount := countPackages(packageCommand)

	asciiArt := strings.Split(getOSLogo(distro), "\n")

	/* 	fmt.Printf(" %s@%s\n", username, hostname)
	   	fmt.Println("", distro)
	   	fmt.Println("󰌽", kernel)
	   	fmt.Println("", wm)
	   	fmt.Printf(" %s\n", shell)
	   	fmt.Printf(" %s (%s)\n", packageCount, packageManager)
	   	fmt.Println("󰅶", uptime)
	   	fmt.Println("", cpu)
	   	fmt.Println("", memory) */

	sysInfo := []string{
		fmt.Sprintf("\033[1;34m\033[0m %s@%s", username, hostname),
		fmt.Sprintf("\033[1;34m\033[0m %s", distro),
		fmt.Sprintf("\033[1;34m󰌽\033[0m %s", kernel),
		fmt.Sprintf("\033[1;34m\033[0m %s", wm),
		fmt.Sprintf("\033[1;34m\033[0m %s", shell),
		fmt.Sprintf("\033[1;34m\033[0m %s (%s)", packageCount, packageManager),
		fmt.Sprintf("\033[1;34m󰅶\033[0m %s", uptime),
		printTerminalColorsInline(),
		//		fmt.Sprintf(" %s", cpu),
		//		fmt.Sprintf(" %s", memory),
	}

	maxLines := len(asciiArt)
	if len(sysInfo) > maxLines {
		maxLines = len(sysInfo)
	}

	for i := 0; i < len(asciiArt) || i < len(sysInfo); i++ {
		var asciiPart, sysPart string

		if i < len(asciiArt) {
			asciiPart = fmt.Sprintf("\033[1;34m%-18s\033[0m", asciiArt[i])
		}

		if i < len(sysInfo) {
			sysPart = sysInfo[i]
		}

		fmt.Printf("%-18s %s\n", asciiPart, sysPart)
	}

}
