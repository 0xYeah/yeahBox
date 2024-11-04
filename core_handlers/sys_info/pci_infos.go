package sys_info

import (
	"github.com/tidwall/gjson"
	"os/exec"
	"runtime"
)

func getPciInfo() []byte {
	as := ""
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("system_profiler", "SPNVMeDataType", "-json", "-detailLevel", "basic")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return nil
		}
		as = gjson.ParseBytes(nob).Get("SPNVMeDataType.0._items.0.device_serial").Str
	}
	return []byte(as)
}
