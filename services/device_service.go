package services

import (
	"ThingsPanel-Go/initialize/psql"
	"ThingsPanel-Go/initialize/redis"
	"ThingsPanel-Go/models"
	cm "ThingsPanel-Go/modules/dataService/mqtt"
	tphttp "ThingsPanel-Go/others/http"
	uuid "ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"

	"github.com/beego/beego/v2/core/logs"
	simplejson "github.com/bitly/go-simplejson"
	"gorm.io/gorm"
)

type DeviceService struct {
}

// Token 获取设备token
func (*DeviceService) Token(id string) (*models.Device, int64) {
	var device models.Device
	result := psql.Mydb.Where("id = ?", id).First(&device)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	return &device, result.RowsAffected
}

// GetDevicesByAssetID 获取设备列表
func (*DeviceService) GetDevicesByAssetID(asset_id string) ([]models.Device, int64) {
	var devices []models.Device
	var count int64
	result := psql.Mydb.Model(&models.Device{}).Where("asset_id = ?", asset_id).Find(&devices)
	psql.Mydb.Model(&models.Device{}).Where("asset_id = ?", asset_id).Count(&count)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	if len(devices) == 0 {
		devices = []models.Device{}
	}
	return devices, count
}

// GetDevicesByAssetID 获取设备列表(business_id string, device_id string, asset_id string, current int, pageSize int,device_type string)
func (*DeviceService) PageGetDevicesByAssetID(business_id string, asset_id string, device_id string, current int, pageSize int, device_type string, token string, name string) ([]map[string]interface{}, int64) {
	sqlWhere := `select (with RECURSIVE ast as 
		( 
		(select aa.id,cast(aa.name as varchar(255)),aa.parent_id  from asset aa where id=a.id) 
		union  
		(select tt.id,cast (kk.name||'/'||tt.name as varchar(255))as name ,kk.parent_id from ast tt inner join asset  kk on kk.id = tt.parent_id )
		)select  name from ast where parent_id='0' limit 1) 
		as asset_name,b.id as business_id ,b."name" as business_name,d.d_id,d.location,a.id as asset_id ,d.id as device ,d."name" as device_name,
		   d."token" as device_token,d."type" as device_type,d.protocol as protocol ,(select ts from ts_kv_latest tkl where tkl.entity_id = d.id order by ts desc limit 1) as latest_ts
		   from device d left join asset a on d.asset_id =  a.id left join business b on b.id = a.business_id  where 1=1 `
	sqlWhereCount := `select count(1) from device d left join asset a on d.asset_id =  a.id left join business b on b.id = a.business_id  where 1=1`
	var values []interface{}
	var where = ""
	if business_id != "" {
		values = append(values, business_id)
		where += " and b.id = ?"
	}
	if asset_id != "" {
		values = append(values, asset_id)
		where += " and a.id = ?"
	}
	if device_id != "" {
		values = append(values, device_id)
		where += " and d.id = ?"
	}
	if device_type != "" {
		values = append(values, device_type)
		where += " and d.type = ?"
	}
	if token != "" {
		values = append(values, token)
		where += " and d.token = ?"
	}
	if name != "" {
		where += " and d.name like '%" + name + "%'"
	}
	sqlWhere += where
	sqlWhereCount += where
	var count int64
	result := psql.Mydb.Raw(sqlWhereCount, values...).Count(&count)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	var offset int = (current - 1) * pageSize
	var limit int = pageSize
	sqlWhere += " offset ? limit ?"
	values = append(values, offset, limit)
	var deviceList []map[string]interface{}
	dataResult := psql.Mydb.Raw(sqlWhere, values...).Scan(&deviceList)
	if dataResult.Error != nil {
		errors.Is(dataResult.Error, gorm.ErrRecordNotFound)
	}
	return deviceList, count
}

