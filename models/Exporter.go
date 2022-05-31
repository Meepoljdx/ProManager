package models

import "time"

type ExporterNode struct {
	// Nodes Exporter的清单的model
	ID       string `json:"node_id" form:"node_id" gorm:"column:node_id;primaryKey"`
	IP       string `json:"node_ip" form:"node_ip" gorm:"column:node_ip"`
	Status   string `json:"node_status" form:"node_status" gorm:"column:node_status"`
	Type     string `json:"node_type" form:"node_type" gorm:"column:node_type;primaryKey"`
	CreateAt time.Time
	UpdateAt time.Time
	URL      string `json:"web_url" form:"web_url" gorm:"column:web_url"`
	Cluster  string `json:"cluster" form:"cluster" gorm:"column:cluster;primaryKey"`
	Groups   string `json:"groups" from:"groups" gorm:"column:groups;primaryKey"`
}

func (n ExporterNode) TableName() string {
	return "exporter_node"
}
