package camera

import (
	_ "embed"
)

const (
	_ errno = iota // 占位

	ErrEnumDeviceFailed           // 枚举设备失败
	ErrDeviceNotFound             // 设备未找到
	ErrDeviceMediaConfigNotFound  // 设备媒体配置未找到
	ErrGetDeviceMediaConfigFailed // 获取设备媒体配置失败
	ErrDeviceRepeatOpening        // 请勿频繁重复打开设备
	ErrDeviceOpenFailed           // 设备打开失败
	ErrGetFrameFailed             // 获取帧失败
	ErrDecodeJpegImageFailed      // 解码JPEG图像失败
	ErrDeviceNotOpen              // 设备未打开
)

// 错误码变量名映射
var errVarName = map[errno]string{
	ErrEnumDeviceFailed:           "ErrEnumDeviceFailed",
	ErrDeviceNotFound:             "ErrDeviceNotFound",
	ErrDeviceMediaConfigNotFound:  "ErrDeviceMediaConfigNotFound",
	ErrGetDeviceMediaConfigFailed: "ErrGetDeviceMediaConfigFailed",
	ErrDeviceRepeatOpening:        "ErrDeviceRepeatOpening",
	ErrDeviceOpenFailed:           "ErrDeviceOpenFailed",
	ErrGetFrameFailed:             "ErrGetFrameFailed",
	ErrDecodeJpegImageFailed:      "ErrDecodeJpegImageFailed",
	ErrDeviceNotOpen:              "ErrDeviceNotOpen",
}
