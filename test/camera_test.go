package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bearki/go-becam"
	"github.com/bearki/go-becam/camera"
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

	var id string = ""
	var info *camera.DeviceConfig
	for i, v := range list {
		fmt.Printf("%d. %s (%s)\n", i+1, v.ID, v.Name)
		// 获取支持的分辨率
		cfgList, err := cameraManage.GetDeviceConfigInfo(v.ID)
		if err != nil {
			t.Fatal(err)
		}
		for j, w := range cfgList {
			fmt.Printf("\t%d. %s %d*%dp (%d)\n", j+1, w.Format, w.Width, w.Height, w.FPS)
			if info == nil && w.Height == 1080 {
				id = v.ID
				info = w
			}
		}
	}

	err = cameraManage.Open(id, *info)
	if err != nil {
		t.Fatal(err)
	}
	defer cameraManage.Close()

	now := time.Now()

	var img []byte
	var imgInfo *camera.DeviceConfig
	for i := 0; i < 100; i++ {
		img, imgInfo, err = cameraManage.GetFrame()
		if err != nil {
			t.Log(err)
		} else {
			t.Logf("Size: %d, PX: %d*%d Fotmat: %s\n", len(img), imgInfo.Width, imgInfo.Height, imgInfo.Format)
			err = os.WriteFile("test."+strings.ToLower(imgInfo.Format.String()), img, 0644)
			if err != nil {
				t.Log(err)
			}
		}
	}

	t.Logf("图像分辨率：%s %d*%dpx，实际帧率：%d\n", imgInfo.Format, imgInfo.Width, imgInfo.Height, 1000/(time.Since(now).Milliseconds()/100))
}
