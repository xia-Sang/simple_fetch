package helper

import "testing"

func TestTime(t *testing.T) {
	ti := "20240816232957.500000+480"
	tim, err := ParseWindowsTime(ti)
	t.Log(tim)
	t.Log(err)
	ans := FormatWindowsTimeCustom(tim)
	t.Log(ans)
}

func TestOsInfo(t *testing.T) {
	osInfo, err := GetOSInfo()
	t.Log(err)
	t.Log(osInfo)
	t.Log(osInfo.FreePhysicalMemory / 1024 / 1024)
	t.Log(osInfo.TotalVisibleMemorySize / 1024 / 1024)
}
func TestOsInfo1(t *testing.T) {
	osInfo, err := GetCPUInfo()
	t.Log(err)
	t.Log(osInfo)
}
func TestOsInfo2(t *testing.T) {
	osInfo, err := GetLogicalDiskUsageInfo()
	t.Log(err)
	t.Log(osInfo)
}
func TestOsInfo3(t *testing.T) {
	osInfo, err := GetPhysicalDiskUsageInfo()
	t.Log(err)
	t.Log(osInfo)
}
