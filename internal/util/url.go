package util

import "regexp"

func IsURL(str string) bool {
	// 简单的 URL 正则 (支持 http/https、域名、端口、路径、参数)
	regex := `^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(:\d+)?(/.*)?$`
	re := regexp.MustCompile(regex)
	return re.MatchString(str)
}
