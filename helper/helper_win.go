//go:build windows
// +build windows

package helper

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func NewPCInfo() (*PCInfo, error) {
	osInfo, err := GetOSInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get OS info: %v", err)
	}

	cpuInfo, err := GetCPUInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %v", err)
	}

	logicalDiskInfo, err := GetLogicalDiskUsageInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get logical disk info: %v", err)
	}

	physicalDiskInfo, err := GetPhysicalDiskUsageInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get physical disk info: %v", err)
	}

	pcInfo := &PCInfo{
		OSInfo:                *osInfo,
		CPUInfo:               *cpuInfo,
		LogicalDiskUsageInfo:  *logicalDiskInfo,
		PhysicalDiskUsageInfo: *physicalDiskInfo,
	}

	return pcInfo, nil
}
func getComputerModel() (string, error) {
	cmd := exec.Command("wmic", "computersystem", "get", "model")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute wmic command: %v", err)
	}

	utf8Output, err := io.ReadAll(transform.NewReader(bytes.NewReader(output), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return "", fmt.Errorf("failed to convert output to UTF-8: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(utf8Output)), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected output format")
	}

	model := strings.TrimSpace(lines[1])
	return model, nil
}

func parseSystemInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}

func getSystemInfo() (string, error) {
	cmd := exec.Command("wmic", "os", "get", "/all", "/format:list")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	utf8Output, err := io.ReadAll(transform.NewReader(bytes.NewReader(output), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(utf8Output)), nil

}

func GetOSInfo() (*OSInfo, error) {
	osInfoStr, err := getSystemInfo()
	if err != nil {
		return nil, err
	}

	osInfoMap := parseSystemInfo(osInfoStr)

	osInfo := &OSInfo{
		UserName:               osInfoMap["CSName"],
		OSInfo:                 osInfoMap["Caption"],
		OSVersion:              osInfoMap["Version"],
		OSArchitecture:         osInfoMap["OSArchitecture"],
		RegisteredUser:         osInfoMap["RegisteredUser"],
		SerialNumber:           osInfoMap["SerialNumber"],
		FreePhysicalMemory:     parseUint64(osInfoMap["FreePhysicalMemory"]),
		TotalVisibleMemorySize: parseUint64(osInfoMap["TotalVisibleMemorySize"]),
	}

	manufacturer, err := getComputerModel()
	if err != nil {
		osInfo.Manufacturer = "unKnown"
	}
	osInfo.Manufacturer = manufacturer
	if installDate, err := ParseWindowsTime(osInfoMap["InstallDate"]); err == nil {
		osInfo.InstallDate = installDate
	}

	if lastBootUpTime, err := ParseWindowsTime(osInfoMap["LastBootUpTime"]); err == nil {
		osInfo.LastBootUpTime = lastBootUpTime
	}

	if localDateTime, err := ParseWindowsTime(osInfoMap["LocalDateTime"]); err == nil {
		osInfo.LocalDateTime = localDateTime
	}

	return osInfo, nil
}

func getCPUInfo() (string, error) {
	cmd := exec.Command("wmic", "cpu", "get", "name,", "NumberOfCores,", "NumberOfLogicalProcessors,", "MaxClockSpeed", "/format:list")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	utf8Output, err := io.ReadAll(transform.NewReader(bytes.NewReader(output), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(utf8Output)), nil
}

func GetCPUInfo() (*CPUInfo, error) {
	cpuInfoStr, err := getCPUInfo()
	if err != nil {
		return nil, err
	}

	cpuInfoMap := parseSystemInfo(cpuInfoStr)

	cpuInfo := &CPUInfo{
		Name:                      cpuInfoMap["Name"],
		NumberOfCores:             parseUint32(cpuInfoMap["NumberOfCores"]),
		NumberOfLogicalProcessors: parseUint32(cpuInfoMap["NumberOfLogicalProcessors"]),
		MaxClockSpeed:             parseUint32(cpuInfoMap["MaxClockSpeed"]),
	}
	return cpuInfo, nil
}
func GetLogicalDiskUsageInfo() (*LogicalDiskUsageInfo, error) {
	cmd := exec.Command("wmic", "logicaldisk", "get", "name,size,freespace,description", "/all", "/format:list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute wmic command: %v", err)
	}

	utf8Output, err := io.ReadAll(transform.NewReader(bytes.NewReader(output), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, fmt.Errorf("failed to decode output: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(utf8Output)), "\n")
	var disks []DiskUsageInfo
	var currentDisk DiskUsageInfo

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Description":
			if currentDisk.Name != "" {
				disks = append(disks, currentDisk)
				currentDisk = DiskUsageInfo{}
			}
			currentDisk.Description = value
		case "FreeSpace":
			currentDisk.FreeSpace, _ = strconv.ParseUint(value, 10, 64)
		case "Name":
			currentDisk.Name = value
		case "Size":
			currentDisk.TotalSize, _ = strconv.ParseUint(value, 10, 64)
		}
	}

	if currentDisk.Name != "" {
		disks = append(disks, currentDisk)
	}

	return &LogicalDiskUsageInfo{Disks: disks}, nil
}

func GetPhysicalDiskUsageInfo() (*PhysicalDiskUsageInfo, error) {
	cmd := exec.Command("wmic", "diskdrive", "get", "name,size,model,caption", "/all", "/format:list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute wmic command: %v", err)
	}

	utf8Output, err := io.ReadAll(transform.NewReader(bytes.NewReader(output), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, fmt.Errorf("failed to decode output: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(utf8Output)), "\n")
	var disks []PhysicalDiskInfo
	var currentDisk PhysicalDiskInfo

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Caption":
			if currentDisk.Name != "" {
				disks = append(disks, currentDisk)
				currentDisk = PhysicalDiskInfo{}
			}
			currentDisk.Caption = value
		case "Model":
			currentDisk.Model = value
		case "Name":
			currentDisk.Name = value
		case "Size":
			currentDisk.Size, _ = strconv.ParseUint(value, 10, 64)
		}
	}

	if currentDisk.Name != "" {
		disks = append(disks, currentDisk)
	}
	return &PhysicalDiskUsageInfo{Disks: disks}, nil
}
