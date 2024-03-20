/**
有备无患
-static-libgcc -static-libstdc++ -Wl,-Bstatic -lstdc++
*/

package dshow

// #cgo windows,386 CFLAGS: -I${SRCDIR}/libbecamdshow_windows_i686/include
// #cgo windows,386 LDFLAGS: -L${SRCDIR}/libbecamdshow_windows_i686/lib -lbecamdshow
// #cgo windows,amd64 CFLAGS: -I${SRCDIR}/libbecamdshow_windows_amd64/include
// #cgo windows,amd64 LDFLAGS: -L${SRCDIR}/libbecamdshow_windows_amd64/lib -lbecamdshow
// #cgo LDFLAGS: -lstrmiids -lole32 -loleaut32 -Wl,-Bstatic -lstdc++
// #include <stdlib.h>
// #include "becam_helper.h"
import "C"
import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"sort"
	"sync"
	"time"
	"unsafe"

	"github.com/bearki/becam/camera"
)

const (
	// MJPEG图像格式（小端序）
	V4L2_PIX_FMT_MJPEG = (uint32('M') | (uint32('J') << 8) | (uint32('P') << 16) | (uint32('G') << 24))
)

// 转换状态码
func convertStatusCode(code C.StatusCode) error {
	switch code {
	case C.STATUS_CODE_SUCCESS:
		return nil
	case C.STATUS_CODE_NOT_FOUND_DEVICE:
		return NotFoundDevice
	case C.STATUS_CODE_ERR_HANDLE_EMPTY:
		return ErrHandleEmpty
	case C.STATUS_CODE_ERR_INPUT_PARAM:
		return ErrInputParam
	case C.STATUS_CODE_ERR_INTERNAL_PARAM:
		return ErrInternalParam
	case C.STATUS_CODE_ERR_INIT_COM:
		return ErrInitCom
	case C.STATUS_CODE_ERR_CREATE_ENUMERATOR:
		return ErrCreateEnumerator
	case C.STATUS_CODE_ERR_DEVICE_ENUM:
		return ErrDeviceEnum
	case C.STATUS_CODE_ERR_GET_DEVICE_PROP:
		return ErrGetDeviceProp
	case C.STATUS_CODE_ERR_GET_STREAM_CAPS:
		return ErrGetStreamCaps
	case C.STATUS_CODE_ERR_NOMATCH_STREAM_CAPS:
		return ErrNomatchStreamCaps
	case C.STATUS_CODE_ERR_SET_MEDIA_TYPE:
		return ErrSetMediaType
	case C.STATUS_CODE_ERR_SELECTED_DEVICE:
		return ErrSelectedDevice
	case C.STATUS_CODE_ERR_CREATE_GRAPH_BUILDER:
		return ErrCreateGraphBuilder
	case C.STATUS_CODE_ERR_ADD_CAPTURE_FILTER:
		return ErrAddCaptureFilter
	case C.STATUS_CODE_ERR_CREATE_SAMPLE_GRABBER:
		return ErrCreateSampleGrabber
	case C.STATUS_CODE_ERR_GET_SAMPLE_GRABBER_INFC:
		return ErrGetSampleGrabberInfc
	case C.STATUS_CODE_ERR_ADD_SAMPLE_GRABBER:
		return ErrAddSampleGrabber
	case C.STATUS_CODE_ERR_CREATE_MEDIA_CONTROL:
		return ErrCreateMediaControl
	case C.STATUS_CODE_ERR_CREATE_NULL_RENDER:
		return ErrCreateNullRender
	case C.STATUS_CODE_ERR_ADD_NULL_RENDER:
		return ErrAddNullRender
	case C.STATUS_CODE_ERR_CAPTURE_GRABBER:
		return ErrCaptureGrabber
	case C.STATUS_CODE_ERR_GRABBER_RENDER:
		return ErrGrabberRender
	case C.STATUS_CODE_ERR_DEVICE_NOT_OPEN:
		return ErrDeviceNotOpen
	case C.STATUS_CODE_ERR_FRAME_EMPTY:
		return ErrFrameEmpty
	case C.STATUS_CODE_ERR_FRAME_NOT_UPDATE:
		return ErrFrameNotUpdate
	default:
		return fmt.Errorf("unknow becam direct show errno: %d", int(code))
	}
}

// Control 相机控制器
type Control struct {
	rwmutex           sync.RWMutex        // 读写锁
	deviceCacheList   camera.DeviceList   // 缓存的相机信息列表
	handle            C.BecamHandle       // 相机库句柄
	deviceInfo        camera.Device       // 当前使用的相机信息
	deviceSupportInfo camera.DeviceConfig // 当前使用的相机支持信息
}

