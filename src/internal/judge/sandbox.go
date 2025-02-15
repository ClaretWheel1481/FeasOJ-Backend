package judge

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"src/config"
	"src/internal/global"
	"src/internal/utils/sql"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// BuildImage 构建Sandbox
func BuildImage() bool {
	ctx := context.Background()

	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panic(err)
		return false
	}

	// 将Dockerfile目录打包成tar格式
	var dir string
	if config.DebugMode {
		dir = global.ParentDir
	} else {
		dir = global.CurrentDir
	}
	tar, err := archive.TarWithOptions(dir, &archive.TarOptions{})
	if err != nil {
		log.Panic(err)
	}

	// 设置镜像构建选项
	buildOptions := types.ImageBuildOptions{
		Context:    tar,                          // 构建上下文
		Dockerfile: "Sandbox",                    // Dockerfile文件名
		Tags:       []string{"judgecore:latest"}, // 镜像标签
	}

	log.Println("[FeasOJ] SandBox is being built...")
	// 构建Docker镜像
	buildResponse, err := cli.ImageBuild(ctx, tar, buildOptions)
	if err != nil {
		log.Panic(err)
		return false
	}
	defer buildResponse.Body.Close()

	// 打印构建响应
	_, err = io.Copy(log.Writer(), buildResponse.Body)
	if err != nil {
		log.Printf("[FeasOJ] Error copying build response: %v", err)
	}

	return true
}

// StartContainer 启动Docker容器
func StartContainer() (string, error) {
	ctx := context.Background()

	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	// 配置容器配置
	containerConfig := &container.Config{
		Image: "judgecore:latest",
		Cmd:   []string{"sh"},
		Tty:   true,
	}

	// 配置主机配置
	hostConfig := &container.HostConfig{
		Resources: config.SandBoxConfig,
		Binds: []string{
			global.CodeDir + ":/workspace", // 挂载文件夹
		},
		AutoRemove: true, // 容器退出后自动删除
	}

	// 创建容器
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}

	// 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

// CompileAndRun 启动容器并编译运行文件、放入输入、捕获输出、对照输出
func CompileAndRun(filename string, containerID string) string {
	ext := filepath.Ext(filename)
	var compileCmd *exec.Cmd

	switch ext {
	case ".cpp":
		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c", fmt.Sprintf("g++ /workspace/%s -o /workspace/%s.out", filename, filename))
		if err := compileCmd.Run(); err != nil {
			TerminateContainer(containerID)
			return "Compile Failed"
		}
	case ".java":
		// 临时重命名为Main.java
		originalName := filename
		tempName := "Main_" + filename + ".java"
		renameCmd := exec.Command("docker", "exec", containerID, "sh", "-c", fmt.Sprintf("mv /workspace/%s /workspace/%s", originalName, tempName))
		if err := renameCmd.Run(); err != nil {
			TerminateContainer(containerID)
			return "Compile Failed"
		}

		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c", fmt.Sprintf("javac /workspace/%s", tempName))
		if err := compileCmd.Run(); err != nil {
			TerminateContainer(containerID)
			return "Compile Failed"
		}

		// 编译完成后改回原名
		renameBackCmd := exec.Command("docker", "exec", containerID, "sh", "-c", fmt.Sprintf("mv /workspace/%s /workspace/%s", tempName, originalName))
		if err := renameBackCmd.Run(); err != nil {
			TerminateContainer(containerID)
			return "Compile Failed"
		}
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 两种方案：1、数据库存放输入输出样例 2、in，out文件存放输入输出样例，数据库存放对应文件路径
	// 数据库获取输入输出样例
	testCases := sql.SelectTestCasesByPid(strings.Split(filename, "_")[1])
	for _, testCase := range testCases {
		var runCmd *exec.Cmd
		switch ext {
		case ".cpp":
			runCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c", fmt.Sprintf("/workspace/%s.out", filename))
		case ".py":
			runCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c", fmt.Sprintf("python /workspace/%s", filename))
		case ".go":
			runCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c", fmt.Sprintf("go run /workspace/%s", filename))
		case ".java":
			runCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c", fmt.Sprintf("java Main_%s", filename))
		default:
			TerminateContainer(containerID)
			return "Failed"
		}

		runCmd.Stdin = strings.NewReader(testCase.InputData)
		output, err := runCmd.CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			TerminateContainer(containerID)
			return "Time Limit Exceeded"
		}
		if err != nil {
			TerminateContainer(containerID)
			return "Failed"
		}
		outputStr := string(output)
		if strings.TrimSpace(outputStr) != strings.TrimSpace(testCase.OutputData) {
			TerminateContainer(containerID)
			return "Wrong Answer"
		}
	}
	// TODO: 第二种方案，读取in，out文件
	TerminateContainer(containerID)
	return "Success"
}

// TerminateContainer 终止并删除Docker容器
func TerminateContainer(containerID string) bool {
	ctx := context.Background()

	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// 终止容器
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		panic(err)
	}

	return true
}
