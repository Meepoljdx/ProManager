package models

import (
	"time"

	"gorm.io/datatypes"
)

type PrometheusConf struct {
	ID      string
	AppConf datatypes.JSON
}

// PrometheusList结构体，记录概要信息
type PrometheusNode struct {
	// Prometheus 的清单model
	ID       string `json:"node_id" form:"node_id" gorm:"column:node_id;primaryKey"` // Promehteus的ID
	IP       string `json:"node_ip" form:"node_ip" gorm:"column:node_ip"`            // Prometheus的IP
	Role     string `json:"node_role" form:"node_role" gorm:"column:node_role"`      // Master/Slave
	Status   string `json:"node_status" form:"node_status" gorm:"column:node_status"`
	CreateAt time.Time
	UpdateAt time.Time
	URL      string `json:"web_url" form:"web_url" gorm:"column:web_url"`
	Cluster  string `json:"cluster" form:"cluster" gorm:"column:cluster;primaryKey"`
	Groups   string `json:"groups" from:"groups" gorm:"column:groups;primaryKey"`
}

func (p PrometheusNode) TableName() string {
	return "prometheus_node" // 数据库中创建的表名
}

func (pc PrometheusConf) TableName() string {
	return "service_config"
}

// Prometheus
