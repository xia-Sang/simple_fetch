package helper

import (
	"fmt"
	"strings"
	"time"
)

type PCInfo struct {
	OSInfo
	CPUInfo
	LogicalDiskUsageInfo
	PhysicalDiskUsageInfo
}

func (pc *PCInfo) String() string {
	var sb strings.Builder

	sb.WriteString("PC Information:\n")
	sb.WriteString("OS Info:\n")
	sb.WriteString(pc.OSInfo.String())
	sb.WriteString("\n")

	sb.WriteString("CPU Info:\n")
	sb.WriteString(pc.CPUInfo.String())
	sb.WriteString("\n")

	sb.WriteString("Logical Disk Usage:\n")
	sb.WriteString(pc.LogicalDiskUsageInfo.String())
	sb.WriteString("\n")

	sb.WriteString("Physical Disk Info:\n")
	sb.WriteString(pc.PhysicalDiskUsageInfo.String())

	return sb.String()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("2006-01-02 15:04:05")
}

type OSInfo struct {
	UserName               string    //用户信息
	Manufacturer           string    //厂商信息
	OSInfo                 string    //存储windows版本信息
	OSVersion              string    //电脑版本
	OSArchitecture         string    //架构 支持多少位
	RegisteredUser         string    //电脑注册用户
	SerialNumber           string    //序列号信息
	InstallDate            time.Time //安装时间
	LastBootUpTime         time.Time //最后一次启动时间
	LocalDateTime          time.Time //本地时间
	FreePhysicalMemory     int64     //空闲物理内存空间
	TotalVisibleMemorySize int64     //总共内存空间
}

func (oi OSInfo) String() string {
	return fmt.Sprintf(
		"User Name: %s\n"+
			"OS Info: %s\n"+
			"PC Manufacturer: %s\n"+
			"OS Version: %s\n"+
			"OS Architecture: %s\n"+
			"Registered User: %s\n"+
			"Serial Number: %s\n"+
			"Install Date: %s\n"+
			"Last Boot Up Time: %s\n"+
			"Local Date Time: %s\n"+
			"Free Physical Memory: %.2f GB\n"+
			"Total Visible Memory Size: %.2f GB",
		oi.UserName,
		oi.OSInfo,
		oi.Manufacturer,
		oi.OSVersion,
		oi.OSArchitecture,
		oi.RegisteredUser,
		oi.SerialNumber,
		formatTime(oi.InstallDate),
		formatTime(oi.LastBootUpTime),
		formatTime(oi.LocalDateTime),
		float64(oi.FreePhysicalMemory)/(1024*1024),
		float64(oi.TotalVisibleMemorySize)/(1024*1024),
	)
}

type CPUInfo struct {
	Name                      string //cpu信息
	NumberOfCores             uint32 //核心
	NumberOfLogicalProcessors uint32 //处理器
	MaxClockSpeed             uint32 //时钟频率
}

func (ci CPUInfo) String() string {
	return fmt.Sprintf(
		"CPU Name: %s\n"+
			"Number of Cores: %d\n"+
			"Number of Logical Processors: %d\n"+
			"Max Clock Speed: %.2f GHz",
		ci.Name,
		ci.NumberOfCores,
		ci.NumberOfLogicalProcessors,
		float64(ci.MaxClockSpeed)/1000,
	)
}

type LogicalDiskUsageInfo struct {
	Disks []DiskUsageInfo
}

type DiskUsageInfo struct {
	Name        string
	Description string
	TotalSize   uint64
	FreeSpace   uint64
}

func (di DiskUsageInfo) String() string {
	return fmt.Sprintf(
		"Disk Name: %s\n"+
			"Description: %s\n"+
			"Total Size: %.2f GB\n"+
			"Free Space: %.2f GB",
		di.Name,
		di.Description,
		float64(di.TotalSize)/(1024*1024*1024),
		float64(di.FreeSpace)/(1024*1024*1024),
	)
}

func (ldi LogicalDiskUsageInfo) String() string {
	var result strings.Builder
	for i, disk := range ldi.Disks {
		if i > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(disk.String())
	}
	return result.String()
}

type PhysicalDiskInfo struct {
	Caption string
	Model   string
	Name    string
	Size    uint64
}

type PhysicalDiskUsageInfo struct {
	Disks []PhysicalDiskInfo
}

func (pdi PhysicalDiskInfo) String() string {
	return fmt.Sprintf(
		"Caption: %s\n"+
			"Model: %s\n"+
			"Name: %s\n"+
			"Size: %.2f GB",
		pdi.Caption,
		pdi.Model,
		pdi.Name,
		float64(pdi.Size)/(1024*1024*1024),
	)
}

func (pdui PhysicalDiskUsageInfo) String() string {
	var result strings.Builder
	for i, disk := range pdui.Disks {
		if i > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(disk.String())
	}
	return result.String()
}
