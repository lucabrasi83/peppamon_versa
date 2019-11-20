// Package initializer contains the environmental data to load before starting the Peppamon Collector
package initializer

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"

	"github.com/lucabrasi83/peppamon_versa/logging"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var (
	Commit  string
	Version string
	BuiltAt string
	BuiltOn string
)

func Initialize() {
	printBanner()
	printReleaseDetails()
	printPlatformDetails()
}

func printBanner() {

	bannerString := `
 ____   ___  ____  ____   ____  ___ ___   ___   ____       __ __    ___  ____    _____  ____ 
|    \ /  _]|    \|    \ /    ||   |   | /   \ |    \     |  |  |  /  _]|    \  / ___/ /    |
|  o  )  [_ |  o  )  o  )  o  || _   _ ||     ||  _  |    |  |  | /  [_ |  D  )(   \_ |  o  |
|   _/    _]|   _/|   _/|     ||  \_/  ||  O  ||  |  |    |  |  ||    _]|    /  \__  ||     |
|  | |   [_ |  |  |  |  |  _  ||   |   ||     ||  |  |    |  :  ||   [_ |    \  /  \ ||  _  |
|  | |     ||  |  |  |  |  |  ||   |   ||     ||  |  |     \   / |     ||  .  \ \    ||  |  |
|__| |_____||__|  |__|  |__|__||___|___| \___/ |__|__|      \_/  |_____||__|\_|  \___||__|__|

`
	buf := new(bytes.Buffer)
	buf.WriteString(bannerString)

	_, err := io.Copy(os.Stdout, buf)
	if err != nil {
		logging.PeppaMonLog("error", "Not able to load banner: ", err.Error())
	}
	fmt.Printf("\n\n")
}

// printReleaseDetails is called as part of init() function and display Vulscano release details such as
// Git Commit, Git tag, build date,...
func printReleaseDetails() {
	fmt.Println(logging.UnderlineText("Peppamon Collector Release:"), logging.InfoMessage(Version))
	fmt.Println(logging.UnderlineText("Github Commit:"), logging.InfoMessage(Commit))

	fmt.Println(logging.UnderlineText(
		"Compiled @"), logging.InfoMessage(BuiltAt),
		"on", logging.InfoMessage(BuiltOn))

	fmt.Printf("\n")
}

// printPlatformDetails is called as part of init() function and display local platform details such as
// CPU info, OS & kernel Version, disk usage on partition "/",...
func printPlatformDetails() {

	platform, err := host.Info()

	if err != nil {
		logging.PeppaMonLog("error", "Unable to fetch platform details:", err.Error())
	} else {
		fmt.Println(
			logging.UnderlineText("Hostname:"),
			logging.InfoMessage(platform.Hostname))
		fmt.Println(
			logging.UnderlineText("Operating System:"),
			logging.InfoMessage(platform.OS),
			logging.InfoMessage(platform.PlatformVersion))
		fmt.Println(logging.UnderlineText("Kernel Version:"), logging.InfoMessage(platform.KernelVersion))
	}

	cpuDetails, err := cpu.Info()
	if err != nil {
		logging.PeppaMonLog("error", "Unable to fetch CPU details:", err.Error())
	} else {
		fmt.Println(logging.UnderlineText("CPU Model:"), logging.InfoMessage(cpuDetails[0].ModelName))
		fmt.Println(logging.UnderlineText("CPU Core(s):"), logging.InfoMessage(runtime.NumCPU()))
		fmt.Println(logging.UnderlineText("OS Architecture:"), logging.InfoMessage(runtime.GOARCH))
	}

	diskUsage, err := disk.Usage("/")

	if err != nil {
		logging.PeppaMonLog("error", "Unable to fetch disk Usage details:", err.Error())
	} else {
		diskUsageRounded := strconv.Itoa(int(math.Round(diskUsage.UsedPercent)))

		fmt.Println(
			logging.UnderlineText("Disk Usage Percentage:"), logging.InfoMessage(diskUsageRounded, "%"))
	}

	memUsage, err := mem.VirtualMemory()

	if err != nil {
		logging.PeppaMonLog("error", "Unable to fetch Memory details:", err.Error())
	} else {
		memUsageRounded := strconv.Itoa(int(math.Round(memUsage.UsedPercent)))
		fmt.Println(
			logging.UnderlineText("Virtual Memory Usage:"), logging.InfoMessage(memUsageRounded, "%"))
	}

	fmt.Printf("\n")

}
