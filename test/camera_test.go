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

	for _, v := range list {
		fmt.Println(v.ID)
	}

	err = cameraManage.Open(list[0].ID, camera.AutoDeviceConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer cameraManage.Close()

	now := time.Now()

	var w, h uint32
	var img []byte
	for i := 0; i < 1000; i++ {
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
