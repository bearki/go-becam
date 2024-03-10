package camera

// Manager 相机管理器
type Manager interface {
	// GetList 获取相机列表
	//
	//	@return 相机列表
	//	@return 错误信息
	GetList() (DeviceList, error)

	// GetDeviceWithID 通过相机ID获取缓存的相机信息
	//
	//	@param	id	相机ID
	//	@return	缓存的相机信息
	//	@return	异常信息
	GetDeviceWithID(id string) (*Device, error)

	// Open 打开相机
	//
	//	@param	id		相机ID
	//	@param	info	分辨率信息
	//	@return	异常信息
	Open(id string, info DeviceConfig) error

	// GetDeviceConfigInfo 获取当前设备配置信息
	//
	//	@return	当前设备信息
	//	@return	当前设备配置信息
	//	@return	异常信息
	GetCurrDeviceConfigInfo() (*Device, *DeviceConfig, error)

	// GetStream 获取图片流
	//
	//	@return	图片流
	//	@return	异常信息
	GetStream() ([]byte, error)

	// Close 关闭已打开的相机
	Close()
}
