package dshow

import (
	_ "embed"
	"fmt"

	goi18n "github.com/bearki/go-i18n/v2"
	"gopkg.in/yaml.v3"
)

// 错误码
type errno uint

// 包名，以便于错误码能正常描述
const packageDescription = "becam camera lib direct show implement"

//go:embed errno_i18n.yml
var errYml []byte

// 错误码描述映射
var errMap = make(map[string]map[goi18n.Code]string)

// 程序启动之前加载国际化配置
func init() {
	_ = yaml.Unmarshal(errYml, errMap)
	errYml = nil // 置空，尝试释放掉内存
}

// Error 实现error接口
func (e errno) Error() string {
	// 是否存在该错误码
	errName, ok := errVarName[e]
	if !ok {
		// 不存在的错误码
		return fmt.Sprintf("unknown %s errno: %d", packageDescription, e)
	}
	// 错误码是否有对应的国际化描述
	if errMsgMap, ok := errMap[errName]; ok && len(errMsgMap) > 0 {
		// 是否有对应语言的描述
		errMsg, ok := errMsgMap[goi18n.GetEnv()]
		if ok {
			// 有描述
			return errMsg
		} else {
			// 无对应语言，但是有其他语言的描述
			for _, msg := range errMsgMap {
				return msg
			}
		}
	}
	// 错误码没有对应的国际化描述
	return fmt.Sprintf("raw %s errno: %d", packageDescription, e)
}
