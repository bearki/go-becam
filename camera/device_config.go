package camera

// DeviceConfig 相机配置信息
type DeviceConfig struct {
	Width  uint32 // 相机支持的分辨率宽度
	Height uint32 // 相机支持的分辨率高度
	FPS    uint32 // 相机在该分辨率下支持的帧率
}

func NewDeviceConfig(w, h, fps uint32) DeviceConfig {
	return DeviceConfig{
		Width:  w,
		Height: h,
		FPS:    fps,
	}
}

func (p *DeviceConfig) Clone() *DeviceConfig {
	if p == nil {
		return nil
	}
	return &DeviceConfig{
		Width:  p.Width,
		Height: p.Height,
		FPS:    p.FPS,
	}
}

func (p *DeviceConfig) Eq(v *DeviceConfig) bool {
	if p == nil && v == nil {
		return true
	}
	if p == nil || v == nil {
		return false
	}
	return p.Width == v.Width &&
		p.Height == v.Height &&
		p.FPS == v.FPS
}

var (
	// DefaultDeviceConfig 默认设备配置信息
	DefaultDeviceConfig = NewDeviceConfig(1280, 720, 60)
	// DefaultDeviceConfig 自动设备配置信息
	AutoDeviceConfig = NewDeviceConfig(0, 0, 0)
)
