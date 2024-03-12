package dshow

// #cgo LDFLAGS: -lkernel32 -lstrmiids -lole32 -loleaut32 -lquartz -lsetupapi
// #include <dshow.h>
// #include <stdint.h>
// #cgo CFLAGS: -I${SRCDIR}
// #include "dshow_windows.hpp"
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

// AsyncThreadResult 异步线程响应结果
type AsyncThreadResult struct {
	Err   error  // 异常信息
	Bytes []byte // 拷贝后的图像流
}

// 全局相机流缓存管道
var streamCacheChannel = make(chan *AsyncThreadResult)

//export imageCallback
func imageCallback(imgBuf *C.char, imgBufSize C.size_t) {
	// 检查图像
	if imgBuf == nil || imgBufSize <= 0 {
		// 发送取流结果信号
		select {
		case streamCacheChannel <- &AsyncThreadResult{Err: camera.ErrGetFrameFailed}:
		case <-time.After(time.Millisecond * 5):
		}
		return
	}
	// 拷贝图像
	copyBuf := C.GoBytes(unsafe.Pointer(imgBuf), C.int(imgBufSize))
	if len(copyBuf) != int(imgBufSize) {
		// 发送取流结果信号
		select {
		case streamCacheChannel <- &AsyncThreadResult{Err: camera.ErrCopyFrameFailed}:
		case <-time.After(time.Millisecond * 5):
		}
		return
	}

	// 发送取流结果信号
	select {
	case streamCacheChannel <- &AsyncThreadResult{Bytes: copyBuf}:
	case <-time.After(time.Millisecond * 5):
	}
}

func init() {
	C.CoInitializeEx(nil, C.COINIT_MULTITHREADED)
}

// Control 相机控制器
type Control struct {
	rwmutex           sync.RWMutex        // 读写锁
	deviceCacheList   camera.DeviceList   // 缓存的相机信息列表
	deviceHandle      *C.camera           // 当前使用的设备
	deviceInfo        camera.Device       // 当前使用的相机信息
	deviceSupportInfo camera.DeviceConfig // 当前使用的相机支持信息
}

// NewControl 创建一个相机控制器
func New() camera.Manager {
	return &Control{}
}

// 获取设备配置信息
//
//	@param	symbolicLink	设备系统路径
func getDeviceConfig(symbolicLink string) (camera.DeviceConfigList, error) {
	// 打开设备
	in := &C.camera{
		path: C.CString(symbolicLink),
	}
	defer C.free(unsafe.Pointer(in.path))

	var errStr *C.char
	if C.getResolution(in, &errStr) != 0 {
		err := errors.New(C.GoString(errStr))
		err = errors.Join(camera.ErrGetDeviceConfigFailed, err)
		return nil, err
	}
	defer C.freeResolution(in)

	// support mjpeg format
	var deviceConfigList camera.DeviceConfigList
	for i := 0; i < int(in.numProps); i++ {
		// 获取分辨率
		p := C.getProp(in, C.int(i))
		// 过滤不支持MJPEG的
		if uint32(p.fcc) != camera.V4L2_PIX_FMT_RGB332 {
			continue
		}
		// 过滤width小于600px的分辨率
		if p.width < 600 {
			continue
		}
		// 丢弃20帧以下的帧率
		if p.fps < 20 {
			continue
		}
		// 添加到支持信息中
		deviceConfigList = append(deviceConfigList, &camera.DeviceConfig{
			Width:  uint32(p.width),
			Height: uint32(p.height),
			FPS:    uint32(p.fps),
		})
	}

	// 检查配置是否为空
	if len(deviceConfigList) == 0 {
		return nil, errors.New("the device does not support MJPEG format")
	}

	// OK
	return deviceConfigList, nil
}

// --------------------------------------------- 实现内部接口 --------------------------------------------- //

