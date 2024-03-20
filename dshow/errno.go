package dshow

import (
	_ "embed"
	"fmt"

	goi18n "github.com/bearki/go-i18n/v2"
	"gopkg.in/yaml.v3"
)

// DirectShow相机异常错误号
type errno uint

//go:embed errno.i18n.zh_cn.yml
var errno_ZH_CN []byte

//go:embed errno.i18n.en_us.yml
var errno_EN_US []byte

const (
	_                       errno = iota // 占位
	NotFoundDevice                       // 未找到设备
	ErrHandleEmpty                       // Becam接口句柄未初始化
	ErrInputParam                        // 传入参数错误
	ErrInternalParam                     // 内部参数错误
	ErrInitCom                           // 初始化COM库失败
	ErrCreateEnumerator                  // 创建设备枚举器失败
	ErrDeviceEnum                        // 设备枚举失败
	ErrGetDeviceProp                     // 获取设备属性失败
	ErrGetStreamCaps                     // 获取设备流能力失败
	ErrNomatchStreamCaps                 // 未匹配到流能力
	ErrSetMediaType                      // 设置媒体类型失败
	ErrSelectedDevice                    // 选择设备失败
	ErrCreateGraphBuilder                // 创建图像构建器失败
	ErrAddCaptureFilter                  // 添加捕获过滤器到图像构建器失败
	ErrCreateSampleGrabber               // 创建样品采集器失败
	ErrGetSampleGrabberInfc              // 获取样品采集器接口失败
	ErrAddSampleGrabber                  // 添加样品采集器到图像构建器失败
	ErrCreateMediaControl                // 创建媒体控制器失败
	ErrCreateNullRender                  // 创建空渲染器失败
	ErrAddNullRender                     // 添加空渲染器到图像构建器失败
	ErrCaptureGrabber                    // 连接捕获器和采集器失败
	ErrGrabberRender                     // 连接采集器和渲染器失败
	ErrDeviceNotOpen                     // 设备未打开
	ErrFrameEmpty                        // 视频帧为空
	ErrFrameNotUpdate                    // 视频帧未更新
	ErrRepeatOpening                     // 请勿频繁重复打开相机
)

// 错误码变量名映射
var errVarName = map[errno]string{
	NotFoundDevice:          "NotFoundDevice",
	ErrHandleEmpty:          "ErrHandleEmpty",
	ErrInputParam:           "ErrInputParam",
	ErrInternalParam:        "ErrInternalParam",
	ErrInitCom:              "ErrInitCom",
	ErrCreateEnumerator:     "ErrCreateEnumerator",
	ErrDeviceEnum:           "ErrDeviceEnum",
	ErrGetDeviceProp:        "ErrGetDeviceProp",
	ErrGetStreamCaps:        "ErrGetStreamCaps",
	ErrNomatchStreamCaps:    "ErrNomatchStreamCaps",
	ErrSetMediaType:         "ErrSetMediaType",
	ErrSelectedDevice:       "ErrSelectedDevice",
	ErrCreateGraphBuilder:   "ErrCreateGraphBuilder",
	ErrAddCaptureFilter:     "ErrAddCaptureFilter",
	ErrCreateSampleGrabber:  "ErrCreateSampleGrabber",
	ErrGetSampleGrabberInfc: "ErrGetSampleGrabberInfc",
	ErrAddSampleGrabber:     "ErrAddSampleGrabber",
	ErrCreateMediaControl:   "ErrCreateMediaControl",
	ErrCreateNullRender:     "ErrCreateNullRender",
	ErrAddNullRender:        "ErrAddNullRender",
	ErrCaptureGrabber:       "ErrCaptureGrabber",
	ErrGrabberRender:        "ErrGrabberRender",
	ErrDeviceNotOpen:        "ErrDeviceNotOpen",
	ErrFrameEmpty:           "ErrFrameEmpty",
	ErrFrameNotUpdate:       "ErrFrameNotUpdate",
	ErrRepeatOpening:        "ErrRepeatOpening",
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
		return fmt.Sprintf("unknown becam direct show errno: %d", e)
	}
	// 错误码是否有对应的国际化描述
	if errMsg, ok := errMap[errName]; ok {
		// 由国际化描述
		return errMsg
	}
	// 错误码没有对应的国际化描述
	return fmt.Sprintf("raw becam direct show errno: %d", e)
}