// NewControl 创建一个相机控制器
func New() *Control {
	// 初始化句柄
	handle := C.BecamNew()
	// OK
	return &Control{
		handle: handle,
	}
}

// 尝试获取帧
func (p *Control) tryGetFrame(parseW, parseH *uint32) ([]byte, error) {
	// 声明响应参数
	var replyData *C.uint8_t
	var replySize C.size_t

	// 循环取100次，有就立即跳出
	for i := 0; i < 100; i++ {
		// 执行取流
		code := C.BecamGetFrame(p.handle, &replyData, &replySize)
		if err := convertStatusCode(code); err != nil {
			if i == 100-1 {
				fmt.Fprintln(os.Stderr, err.Error())
				return nil, errors.Join(camera.ErrGetFrameFailed, err)
			}
		} else {
			break
		}
	}
	// 延迟释放
	defer C.BecamFreeFrame(p.handle, &replyData)

	// 拷贝帧
	data := C.GoBytes(unsafe.Pointer(replyData), C.int(replySize))

	// 是否需要解码图像
	if parseW != nil && parseH != nil {
		// 解码图像
		imgConf, err := jpeg.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			// 返回异常
			return nil, err
		}
		// 赋值宽高
		*parseW = uint32(imgConf.Width)
		*parseH = uint32(imgConf.Height)
	}

	// 获取帧成功
	return data, nil
}

// --------------------------------------------- 实现Manager接口 --------------------------------------------- //

// 释放所有相机资源
func (p *Control) Free() {
	// 关闭已打开的相机
	p.Close()

	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 释放句柄
	C.BecamFree(&p.handle)
	p.handle = nil
	p.deviceCacheList = nil
	p.deviceInfo = camera.Device{}
	p.deviceSupportInfo = camera.DeviceConfig{}
}

// GetList 获取相机列表
//
//	@return 相机列表
//	@return 错误信息
func (p *Control) GetList() (camera.DeviceList, error) {
	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 清空缓存列表
	p.deviceCacheList = nil

	// 调用C接口获取相机列表
	var reply C.GetDeviceListReply
	code := C.BecamGetDeviceList(p.handle, &reply)
	if err := convertStatusCode(code); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, errors.Join(camera.ErrGetDeviceConfigFailed, err)
	}
	defer C.BecamFreeDeviceList(p.handle, &reply)

	// 遍历相机列表
	for i := 0; i < int(reply.deviceInfoListSize); i++ {
		// 使用C助手函数获取设备信息
		device := C.getDeviceInfoListItem(reply.deviceInfoList, C.size_t(i))

		// 获取设备名称
		name := C.GoString(device.name)
		// 获取设备路径
		devicePath := C.GoString(device.devicePath)
		// 获取设备位置信息
		locationInfo := C.GoString(device.locationInfo)
		// 计算设备唯一ID
		idData := md5.Sum([]byte(devicePath + name))
		id := hex.EncodeToString(idData[:])
		// 构建设备信息
		info := &camera.Device{
			ID:           id,
			Name:         name,
			SymbolicLink: devicePath,
			LocationInfo: locationInfo,
			ConfigList:   camera.DeviceConfigList{camera.AutoDeviceConfig.Clone()},
		}

		// 遍历设备支持的视频帧
		var deviceConfigList camera.DeviceConfigList
		for j := 0; j < int(device.frameInfoListSize); j++ {
			// 使用C助手函数获取设备配置信息
			frameInfo := C.getFrameInfoListItem(device.frameInfoList, C.size_t(j))
			// 跳过不要的格式和分辨率
			if uint32(frameInfo.format) != V4L2_PIX_FMT_MJPEG || frameInfo.width < 600 {
				continue
			}
			// 追加配置信息
			deviceConfigList = append(deviceConfigList, &camera.DeviceConfig{
				Width:  uint32(frameInfo.width),
				Height: uint32(frameInfo.height),
				FPS:    uint32(frameInfo.fps),
			})
		}

		// 对支持信息进行排序
		sort.Slice(deviceConfigList, func(i, j int) bool {
			// 宽度是否相等
			if deviceConfigList[i].Width == deviceConfigList[j].Width {
				// 高度是否相等
				if deviceConfigList[i].Height == deviceConfigList[j].Height {
					// 按帧率从大到小排序
					return deviceConfigList[i].FPS >= deviceConfigList[j].FPS
				}
				// 按高度从大到小排序
				return deviceConfigList[i].Height > deviceConfigList[j].Height
			}
			// 按宽度从大到小排序
			return deviceConfigList[i].Width > deviceConfigList[j].Width
		})

		// 追加到相机配置
		info.ConfigList = append(info.ConfigList, deviceConfigList...)
		// 追加到相机列表
		p.deviceCacheList = append(p.deviceCacheList, info)
	}

	// 返回相机克隆列表
	return p.deviceCacheList.Clone(), nil
}

