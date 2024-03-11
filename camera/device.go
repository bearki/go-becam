package camera

// Device 相机设备信息
type Device struct {
	ID           string           // 相机ID
	Name         string           // 相机名称
	SymbolicLink string           // 相机系统路径
	LocationInfo string           // 相机固定编码
	ConfigList   DeviceConfigList // 相机支持分辨率列表
}

func (p *Device) Clone() *Device {
	if p == nil {
		return nil
	}
	return &Device{
		ID:           p.ID,
		Name:         p.Name,
		SymbolicLink: p.SymbolicLink,
		LocationInfo: p.LocationInfo,
		ConfigList:   p.ConfigList.Clone(),
	}
}
