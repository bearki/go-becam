package becam

import (
	"github.com/bearki/go-becam/camera"
	"github.com/bearki/go-becam/dshow"
)

// New 创建相机管理器
func New() camera.Manager {
	return dshow.New()
}