// GetDevicesByAssetID 获取设备列表(business_id string, device_id string, asset_id string, current int, pageSize int,device_type string)
func (*DeviceService) PageGetDevicesByAssetIDTree(req valid.DevicePageListValidate) ([]map[string]interface{}, int64) {
	sqlWhere := `select (with RECURSIVE ast as 
		( 
		(select aa.id,cast(aa.name as varchar(255)),aa.parent_id  from asset aa where id=a.id) 
		union  
		(select tt.id,cast (kk.name||'/'||tt.name as varchar(255))as name ,kk.parent_id from ast tt inner join asset  kk on kk.id = tt.parent_id )
		)select  name from ast where parent_id='0' limit 1) 
		as asset_name,b.id as business_id ,b."name" as business_name,d.d_id,d.location,a.id as asset_id ,d.id as device ,d."name" as device_name,d.device_type as device_type,d.parent_id as parent_id,
		   d."token" as device_token,d."type" as "type",d.protocol as protocol ,(select ts from ts_kv_latest tkl where tkl.entity_id = d.id order by ts desc limit 1) as latest_ts
		   from device d left join asset a on d.asset_id =  a.id left join business b on b.id = a.business_id  where 1=1  and d.device_type != '3'`
	sqlWhereCount := `select count(1) from device d left join asset a on d.asset_id =  a.id left join business b on b.id = a.business_id  where 1=1 and d.device_type != '3'`
	var values []interface{}
	var where = ""
	if req.BusinessId != "" {
		values = append(values, req.BusinessId)
		where += " and b.id = ?"
	}
	if req.AssetId != "" {
		values = append(values, req.AssetId)
		where += " and a.id = ?"
	}
	if req.DeviceId != "" {
		values = append(values, req.DeviceId)
		where += " and d.id = ?"
	}
	if req.DeviceType != "" {
		values = append(values, req.DeviceType)
		where += " and d.type = ?"
	}
	if req.Token != "" {
		values = append(values, req.Token)
		where += " and d.token = ?"
	}
	if req.Name != "" {
		where += " and d.name like '%" + req.Name + "%'"
	}
	sqlWhere += where
	sqlWhereCount += where
	var count int64
	result := psql.Mydb.Raw(sqlWhereCount, values...).Count(&count)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	var offset int = (req.CurrentPage - 1) * req.PerPage
	var limit int = req.PerPage
	sqlWhere += " offset ? limit ?"
	values = append(values, offset, limit)
	var deviceList []map[string]interface{}
	dataResult := psql.Mydb.Raw(sqlWhere, values...).Scan(&deviceList)
	if dataResult.Error != nil {
		errors.Is(dataResult.Error, gorm.ErrRecordNotFound)
	} else {
		for _, device := range deviceList {
			fmt.Println("=====================================")
			fmt.Println(device)
			fmt.Println(device["device_type"])
			if device["device_type"].(string) == "2" { // 网关设备需要查询子设备
				var subDeviceList []map[string]interface{}
				sql := `select (with RECURSIVE ast as 
					( 
					(select aa.id,cast(aa.name as varchar(255)),aa.parent_id  from asset aa where id=a.id) 
					union  
					(select tt.id,cast (kk.name||'/'||tt.name as varchar(255))as name ,kk.parent_id from ast tt inner join asset  kk on kk.id = tt.parent_id )
					)select  name from ast where parent_id='0' limit 1) 
					as asset_name,b.id as business_id ,b."name" as business_name,d.d_id,d.location,a.id as asset_id ,d.id as device ,d."name" as device_name,d.device_type  as device_type,d.parent_id as parent_id,
					   d."token" as device_token,d."type" as "type",d.protocol as protocol ,(select ts from ts_kv_latest tkl where tkl.entity_id = d.id order by ts desc limit 1) as latest_ts
					   from device d left join asset a on d.asset_id =  a.id left join business b on b.id = a.business_id  where 1=1  and d.device_type = '3' and d.parent_id = '` + device["device"].(string) + `'`
				result := psql.Mydb.Raw(sql).Scan(&subDeviceList)
				if result.Error != nil {
					errors.Is(result.Error, gorm.ErrRecordNotFound)
				} else {
					device["children"] = subDeviceList
				}
			}
		}
	}
	return deviceList, count
}

// GetDevicesByBusinessID 根据业务ID获取设备列表
// return []设备,设备数量
// 2022-04-18新增
func (*DeviceService) GetDevicesByBusinessID(business_id string) ([]models.Device, int64) {
	var devices []models.Device
	SQL := `select device.id,device.asset_id ,device.additional_info,device."type" ,device."location",device."d_id",device."name",device."label",device.protocol from device left join asset on device.asset_id  = asset.id where asset.business_id =?`
	if err := psql.Mydb.Raw(SQL, business_id).Scan(&devices).Error; err != nil {
		log.Println(err.Error())
	}
	if len(devices) == 0 {
		devices = []models.Device{}
	}
	return devices, int64(len(devices))
}

// GetDevicesByBusinessID 根据业务ID获取设备列表
// return []设备,设备数量
// 2022-04-18新增
func (*DeviceService) GetDevicesInfoAndCurrentByAssetID(asset_id string) ([]models.Device, int64) {
	var devices []models.Device
	SQL := `select device.id,device.asset_id ,device.additional_info,device."type" ,device."location",device."d_id",device."name",device."label",device.protocol from device left join asset on device.asset_id  = asset.id where asset.id =?`
	if err := psql.Mydb.Raw(SQL, asset_id).Scan(&devices).Error; err != nil {
		log.Println(err.Error())
	}
	if len(devices) == 0 {
		devices = []models.Device{}
	}
	return devices, int64(len(devices))
}

