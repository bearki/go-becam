package becam

import (
	"github.com/bearki/becam/camera"
	"github.com/bearki/becam/v4l2"
)

// New 创建相机管理器
func New() camera.Manager {
	return v4l2.New()
}
