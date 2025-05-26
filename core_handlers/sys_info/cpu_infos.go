package sys_info

import (
	"fmt"
	"github.com/0xYeah/yeahBox/core_handlers/brand_flag"
	"github.com/shirou/gopsutil/v4/cpu"
	"strconv"
	"strings"
	"time"
)

type CpuInfo struct {
	SocketAt         int             `json:"socket_at"` // cpu 插槽
	Brand            brandFlag.Brand `json:"brand"`     // 厂商品牌
	Class            string          `json:"class"`     // 系列分类
	ClassRd          int64           `json:"class_rd"`  // 系列分类代数
	Cores            int64
	Threads          int64
	ModelName        string  `json:"model_name"`      // CPU 型号
	FrequencyMainMHz float64 `json:"frequency_main"`  // CPU 主频 单位MHz
	Percentage       float64 `json:"percentage"`      // 使用率
	PercentageShow   string  `json:"percentage_show"` //使用率显示
}

func getCpuInfos() (map[string]CpuInfo, error) {
	aCpuInfos, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	infos := map[string]CpuInfo{}

	for i, cInfo := range aCpuInfos {
		cuts := strings.Split(cInfo.ModelName, " ")
		cRd, _ := strconv.ParseInt(cInfo.Family, 10, 64)
		cpuInfo := CpuInfo{
			SocketAt:         i,
			Brand:            brandFlag.GetBrandString(cInfo.ModelName),
			Class:            cuts[1],
			ClassRd:          cRd,
			ModelName:        cInfo.ModelName,
			FrequencyMainMHz: cInfo.Mhz,
		}

		// 获取CPU使用率，参数interval指定计算使用率的时间间隔，返回值是每个CPU核心的使用率列表
		percentages, c_err := cpu.Percent(time.Second, false) // false表示获取总的CPU使用率
		if c_err != nil {
			fmt.Printf("get cpu info error: %v", err)
			continue
		}
		cpuInfo.Percentage = percentages[i]
		cpuInfo.PercentageShow = fmt.Sprintf("%.2f%%", percentages[i])
		infos[fmt.Sprintf("%d", i)] = cpuInfo
	}

	return infos, nil
}
