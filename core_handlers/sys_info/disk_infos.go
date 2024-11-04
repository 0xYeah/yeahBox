package sys_info

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskInfo struct {
	Partitions []MountInfo
}

type MountInfo struct {
	Path      string
	Partition disk.PartitionStat
	Usage     disk.UsageStat
}

func getDiskInfos() (DiskInfo, error) {
	//disk_info.GetDiskNumber()

	aDiskInfo, err := disk.Partitions(true)
	if err != nil {
		return DiskInfo{}, err
	}
	var mounts []MountInfo
	for _, a := range aDiskInfo {
		u, err := disk.Usage(a.Mountpoint)
		if err != nil {
			fmt.Printf("get path useage error: %v", err)
			continue
		}

		mounts = append(mounts, MountInfo{Partition: a, Path: a.Mountpoint, Usage: *u})
	}

	return DiskInfo{mounts}, nil

}
