package internal

/*
#cgo linux pkg-config: becamv4l2
#cgo windows pkg-config: becamdshow
#include <stdlib.h>
#include "becam_helper.h"
*/
import "C"
import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
	"unsafe"

	"github.com/bearki/go-becam/camera"
)

// 转换状态码
func convertStatusCode(code C.StatusCode) error {
	switch code {
	case C.STATUS_CODE_SUCCESS:
		return nil
	case C.STATUS_CODE_ERR_HANDLE_EMPTY:
		return STATUS_CODE_ERR_HANDLE_EMPTY
	case C.STATUS_CODE_ERR_INPUT_PARAM:
		return STATUS_CODE_ERR_INPUT_PARAM
	case C.STATUS_CODE_ERR_DEVICE_ENUM_FAILED:
		return STATUS_CODE_ERR_DEVICE_ENUM_FAILED
	case C.STATUS_CODE_ERR_DEVICE_NOT_FOUND:
		return STATUS_CODE_ERR_DEVICE_NOT_FOUND
	case C.STATUS_CODE_ERR_DEVICE_OPEN_FAILED:
		return STATUS_CODE_ERR_DEVICE_OPEN_FAILED
	case C.STATUS_CODE_ERR_DEVICE_NOT_OPEN:
		return STATUS_CODE_ERR_DEVICE_NOT_OPEN
	case C.STATUS_CODE_ERR_DEVICE_FRAME_FMT_NOT_FOUND:
		return STATUS_CODE_ERR_DEVICE_FRAME_FMT_NOT_FOUND
	case C.STATUS_CODE_ERR_DEVICE_FRAME_FMT_SET_FAILED:
		return STATUS_CODE_ERR_DEVICE_FRAME_FMT_SET_FAILED
	case C.STATUS_CODE_ERR_DEVICE_RUN_FAILED:
		return STATUS_CODE_ERR_DEVICE_RUN_FAILED
	case C.STATUS_CODE_ERR_DEVICE_NOT_RUN:
		return STATUS_CODE_ERR_DEVICE_NOT_RUN
	case C.STATUS_CODE_ERR_GET_FRAME_FAILED:
		return STATUS_CODE_ERR_GET_FRAME_FAILED
	case C.STATUS_CODE_ERR_GET_FRAME_EMPTY:
		return STATUS_CODE_ERR_GET_FRAME_EMPTY
	case C.STATUS_CODE_DSHOW_ERR_INTERNAL_PARAM:
		return STATUS_CODE_DSHOW_ERR_INTERNAL_PARAM
	case C.STATUS_CODE_DSHOW_ERR_INIT_COM:
		return STATUS_CODE_DSHOW_ERR_INIT_COM
	case C.STATUS_CODE_DSHOW_ERR_CREATE_ENUMERATOR:
		return STATUS_CODE_DSHOW_ERR_CREATE_ENUMERATOR
	case C.STATUS_CODE_DSHOW_ERR_GET_DEVICE_PROP:
		return STATUS_CODE_DSHOW_ERR_GET_DEVICE_PROP
	case C.STATUS_CODE_DSHOW_ERR_GET_STREAM_CAPS:
		return STATUS_CODE_DSHOW_ERR_GET_STREAM_CAPS
	case C.STATUS_CODE_DSHOW_ERR_CREATE_GRAPH_BUILDER:
		return STATUS_CODE_DSHOW_ERR_CREATE_GRAPH_BUILDER
	case C.STATUS_CODE_DSHOW_ERR_ADD_CAPTURE_FILTER:
		return STATUS_CODE_DSHOW_ERR_ADD_CAPTURE_FILTER
	case C.STATUS_CODE_DSHOW_ERR_CREATE_SAMPLE_GRABBER:
		return STATUS_CODE_DSHOW_ERR_CREATE_SAMPLE_GRABBER
	case C.STATUS_CODE_DSHOW_ERR_GET_SAMPLE_GRABBER_INFC:
		return STATUS_CODE_DSHOW_ERR_GET_SAMPLE_GRABBER_INFC
	case C.STATUS_CODE_DSHOW_ERR_ADD_SAMPLE_GRABBER:
		return STATUS_CODE_DSHOW_ERR_ADD_SAMPLE_GRABBER
	case C.STATUS_CODE_DSHOW_ERR_CREATE_MEDIA_CONTROL:
		return STATUS_CODE_DSHOW_ERR_CREATE_MEDIA_CONTROL
	case C.STATUS_CODE_DSHOW_ERR_CREATE_NULL_RENDER:
		return STATUS_CODE_DSHOW_ERR_CREATE_NULL_RENDER
	case C.STATUS_CODE_DSHOW_ERR_ADD_NULL_RENDER:
		return STATUS_CODE_DSHOW_ERR_ADD_NULL_RENDER
	case C.STATUS_CODE_DSHOW_ERR_CAPTURE_GRABBER:
		return STATUS_CODE_DSHOW_ERR_CAPTURE_GRABBER
	case C.STATUS_CODE_DSHOW_ERR_GRABBER_RENDER:
		return STATUS_CODE_DSHOW_ERR_GRABBER_RENDER
	case C.STATUS_CODE_DSHOW_ERR_FRAME_NOT_UPDATE:
		return STATUS_CODE_DSHOW_ERR_FRAME_NOT_UPDATE
	case C.STATUS_CODE_MF_ERR_CREATE_ATTR_STORE:
		return STATUS_CODE_MF_ERR_CREATE_ATTR_STORE
	case C.STATUS_CODE_MF_ERR_SET_ATTR_STORE:
		return STATUS_CODE_MF_ERR_SET_ATTR_STORE
	case C.STATUS_CODE_MF_ERR_CREATE_PRESENT_DESC:
		return STATUS_CODE_MF_ERR_CREATE_PRESENT_DESC
	case C.STATUS_CODE_MF_ERR_GET_STREAM_DESC:
		return STATUS_CODE_MF_ERR_GET_STREAM_DESC
	case C.STATUS_CODE_MF_ERR_GET_MEDIA_TYPE_HANDLER:
		return STATUS_CODE_MF_ERR_GET_MEDIA_TYPE_HANDLER
	case C.STATUS_CODE_MF_ERR_GET_MEDIA_TYPE_COUNT:
		return STATUS_CODE_MF_ERR_GET_MEDIA_TYPE_COUNT
	case C.STATUS_CODE_MF_ERR_GET_MEDIA_TYPE:
		return STATUS_CODE_MF_ERR_GET_MEDIA_TYPE
	case C.STATUS_CODE_MF_ERR_CONVERT_FRAME_BUFFER:
		return STATUS_CODE_MF_ERR_CONVERT_FRAME_BUFFER
	case C.STATUS_CODE_MF_ERR_LOCK_FRAME_BUFFER:
		return STATUS_CODE_MF_ERR_LOCK_FRAME_BUFFER
	case C.STATUS_CODE_V4L2_ERR_REQUEST_BUF:
		return STATUS_CODE_V4L2_ERR_REQUEST_BUF
	case C.STATUS_CODE_V4L2_ERR_QUERY_BUF:
		return STATUS_CODE_V4L2_ERR_QUERY_BUF
	case C.STATUS_CODE_V4L2_ERR_MMAP_BUF:
		return STATUS_CODE_V4L2_ERR_MMAP_BUF
	case C.STATUS_CODE_V4L2_ERR_LOCK_BUF:
		return STATUS_CODE_V4L2_ERR_LOCK_BUF
	case C.STATUS_CODE_V4L2_ERR_UNLOCK_BUF:
		return STATUS_CODE_V4L2_ERR_UNLOCK_BUF
	default:
		return fmt.Errorf("unknow becam errno %d", int(code))
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
//
//	@return 帧数据
//	@return 帧信息
//	@return 错误信息
func (p *Control) tryGetFrame() ([]byte, *camera.DeviceConfig, error) {
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
				return nil, nil, errors.Join(camera.ErrGetFrameFailed, err)
			}
		} else {
			break
		}
	}
	// 延迟释放
	defer C.BecamFreeFrame(p.handle, &replyData)

	// 拷贝帧
	data := C.GoBytes(unsafe.Pointer(replyData), C.int(replySize))
	// 获取帧成功
	return data, p.deviceSupportInfo.Clone(), nil
}

