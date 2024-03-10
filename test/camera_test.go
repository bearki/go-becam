package test

import (
	"fmt"
	"os"
	"testing"

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

	err = cameraManage.Open(list[0].ID, camera.DefaultDeviceConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer cameraManage.Close()

	img, err := cameraManage.GetStream()
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("test.jpg", img, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
