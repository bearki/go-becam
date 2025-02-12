package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bearki/go-becam"
	"github.com/bearki/go-becam/camera"
)

func main() {
	cameraManage := becam.New()
	defer cameraManage.Free()

	for {
		func() {
			list, err := cameraManage.GetList()
			if err != nil {
				log.Fatal(err)
			}

			if len(list) == 0 {
				log.Fatal("未找到相机")
			}

			var id string = ""
			var info *camera.DeviceConfig
			for i, v := range list {
				fmt.Printf("%d. %s (%s)\n", i+1, v.ID, v.Name)
				// 获取支持的分辨率
				cfgList, err := cameraManage.GetDeviceConfigInfo(v.ID)
				if err != nil {
					log.Fatal(err)
				}
				for j, w := range cfgList {
					fmt.Printf("\t%d. %d*%dp (%d)\n", j+1, w.Width, w.Height, w.FPS)
					if info == nil && w.Width == 1920 {
						id = v.ID
						info = w
					}
				}
			}

			err = cameraManage.Open(id, *info)
			if err != nil {
				log.Fatal(err)
			}
			defer cameraManage.Close()

			now := time.Now()

			var img []byte
			var imgInfo *camera.DeviceConfig
			for i := 0; i < 100; i++ {
				img, imgInfo, err = cameraManage.GetFrame()
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Size: %d, PX: %d*%d Fotmat: %s\n", len(img), imgInfo.Width, imgInfo.Height, imgInfo.Format)
					err = os.WriteFile("test."+imgInfo.Format.String(), img, 0644)
					if err != nil {
						log.Println(err)
					}
				}
			}

			log.Printf("图像分辨率：%s %d*%dpx，实际帧率：%d\n", imgInfo.Format, imgInfo.Width, imgInfo.Height, 1000/time.Since(now).Milliseconds()/100)
		}()
	}

}
