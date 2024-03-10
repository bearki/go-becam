package camera

import (
	"fmt"

	goi18n "github.com/bearki/go-i18n/v2"
)

// 相机异常错误号
type Errno uint

const (
	ErrQueryListFail                 Errno = iota + 1 // 查询相机列表失败
	ErrNotFoundMatchCamera                            // 未查询到匹配的相机
	ErrOpenFail                                       // 打开相机失败
	ErrStartStreamingFail                             // 启动取流线程失败
	ErrSetSupportInfoFail                             // 设置相机分辨率失败
	ErrNotOpen                                        // 相机未打开
	ErrWriteChannelSendTimeout                        // 最新图像写信号发送超时
	ErrWriteChannelRecvTimeout                        // 最新图像写信号接收超时
	ErrReadChannelSendTimeout                         // 最新图像读信号发送超时
	ErrReadChannelRecvTimeout                         // 最新图像读信号接收超时
	ErrGetStreamFail                                  // 取流失败
	ErrCopyStreamFail                                 // 拷贝流失败
	ErrNotFoundMatchDeviceConfigInfo                  // 未查询到匹配的设备配置信息
	ErrGetCurrDeviceInfoFail                          // 获取当前相机的设备信息失败
	ErrGetCurrConfigInfoFail                          // 获取当前相机的配置信息失败
	ErrCloseStartChannelSendTimeout                   // 相机关闭开始信号发送超时
	ErrCloseEndChannelSendTimeout                     // 相机关闭结束信号发送超时
	ErrCloseEndChannelRecvTimeout                     // 相机关闭结束信号接收超时
)

var errLangMap = map[Errno]map[goi18n.Code]string{
	ErrQueryListFail: {
		goi18n.EN_US: "Failed to query camera list",
		goi18n.ZH_CN: "查询相机列表失败",
	},
	ErrNotFoundMatchCamera: {
		goi18n.EN_US: "No matching camera found",
		goi18n.ZH_CN: "未查询到匹配的相机",
	},
	ErrOpenFail: {
		goi18n.EN_US: "Failed to open camera",
		goi18n.ZH_CN: "打开相机失败",
	},
	ErrStartStreamingFail: {
		goi18n.EN_US: "Failed to start streaming thread",
		goi18n.ZH_CN: "启动取流线程失败",
	},
	ErrSetSupportInfoFail: {
		goi18n.EN_US: "Failed to set camera resolution",
		goi18n.ZH_CN: "设置相机分辨率失败",
	},
	ErrNotOpen: {
		goi18n.EN_US: "camera not on",
		goi18n.ZH_CN: "相机未打开",
	},
	ErrWriteChannelSendTimeout: {
		goi18n.EN_US: "The latest image write signal sending timeout",
		goi18n.ZH_CN: "最新图像写信号发送超时",
	},
	ErrWriteChannelRecvTimeout: {
		goi18n.EN_US: "The latest image write signal receiving timeout",
		goi18n.ZH_CN: "最新图像写信号接收超时",
	},
	ErrReadChannelSendTimeout: {
		goi18n.EN_US: "The latest image read signal sending timeout",
		goi18n.ZH_CN: "最新图像读信号发送超时",
	},
	ErrReadChannelRecvTimeout: {
		goi18n.EN_US: "The latest image read signal reception timeout",
		goi18n.ZH_CN: "最新图像读信号接收超时",
	},
	ErrGetStreamFail: {
		goi18n.EN_US: "Failed to fetch stream",
		goi18n.ZH_CN: "取流失败",
	},
	ErrCopyStreamFail: {
		goi18n.EN_US: "copy stream failed",
		goi18n.ZH_CN: "拷贝流失败",
	},
	ErrNotFoundMatchDeviceConfigInfo: {
		goi18n.EN_US: "No matching device config information found",
		goi18n.ZH_CN: "未查询到匹配的设备配置信息",
	},
	ErrGetCurrDeviceInfoFail: {
		goi18n.EN_US: "Failed to get the device information of the current camera",
		goi18n.ZH_CN: "获取当前相机的设备信息失败",
	},
	ErrGetCurrConfigInfoFail: {
		goi18n.EN_US: "Failed to get the configuration information of the current camera",
		goi18n.ZH_CN: "获取当前相机的配置信息失败",
	},
	ErrCloseStartChannelSendTimeout: {
		goi18n.EN_US: "Camera off start signal send timed out",
		goi18n.ZH_CN: "相机关闭开始信号发送超时",
	},
	ErrCloseEndChannelSendTimeout: {
		goi18n.EN_US: "Camera close end signal sending timeout",
		goi18n.ZH_CN: "相机关闭结束信号发送超时",
	},
	ErrCloseEndChannelRecvTimeout: {
		goi18n.EN_US: "Camera off end signal reception timeout",
		goi18n.ZH_CN: "相机关闭结束信号接收超时",
	},
}

func (e Errno) Error() string {
	// 是否存在该错误
	errLang, ok := errLangMap[e]
	if !ok {
		return fmt.Sprintf("unknown camera errno: %d", e)
	} else if errLang == nil {
		return fmt.Sprintf("raw camera errno: %d", e)
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
	return fmt.Sprintf("raw camera errno: %d", e)
}
