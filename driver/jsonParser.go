package driver

/* Device Struct
*
* Data structure for Device info
//*
// */
//type Device struct {
//	Id       string `json:"id"`
//	Created  int64  `json:"created"`
//	Origin   int64  `json:"origin"`
//	Modified int64  `json:"modified"`
//	Device   string `json:"device"`
//	Name     string `json:"name"`
//	Value    int32  `json:"value"`
//}
//
///* NewTransaction()
//*
//* To return new transaction data
//*
// */
//func NewDevice(id string, created int64, origin int64, modified int64, device string, name string, value int32) Device {
//	return Device{
//		Id:       id,
//		Created:  created,
//		Origin:   origin,
//		Modified: modified,
//		Device:   device,
//		Name:     name,
//		Value:    value,
//	}
//}
//
///* EncodeToJson()
//*
//* To Encode Device info from json format
//*
// */
//func (device *Device) EncodeToJson() (string, error) {
//	jsonBytes, error := json.Marshal(device)
//	return string(jsonBytes), error
//}
//
///* DecodeFromJson()
//*
//* To Decode Device info from json format
//*
// */
//func (device *Device) DecodeFromJson(jsonString string) error {
//	return json.Unmarshal([]byte(jsonString), device)
//}
//
///* GetId
//*
//* To Get device Id
//*
// */
//func (device *Device) getId() string {
//	return device.Id
//}
//
///* getCreated
//*
//* To Get device Created date
//*
// */
//func (device *Device) getCreated() int64 {
//	return device.Created
//}
//
///* getOrigin
//*
//* To Get device Origin date
//*
// */
//func (device *Device) getOrigin() int64 {
//	return device.Origin
//}
//
///* getModified
//*
//* To Get Modified date
//*
// */
//func (device *Device) getModified() int64 {
//	return device.Modified
//}
//
///* getDevice
//*
//* To Get Device info
//*
// */
//func (device *Device) getDevice() string {
//	return device.Device
//}
//
///* getDeviceName
//*
//* To Get Device Name
//*
// */
//func (device *Device) getDeviceName() string {
//	return device.Name
//}
//
///* getDeviceValue
//*
//* To Get Device Value
//*
// */
//func (device *Device) getDeviceValue() int32 {
//	return device.Value
//}
//
///* printDeviceInfo
//*
//* To print device information
//*
// */
//func (device *Device) printDeviceInfo() {
//
//	fmt.Print("Device ID: ", device.getId(), "Device Created Time: ", device.getCreated(), "Device Origin: ", device.getOrigin(), "Device Modified: ", device.getModified(), "Device Device: ", device.getDevice(), "Device Device Name: ", device.getDeviceName(), "Device Value: ", device.getDeviceValue())
//}
//
//func DevicePurchasedEnergy() {
//	fmt.Print("List of device(s) need energy!!")
//}
//
//func DeviceNeedEnergy() {
//	fmt.Print("List of device(s) purchased energy!!")
//}
//
//func DeviceSoldEnergy() {
//	fmt.Print("List of device(s) sold energy!!")
//}
