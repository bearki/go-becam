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
//
//	@param	target			期望目标配置
//	@param	standbyTargets	备用期望目标配置（依次查询）
func (s DeviceConfigList) GetMostSimilar(target DeviceConfig, standbyTargets []DeviceConfig) (*DeviceConfig, error) {
	// 只有一个分辨率时直接返回
	if len(s) == 1 {
		return s[0].Clone(), nil
	}

	// 精准查询
	res, err := s.Get(target)
	if err == nil {
		return res, nil
	}

	// 备用格式转MAP
	standbyTargetMap := make(map[Fourcc]DeviceConfig)
	for _, v := range standbyTargets {
		standbyTargetMap[v.Format] = v
	}

	// 拷贝一份用于排序
	tmpS := s.Clone()

	// 格式分离映射表
	// Map映射的三个列表
	const EQ = 0 // 分辨率一致帧率不一致
	const GT = 1 // 大于目标的支持信息
	const LT = 2 // 小于目标的支持信息
	fmtSplitMap := make(map[Fourcc]map[int]DeviceConfigList)
	for _, v1 := range tmpS {
		// 格式是否一致
		if v1.Format == target.Format {
			// 格式映射是否存在
			if _, ok := fmtSplitMap[v1.Format]; !ok {
				fmtSplitMap[v1.Format] = make(map[int]DeviceConfigList, 3)
			}
			// 对分辨率进行分离
			if v1.Width == target.Width && v1.Height == target.Height {
				// 分辨率一致
				fmtSplitMap[v1.Format][EQ] = append(fmtSplitMap[v1.Format][EQ], v1)
			} else if v1.Width > target.Width && v1.Height > target.Height {
				// 分辨率大于目标分辨率
				fmtSplitMap[v1.Format][GT] = append(fmtSplitMap[v1.Format][GT], v1)
			} else {
				// 分辨率小于目标分辨率
				fmtSplitMap[v1.Format][LT] = append(fmtSplitMap[v1.Format][LT], v1)
			}
		} else {
			// 是否在备用格式中
			if standbyTarget, ok := standbyTargetMap[v1.Format]; ok {
				// 格式映射是否存在
				if _, ok := fmtSplitMap[v1.Format]; !ok {
					fmtSplitMap[v1.Format] = make(map[int]DeviceConfigList, 3)
				}
				// 对分辨率进行分离
				if v1.Width == standbyTarget.Width && v1.Height == standbyTarget.Height {
					// 分辨率一致
					fmtSplitMap[v1.Format][EQ] = append(fmtSplitMap[v1.Format][EQ], v1)
				} else if v1.Width > standbyTarget.Width && v1.Height > standbyTarget.Height {
					// 分辨率大于目标分辨率
					fmtSplitMap[v1.Format][GT] = append(fmtSplitMap[v1.Format][GT], v1)
				} else {
					// 分辨率小于目标分辨率
					fmtSplitMap[v1.Format][LT] = append(fmtSplitMap[v1.Format][LT], v1)
				}
			}
			// 不在备用格式中就忽略吧，不支持这样的格式了
		}
	}

	// 优先处理格式一致的情况
	if targetMap, ok := fmtSplitMap[target.Format]; ok {
		// 找到相同分辨率的不同帧率时优先返回
		if len(targetMap[EQ]) > 0 {
			return targetMap[EQ][0].Clone(), nil
		}
		// 优先返回分辨率高于目标分辨率的最小分辨率的最高帧率
		if len(targetMap[GT]) > 0 {
			// 对支持信息进行排序
			sort.Slice(targetMap[GT], func(i, j int) bool {
				// 宽度是否相等
				if targetMap[GT][i].Width == targetMap[GT][j].Width {
					// 高度是否相等
					if targetMap[GT][i].Height == targetMap[GT][j].Height {
						// 按帧率从小到大排序
						return targetMap[GT][i].FPS < targetMap[GT][j].FPS
					}
					// 按高度从大到小排序
					return targetMap[GT][i].Height > targetMap[GT][j].Height
				}
				// 按宽度从大到小排序
				return targetMap[GT][i].Width > targetMap[GT][j].Width
			})
			// 返回结果
			return targetMap[GT][len(targetMap[GT])-1].Clone(), nil
		}
		// 优先返回分辨率低于目标分辨率的最大分辨率的最高帧率
		if len(targetMap[LT]) > 0 {
			return targetMap[LT][0].Clone(), nil
		}
	}

	// 再遍历备用目标格式
	for _, standbyTarget := range standbyTargets {
		// 是否存在
		if standbyTargetMap, ok := fmtSplitMap[standbyTarget.Format]; ok {
			// 找到相同分辨率的不同帧率时优先返回
			if len(standbyTargetMap[EQ]) > 0 {
				return standbyTargetMap[EQ][0].Clone(), nil
			}
			// 优先返回分辨率高于目标分辨率的最小分辨率的最高帧率
			if len(standbyTargetMap[GT]) > 0 {
				// 对支持信息进行排序
				sort.Slice(standbyTargetMap[GT], func(i, j int) bool {
					// 宽度是否相等
					if standbyTargetMap[GT][i].Width == standbyTargetMap[GT][j].Width {
						// 高度是否相等
						if standbyTargetMap[GT][i].Height == standbyTargetMap[GT][j].Height {
							// 按帧率从小到大排序
							return standbyTargetMap[GT][i].FPS < standbyTargetMap[GT][j].FPS
						}
						// 按高度从大到小排序
						return standbyTargetMap[GT][i].Height > standbyTargetMap[GT][j].Height
					}
					// 按宽度从大到小排序
					return standbyTargetMap[GT][i].Width > standbyTargetMap[GT][j].Width
				})
				// 返回结果
				return standbyTargetMap[GT][len(standbyTargetMap[GT])-1].Clone(), nil
			}
			// 优先返回分辨率低于目标分辨率的最大分辨率的最高帧率
			if len(standbyTargetMap[LT]) > 0 {
				return standbyTargetMap[LT][0].Clone(), nil
			}
		}
	}

	// 默认未查询到分辨率
	return nil, ErrDeviceMediaConfigNotFound
}
