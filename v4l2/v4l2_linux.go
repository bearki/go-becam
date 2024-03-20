package v4l2

// #include <linux/videodev2.h>
import "C"
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/bearki/go-becam/camera"
	"github.com/blackjack/webcam"
)

// AsyncThreadResult 异步线程响应结果
type AsyncThreadResult struct {
	Err   error  // 异常信息
	Bytes []byte // 拷贝后的图像流
}

// Control 相机控制器
type Control struct {
	rwmutex            sync.RWMutex            // 读写锁
	deviceCacheList    camera.DeviceList       // 缓存的相机信息列表
	deviceHandle       *webcam.Webcam          // 当前使用的设备
	deviceInfo         camera.Device           // 当前使用的相机信息
	deviceSupportInfo  camera.DeviceConfig     // 当前使用的相机支持信息
	flushStreamWG      sync.WaitGroup          // 刷流线程等待组
	streamCacheChannel chan *AsyncThreadResult // 相机流缓存管道（必须为无缓冲区的管道）
	closeSignal        context.Context         // 相机关闭请求信号
	closeSignalFunc    context.CancelFunc      // 相机关闭请求方法
}

// New 创建一个相机控制器
func New() camera.Manager {
	return &Control{
		streamCacheChannel: make(chan *AsyncThreadResult),
	}
}

