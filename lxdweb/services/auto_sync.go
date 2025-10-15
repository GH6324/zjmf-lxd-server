package services

import (
	"log"
	"lxdweb/database"
	"lxdweb/models"
	"time"
)

func StartAutoSyncService() {
	log.Println("[AUTO-SYNC] 自动同步服务已禁用")
}

// SyncAllNodesFullAsync 完整同步所有节点的所有数据
func SyncAllNodesFullAsync() {
	log.Println("[AUTO-SYNC] 开始执行完整实时同步任务")

	var nodes []models.Node
	if err := database.DB.Where("status = ?", "active").Find(&nodes).Error; err != nil {
		log.Printf("[AUTO-SYNC] 查询节点失败: %v", err)
		return
	}

	if len(nodes) == 0 {
		log.Println("[AUTO-SYNC] 没有活跃的节点需要同步")
		return
	}

	log.Printf("[AUTO-SYNC] 找到 %d 个活跃节点，开始实时同步", len(nodes))

	for i, node := range nodes {
		log.Printf("[AUTO-SYNC] 处理节点 %d/%d: %s", i+1, len(nodes), node.Name)
		syncNodeFull(node)

		if i < len(nodes)-1 {
			interval := time.Duration(node.BatchInterval) * time.Second
			log.Printf("[AUTO-SYNC] 等待 %v 后处理下一个节点", interval)
			time.Sleep(interval)
		}
	}

	log.Println("[AUTO-SYNC] 所有节点完整同步任务完成")
}

// syncNodeFull 完整同步单个节点的所有数据类型
func syncNodeFull(n models.Node) {
	log.Printf("[AUTO-SYNC] 实时同步节点: %s (ID: %d)", n.Name, n.ID)

	if err := SyncNodeContainers(n.ID, false); err != nil {
		log.Printf("[AUTO-SYNC] 节点 %s 容器同步失败: %v", n.Name, err)
	}

	time.Sleep(1 * time.Second)

	if err := SyncNodeNATRules(n.ID, false); err != nil {
		log.Printf("[AUTO-SYNC] 节点 %s NAT规则同步失败: %v", n.Name, err)
	}

	time.Sleep(1 * time.Second)

	if err := SyncNodeIPv6Bindings(n.ID); err != nil {
		log.Printf("[AUTO-SYNC] 节点 %s IPv6绑定同步失败: %v", n.Name, err)
	}

	time.Sleep(1 * time.Second)

	if err := SyncNodeProxyConfigs(n.ID); err != nil {
		log.Printf("[AUTO-SYNC] 节点 %s Proxy配置同步失败: %v", n.Name, err)
	}

	log.Printf("[AUTO-SYNC] 节点 %s 同步完成", n.Name)
}

func EnableAutoSync() {
	log.Println("[AUTO-SYNC] 自动同步功能已禁用")
}

func DisableAutoSync() {
	log.Println("[AUTO-SYNC] 自动同步功能已禁用")
}

func IsAutoSyncEnabled() bool {
	return false
}

