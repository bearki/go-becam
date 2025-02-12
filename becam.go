package becam

import (
	"github.com/bearki/go-becam/camera"
	"github.com/bearki/go-becam/internal"
)

// New 创建相机管理器
func New() camera.Manager {
	return internal.New()
}
