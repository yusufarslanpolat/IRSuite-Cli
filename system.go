package main

import (
	"encoding/json"
	"fmt"
	"strings"

	//"fmt"
	//"github.com/StackExchange/wmi"
	"github.com/atotto/clipboard"
	"github.com/botherder/go-autoruns"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
)

//type systemInfo struct {
//	Time         string
//	ComputerName string
//	ComputerUser string
//	DNSHostName string
//	Domain      string
//	DomainRole  int
//	Model       string
//	Manufacturer	string
//	BuildNumber	string
//	BuildType	string
//	CountryCode	string
//	CurrentTimeZone	int
//	InstallDate	string
//	Locale	string
//	OSArchitecture	string
//	OSLanguage	int
//	SystemDirectory	string
//	LastBootUpTime	time.Time
//}
//type Win32_ComputerSystem struct {
//	DNSHostName string
//	Domain      string
//	DomainRole  int
//	Model       string
//	Name		string
//	UserName	string
//	Manufacturer	string
//}
//
//type Win32_OperatingSystem  struct {
//	BuildNumber	string
//	BuildType	string
//	CountryCode	string
//	CurrentTimeZone	int
//	InstallDate	string
//	Locale	string
//	OSArchitecture	string
//	OSLanguage	int
//	SystemDirectory	string
//	LastBootUpTime	time.Time
//}
func getClipboard(caseName string) error {
	clip, err := clipboard.ReadAll()
	if err != nil {
		return errors.Wrap(err, "Copy Clipboard Failed.")
	}
	return WriteToFile(filepath.Join(caseName, "Clipboard.txt"), clip)
}

func getAutoruns(caseName string) error {
	autorunList := autoruns.Autoruns()
	autorunsJSONPath := filepath.Join(caseName, "Autoruns.json")
	autorunsJSON, err := os.Create(autorunsJSONPath)
	if err != nil {
		return errors.Wrap(err, "Autoruns Failed.")
	}
	defer autorunsJSON.Close()

	// Encoding into json.
	buf, _ := json.MarshalIndent(autorunList, "", "    ")

	autorunsJSON.WriteString(string(buf[:]))
	autorunsJSON.Sync()
	return nil
}

func getSystemInfo(caseName string) {
	//var info systemInfo
	//var result []Win32_ComputerSystem
	//if err := wmi.Query(wmi.CreateQuery(&result, ""), &result); err != nil {
	//	errors.New("Error in system info.")
	//}
	//var resultOs []Win32_OperatingSystem
	//if err := wmi.Query(wmi.CreateQuery(&resultOs, ""), &resultOs); err != nil {
	//	errors.New("Error in OS system info.")
	//}
	//info.DNSHostName = result[0].DNSHostName
	//info.Domain = result[0].Domain
	//info.DomainRole = result[0].DomainRole
	//info.Model = result[0].Model
	//info.Time = time.Now().UTC().Format(time.RFC822)
	//info.ComputerName = result[0].Name
	//info.ComputerUser = result[0].UserName
	//info.Manufacturer = result[0].Manufacturer
	//info.BuildNumber = resultOs[0].BuildNumber
	//info.BuildType = resultOs[0].BuildType
	//info.CountryCode = resultOs[0].CountryCode
	//info.CurrentTimeZone = resultOs[0].CurrentTimeZone
	//info.InstallDate = resultOs[0].InstallDate
	//info.Locale = resultOs[0].Locale
	//info.OSArchitecture = resultOs[0].OSArchitecture
	//info.OSLanguage = resultOs[0].OSLanguage
	//info.SystemDirectory = resultOs[0].SystemDirectory
	//info.LastBootUpTime = resultOs[0].LastBootUpTime
	//fmt.Println(info.ComputerUser)
	//fmt.Println(info.ComputerName)
	//fmt.Println(info.Time)
	//fmt.Println(info.Manufacturer)
	//fmt.Println(info.Model)
	//fmt.Println(info.DomainRole)
	//fmt.Println(info.Domain)
	//fmt.Println(info.DNSHostName)
	//fmt.Println(info.BuildNumber)
	//fmt.Println(info.BuildType)
	//fmt.Println(info.CountryCode)
	//fmt.Println(info.CurrentTimeZone)
	//fmt.Println(info.InstallDate)
	//fmt.Println(info.Locale)
	//fmt.Println(info.OSArchitecture)
	//fmt.Println(info.OSLanguage)
	//fmt.Println(info.SystemDirectory)
	//fmt.Println(info.LastBootUpTime)
	//fmt.Printf("%#v", info)
	out, err := exec.Command("cmd", "/c", "systeminfo").Output()
	if err != nil {
		errors.Wrap(err, "Failed to get systeminfo.")
	}
	var data = strings.TrimSpace(string(out))
	WriteToFile(filepath.Join(caseName, "SystemInfo.txt"), string(data))
}

func getArpTable(caseName string) {
	out, err := exec.Command("cmd", "/c", "arp -a").Output()
	if err != nil {
		errors.Wrap(err, "Failed to get arp table.")
	}
	var data = strings.TrimSpace(string(out))
	WriteToFile(filepath.Join(caseName, "ArpTable.txt"), string(data))
}

func getIPConfig(caseName string) {
	out, err := exec.Command("cmd", "/c", "ipconfig /all").Output()
	if err != nil {
		errors.Wrap(err, "Failed to get ip configuration.")
	}
	var data = strings.TrimSpace(string(out))
	WriteToFile(filepath.Join(caseName, "IPConfiguration.txt"), string(data))
}
func getProxyConfig(caseName string) {
	out, err := exec.Command("cmd", "/c", "netsh winhttp show proxy").Output()
	if err != nil {
		errors.Wrap(err, "Failed to get proxy configuration.")
	}
	var data = strings.TrimSpace(string(out))
	WriteToFile(filepath.Join(caseName, "ProxyConfiguration.txt"), string(data))
}
func getVSSConfig(caseName string) {
	out, err := exec.Command("powershell", "-command", fmt.Sprintf("Get-ComputerRestorePoint | Format-Table -AutoSize")).Output()
	if err != nil {
		errors.Wrap(err, "Failed to get system restore points configuration.")
	}
	var data = strings.TrimSpace(string(out))
	WriteToFile(filepath.Join(caseName, "SystemRestorePoints.txt"), string(data))
}
func getRecyclebinitems(caseName string) {
	out, err := exec.Command("powershell", "-command", fmt.Sprintf("(New-Object -ComObject Shell.Application).NameSpace(0x0a).Items() | Format-Table -AutoSize")).Output()
	if err != nil {
		errors.Wrap(err, "Failed to get recyclebin items  configuration.")
	}
	var data = strings.TrimSpace(string(out))
	WriteToFile(filepath.Join(caseName, "RecyclebinItems.txt"), string(data))
}