// GetDevicesByAssetIDs 获取设备列表
func (*DeviceService) GetDevicesByAssetIDs(asset_ids []string) (devices []models.Device, err error) {
	err = psql.Mydb.Model(&models.Device{}).Where("asset_id IN ?", asset_ids).Find(&devices).Error
	if err != nil {
		return devices, err
	}
	return devices, nil
}

// GetAllDevicesByID 获取所有设备
func (*DeviceService) GetAllDeviceByID(id string) ([]models.Device, int64) {
	var devices []models.Device
	var count int64
	result := psql.Mydb.Model(&models.Device{}).Where("id = ?", id).Find(&devices)
	psql.Mydb.Model(&models.Device{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	if len(devices) == 0 {
		devices = []models.Device{}
	}
	return devices, count
}

// GetDevicesByID 获取设备
func (*DeviceService) GetDeviceByID(id string) (*models.Device, int64) {
	var device models.Device
	result := psql.Mydb.Where("id = ?", id).First(&device)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	return &device, result.RowsAffected
}

// Delete 根据ID删除Device
func (*DeviceService) Delete(id string) bool {
	var device models.Device
	psql.Mydb.Where("id = ?", id).First(&device)
	result := psql.Mydb.Where("id = ?", id).Delete(&models.Device{})
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return false
	}
	if device.Token != "" {
		redis.DelKey("token" + device.Token)
		tphttp.Delete(viper.GetString("api.delete")+device.Token, "{}")
	}
	return true
}

// 获取全部Device
func (*DeviceService) All() ([]models.Device, int64) {
	var devices []models.Device
	var count int64
	result := psql.Mydb.Model(&devices).Count(&count)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	if len(devices) == 0 {
		devices = []models.Device{}
	}
	return devices, count
}

// 判断token是否存在
func (*DeviceService) IsToken(token string) bool {
	var devices []models.Device
	var count int64
	result := psql.Mydb.Model(&devices).Where("token = ?", token).Count(&count)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return false
	}
	return int(count) > 0
}

// 根据ID编辑Device的Token
func (*DeviceService) Edit(deviceModel valid.EditDevice) bool {
	var device models.Device
	psql.Mydb.Where("id = ?", deviceModel.ID).First(&device)
	result := psql.Mydb.Model(&models.Device{}).Where("id = ?", deviceModel.ID).Updates(models.Device{
		Token:     deviceModel.Token,
		Protocol:  deviceModel.Protocol,
		Port:      deviceModel.Port,
		Publish:   deviceModel.Publish,
		Subscribe: deviceModel.Subscribe,
		Username:  deviceModel.Username,
		Password:  deviceModel.Password,
		AssetID:   deviceModel.AssetID,
		Type:      deviceModel.Type,
		Name:      deviceModel.Name,
	})
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return false
	}
	if deviceModel.Token != "" {
		if device.Token != "" {
			redis.DelKey("token" + device.Token)
			tphttp.Delete(viper.GetString("api.delete")+device.Token, "{}")
		}
		redis.SetStr("token"+deviceModel.Token, deviceModel.ID, 3600*time.Second)
		tphttp.Post(viper.GetString("api.add")+device.Token, "{\"password\":\""+device.Password+"\"}")
	}

	return true
}

func (*DeviceService) Add(device models.Device) (bool, string) {

	var uuid = uuid.GetUuid()
	device.ID = uuid
	result := psql.Mydb.Create(&device)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return false, ""
	}
	if device.Token != "" {
		redis.SetStr("token"+device.Token, uuid, 3600*time.Second)
		tphttp.Post(viper.GetString("api.add")+device.Token, "{\"password\":\"\"}")
	}

	return true, uuid
}

// 向mqtt发送控制指令
func (*DeviceService) OperatingDevice(deviceId string, field string, value interface{}) bool {
	//reqMap := make(map[string]interface{})
	valueMap := make(map[string]interface{})
	logs.Info("通过设备id获取设备token")
	var DeviceService DeviceService
	device, _ := DeviceService.Token(deviceId)
	if device == nil {
		logs.Info("没有匹配的token")
		return false
	}
	//reqMap["token"] = device.Token
	logs.Info("token-%s", device.Token)
	// logs.Info("把field字段映射回设备端字段")
	// var fieldMappingService FieldMappingService
	// deviceField := fieldMappingService.TransformByDeviceid(deviceId, field)
	// if deviceField != "" {
	// 	valueMap[deviceField] = value
	// }
	valueMap[field] = value
	//reqMap["values"] = valueMap
	logs.Info("将map转json")
	mjson, _ := json.Marshal(valueMap)
	logs.Info("json-%s", string(mjson))
	err := cm.Send(mjson, device.Token)
	if err == nil {
		logs.Info("发送到mqtt成功")
		return true
	} else {
		logs.Info(err.Error())
		return false
	}
}

