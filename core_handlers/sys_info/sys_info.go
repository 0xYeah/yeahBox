package sys_info

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v4/host"
)

var CurrentSysInfo = &SysInfo{}

type SysInfo struct {
	CpuInfos   map[string]CpuInfo `json:"cpu"`
	MemoryInfo MemoryInfo         `json:"mem"`
	DiskInfo   DiskInfo           `json:"disk"`
	NetInfo    NetInfo            `json:"net"`
	Host       host.InfoStat      `json:"host"`
}

func RefreshSysInfos() (*SysInfo, error) {
	var err error

	CurrentSysInfo.CpuInfos, err = getCpuInfos()

	d, err := getDiskInfos()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("disk error: %v", err))
	}
	CurrentSysInfo.DiskInfo = d

	n, err := getNetInfos()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("net error: %v", err))

	}
	CurrentSysInfo.NetInfo = n

	CurrentSysInfo.MemoryInfo = *getMemoryInfo()

	h, err := host.Info()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("host error: %v", err))

	}
	CurrentSysInfo.Host = *h

	return CurrentSysInfo, nil
}