// 获取相机信息
//
//	@param	path	相机系统路径
//	@return	设备信息
//	@return	异常信息
func getDeviceInfo(path string) (*camera.Device, error) {
	// 打开设备
	webcamCam, err := webcam.Open(path)
	if err != nil {
		return nil, err
	}
	defer webcamCam.Close()

	// 检查设备是否支持mjpeg
	formats := webcamCam.GetSupportedFormats()
	if _, ok := formats[C.V4L2_PIX_FMT_MJPEG]; !ok {
		return nil, camera.ErrDeviceUnsupportMjpegFormat
	}

	// 获取相机信息
	name, _ := webcamCam.GetName()       // 相机名称
	busInfo, _ := webcamCam.GetBusInfo() // 相机USB信息
	// 构建相机基础信息
	device := &camera.Device{
		ID:           path,
		Name:         name,
		LocationInfo: busInfo,
		SymbolicLink: path,
		ConfigList: camera.DeviceConfigList{
			camera.AutoDeviceConfig.Clone(), // 预制一个自动分辨率在顶部
		},
	}

	// 过滤重复的相机信息
	filter := make(map[string]*camera.DeviceConfig)
	// 设备配置信息列表
	deviceConfigList := make(camera.DeviceConfigList, 0, 30)
	// 提取MJPEG支持的所有分辨率
	sizes := webcamCam.GetSupportedFrameSizes(C.V4L2_PIX_FMT_MJPEG)
	for _, item := range sizes {
		// 过滤width小于600px的分辨率
		if item.MaxWidth < 600 {
			continue
		}
		// 获取分辨率支持的帧率
		rates := webcamCam.GetSupportedFramerates(C.V4L2_PIX_FMT_MJPEG, item.MaxWidth, item.MaxHeight)
		for _, v := range rates {
			// 丢弃20帧以下的帧率
			if v.MaxDenominator/v.MaxNumerator < 20 {
				continue
			}
			// 是否已存在
			key := fmt.Sprintf("%d-%d-%d", item.MaxWidth, item.MaxHeight, v.MaxDenominator/v.MaxNumerator)
			if _, ok := filter[key]; ok {
				continue
			}
			// 构建设备配置信息
			deviceConf := &camera.DeviceConfig{
				Width:  item.MaxWidth,
				Height: item.MaxHeight,
				FPS:    v.MaxDenominator / v.MaxNumerator,
			}
			// 缓存设备配置信息
			filter[key] = deviceConf
			// 追加支持信息
			deviceConfigList = append(deviceConfigList, deviceConf)
		}
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

	// 追加分辨率信息到设备
	device.ConfigList = append(device.ConfigList, deviceConfigList...)

	// OK
	return device, nil
}

// 发现（查找）设备
//
//	@param	list		已缓存的相机列表
//	@param	discovered	重复设备过滤器
//	@param	pattern		设备查找范围
//	@param	新的相机缓存列表
func discover(list camera.DeviceList, discovered map[string]struct{}, pattern string) camera.DeviceList {
	// 匹配设备列表
	devices, err := filepath.Glob(pattern)
	if err != nil {
		// 没有匹配到，返回原列表
		return list
	}

	// 遍历匹配到的设备路径
	for _, device := range devices {
		// 提取设备真实链接文件名
		reallinkBaseName := filepath.Base(device)
		// 尝试获取设备的真实路径（防止路径是软链）
		reallink, err := os.Readlink(device)
		if err == nil {
			// 再次赋值真实链接文件名
			reallinkBaseName = filepath.Base(reallink)
		}
		// 检查设备是否重复
		if _, ok := discovered[reallinkBaseName]; ok {
			continue
		}
		// 缓存设备到过滤器
		discovered[reallinkBaseName] = struct{}{}
		// 获取设备配置信息
		devicess, err := getDeviceInfo(device)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		// 追加到设备列表
		list = append(list, devicess)
	}

	// OK
	return list
}

// --------------------------------------------- 实现内部接口 --------------------------------------------- //

// 拷贝帧
func (p *Control) copyFrame(handle *webcam.Webcam) ([]byte, error) {
	// 读取一帧
	buf, index, err := handle.GetFrame()
	// 取流是否异常
	if err != nil {
		return nil, errors.Join(camera.ErrGetFrameFailed, err)
	}
	// 延迟释放
	defer handle.ReleaseFrame(index)

	// 拷贝流
	copyBuf := make([]byte, len(buf))
	wn := copy(copyBuf, buf)
	if wn != len(buf) {
		return nil, camera.ErrCopyFrameFailed
	}

	// OK
	return copyBuf, nil
}

// 刷新相机帧
func (p *Control) flushFrame(handle *webcam.Webcam) {
	// 捕获异常
	defer func() {
		if e := recover(); e != nil {
			// 打印异常
			fmt.Fprintln(os.Stderr, "FlushStream Panic: ", e)
		}
		// 移除一个组员
		p.flushStreamWG.Done()
	}()

	// 死循环刷流
	for {
		// 等待帧
		err := handle.WaitForFrame(5)
		// 等待信号
		select {
		// 请求关闭相机
		case <-p.closeSignal.Done():
			// 结束循环
			return

		// 默认处理方式
		default:
			// 预声明
			var reply *AsyncThreadResult
			// 取流是否异常
			if err != nil {
				// 包装异常
				err = errors.Join(camera.ErrWaitForFrameFailed, err)
				// 打印异常
				fmt.Fprintln(os.Stderr, err.Error())
				// 构建响应
				reply = &AsyncThreadResult{
					Err:   err,
					Bytes: nil,
				}
			} else {
				// 拷贝流
				data, err := p.copyFrame(handle)
				// 构建响应
				reply = &AsyncThreadResult{
					Err:   err,
					Bytes: data,
				}
			}

			// 管道处理
			select {
			// 尝试发送
			case p.streamCacheChannel <- reply:
				// 发送成功啥也不做
			// 发送失败直接丢弃
			default:
				// 默认丢弃
			}
		}
	}
}

// 尝试获取帧
func (p *Control) tryGetFrame(parseW, parseH *uint32) ([]byte, error) {
	// 执行取流，并做好取流失败重试准备
	for i := 1; i <= camera.GetFrameRetryCount; i++ {
		// 等待取流完成或者取流超时
		select {
		// 取流完成信号
		case res := <-p.streamCacheChannel:
			// 是否异常
			if res.Err != nil {
				// 是否需要尝试再次取流
				if errors.Is(res.Err, new(webcam.Timeout)) && i < camera.GetFrameRetryCount {
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

// 释放所有相机资源
func (p *Control) Free() {
	// 先关闭已打开的相机
	p.Close()
	// 加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()
	// 执行资源释放
	p.deviceCacheList = nil
	close(p.streamCacheChannel)
}

// GetList 获取相机列表
//
//	@return 相机列表
//	@return 错误信息
func (p *Control) GetList() (camera.DeviceList, error) {
	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 清除缓存
	p.deviceCacheList = nil
	// 过滤重复相机
	discovered := make(map[string]struct{})
	// 查询设备列表
	p.deviceCacheList = discover(p.deviceCacheList, discovered, "/dev/v4l/by-id/*")
	p.deviceCacheList = discover(p.deviceCacheList, discovered, "/dev/v4l/by-path/*")
	p.deviceCacheList = discover(p.deviceCacheList, discovered, "/dev/video*")

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
	// 关闭已打开的相机
	p.Close()

	// 操作加锁
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

	// 查询ID对应的相机信息
	cameraInfo, err := p.deviceCacheList.Get(id)
	if err != nil {
		return err
	}

	// 打开新的相机
	tmpHandle, err := webcam.Open(cameraInfo.SymbolicLink)
	if err != nil {
		return errors.Join(camera.ErrDeviceOpenFailed, err)
	}
	defer func() {
		// 后续操作是否存在异常
		if err != nil {
			// 关闭相机
			tmpHandle.Close()
		}
	}()

	// 是否需要修改分辨率
	if info.Width > 0 && info.Height > 0 && info.FPS > 0 {
		// 筛选分辨率
		sInfo, err := cameraInfo.ConfigList.Get(info)
		if err != nil {
			return err
		}
		// 赋值选择的分辨率
		_, w, h, err := tmpHandle.SetImageFormat(C.V4L2_PIX_FMT_MJPEG, sInfo.Width, sInfo.Height)
		if err != nil {
			return errors.Join(camera.ErrSetDeviceConfigFailed, err)
		}
		err = tmpHandle.SetFramerate(float32(sInfo.FPS))
		if err != nil {
			return errors.Join(camera.ErrSetDeviceConfigFailed, err)
		}
		fps, err := tmpHandle.GetFramerate()
		if err != nil {
			return errors.Join(camera.ErrGetDeviceConfigFailed, err)
		}
		// 赋值配置后的结果
		info = camera.DeviceConfig{
			Width:  w,
			Height: h,
			FPS:    uint32(fps),
		}
	} else if len(cameraInfo.ConfigList) > 0 {
		// 选中与默认分辨率最相似的分辨率
		sInfo, err := cameraInfo.ConfigList.GetMostSimilar(camera.DefaultDeviceConfig)
		if err != nil {
			return err
		}
		// 赋值选择的分辨率
		_, w, h, err := tmpHandle.SetImageFormat(C.V4L2_PIX_FMT_MJPEG, sInfo.Width, sInfo.Height)
		if err != nil {
			return errors.Join(camera.ErrSetDeviceConfigFailed, err)
		}
		err = tmpHandle.SetFramerate(float32(sInfo.FPS))
		if err != nil {
			return errors.Join(camera.ErrSetDeviceConfigFailed, err)
		}
		fps, err := tmpHandle.GetFramerate()
		if err != nil {
			return errors.Join(camera.ErrGetDeviceConfigFailed, err)
		}
		// 赋值配置后的结果
		info = camera.DeviceConfig{
			Width:  w,
			Height: h,
			FPS:    uint32(fps),
		}
	}

	// 开启取流线程
	err = tmpHandle.StartStreaming()
	if err != nil {
		return errors.Join(camera.ErrRunStreamingFailed, err)
	}

	// 赋值新的相机关闭信号管道
	p.closeSignal, p.closeSignalFunc = context.WithCancel(context.Background())
	// 赋值句柄
	p.deviceHandle = tmpHandle
	// 赋值当前使用的相机信息
	p.deviceInfo = *cameraInfo
	// 忘刷流线程分组增加一个组员
	p.flushStreamWG.Add(1)
	// 异步循环刷新缓冲区
	go p.flushFrame(tmpHandle)

	// 获取一帧图像，提取其分辨率
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
	if p.deviceHandle == nil || p.closeSignal == nil || p.closeSignalFunc == nil {
		return
	}

	// 通知需要关闭相机
	p.closeSignalFunc()
	// 等待所有刷流线程关闭
	p.flushStreamWG.Wait()

	// 关闭已打开的相机
	err := p.deviceHandle.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	// 给内核一点时间
	time.Sleep(time.Millisecond * 100)
	// 清空句柄
	p.deviceHandle = nil
	// 清除当前使用的相机信息
	p.deviceInfo = camera.Device{}
	p.deviceSupportInfo = camera.DeviceConfig{}
}
