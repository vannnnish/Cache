package main

import (
	"fmt"
	"stathat.com/c/consistent"
)

func main() {
	// 创建一个一致性哈希的环实例
	circle := consistent.New()

	// 设置虚拟节点个数
	circle.NumberOfReplicas = 1024

	// 设置物理节点
	circle.Set([]string{
		"A", "B", "C", "D",
	})
	// 获取 key1 应该在哪台机器上
	machine, err := circle.Get("key1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("key1 在机器 %s 上！", machine)
}
