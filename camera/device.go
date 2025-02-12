package camera

// Device 相机设备信息
type Device struct {
	ID           string // 相机ID
	Name         string // 相机名称
	SymbolicLink string // 相机系统路径
}

// Clone 克隆相机设备信息
func (p *Device) Clone() *Device {
	if p == nil {
		return nil
	}
	return &Device{
		ID:           p.ID,
		Name:         p.Name,
		SymbolicLink: p.SymbolicLink,
	}
}