// GetDeviceWithID 通过相机ID获取缓存的相机信息
//
//	@param	id	相机ID
//	@return	缓存的相机信息
//	@return	异常信息
func (p *Control) GetDeviceWithID(id string) (*camera.Device, error) {
	// 加读锁
	p.rwmutex.RLock()
	defer p.rwmutex.RUnlock()
	// 执行查找
	return p.deviceCacheList.Get(id)
}

// Open 打开相机
//
//	@param	id		相机ID
//	@param	info	分辨率信息
//	@return	异常信息
func (p *Control) Open(id string, info camera.DeviceConfig) error {
	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 查询ID对应的相机信息
	cameraInfo, err := p.deviceCacheList.Get(id)
	if err != nil {
		return err
	}

	// 是否需要修改分辨率
	if info.Width > 0 && info.Height > 0 && info.FPS > 0 {
		// 筛选分辨率
		sInfo, err := cameraInfo.ConfigList.Get(info)
		if err != nil {
			return err
		}
		// 赋值选择的分辨率
		info.Width = sInfo.Width
		info.Height = sInfo.Height
		info.FPS = sInfo.FPS
	} else if len(cameraInfo.ConfigList) > 0 {
		// 选中与默认分辨率最相似的分辨率
		sInfo, err := cameraInfo.ConfigList.GetMostSimilar(camera.DefaultDeviceConfig)
		if err != nil {
			return err
		}
		// 赋值选择的分辨率
		info.Width = sInfo.Width
		info.Height = sInfo.Height
		info.FPS = sInfo.FPS
	}

	// 转换设备路径
	devicePath := C.CString(cameraInfo.SymbolicLink)
	defer C.free(unsafe.Pointer(devicePath))
	// 转换配置信息
	var frameInfo C.VideoFrameInfo
	frameInfo.width = C.uint32_t(info.Width)
	frameInfo.height = C.uint32_t(info.Height)
	frameInfo.fps = C.uint32_t(info.FPS)
	frameInfo.format = C.uint32_t(V4L2_PIX_FMT_MJPEG)
	// 执行打开
	code := C.BecamOpenDevice(p.handle, devicePath, &frameInfo)
	if err := convertStatusCode(code); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return errors.Join(camera.ErrDeviceOpenFailed, err)
	}

	// 赋值当前使用的相机信息
	p.deviceInfo = *cameraInfo

	// 尝试获取帧，计算其分辨率
	_, err = p.tryGetFrame(&p.deviceSupportInfo.Width, &p.deviceSupportInfo.Height)
	if err != nil {
		return err
	}

	// 赋值帧率
	p.deviceSupportInfo.FPS = info.FPS

	// OK
	return nil
}

// GetDeviceConfigInfo 获取当前设备配置信息
//
//	@return	当前设备信息
//	@return	当前设备配置信息
//	@return	异常信息
func (p *Control) GetCurrDeviceConfigInfo() (*camera.Device, *camera.DeviceConfig, error) {
	// 操作加读锁
	p.rwmutex.RLock()
	defer p.rwmutex.RUnlock()

	// 检查相机是否已打开
	if p.handle == nil {
		return nil, nil, camera.ErrDeviceNotOpen
	}

	// 返回结果
	return p.deviceInfo.Clone(), p.deviceSupportInfo.Clone(), nil
}

// GetStream 获取帧
//
//	@param	outWidth	需要解析图像宽高时请传入地址，以便于内部赋值
//	@param	outHeight	需要解析图像宽高时请传入地址，以便于内部赋值
//	@return	图片流
//	@return	异常信息
func (p *Control) GetFrame(outWidth, outHeight *uint32) ([]byte, error) {
	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 检查相机是否已打开
	if p.handle == nil {
		return nil, camera.ErrDeviceNotOpen
	}

	// 尝试获取帧
	return p.tryGetFrame(outWidth, outHeight)
}

// Close 关闭已打开的相机
func (p *Control) Close() {
	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 是否存在已打开的相机
	if p.handle == nil {
		return
	}

	// 释放相机内存
	C.BecamCloseDevice(p.handle)
	// 给内核一点时间
	time.Sleep(time.Millisecond * 100)
	// 清除当前使用的相机信息
	p.deviceInfo = camera.Device{}
	p.deviceSupportInfo = camera.DeviceConfig{}
}
