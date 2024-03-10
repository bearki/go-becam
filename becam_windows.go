package becam

import (
	"github.com/bearki/becam/camera"
	"github.com/bearki/becam/dshow"
)

func New() camera.Manager {
	return dshow.NewControl()
}
