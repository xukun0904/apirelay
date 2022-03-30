package util

func IsBlank(str string) bool {
	return str == ""
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

func BlankToDefault(str string, defaultStr string) string {
	if IsBlank(str) {
		return defaultStr
	}
	return str
}

func ZeroToDefault(num uint64, defaultNum uint64) uint64 {
	if num == 0 {
		return defaultNum
	}
	return num
}
