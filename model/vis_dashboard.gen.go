// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameVisDashboard = "vis_dashboard"

// VisDashboard mapped from table <vis_dashboard>
type VisDashboard struct {
	ID            string     `gorm:"column:id;primaryKey" json:"id"`
	RelationID    *string    `gorm:"column:relation_id" json:"relation_id"`
	JSONDatum     *string    `gorm:"column:json_data;default:{}" json:"json_data"`
	DashboardName *string    `gorm:"column:dashboard_name" json:"dashboard_name"`
	CreateAt      *time.Time `gorm:"column:create_at" json:"create_at"`
	Sort          *int32     `gorm:"column:sort;comment:排序" json:"sort"` // 排序
	Remark        *string    `gorm:"column:remark" json:"remark"`
	TenantID      *string    `gorm:"column:tenant_id" json:"tenant_id"`
	ShareID       *string    `gorm:"column:share_id;comment:分享id" json:"share_id"` // 分享id
}

// TableName VisDashboard's table name
func (*VisDashboard) TableName() string {
	return TableNameVisDashboard
}