// 尝试获取帧
func (p *Control) tryGetFrame(parseW, parseH *uint32) ([]byte, error) {
	// 执行取流，并做好取流失败重试准备
	for i := 1; i <= camera.GetFrameRetryCount; i++ {
		// 等待取流完成或者取流超时
		select {
		// 取流完成信号
		case res := <-streamCacheChannel:
			// 是否异常
			if res.Err != nil {
				// 是否需要尝试再次取流
				if i < camera.GetFrameRetryCount {
					continue
				}
				// 返回异常
				return nil, res.Err
			}
			// 是否需要解码图像
			if parseW != nil && parseH != nil {
				// 解码图像
				imgConf, err := jpeg.DecodeConfig(bytes.NewReader(res.Bytes))
				if err != nil {
					if i < camera.GetFrameRetryCount {
						continue
					}
					// 返回异常
					return nil, err
				}
				// 赋值宽高
				*parseW = uint32(imgConf.Width)
				*parseH = uint32(imgConf.Height)
			}
			// 获取帧成功
			return res.Bytes, nil

		// 取流超时信号
		case <-time.After(time.Millisecond * 50):
			// 继续剩余次数
			continue
		}
	}

	// 超时了
	return nil, camera.ErrGetFrameTimout
}

// --------------------------------------------- 实现Manager接口 --------------------------------------------- //

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
	var list C.cameraList
	var errStr *C.char
	if C.listCamera(&list, &errStr) != 0 {
		err := errors.New("enum camera device list error, " + C.GoString(errStr))
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, errors.Join(camera.ErrGetDeviceConfigFailed, err)
	}
	defer C.freeCameraList(&list, &errStr)

	// 遍历相机列表
	for i := 0; i < int(list.num); i++ {
		// 获取设备信息
		path := C.GoString(C.getPath(&list, C.int(i)))
		devicePath := C.GoString(C.getLocationInfo(&list, C.int(i)))
		name := C.GoString(C.getName(&list, C.int(i)))
		// 计算相机ID
		idData := md5.Sum([]byte(devicePath + name))

		// 构建设备信息
		info := &camera.Device{
			ID:           hex.EncodeToString(idData[:]),
			Name:         name,
			SymbolicLink: path,
			LocationInfo: devicePath,
			ConfigList:   camera.DeviceConfigList{camera.AutoDeviceConfig.Clone()},
		}

		// 获取设备配置信息
		deviceConfigList, err := getDeviceConfig(info.SymbolicLink)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
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

	// 构建新的相机打开信息
	tmpHandle := &C.camera{
		path: C.CString(cameraInfo.SymbolicLink),
	}

	// 是否需要修改分辨率
	if info.Width > 0 && info.Height > 0 && info.FPS > 0 {
		// 筛选分辨率
		sInfo, err := cameraInfo.ConfigList.Get(info)
		if err != nil {
			return err
		}
		// 赋值选择的分辨率
		info.Width, tmpHandle.width = sInfo.Width, C.int(sInfo.Width)
		info.Height, tmpHandle.height = sInfo.Height, C.int(sInfo.Height)
		info.FPS, tmpHandle.fps = sInfo.FPS, C.int(sInfo.FPS)
	} else if len(cameraInfo.ConfigList) > 0 {
		// 选中与默认分辨率最相似的分辨率
		sInfo, err := cameraInfo.ConfigList.GetMostSimilar(camera.DefaultDeviceConfig)
		if err != nil {
			return err
		}
		// 赋值选择的分辨率
		info.Width, tmpHandle.width = sInfo.Width, C.int(sInfo.Width)
		info.Height, tmpHandle.height = sInfo.Height, C.int(sInfo.Height)
		info.FPS, tmpHandle.fps = sInfo.FPS, C.int(sInfo.FPS)
	}

	// 执行打开
	var errStr *C.char
	res := C.openCamera(tmpHandle, &errStr)
	if res != 0 {
		C.free(unsafe.Pointer(tmpHandle.path))
		return errors.Join(camera.ErrDeviceOpenFailed, errors.New(C.GoString(errStr)))
	}

	// 缓存句柄
	p.deviceHandle = tmpHandle
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
	if p.deviceHandle == nil {
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
	if p.deviceHandle == nil {
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
	if p.deviceHandle == nil {
		return
	}

	// 释放申请的路径内存
	C.free(unsafe.Pointer(p.deviceHandle.path))
	// 释放相机内存
	C.freeCamera(p.deviceHandle)
	// 给内核一点时间
	time.Sleep(time.Millisecond * 100)
	// 清空句柄
	p.deviceHandle = nil
	// 清除当前使用的相机信息
	p.deviceInfo = camera.Device{}
	p.deviceSupportInfo = camera.DeviceConfig{}
}
