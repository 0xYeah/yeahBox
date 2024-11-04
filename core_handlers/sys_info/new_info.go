package sys_info

import (
	"github.com/shirou/gopsutil/v4/net"
)

type NetInfo struct {
	Stat []net.IOCountersStat
}

func getNetInfos() (NetInfo, error) {
	n, err := net.IOCounters(true)
	if err != nil {
		return NetInfo{}, err
	}

	return NetInfo{n}, nil
}
