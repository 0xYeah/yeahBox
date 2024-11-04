package brandFlag

import "strings"

type Brand int

const (
	BrandNone Brand = iota
	BrandAMD
	BrandIntel
	BrandNVIDIA
)

func GetBrandString(str string) Brand {

	if strings.Contains(str, "intel") || strings.Contains(str, "Intel") {
		return BrandIntel
	} else if strings.Contains(str, "amd") || strings.Contains(str, "AMD") {
		return BrandAMD
	} else if strings.Contains(str, "nvidia") || strings.Contains(str, "NVIDIA") || strings.Contains(str, "Nvidia") {
		return BrandNVIDIA
	}

	return BrandNone
}
