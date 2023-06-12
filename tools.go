package push

import (
	"fmt"
	"strings"
)

// 解析模板
func parserTemplate(content string, templateParams map[string]string) string {
	for key, value := range templateParams {
		content = strings.ReplaceAll(content, fmt.Sprintf("{%s}", key), value)
	}
	return content
}
