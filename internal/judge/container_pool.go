package judge

import (
	"log"
	"os/exec"
	"sync"
)

// 全局容器池（存储容器ID）
var containerPool chan string

// 用于保护池内容器的更新
var poolMutex sync.Mutex

// InitializeContainerPool 预热容器池
func InitializeContainerPool(n int) {
	containerPool = make(chan string, n)
	for range n {
		containerID, err := StartContainer()
		if err != nil {
			log.Printf("[FeasOJ] Error starting container during preheat: %v", err)
			continue
		}
		containerPool <- containerID
	}
	log.Printf("[FeasOJ] Preheated %d containers", len(containerPool))
}

// AcquireContainer 从池中获取一个空闲容器（若池为空则阻塞等待）
func AcquireContainer() string {
	containerID := <-containerPool
	log.Printf("[FeasOJ] Acquired container %s", containerID)
	return containerID
}

// ResetContainerForPool 用于在归还容器到池中前清理所有残留的任务目录
func ResetContainerForPool(containerID string) error {
	// 使用 find 命令删除 /workspace 下所有以 task_ 开头的目录
	resetCmd := exec.Command("docker", "exec", containerID, "sh", "-c", "find /workspace -maxdepth 1 -type d -name 'task_*' -exec rm -rf {} +")
	if err := resetCmd.Run(); err != nil {
		log.Printf("[FeasOJ] Error resetting container %s: %v", containerID, err)
		return err
	}
	return nil
}

// ReleaseContainer 将容器归还到池中（执行全局环境重置，若重置失败则重新启动新容器）
func ReleaseContainer(containerID string) {
	// 尝试清理容器中所有残留的任务目录
	if err := ResetContainerForPool(containerID); err != nil {
		log.Printf("[FeasOJ] Reset failed for container %s: %v, terminating it", containerID, err)
		// 清理失败，则终止当前容器
		TerminateContainer(containerID)
		// 尝试启动一个新容器替换
		newContainerID, err := StartContainer()
		if err != nil {
			log.Printf("[FeasOJ] Failed to start new container: %v", err)
			// 启动失败，则直接返回，不归还任何容器
			return
		}
		containerID = newContainerID
	}
	// 将容器归还到池中
	poolMutex.Lock()
	containerPool <- containerID
	poolMutex.Unlock()
	log.Printf("[FeasOJ] Released container %s back to pool", containerID)
}

// ShutdownContainerPool 在服务关闭时终止池中所有容器
func ShutdownContainerPool() {
	poolMutex.Lock()
	close(containerPool)
	for containerID := range containerPool {
		TerminateContainer(containerID)
		log.Printf("[FeasOJ] Terminated container %s", containerID)
	}
	poolMutex.Unlock()
}