// --------------------------------------------- 实现Manager接口 --------------------------------------------- //

// 获取相机列表（无锁）
//
//	@return 相机列表
//	@return 错误信息
func (p *Control) getList() (camera.DeviceList, error) {
	// 清空缓存列表
	p.deviceCacheList = nil

	// 调用C接口获取相机列表
	var reply C.GetDeviceListReply
	code := C.BecamGetDeviceList(p.handle, &reply)
	if err := convertStatusCode(code); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, errors.Join(camera.ErrEnumDeviceFailed, err)
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
		// 计算设备唯一ID
		idData := md5.Sum([]byte(devicePath + name))
		id := hex.EncodeToString(idData[:])
		// 追加到相机列表
		p.deviceCacheList = append(p.deviceCacheList, &camera.Device{
			ID:           id,
			Name:         name,
			SymbolicLink: devicePath,
		})
	}

	// 返回相机克隆列表
	return p.deviceCacheList.Clone(), nil
}

// GetList 获取相机列表（有锁）
//
//	@return 相机列表
//	@return 错误信息
func (p *Control) GetList() (camera.DeviceList, error) {
	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 调用内部实现
	return p.getList()
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

// 通过相机系统路径获取设备的配置信息（无锁）
//
//	@param	devicePath	相机系统路径
//	@return	设备配置信息
//	@return	异常信息
func (p *Control) getDeviceConfigInfo(devicePath string) (camera.DeviceConfigList, error) {
	// 转换设备路径
	devicePathPtr := C.CString(devicePath)
	defer C.free(unsafe.Pointer(devicePathPtr))

	// 获取支持的配置
	var reply C.GetDeviceConfigListReply
	code := C.BecamGetDeviceConfigList(p.handle, devicePathPtr, &reply)
	if err := convertStatusCode(code); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, errors.Join(camera.ErrGetDeviceMediaConfigFailed, err)
	}
	defer C.BecamFreeDeviceConfigList(p.handle, &reply)

	// 拷贝配置列表
	var deviceConfigList camera.DeviceConfigList
	for j := 0; j < int(reply.videoFrameInfoListSize); j++ {
		// 使用C助手函数获取设备配置信息
		frameInfo := C.getFrameInfoListItem(reply.videoFrameInfoList, C.size_t(j))
		// 追加配置信息
		deviceConfigList = append(deviceConfigList, &camera.DeviceConfig{
			Width:  uint32(frameInfo.width),
			Height: uint32(frameInfo.height),
			FPS:    uint32(frameInfo.fps),
			Format: camera.NewFourccFromNumber(uint32(frameInfo.format)),
		})
	}

	// 对支持信息进行排序
	sort.Slice(deviceConfigList, func(i, j int) bool {
		// 格式是否一致
		if deviceConfigList[i].Format == deviceConfigList[j].Format {
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
		}
		// 按格式随便
		return deviceConfigList[i].Format < deviceConfigList[j].Format
	})

	// 返回配置信息列表
	return deviceConfigList, nil
}

