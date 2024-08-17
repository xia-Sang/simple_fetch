package helper

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func NewPCInfo() (*PCInfo, error) {
	info := &PCInfo{}

	var err error
	info.OSInfo, err = GetOSInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get OS info: %v", err)
	}

	info.CPUInfo, err = GetCPUInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %v", err)
	}

	info.LogicalDiskUsageInfo, err = GetLogicalDiskUsageInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get logical disk usage info: %v", err)
	}

	info.PhysicalDiskUsageInfo, err = GetPhysicalDiskUsageInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get physical disk usage info: %v", err)
	}

	return info, nil
}

func GetOSInfo() (OSInfo, error) {
	var info OSInfo

	info.UserName = os.Getenv("USER")

	out, err := exec.Command("lsb_release", "-a").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Description:") {
				info.OSInfo = strings.TrimSpace(strings.TrimPrefix(line, "Description:"))
			}
		}
	}

	out, err = exec.Command("uname", "-r").Output()
	if err == nil {
		info.OSVersion = strings.TrimSpace(string(out))
	}

	out, err = exec.Command("uname", "-m").Output()
	if err == nil {
		info.OSArchitecture = strings.TrimSpace(string(out))
	}

	info.RegisteredUser = info.UserName
	info.LocalDateTime = time.Now()

	out, err = exec.Command("free", "-b").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) > 1 {
				total, _ := strconv.ParseInt(fields[1], 10, 64)
				info.TotalVisibleMemorySize = total
			}
			if len(fields) > 3 {
				free, _ := strconv.ParseInt(fields[3], 10, 64)
				info.FreePhysicalMemory = free
			}
		}
	}

	return info, nil
}

func GetCPUInfo() (CPUInfo, error) {
	var info CPUInfo

	out, err := exec.Command("cat", "/proc/cpuinfo").Output()
	if err != nil {
		return info, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "model name") {
			info.Name = strings.TrimSpace(strings.TrimPrefix(line, "model name	:"))
		} else if strings.HasPrefix(line, "cpu cores") {
			cores, _ := strconv.ParseUint(strings.TrimSpace(strings.TrimPrefix(line, "cpu cores	:")), 10, 32)
			info.NumberOfCores = uint32(cores)
		} else if strings.HasPrefix(line, "siblings") {
			processors, _ := strconv.ParseUint(strings.TrimSpace(strings.TrimPrefix(line, "siblings	:")), 10, 32)
			info.NumberOfLogicalProcessors = uint32(processors)
		} else if strings.HasPrefix(line, "cpu MHz") {
			mhz, _ := strconv.ParseFloat(strings.TrimSpace(strings.TrimPrefix(line, "cpu MHz		:")), 64)
			info.MaxClockSpeed = uint32(mhz)
		}
	}

	return info, nil
}

func GetLogicalDiskUsageInfo() (LogicalDiskUsageInfo, error) {
	var info LogicalDiskUsageInfo

	out, err := exec.Command("df", "-B1").Output()
	if err != nil {
		return info, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines[1:] { // Skip header
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			size, _ := strconv.ParseUint(fields[1], 10, 64)
			free, _ := strconv.ParseUint(fields[3], 10, 64)
			info.Disks = append(info.Disks, DiskUsageInfo{
				Name:        fields[0],
				Description: fields[5],
				TotalSize:   size,
				FreeSpace:   free,
			})
		}
	}

	return info, nil
}

func GetPhysicalDiskUsageInfo() (PhysicalDiskUsageInfo, error) {
	var info PhysicalDiskUsageInfo

	out, err := exec.Command("lsblk", "-ndo", "NAME,SIZE,MODEL").Output()
	if err != nil {
		return info, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			size, _ := parseSize(fields[1])
			model := ""
			if len(fields) >= 3 {
				model = strings.Join(fields[2:], " ")
			}
			info.Disks = append(info.Disks, PhysicalDiskInfo{
				Name:    fields[0],
				Size:    size,
				Model:   model,
				Caption: fields[0],
			})
		}
	}

	return info, nil
}

func parseSize(sizeStr string) (uint64, error) {
	multiplier := uint64(1)
	if strings.HasSuffix(sizeStr, "G") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "G")
	} else if strings.HasSuffix(sizeStr, "T") {
		multiplier = 1024 * 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "T")
	}

	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, err
	}

	return uint64(size * float64(multiplier)), nil
}
