package camera

// DeviceList 相机设备列表
type DeviceList []*Device

func (s DeviceList) Clone() DeviceList {
	if len(s) == 0 {
		return nil
	}
	res := make(DeviceList, 0, len(s))
	for _, v := range s {
		res = append(res, v.Clone())
	}
	return res
}

func (s DeviceList) Get(id string) (*Device, error) {
	// 遍历全部相机信息
	for _, item := range s {
		if item.ID != id {
			continue
		}
		// 返回匹配到的相机
		return item.Clone(), nil
	}

	// 默认为找不到匹配的相机
	return nil, ErrDeviceNotFound
}
