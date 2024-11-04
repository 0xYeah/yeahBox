package sys_info

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/wmyeah/yeah_box/unit_tools"
	"math/big"
)

type MemoryInfo struct {
	Total           uint64
	TotalShow       string
	Used            uint64
	usedShow        string
	UsedPercent     float64
	UsedPercentShow string
	Free            uint64
	FreeShow        string
}

func (mem *MemoryInfo) handleShow() {
	mem.TotalShow = unit_tools.UnitFormatWith1024(big.NewFloat(float64(mem.Total)), 2)
	mem.usedShow = unit_tools.UnitFormatWith1024(big.NewFloat(float64(mem.Used)), 2)
	mem.FreeShow = unit_tools.UnitFormatWith1024(big.NewFloat(float64(mem.Free)), 2)
	mem.UsedPercentShow = fmt.Sprintf("%.2f%s", mem.UsedPercent, "%")

}

func getMemoryInfo() *MemoryInfo {
	aMemInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil
	}
	aMem := &MemoryInfo{
		Total:       aMemInfo.Total,
		Used:        aMemInfo.Used,
		UsedPercent: aMemInfo.UsedPercent,
		Free:        aMemInfo.Free,
	}
	aMem.handleShow()
	return aMem
}
