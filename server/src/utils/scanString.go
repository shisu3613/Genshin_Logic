package utils

import (
	"bufio"
	"os"
	"strings"
)

/**
    @author: WangYuding
    @since: 2022/5/6
    @desc: //添加扫描字符串的辅助函数
**/

func ScanString() string {
	var msgStr string
	for msgStr == "" {
		msgStr, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		msgStr = strings.TrimSpace(msgStr)
	}
	//fmt.Println(msgStr)
	return msgStr
}
