// @Title: Clusters.go
// @Description: 集群models，把组和集群存在一起
// @Author: 李嘉栋
package models

import "time"

type ClusterMessage struct {
	// 集群信息
	Cluster  string `gorm:"primaryKey" json:"cluster"`
	Groups   string `gorm:"primaryKey" json:"groups"`
	NodeNums int    `gorm:"column:node_nums"` // 节点数量
	CreateAt time.Time
	UpdateAt time.Time
}

func (c ClusterMessage) TableName() string {
	return "clusters_manager"
}
