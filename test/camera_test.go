package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bearki/becam"
	"github.com/bearki/becam/camera"
)

func TestCamera(t *testing.T) {
	cameraManage := becam.New()

	list, err := cameraManage.GetList()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("未找到相机")
	}

	var index int = 0
	var info *camera.DeviceConfig
	for i, v := range list {
		fmt.Printf("%d. %s (%s)\n", i+1, v.ID, v.Name)
		for j, w := range v.ConfigList {
			fmt.Printf("\t%d. %d*%dp (%d)\n", j+1, w.Width, w.Height, w.FPS)
			if info == nil && w.Width == 1920 {
				index = i
				info = &camera.DeviceConfig{
					Width:  w.Width,
					Height: w.Height,
					FPS:    w.FPS,
				}
			}
		}
	}

	err = cameraManage.Open(list[index].ID, *info)
	if err != nil {
		t.Fatal(err)
	}
	defer cameraManage.Close()

	now := time.Now()

	var w, h uint32
	var img []byte
	for i := 0; i < 10000000; i++ {
		img, err = cameraManage.GetFrame(&w, &h)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile("test.jpg", img, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Logf("图像分辨率：%d*%dpx，实际帧率：%d", w, h, time.Since(now).Milliseconds()/1000)
}