//自动化发送控制
func (*DeviceService) ApplyControl(res *simplejson.Json) {
	logs.Info("执行控制开始")
	//"apply":[{"asset_id":"xxx","field":"hum","device_id":"xxx","value":"1"}]}
	applyRows, _ := res.Get("apply").Array()
	logs.Info("applyRows-", applyRows)
	for _, applyRow := range applyRows {
		logs.Info("applyRow-", applyRow)
		if applyMap, ok := applyRow.(map[string]interface{}); ok {
			logs.Info(applyMap)
			// 如果有“或者，并且”操作符，就给code加上操作符
			if applyMap["field"] != nil && applyMap["value"] != nil {
				logs.Info("准备执行控制发送函数")
				var s = ""
				switch applyMap["value"].(type) {
				case string:
					s = applyMap["value"].(string)
				case json.Number:
					s = applyMap["value"].(json.Number).String()
				}
				ConditionsLog := models.ConditionsLog{
					DeviceId:      applyMap["device_id"].(string),
					OperationType: "3",
					Instruct:      applyMap["field"].(string) + ":" + s,
					ProtocolType:  "mqtt",
					CteateTime:    time.Now().Format("2006-01-02 15:04:05"),
				}
				var DeviceService DeviceService
				reqFlag := DeviceService.OperatingDevice(applyMap["device_id"].(string), applyMap["field"].(string), applyMap["value"])
				if reqFlag {
					logs.Info("成功发送控制")
					ConditionsLog.SendResult = "1"
				} else {
					logs.Info("成功发送失败")
					ConditionsLog.SendResult = "2"
				}
				// 记录日志
				var ConditionsLogService ConditionsLogService
				ConditionsLogService.Insert(&ConditionsLog)
			}
		}
	}
}

// func (*DeviceService) ApplyControl(res *simplejson.Json) {
// 	logs.Info("执行控制开始")
// 	//"apply":[{"asset_id":"xxx","field":"hum","device_id":"xxx","value":"1"}]}
// 	applyRows, _ := res.Get("apply").Array()
// 	logs.Info("applyRows-", applyRows)
// 	for _, applyRow := range applyRows {
// 		logs.Info("applyRow-", applyRow)
// 		if applyMap, ok := applyRow.(map[string]interface{}); ok {
// 			logs.Info(applyMap)
// 			// 如果有“或者，并且”操作符，就给code加上操作符
// 			if applyMap["field"] != nil && applyMap["value"] != nil {
// 				logs.Info("准备执行控制发送函数")
// 				var s = ""
// 				switch applyMap["value"].(type) {
// 				case string:
// 					s = applyMap["value"].(string)
// 				case json.Number:
// 					s = applyMap["value"].(json.Number).String()
// 				}
// 				ConditionsLog := models.ConditionsLog{
// 					DeviceId:      applyMap["device_id"].(string),
// 					OperationType: "3",
// 					Instruct:      applyMap["field"].(string) + ":" + s,
// 					ProtocolType:  "mqtt",
// 					CteateTime:    time.Now().Format("2006-01-02 15:04:05"),
// 				}
// 				var DeviceService DeviceService
// 				reqFlag := DeviceService.OperatingDevice(applyMap["device_id"].(string), applyMap["field"].(string), applyMap["value"])
// 				if reqFlag {
// 					logs.Info("成功发送控制")
// 					ConditionsLog.SendResult = "1"
// 				} else {
// 					logs.Info("成功发送失败")
// 					ConditionsLog.SendResult = "2"
// 				}
// 				// 记录日志
// 				var ConditionsLogService ConditionsLogService
// 				ConditionsLogService.Insert(&ConditionsLog)
// 			}
// 		}
// 	}
// }

// 根据token获取网关设备和子设备的配置
func (*DeviceService) GetConfigByToken(token string) map[string]interface{} {
	var GatewayConfigMap = make(map[string]interface{})
	var device models.Device
	result := psql.Mydb.First(&device, "token = ?", token)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return GatewayConfigMap
	}
	var sub_devices []models.Device
	sub_result := psql.Mydb.Find(&sub_devices, "parent_id = ?", device.ID)
	if sub_result.Error != nil {
		errors.Is(sub_result.Error, gorm.ErrRecordNotFound)
	} else {

		GatewayConfigMap["GatewayId"] = device.ID
		GatewayConfigMap["ProtocolType"] = device.Protocol
		GatewayConfigMap["AccessToken"] = token
		for _, sub_device := range sub_devices {
			var m = make(map[string]interface{})
			err := json.Unmarshal([]byte(sub_device.ProtocolConfig), &m)
			if err != nil {
				fmt.Println("Unmarshal failed:", err)
			}
			GatewayConfigMap["SubDevice"] = m
		}
		return GatewayConfigMap
	}
	return GatewayConfigMap
}