// GetDeviceConfigInfo 通过相机ID获取设备的配置信息（有锁）
//
//	@param	id	相机ID
//	@return	设备配置信息
//	@return	异常信息
func (p *Control) GetDeviceConfigInfo(id string) (camera.DeviceConfigList, error) {
	// 获取设备
	dev, err := p.GetDeviceWithID(id)
	if err != nil {
		return nil, err
	}

	// 加读锁
	p.rwmutex.RLock()
	defer p.rwmutex.RUnlock()

	// 调用内部方法获取
	return p.getDeviceConfigInfo(dev.SymbolicLink)
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

// Open 打开相机
//
//	@param	id		相机ID
//	@param	info	分辨率信息
//	@return	异常信息
func (p *Control) Open(id string, info camera.DeviceConfig) error {
	// 操作尝试加锁
	ok := p.rwmutex.TryLock()
	if !ok {
		return camera.ErrDeviceRepeatOpening
	}
	defer p.rwmutex.Unlock()

	// 相机列表为空时获取相机列表
	if len(p.deviceCacheList) == 0 {
		// 获取相机列表（必须使用无锁）
		_, err := p.getList()
		if err != nil {
			return err
		}
	}

	// 查询ID对应的相机信息
	cameraInfo, err := p.deviceCacheList.Get(id)
	if err != nil {
		return err
	}

	// 确认输入的配置是否在列表中
	configList, err := p.getDeviceConfigInfo(cameraInfo.SymbolicLink)
	if err != nil {
		return err
	}
	yesInfo, err := configList.Get(info)
	if err != nil {
		return err
	}

	// 关闭已打开的相机
	p.close()

	// 转换设备路径
	devicePath := C.CString(cameraInfo.SymbolicLink)
	defer C.free(unsafe.Pointer(devicePath))
	// 转换配置信息
	var frameInfo C.VideoFrameInfo
	frameInfo.width = C.uint32_t(yesInfo.Width)
	frameInfo.height = C.uint32_t(yesInfo.Height)
	frameInfo.fps = C.uint32_t(yesInfo.FPS)
	frameInfo.format = C.uint32_t(yesInfo.Format.Number())
	// 执行打开
	code := C.BecamOpenDevice(p.handle, devicePath, &frameInfo)
	if err := convertStatusCode(code); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return errors.Join(camera.ErrDeviceOpenFailed, err)
	}

	// 赋值当前使用的相机信息
	p.deviceInfo = *cameraInfo
	p.deviceSupportInfo = *yesInfo

	// 尝试获取帧
	_, _, err = p.tryGetFrame()
	return err
}

// GetStream 获取帧
//
//	@return	帧数据
//	@return	帧信息
//	@return	异常信息
func (p *Control) GetFrame() ([]byte, *camera.DeviceConfig, error) {
	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 检查相机是否已打开
	if p.handle == nil {
		return nil, nil, camera.ErrDeviceNotOpen
	}

	// 尝试获取帧
	return p.tryGetFrame()
}

// 关闭已打开的相机（无锁）
func (p *Control) close() {
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

// Close 关闭已打开的相机
func (p *Control) Close() {
	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 调用内部实现
	p.close()
}

// 释放所有相机资源
func (p *Control) Free() {
	// 关闭已打开的相机
	p.Close()

	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 是否需要释放句柄
	if p.handle != nil {
		// 释放句柄
		C.BecamFree(&p.handle)
		p.handle = nil
	}
	p.deviceCacheList = nil
	p.deviceInfo = camera.Device{}
	p.deviceSupportInfo = camera.DeviceConfig{}
}
