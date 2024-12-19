package camera

import (
	"sort"
)

// DeviceConfigList 相机配置信息列表
type DeviceConfigList []*DeviceConfig

// Clone 克隆相机配置信息列表
func (s DeviceConfigList) Clone() DeviceConfigList {
	res := make(DeviceConfigList, 0, len(s))
	for _, v := range s {
		res = append(res, v.Clone())
	}
	return res
}

// Get 从列表中查询目标配置信息（通常用于检测目标配置是否存在）
func (s DeviceConfigList) Get(val DeviceConfig) (*DeviceConfig, error) {
	for _, v := range s {
		if v.Eq(&val) {
			return v.Clone(), nil
		}
	}

	// 默认未查询到
	return nil, ErrDeviceMediaConfigNotFound
}

// GetMostSimilar 查找与目标配置信息最相似的配置信息
func (s DeviceConfigList) GetMostSimilar(val DeviceConfig) (*DeviceConfig, error) {
	// 只有一个分辨率时直接返回
	if len(s) == 1 {
		return s[0].Clone(), nil
	}

	// 精准查询
	res, err := s.Get(val)
	if err == nil {
		return res, nil
	}

	// 拷贝一份用于排序
	tmpS := s.Clone()
	var (
		// 相同格式：大于目标的支持信息
		fmtGtList DeviceConfigList
		// 相同格式：分辨率一致帧率不一致
		fmtEqList DeviceConfigList
		// 相同格式：小于目标的支持信息
		fmtLtList DeviceConfigList
	)
	var (
		// 其他格式：大于目标的支持信息
		othGtList DeviceConfigList
		// 其他格式：分辨率一致帧率不一致
		othEqList DeviceConfigList
		// 其他格式：小于目标的支持信息
		othLtList DeviceConfigList
	)
	// 筛选分辨率
	for _, v := range tmpS {
		// 比较格式
		if v.Format == val.Format {
			// 格式一致的
			// 对分辨率进行比较
			if v.Width == val.Width && v.Height == val.Height { // 分辨率一致
				fmtEqList = append(fmtEqList, v)
			} else if v.Width < val.Width || v.Height < val.Height { // 分辨率小于目标分辨率
				fmtLtList = append(fmtLtList, v)
			} else { // 分辨率大于目标分辨率
				fmtGtList = append(fmtGtList, v)
			}
		} else {
			// 格式不一致的
			// 对分辨率进行比较
			if v.Width == val.Width && v.Height == val.Height { // 分辨率一致
				othEqList = append(othEqList, v)
			} else if v.Width < val.Width || v.Height < val.Height { // 分辨率小于目标分辨率
				othLtList = append(othLtList, v)
			} else { // 分辨率大于目标分辨率
				othGtList = append(othGtList, v)
			}
		}
	}

	// 优先处理格式一致的情况
	{
		// 找到相同分辨率的不同帧率时优先返回
		if len(fmtEqList) > 0 {
			return fmtEqList[0].Clone(), nil
		}
		// 优先返回分辨率高于目标分辨率的最小分辨率的最高帧率
		if len(fmtGtList) > 0 {
			// 对支持信息进行排序
			sort.Slice(fmtGtList, func(i, j int) bool {
				// 宽度是否相等
				if fmtGtList[i].Width == fmtGtList[j].Width {
					// 高度是否相等
					if fmtGtList[i].Height == fmtGtList[j].Height {
						// 按帧率从小到大排序
						return fmtGtList[i].FPS < fmtGtList[j].FPS
					}
					// 按高度从大到小排序
					return fmtGtList[i].Height > fmtGtList[j].Height
				}
				// 按宽度从大到小排序
				return fmtGtList[i].Width > fmtGtList[j].Width
			})
			// 返回结果
			return fmtGtList[len(fmtGtList)-1].Clone(), nil
		}
		// 优先返回分辨率低于目标分辨率的最大分辨率的最高帧率
		if len(fmtLtList) > 0 {
			return fmtLtList[0].Clone(), nil
		}
	}

	// 再处理格式不一致的情况
	{
		// 找到相同分辨率的不同帧率时优先返回
		if len(othEqList) > 0 {
			return othEqList[0].Clone(), nil
		}
		// 优先返回分辨率高于目标分辨率的最小分辨率的最高帧率
		if len(othGtList) > 0 {
			// 对支持信息进行排序
			sort.Slice(othGtList, func(i, j int) bool {
				// 宽度是否相等
				if othGtList[i].Width == othGtList[j].Width {
					// 高度是否相等
					if othGtList[i].Height == othGtList[j].Height {
						// 按帧率从小到大排序
						return othGtList[i].FPS < othGtList[j].FPS
					}
					// 按高度从大到小排序
					return othGtList[i].Height > othGtList[j].Height
				}
				// 按宽度从大到小排序
				return othGtList[i].Width > othGtList[j].Width
			})
			// 返回结果
			return othGtList[len(othGtList)-1].Clone(), nil
		}
		// 优先返回分辨率低于目标分辨率的最大分辨率的最高帧率
		if len(othLtList) > 0 {
			return othLtList[0].Clone(), nil
		}
	}

	// 默认未查询到分辨率
	return nil, ErrDeviceMediaConfigNotFound
}
