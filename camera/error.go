package camera

import (
	"fmt"

	goi18n "github.com/bearki/go-i18n/v2"
)

// 相机异常错误号
type errno uint

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

// 错误码描述映射
var errMap = map[errno]map[goi18n.Code]string{}

func (e errno) Error() string {
	// 是否存在该错误
	errLang, ok := errMap[e]
	if !ok {
		return fmt.Sprintf("unknown becam errno: %d", e)
	} else if errLang == nil {
		return fmt.Sprintf("raw becam errno: %d", e)
	}
	// 是否存在对应语言
	errStr, ok := errLang[goi18n.GetEnv()]
	if ok {
		return errStr
	}
	// 优先使用英文
	errStr, ok = errLang[goi18n.EN_US]
	if ok {
		return errStr
	}
	// 次优先使用中文
	errStr, ok = errLang[goi18n.ZH_CN]
	if ok {
		return errStr
	}
	// 使用随机语言
	for _, v := range errLang {
		return v
	}
	// 不存在语言时
	return fmt.Sprintf("raw becam errno: %d", e)
}
