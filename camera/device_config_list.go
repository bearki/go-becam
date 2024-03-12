package camera

import (
	"sort"
)

// DeviceConfigList 相机配置信息列表
type DeviceConfigList []*DeviceConfig

func (s DeviceConfigList) Clone() DeviceConfigList {
	res := make(DeviceConfigList, 0, len(s))
	for _, v := range s {
		res = append(res, v.Clone())
	}
	return res
}

func (s DeviceConfigList) Get(val DeviceConfig) (*DeviceConfig, error) {
	for _, v := range s {
		if v.Width == val.Width && v.Height == val.Height && v.FPS == val.FPS {
			return v.Clone(), nil
		}
	}

	// 默认未查询到分辨率
	return nil, ErrDeviceConfigNotFound
}

// GetMostSimilar 查找与目标配置信息最相似的配置信息
func (s DeviceConfigList) GetMostSimilar(val DeviceConfig) (*DeviceConfig, error) {
	// 从列表中过滤掉自动分辨率
	vailds := make(DeviceConfigList, 0, len(s))
	for _, v := range s {
		if !v.Eq(AutoDeviceConfig.Clone()) {
			vailds = append(vailds, v)
		}
	}

	// 只有一个分辨率时直接返回
	if len(vailds) == 1 {
		return vailds[0].Clone(), nil
	}

	// 精准查询
	res, err := vailds.Get(val)
	if err == nil {
		return res, nil
	}

	// 拷贝一份用于排序
	tmpS := vailds.Clone()
	// 大于目标的支持信息
	var gtList DeviceConfigList
	// 分辨率一致帧率不一致
	var eqList DeviceConfigList
	// 小于目标的支持信息
	var ltList DeviceConfigList
	// 筛选分辨率
	for _, v := range tmpS {
		// 过滤掉默认分辨率
		if v.Eq(&DefaultDeviceConfig) {
			continue
		}
		if v.Width == val.Width && v.Height == val.Height { // 分辨率一致
			eqList = append(eqList, v)
		} else if v.Width < val.Width || v.Height < val.Height { // 分辨率小于目标分辨率
			ltList = append(ltList, v)
		} else { // 分辨率大于目标分辨率
			gtList = append(gtList, v)
		}
	}

	// 找到相同分辨率的不同帧率时优先返回
	if len(eqList) > 0 {
		return eqList[0].Clone(), nil
	}

	// 优先返回分辨率高于目标分辨率的最小分辨率的最高帧率
	if len(gtList) > 0 {
		// 对支持信息进行排序
		sort.Slice(gtList, func(i, j int) bool {
			// 宽度是否相等
			if gtList[i].Width == gtList[j].Width {
				// 高度是否相等
				if gtList[i].Height == gtList[j].Height {
					// 按帧率从小到大排序
					return gtList[i].FPS < gtList[j].FPS
				}
				// 按高度从大到小排序
				return gtList[i].Height > gtList[j].Height
			}
			// 按宽度从大到小排序
			return gtList[i].Width > gtList[j].Width
		})
		// 返回结果
		return gtList[len(gtList)-1].Clone(), nil
	}

	// 优先返回分辨率低于目标分辨率的最大分辨率的最高帧率
	if len(ltList) > 0 {
		return ltList[0].Clone(), nil
	}

	// 默认未查询到分辨率
	return nil, ErrDeviceConfigNotFound
}
