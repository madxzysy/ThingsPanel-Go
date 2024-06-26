// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameOtaUpgradeTask = "ota_upgrade_tasks"

// OtaUpgradeTask mapped from table <ota_upgrade_tasks>
type OtaUpgradeTask struct {
	ID                  string    `gorm:"column:id;primaryKey;comment:Id" json:"id"`                                                   // Id
	Name                string    `gorm:"column:name;not null;comment:任务名称" json:"name"`                                               // 任务名称
	OtaUpgradePackageID string    `gorm:"column:ota_upgrade_package_id;not null;comment:升级包id（外键，关联删除）" json:"ota_upgrade_package_id"` // 升级包id（外键，关联删除）
	Description         *string   `gorm:"column:description;comment:描述" json:"description"`                                            // 描述
	CreatedAt           time.Time `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`                                   // 创建时间
	Remark              *string   `gorm:"column:remark;comment:备注" json:"remark"`                                                      // 备注
}

// TableName OtaUpgradeTask's table name
func (*OtaUpgradeTask) TableName() string {
	return TableNameOtaUpgradeTask
}
