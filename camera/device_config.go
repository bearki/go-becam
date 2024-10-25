package camera

// Fourcc 4字节字符编码(小端序)
type Fourcc string

// NewFourccFromNumber 从数值创建FOURCC
func NewFourccFromNumber(num uint32) Fourcc {
	return Fourcc(string([]byte{byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24)}))
}

// Number 返回FOURCC的数值
func (p Fourcc) Number() uint32 {
	if len(p) != 4 {
		return 0
	}
	// 小端序
	return uint32(p[0]) | (uint32(p[1]) << 8) | (uint32(p[2]) << 16) | (uint32(p[3]) << 24)
}

// String 返回FOURCC的字符串
func (p Fourcc) String() string {
	if len(p) != 4 {
		return ""
	}
	return string(p)
}

const (
	// RGB formats (1 or 2 bytes per pixel)
	FOURCC_RGB332  = Fourcc("RGB1") //  8  RGB-3-3-2
	FOURCC_RGB444  = Fourcc("R444") // 16  xxxxrrrr ggggbbbb
	FOURCC_ARGB444 = Fourcc("AR12") // 16  aaaarrrr ggggbbbb
	FOURCC_XRGB444 = Fourcc("XR12") // 16  xxxxrrrr ggggbbbb
	FOURCC_RGBA444 = Fourcc("RA12") // 16  rrrrgggg bbbbaaaa
	FOURCC_RGBX444 = Fourcc("RX12") // 16  rrrrgggg bbbbxxxx
	FOURCC_ABGR444 = Fourcc("AB12") // 16  aaaabbbb ggggrrrr
	FOURCC_XBGR444 = Fourcc("XB12") // 16  xxxxbbbb ggggrrrr
	FOURCC_BGRA444 = Fourcc("GA12") // 16  bbbbgggg rrrraaaa
	FOURCC_BGRX444 = Fourcc("BX12") // 16  bbbbgggg rrrrxxxx
	FOURCC_RGB555  = Fourcc("RGBO") // 16  RGB-5-5-5
	FOURCC_ARGB555 = Fourcc("AR15") // 16  ARGB-1-5-5-5
	FOURCC_XRGB555 = Fourcc("XR15") // 16  XRGB-1-5-5-5
	FOURCC_RGBA555 = Fourcc("RA15") // 16  RGBA-5-5-5-1
	FOURCC_RGBX555 = Fourcc("RX15") // 16  RGBX-5-5-5-1
	FOURCC_ABGR555 = Fourcc("AB15") // 16  ABGR-1-5-5-5
	FOURCC_XBGR555 = Fourcc("XB15") // 16  XBGR-1-5-5-5
	FOURCC_BGRA555 = Fourcc("BA15") // 16  BGRA-5-5-5-1
	FOURCC_BGRX555 = Fourcc("BX15") // 16  BGRX-5-5-5-1
	FOURCC_RGB565  = Fourcc("RGBP") // 16  RGB-5-6-5
	FOURCC_RGB555X = Fourcc("RGBQ") // 16  RGB-5-5-5 BE
	// FOURCC_ARGB555X = Fourcc_be("AR15") // 16  ARGB-5-5-5 BE
	// FOURCC_XRGB555X = Fourcc_be("XR15") // 16  XRGB-5-5-5 BE
	FOURCC_RGB565X = Fourcc("RGBR") // 16  RGB-5-6-5 BE

	// RGB formats (3 or 4 bytes per pixel)
	FOURCC_BGR666 = Fourcc("BGRH") // 18  BGR-6-6-6
	FOURCC_BGR24  = Fourcc("BGR3") // 24  BGR-8-8-8
	FOURCC_RGB24  = Fourcc("RGB3") // 24  RGB-8-8-8
	FOURCC_BGR32  = Fourcc("BGR4") // 32  BGR-8-8-8-8
	FOURCC_ABGR32 = Fourcc("AR24") // 32  BGRA-8-8-8-8
	FOURCC_XBGR32 = Fourcc("XR24") // 32  BGRX-8-8-8-8
	FOURCC_BGRA32 = Fourcc("RA24") // 32  ABGR-8-8-8-8
	FOURCC_BGRX32 = Fourcc("RX24") // 32  XBGR-8-8-8-8
	FOURCC_RGB32  = Fourcc("RGB4") // 32  RGB-8-8-8-8
	FOURCC_RGBA32 = Fourcc("AB24") // 32  RGBA-8-8-8-8
	FOURCC_RGBX32 = Fourcc("XB24") // 32  RGBX-8-8-8-8
	FOURCC_ARGB32 = Fourcc("BA24") // 32  ARGB-8-8-8-8
	FOURCC_XRGB32 = Fourcc("BX24") // 32  XRGB-8-8-8-8

	// Grey formats
	FOURCC_GREY = Fourcc("GREY") //  8  Greyscale
	FOURCC_Y4   = Fourcc("Y04 ") //  4  Greyscale
	FOURCC_Y6   = Fourcc("Y06 ") //  6  Greyscale
	FOURCC_Y10  = Fourcc("Y10 ") // 10  Greyscale
	FOURCC_Y12  = Fourcc("Y12 ") // 12  Greyscale
	FOURCC_Y14  = Fourcc("Y14 ") // 14  Greyscale
	FOURCC_Y16  = Fourcc("Y16 ") // 16  Greyscale
	// FOURCC_Y16_BE = Fourcc_be("Y16 ") // 16  Greyscale BE

	// Grey bit-packed formats
	FOURCC_Y10BPACK = Fourcc("Y10B") // 10  Greyscale bit-packed
	FOURCC_Y10P     = Fourcc("Y10P") // 10  Greyscale, MIPI RAW10 packed

	// Palette formats
	FOURCC_PAL8 = Fourcc("PAL8") //  8  8-bit palette

	// Chrominance formats
	FOURCC_UV8 = Fourcc("UV8 ") //  8  UV 4:4

	// Luminance+Chrominance formats
	FOURCC_YUYV   = Fourcc("YUYV") // 16  YUV 4:2:2
	FOURCC_YYUV   = Fourcc("YYUV") // 16  YUV 4:2:2
	FOURCC_YVYU   = Fourcc("YVYU") // 16 YVU 4:2:2
	FOURCC_UYVY   = Fourcc("UYVY") // 16  YUV 4:2:2
	FOURCC_VYUY   = Fourcc("VYUY") // 16  YUV 4:2:2
	FOURCC_Y41P   = Fourcc("Y41P") // 12  YUV 4:1:1
	FOURCC_YUV444 = Fourcc("Y444") // 16  xxxxyyyy uuuuvvvv
	FOURCC_YUV555 = Fourcc("YUVO") // 16  YUV-5-5-5
	FOURCC_YUV565 = Fourcc("YUVP") // 16  YUV-5-6-5
	FOURCC_YUV24  = Fourcc("YUV3") // 24  YUV-8-8-8
	FOURCC_YUV32  = Fourcc("YUV4") // 32  YUV-8-8-8-8
	FOURCC_AYUV32 = Fourcc("AYUV") // 32  AYUV-8-8-8-8
	FOURCC_XYUV32 = Fourcc("XYUV") // 32  XYUV-8-8-8-8
	FOURCC_VUYA32 = Fourcc("VUYA") // 32  VUYA-8-8-8-8
	FOURCC_VUYX32 = Fourcc("VUYX") // 32  VUYX-8-8-8-8
	FOURCC_M420   = Fourcc("M420") // 12  YUV 4:2:0 2 lines y, 1 line uv interleaved

	// two planes -- one Y, one Cr + Cb interleaved
	FOURCC_NV12 = Fourcc("NV12") // 12  Y/CbCr 4:2:0
	FOURCC_NV21 = Fourcc("NV21") // 12  Y/CrCb 4:2:0
	FOURCC_NV16 = Fourcc("NV16") // 16  Y/CbCr 4:2:2
	FOURCC_NV61 = Fourcc("NV61") // 16  Y/CrCb 4:2:2
	FOURCC_NV24 = Fourcc("NV24") // 24  Y/CbCr 4:4:4
	FOURCC_NV42 = Fourcc("NV42") // 24  Y/CrCb 4:4:4
	FOURCC_HM12 = Fourcc("HM12") //  8  YUV 4:2:0 16x16 macroblocks

	// two non contiguous planes - one Y, one Cr + Cb interleaved
	FOURCC_NV12M        = Fourcc("NM12") // 12  Y/CbCr 4:2:0
	FOURCC_NV21M        = Fourcc("NM21") // 21  Y/CrCb 4:2:0
	FOURCC_NV16M        = Fourcc("NM16") // 16  Y/CbCr 4:2:2
	FOURCC_NV61M        = Fourcc("NM61") // 16  Y/CrCb 4:2:2
	FOURCC_NV12MT       = Fourcc("TM12") // 12  Y/CbCr 4:2:0 64x32 macroblocks
	FOURCC_NV12MT_16X16 = Fourcc("VM12") // 12  Y/CbCr 4:2:0 16x16 macroblocks

	// three planes - Y Cb, Cr
	FOURCC_YUV410  = Fourcc("YUV9") //  9  YUV 4:1:0
	FOURCC_YVU410  = Fourcc("YVU9") //  9  YVU 4:1:0
	FOURCC_YUV411P = Fourcc("411P") // 12  YVU411 planar
	FOURCC_YUV420  = Fourcc("YU12") // 12  YUV 4:2:0
	FOURCC_YVU420  = Fourcc("YV12") // 12  YVU 4:2:0
	FOURCC_YUV422P = Fourcc("422P") // 16  YVU422 planar

	// three non contiguous planes - Y, Cb, Cr
	FOURCC_YUV420M = Fourcc("YM12") // 12  YUV420 planar
	FOURCC_YVU420M = Fourcc("YM21") // 12  YVU420 planar
	FOURCC_YUV422M = Fourcc("YM16") // 16  YUV422 planar
	FOURCC_YVU422M = Fourcc("YM61") // 16  YVU422 planar
	FOURCC_YUV444M = Fourcc("YM24") // 24  YUV444 planar
	FOURCC_YVU444M = Fourcc("YM42") // 24  YVU444 planar

	// Bayer formats - see http://www.siliconimaging.com/RGB%20Bayer.htm
	FOURCC_SBGGR8  = Fourcc("BA81") //  8  BGBG.. GRGR..
	FOURCC_SGBRG8  = Fourcc("GBRG") //  8  GBGB.. RGRG..
	FOURCC_SGRBG8  = Fourcc("GRBG") //  8  GRGR.. BGBG..
	FOURCC_SRGGB8  = Fourcc("RGGB") //  8  RGRG.. GBGB..
	FOURCC_SBGGR10 = Fourcc("BG10") // 10  BGBG.. GRGR..
	FOURCC_SGBRG10 = Fourcc("GB10") // 10  GBGB.. RGRG..
	FOURCC_SGRBG10 = Fourcc("BA10") // 10  GRGR.. BGBG..
	FOURCC_SRGGB10 = Fourcc("RG10") // 10  RGRG.. GBGB..
	// 10bit raw bayer packed, 5 bytes for every 4 pixels
	FOURCC_SBGGR10P = Fourcc("pBAA")
	FOURCC_SGBRG10P = Fourcc("pGAA")
	FOURCC_SGRBG10P = Fourcc("pgAA")
	FOURCC_SRGGB10P = Fourcc("pRAA")
	// 10bit raw bayer a-law compressed to 8 bits
	FOURCC_SBGGR10ALAW8 = Fourcc("aBA8")
	FOURCC_SGBRG10ALAW8 = Fourcc("aGA8")
	FOURCC_SGRBG10ALAW8 = Fourcc("agA8")
	FOURCC_SRGGB10ALAW8 = Fourcc("aRA8")
	// 10bit raw bayer DPCM compressed to 8 bits
	FOURCC_SBGGR10DPCM8 = Fourcc("bBA8")
	FOURCC_SGBRG10DPCM8 = Fourcc("bGA8")
	FOURCC_SGRBG10DPCM8 = Fourcc("BD10")
	FOURCC_SRGGB10DPCM8 = Fourcc("bRA8")
	FOURCC_SBGGR12      = Fourcc("BG12") // 12  BGBG.. GRGR..
	FOURCC_SGBRG12      = Fourcc("GB12") // 12  GBGB.. RGRG..
	FOURCC_SGRBG12      = Fourcc("BA12") // 12  GRGR.. BGBG..
	FOURCC_SRGGB12      = Fourcc("RG12") // 12  RGRG.. GBGB..
	// 12bit raw bayer packed, 6 bytes for every 4 pixels
	FOURCC_SBGGR12P = Fourcc("pBCC")
	FOURCC_SGBRG12P = Fourcc("pGCC")
	FOURCC_SGRBG12P = Fourcc("pgCC")
	FOURCC_SRGGB12P = Fourcc("pRCC")
	FOURCC_SBGGR14  = Fourcc("BG14") // 14  BGBG.. GRGR..
	FOURCC_SGBRG14  = Fourcc("GB14") // 14  GBGB.. RGRG..
	FOURCC_SGRBG14  = Fourcc("GR14") // 14  GRGR.. BGBG..
	FOURCC_SRGGB14  = Fourcc("RG14") // 14  RGRG.. GBGB..
	// 14bit raw bayer packed, 7 bytes for every 4 pixels
	FOURCC_SBGGR14P = Fourcc("pBEE")
	FOURCC_SGBRG14P = Fourcc("pGEE")
	FOURCC_SGRBG14P = Fourcc("pgEE")
	FOURCC_SRGGB14P = Fourcc("pREE")
	FOURCC_SBGGR16  = Fourcc("BYR2") // 16  BGBG.. GRGR..
	FOURCC_SGBRG16  = Fourcc("GB16") // 16  GBGB.. RGRG..
	FOURCC_SGRBG16  = Fourcc("GR16") // 16  GRGR.. BGBG..
	FOURCC_SRGGB16  = Fourcc("RG16") // 16  RGRG.. GBGB..

	// HSV formats
	FOURCC_HSV24 = Fourcc("HSV3")
	FOURCC_HSV32 = Fourcc("HSV4")

	// compressed formats
	FOURCC_MJPEG          = Fourcc("MJPG") // Motion-JPEG
	FOURCC_JPEG           = Fourcc("JPEG") // JFIF JPEG
	FOURCC_DV             = Fourcc("dvsd") // 1394
	FOURCC_MPEG           = Fourcc("MPEG") // MPEG-1/2/4 Multiplexed
	FOURCC_H264           = Fourcc("H264") // H264 with start codes
	FOURCC_H264_NO_SC     = Fourcc("AVC1") // H264 without start codes
	FOURCC_H264_MVC       = Fourcc("M264") // H264 MVC
	FOURCC_H263           = Fourcc("H263") // H263
	FOURCC_MPEG1          = Fourcc("MPG1") // MPEG-1 ES
	FOURCC_MPEG2          = Fourcc("MPG2") // MPEG-2 ES
	FOURCC_MPEG2_SLICE    = Fourcc("MG2S") // MPEG-2 parsed slice data
	FOURCC_MPEG4          = Fourcc("MPG4") // MPEG-4 part 2 ES
	FOURCC_XVID           = Fourcc("XVID") // Xvid
	FOURCC_VC1_ANNEX_G    = Fourcc("VC1G") // SMPTE 421M Annex G compliant stream
	FOURCC_VC1_ANNEX_L    = Fourcc("VC1L") // SMPTE 421M Annex L compliant stream
	FOURCC_VP8            = Fourcc("VP80") // VP8
	FOURCC_VP8_FRAME      = Fourcc("VP8F") // VP8 parsed frame
	FOURCC_VP9            = Fourcc("VP90") // VP9
	FOURCC_HEVC           = Fourcc("HEVC") // HEVC aka H.265
	FOURCC_FWHT           = Fourcc("FWHT") // Fast Walsh Hadamard Transform (vicodec)
	FOURCC_FWHT_STATELESS = Fourcc("SFWH") // Stateless FWHT (vicodec)
	FOURCC_H264_SLICE     = Fourcc("S264") // H264 parsed slices

	//  Vendor-specific formats
	FOURCC_CPIA1            = Fourcc("CPIA") // cpia1 YUV
	FOURCC_WNVA             = Fourcc("WNVA") // Winnov hw compress
	FOURCC_SN9C10X          = Fourcc("S910") // SN9C10x compression
	FOURCC_SN9C20X_I420     = Fourcc("S920") // SN9C20x YUV 4:2:0
	FOURCC_PWC1             = Fourcc("PWC1") // pwc older webcam
	FOURCC_PWC2             = Fourcc("PWC2") // pwc newer webcam
	FOURCC_ET61X251         = Fourcc("E625") // ET61X251 compression
	FOURCC_SPCA501          = Fourcc("S501") // YUYV per line
	FOURCC_SPCA505          = Fourcc("S505") // YYUV per line
	FOURCC_SPCA508          = Fourcc("S508") // YUVY per line
	FOURCC_SPCA561          = Fourcc("S561") // compressed GBRG bayer
	FOURCC_PAC207           = Fourcc("P207") // compressed BGGR bayer
	FOURCC_MR97310A         = Fourcc("M310") // compressed BGGR bayer
	FOURCC_JL2005BCD        = Fourcc("JL20") // compressed RGGB bayer
	FOURCC_SN9C2028         = Fourcc("SONX") // compressed GBRG bayer
	FOURCC_SQ905C           = Fourcc("905C") // compressed RGGB bayer
	FOURCC_PJPG             = Fourcc("PJPG") // Pixart 73xx JPEG
	FOURCC_OV511            = Fourcc("O511") // ov511 JPEG
	FOURCC_OV518            = Fourcc("O518") // ov518 JPEG
	FOURCC_STV0680          = Fourcc("S680") // stv0680 bayer
	FOURCC_TM6000           = Fourcc("TM60") // tm5600/tm60x0
	FOURCC_CIT_YYVYUY       = Fourcc("CITV") // one line of Y then 1 line of VYUY
	FOURCC_KONICA420        = Fourcc("KONI") // YUV420 planar in blocks of 256 pixels
	FOURCC_JPGL             = Fourcc("JPGL") // JPEG-Lite
	FOURCC_SE401            = Fourcc("S401") // se401 janggu compressed rgb
	FOURCC_S5C_UYVY_JPG     = Fourcc("S5CI") // S5C73M3 interleaved UYVY/JPEG
	FOURCC_Y8I              = Fourcc("Y8I ") // Greyscale 8-bit L/R interleaved
	FOURCC_Y12I             = Fourcc("Y12I") // Greyscale 12-bit L/R interleaved
	FOURCC_Z16              = Fourcc("Z16 ") // Depth data 16-bit
	FOURCC_MT21C            = Fourcc("MT21") // Mediatek compressed block mode
	FOURCC_INZI             = Fourcc("INZI") // Intel Planar Greyscale 10-bit and Depth 16-bit
	FOURCC_SUNXI_TILED_NV12 = Fourcc("ST12") // Sunxi Tiled NV12 Format
	FOURCC_CNF4             = Fourcc("CNF4") // Intel 4-bit packed depth confidence information
	FOURCC_HI240            = Fourcc("HI24") // BTTV 8-bit dithered RGB
)

// DeviceConfig 相机配置信息
type DeviceConfig struct {
	Width  uint32 // 相机支持的分辨率宽度
	Height uint32 // 相机支持的分辨率高度
	FPS    uint32 // 相机在该分辨率下支持的帧率
	Format Fourcc // 相机支持的格式
}

func NewDeviceConfig(w, h, fps uint32, format Fourcc) DeviceConfig {
	return DeviceConfig{
		Width:  w,
		Height: h,
		FPS:    fps,
		Format: format,
	}
}

func (p *DeviceConfig) Clone() *DeviceConfig {
	if p == nil {
		return nil
	}
	return &DeviceConfig{
		Width:  p.Width,
		Height: p.Height,
		FPS:    p.FPS,
		Format: p.Format,
	}
}

func (p *DeviceConfig) Eq(v *DeviceConfig) bool {
	if p == nil && v == nil {
		return true
	}
	if p == nil || v == nil {
		return false
	}
	return p.Width == v.Width &&
		p.Height == v.Height &&
		p.FPS == v.FPS &&
		p.Format == v.Format
}
