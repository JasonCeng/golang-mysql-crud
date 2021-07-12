package tools

import (
	"regexp"
	"strings"
)

/**
功能：从sql中解析获取密文的cisno
返回值：cisno string
 */
func getCisno(sql string) (cisno string) {
	//1、去除空格
	sql = strings.Replace(sql, " ", "", -1)

	//2、匹配'***'，取中间字符串作为cisno密文值
	reg := regexp.MustCompile(`'(.*?)'`)
	preStr := reg.FindString(sql)
	cisno = strings.Trim(preStr, "'")

	return cisno
}
