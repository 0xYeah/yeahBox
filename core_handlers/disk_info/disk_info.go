package disk_info

import (
	"github.com/tidwall/gjson"
	"os/exec"
	"runtime"
	"strings"
)

// GetDiskNumberWithMacOS 获取MacOS 硬盘序列号
func GetDiskNumberWithMacOS() string {
	as := ""
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("system_profiler", "SPNVMeDataType", "-json", "-detailLevel", "basic")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return ""
		}
		as = gjson.ParseBytes(nob).Get("SPNVMeDataType.0._items.0.device_serial").Str
	}
	return as
}

func GetDiskNumber() (disk_seria_number string) {
	diskNo := ""
	if runtime.GOOS == "darwin" {
		diskNo = GetDiskNumberWithMacOS()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "diskdrive", "get", "serialnumber")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return
		}
		diskNo = string(nob)
		diskNo = strings.ReplaceAll(diskNo, "\n", "")
		diskNo = strings.ReplaceAll(diskNo, "\r", "")
		diskNo = strings.ReplaceAll(diskNo, " ", "")
		diskNo = strings.ReplaceAll(diskNo, "SerialNumber", "")
		diskNo = strings.ReplaceAll(diskNo, ".", "")
	}

	return diskNo
}
