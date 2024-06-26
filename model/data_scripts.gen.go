// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameDataScript = "data_scripts"

// DataScript mapped from table <data_scripts>
type DataScript struct {
	ID              string     `gorm:"column:id;primaryKey;comment:Id" json:"id"`                                                        // Id
	Name            string     `gorm:"column:name;not null;comment:名称" json:"name"`                                                      // 名称
	DeviceConfigID  string     `gorm:"column:device_config_id;not null;comment:设备配置id 关联删除" json:"device_config_id"`                     // 设备配置id 关联删除
	EnableFlag      string     `gorm:"column:enable_flag;not null;comment:启用标志Y-启用 N-停用 默认启用" json:"enable_flag"`                        // 启用标志Y-启用 N-停用 默认启用
	Content         *string    `gorm:"column:content;comment:内容" json:"content"`                                                         // 内容
	ScriptType      string     `gorm:"column:script_type;not null;comment:脚本类型 A-遥测上报预处理B-遥测下发预处理C-属性上报预处理D-属性下发预处理" json:"script_type"` // 脚本类型 A-遥测上报预处理B-遥测下发预处理C-属性上报预处理D-属性下发预处理
	LastAnalogInput *string    `gorm:"column:last_analog_input;comment:上次模拟输入" json:"last_analog_input"`                                 // 上次模拟输入
	Description     *string    `gorm:"column:description;comment:描述" json:"description"`                                                 // 描述
	CreatedAt       *time.Time `gorm:"column:created_at;comment:创建时间" json:"created_at"`                                                 // 创建时间
	UpdatedAt       *time.Time `gorm:"column:updated_at;comment:更新时间" json:"updated_at"`                                                 // 更新时间
	Remark          *string    `gorm:"column:remark;comment:备注" json:"remark"`                                                           // 备注
}

// TableName DataScript's table name
func (*DataScript) TableName() string {
	return TableNameDataScript
}
