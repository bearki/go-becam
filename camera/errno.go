package camera

import (
	_ "embed"
	"fmt"

	goi18n "github.com/bearki/go-i18n/v2"
	"gopkg.in/yaml.v3"
)

// 相机异常错误号
type errno uint

//go:embed errno.i18n.zh_cn.yml
var errno_ZH_CN []byte

//go:embed errno.i18n.en_us.yml
var errno_EN_US []byte

const (
	_                             errno = iota // 占位
	ErrDeviceUnsupportMjpegFormat              // 设备不支持MJPEG格式
	ErrDeviceNotFound                          // 设备未找到
	ErrWaitForFrameFailed                      // 等待帧失败
	ErrGetFrameFailed                          // 获取帧失败
	ErrCopyFrameFailed                         // 拷贝帧失败
	ErrGetFrameTimout                          // 获取帧超时
	ErrDeviceOpenFailed                        // 设备打开失败
	ErrDeviceConfigNotFound                    // 设备配置未找到
	ErrSetDeviceConfigFailed                   // 修改设备配置失败
	ErrGetDeviceConfigFailed                   // 获取设备配置失败
	ErrRunStreamingFailed                      // 运行取流线程失败
	ErrDecodeJpegImageFailed                   // 解码JPEG图像失败
	ErrDeviceNotOpen                           // 设备未打开
)

// 错误码变量名映射
var errVarName = map[errno]string{
	ErrDeviceUnsupportMjpegFormat: "ErrDeviceUnsupportMjpegFormat",
	ErrDeviceNotFound:             "ErrDeviceNotFound",
	ErrWaitForFrameFailed:         "ErrWaitForFrameFailed",
	ErrGetFrameFailed:             "ErrGetFrameFailed",
	ErrCopyFrameFailed:            "ErrCopyFrameFailed",
	ErrGetFrameTimout:             "ErrGetFrameTimout",
	ErrDeviceOpenFailed:           "ErrDeviceOpenFailed",
	ErrDeviceConfigNotFound:       "ErrDeviceConfigNotFound",
	ErrSetDeviceConfigFailed:      "ErrSetDeviceConfigFailed",
	ErrGetDeviceConfigFailed:      "ErrGetDeviceConfigFailed",
	ErrRunStreamingFailed:         "ErrRunStreamingFailed",
	ErrDecodeJpegImageFailed:      "ErrDecodeJpegImageFailed",
	ErrDeviceNotOpen:              "ErrDeviceNotOpen",
}

// 错误码描述映射
var errMap = make(map[string]string)

func init() {
	// 根据环境变量初始化资源
	switch goi18n.GetEnv() {
	case goi18n.ZH_CN: // 中文简体（中国大陆）
		_ = yaml.Unmarshal(errno_ZH_CN, errMap)
	case goi18n.EN_US: // 英语（美国）
		_ = yaml.Unmarshal(errno_EN_US, errMap)
	}
}

func (e errno) Error() string {
	// 是否存在该错误码
	errName, ok := errVarName[e]
	if !ok {
		// 不存在的错误码
		return fmt.Sprintf("unknown becam errno: %d", e)
	}
	// 错误码是否有对应的国际化描述
	if errMsg, ok := errMap[errName]; ok {
		// 由国际化描述
		return errMsg
	}
	// 错误码没有对应的国际化描述
	return fmt.Sprintf("raw becam errno: %d", e)
}
