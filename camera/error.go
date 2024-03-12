package camera

import (
	"fmt"

	goi18n "github.com/bearki/go-i18n/v2"
)

// 相机异常错误号
type errno uint

const (
	_                             errno = iota // 占位
	ErrDeviceUnsupportMjpegFormat              // 设备不支持MJPEG格式
	ErrDeviceNotFound                          // 设备未找到
	ErrWaitForFrameFailed                      // 等待帧失败
	ErrGetFrameFailed                          // 获取帧失败
	ErrCopyFrameFailed                         // 拷贝帧失败
	ErrGetFrameTimout                          // 获取帧超时
	ErrDeviceOpenFailed                        // 设备打开失败
	ErrDeviceConfigNotFound                    // 设备配置未找到
	ErrSetDeviceConfigFailed                   // 修改设备配置失败
	ErrGetDeviceConfigFailed                   // 获取设备配置失败
	ErrRunStreamingFailed                      // 运行取流线程失败
	ErrDecodeJpegImageFailed                   // 解码JPEG图像失败
	ErrDeviceNotOpen                           // 设备未打开
)

// 错误码描述映射
var errMap = map[errno]map[goi18n.Code]string{
	ErrDeviceUnsupportMjpegFormat: {
		goi18n.ZH_CN: "设备不支持MJPEG格式",
		goi18n.ZH_TW: "裝置不支援MJPEG格式",
		goi18n.ZH_HK: "裝置不支援MJPEG格式",
		goi18n.EN_HK: "Device does not support MJPEG format",
		goi18n.EN_US: "Device does not support MJPEG format",
		goi18n.EN_GB: "Device does not support MJPEG format",
		goi18n.EN_WW: "Device does not support MJPEG format",
		goi18n.EN_CA: "Device does not support MJPEG format",
		goi18n.EN_AU: "Device does not support MJPEG format",
		goi18n.EN_IE: "Device does not support MJPEG format",
		goi18n.EN_FI: "Laitteen ei tueta MJPEG-muotoa",
		goi18n.FI_FI: "Laitteella ei ole tukea MJPEG-muodolle",
		goi18n.EN_DK: "Enhed understøtter ikke MJPEG-format",
		goi18n.DA_DK: "Enheden understøtter ikke MJPEG-format",
		goi18n.EN_IL: "The device does not support the MJPEG format",
		goi18n.HE_IL: "הустройство לא תומך בפורמט MJPEG",
		goi18n.EN_ZA: "Device does not support MJPEG format",
		goi18n.EN_IN: "Device does not support MJPEG format",
		goi18n.EN_NO: "Enheten støtter ikke MJPEG-formatet",
		goi18n.EN_SG: "Device does not support MJPEG format",
		goi18n.EN_NZ: "Device does not support MJPEG format",
		goi18n.EN_ID: "Perangkat tidak mendukung format MJPEG",
		goi18n.EN_PH: "Ang device ay hindi sumusuporta sa format ng MJPEG",
		goi18n.EN_TH: "อุปกรณ์ไม่รองรับรูปแบบ MJPEG",
		goi18n.EN_MY: "Peralatan tidak menyokong format MJPEG",
		goi18n.EN_XA: "Device does not support MJPEG format",
		goi18n.KO_KR: "장치가 MJPEG 형식을 지원하지 않습니다.",
		goi18n.JA_JP: "デバイスはMJPEG形式をサポートしていません",
		goi18n.NL_NL: "Apparaat ondersteunt geen MJPEG-formaat",
		goi18n.NL_BE: "Toestel biedt geen ondersteuning voor het MJPEG-formaat",
		goi18n.PT_PT: "Dispositivo não suporta formato MJPEG",
		goi18n.PT_BR: "Dispositivo não suporta o formato MJPEG",
		goi18n.FR_FR: "Le dispositif ne prend pas en charge le format MJPEG",
		goi18n.FR_LU: "Le dispositif ne prend pas en charge le format MJPEG",
		goi18n.FR_CH: "Le dispositif ne prend pas en charge le format MJPEG",
		goi18n.FR_BE: "Le dispositif ne prend pas en charge le format MJPEG",
		goi18n.FR_CA: "Le périphérique ne prend pas en charge le format MJPEG",
		goi18n.ES_LA: "El dispositivo no admite el formato MJPEG",
		goi18n.ES_ES: "El dispositivo no admite el formato MJPEG",
		goi18n.ES_AR: "El dispositivo no soporta el formato MJPEG",
		goi18n.ES_US: "El dispositivo no admite el formato MJPEG",
		goi18n.ES_MX: "El dispositivo no admite el formato MJPEG",
		goi18n.ES_CO: "El dispositivo no admite el formato MJPEG",
		goi18n.ES_PR: "El dispositivo no admite el formato MJPEG",
		goi18n.DE_DE: "Gerät unterstützt kein MJPEG-Format",
		goi18n.DE_AT: "Gerät unterstützt kein MJPEG-Format",
		goi18n.DE_CH: "Gerät unterstützt kein MJPEG-Format",
		goi18n.RU_RU: "Устройство не поддерживает формат MJPEG",
		goi18n.IT_IT: "Il dispositivo non supporta il formato MJPEG",
		goi18n.EL_GR: "Ο συσκευαστής δεν υποστηρίζει τη μορφή MJPEG",
		goi18n.NO_NO: "Enheten støtter ikke MJPEG-format",
		goi18n.HU_HU: "Az eszköz nem támogatja a MJPEG formátumot",
		goi18n.TR_TR: "Cihaz MJPEG formatını desteklemiyor",
		goi18n.CS_CZ: "Zařízení nepodporuje formát MJPEG",
		goi18n.SL_SL: "Naprava ne podpira format MJPEG",
		goi18n.PL_PL: "Urządzenie nie obsługuje formatu MJPEG",
		goi18n.SV_SE: "Enheten har inte stöd för MJPEG-formatet",
	},
	ErrDeviceNotFound: {
		goi18n.ZH_CN: "设备未找到",
		goi18n.ZH_TW: "裝置未找到",
		goi18n.ZH_HK: "裝置未找到",
		goi18n.EN_HK: "Device not found",
		goi18n.EN_US: "Device not found",
		goi18n.EN_GB: "Device not found",
		goi18n.EN_WW: "Device not found",
		goi18n.EN_CA: "Device not found",
		goi18n.EN_AU: "Device not found",
		goi18n.EN_IE: "Device not found",
		goi18n.EN_FI: "Laitetta ei löytynyt",
		goi18n.FI_FI: "Laitetta ei löydy",
		goi18n.EN_DK: "Enhed ikke fundet",
		goi18n.DA_DK: "Enheden blev ikke fundet",
		goi18n.EN_IL: "The device was not found",
		goi18n.HE_IL: "הустройство לא נמצא",
		goi18n.EN_ZA: "Device not found",
		goi18n.EN_IN: "Device not found",
		goi18n.EN_NO: "Enheten ble ikke funnet",
		goi18n.EN_SG: "Device not found",
		goi18n.EN_NZ: "Device not found",
		goi18n.EN_ID: "Perangkat tidak ditemukan",
		goi18n.EN_PH: "Ang device ay hindi nakita",
		goi18n.EN_TH: "ไม่พบอุปกรณ์",
		goi18n.EN_MY: "Peralatan tidak dijumpai",
		goi18n.EN_XA: "Device not found",
		goi18n.KO_KR: "장치를 찾을 수 없습니다.",
		goi18n.JA_JP: "デバイスが見つかりません",
		goi18n.NL_NL: "Apparaat niet gevonden",
		goi18n.NL_BE: "Toestel niet gevonden",
		goi18n.PT_PT: "Dispositivo não encontrado",
		goi18n.PT_BR: "Dispositivo não encontrado",
		goi18n.FR_FR: "Le dispositif est introuvable",
		goi18n.FR_LU: "Le dispositif est introuvable",
		goi18n.FR_CH: "Le dispositif est introuvable",
		goi18n.FR_BE: "Le dispositif est introuvable",
		goi18n.FR_CA: "Le périphérique est introuvable",
		goi18n.ES_LA: "No se encontró el dispositivo",
		goi18n.ES_ES: "Dispositivo no encontrado",
		goi18n.ES_AR: "Dispositivo no encontrado",
		goi18n.ES_US: "Dispositivo no encontrado",
		goi18n.ES_MX: "Dispositivo no encontrado",
		goi18n.ES_CO: "Dispositivo no encontrado",
		goi18n.ES_PR: "Dispositivo no encontrado",
		goi18n.DE_DE: "Gerät nicht gefunden",
		goi18n.DE_AT: "Gerät nicht gefunden",
		goi18n.DE_CH: "Gerät nicht gefunden",
		goi18n.RU_RU: "Устройство не найдено",
		goi18n.IT_IT: "Dispositivo non trovato",
		goi18n.EL_GR: "Ο συσκευαστής δεν βρέθηκε",
		goi18n.NO_NO: "Enheten ble ikke funnet",
		goi18n.HU_HU: "Az eszköz nem található",
		goi18n.TR_TR: "Cihaz bulunamadı",
		goi18n.CS_CZ: "Zařízení nebylo nalezeno",
		goi18n.SL_SL: "Naprave ni bilo mogoče najti",
		goi18n.PL_PL: "Urządzenie nie zostało znalezione",
		goi18n.SV_SE: "Enheten kunde inte hittas",
	},
}

func (e errno) Error() string {
	// 是否存在该错误
	errLang, ok := errMap[e]
	if !ok {
		return fmt.Sprintf("unknown becam errno: %d", e)
	} else if errLang == nil {
		return fmt.Sprintf("raw becam errno: %d", e)
	}
	// 是否存在对应语言
	errStr, ok := errLang[goi18n.GetEnv()]
	if ok {
		return errStr
	}
	// 优先使用英文
	errStr, ok = errLang[goi18n.EN_US]
	if ok {
		return errStr
	}
	// 次优先使用中文
	errStr, ok = errLang[goi18n.ZH_CN]
	if ok {
		return errStr
	}
	// 使用随机语言
	for _, v := range errLang {
		return v
	}
	// 不存在语言时
	return fmt.Sprintf("raw becam errno: %d", e)
}
