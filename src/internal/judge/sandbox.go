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
		CapDrop:    []string{"ALL"},
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

// ResetContainer 只清理任务专属的目录，而不影响其他任务
func ResetContainer(containerID, taskDir string) error {
	resetCmd := exec.Command("docker", "exec", containerID, "sh", "-c", fmt.Sprintf("rm -rf %s", taskDir))
	if err := resetCmd.Run(); err != nil {
		log.Printf("[FeasOJ] Error cleaning task directory %s in container %s: %v", taskDir, containerID, err)
		return err
	}
	return nil
}

// CompileAndRun 编译并运行代码
func CompileAndRun(filename string, containerID string) string {
	// 生成唯一任务目录（使用当前时间戳纳秒值）
	taskDir := fmt.Sprintf("/workspace/task_%d", time.Now().UnixNano())

	// 在容器内创建任务目录
	mkdirCmd := exec.Command("docker", "exec", containerID, "mkdir", "-p", taskDir)
	if err := mkdirCmd.Run(); err != nil {
		return "Internal Error: cannot create task dir"
	}

	// 将代码文件从挂载的workspace目录复制到任务目录中
	copyCmd := exec.Command("docker", "exec", containerID, "cp", fmt.Sprintf("/workspace/%s", filename), taskDir)
	if err := copyCmd.Run(); err != nil {
		return "Internal Error: cannot copy file"
	}

	// 确保任务结束后清理任务目录
	defer func() {
		if err := ResetContainer(containerID, taskDir); err != nil {
			log.Printf("Reset task dir %s error: %v", taskDir, err)
		}
	}()

	ext := filepath.Ext(filename)
	var compileCmd *exec.Cmd

	switch ext {
	case ".cpp":
		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("g++ %s/%s -o %s/%s.out", taskDir, filename, taskDir, filename))
		if err := compileCmd.Run(); err != nil {
			return "Compile Failed"
		}
	case ".java":
		renameCmd := exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("mv %s/%s %s/Main.java", taskDir, filename, taskDir))
		if err := renameCmd.Run(); err != nil {
			return "Compile Failed"
		}
		// 编译Java代码
		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("javac %s/Main.java", taskDir))
		if err := compileCmd.Run(); err != nil {
			return "Compile Failed"
		}
	default:

	}

	// 设置超时上下文用于运行测试用例
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 从数据库中获取输入输出样例
	testCases := sql.SelectTestCasesByPid(strings.Split(filename, "_")[1])
	for _, testCase := range testCases {
		var runCmd *exec.Cmd
		switch ext {
		case ".cpp":
			// 执行编译生成的C++可执行文件
			runCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c",
				fmt.Sprintf("%s/%s.out", taskDir, filename))
		case ".py":
			runCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c",
				fmt.Sprintf("python %s/%s", taskDir, filename))
		case ".go":
			runCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c",
				fmt.Sprintf("go run %s/%s", taskDir, filename))
		case ".java":
			// 运行Java程序：指定任务目录作为classpath，执行Main类
			runCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c",
				fmt.Sprintf("java -cp %s Main", taskDir))
		default:
			return "Failed"
		}

		// 将测试用例的输入数据传入命令
		runCmd.Stdin = strings.NewReader(testCase.InputData)
		output, err := runCmd.CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			return "Time Limit Exceeded"
		}
		if err != nil {
			return "Failed"
		}
		outputStr := string(output)
		if strings.TrimSpace(outputStr) != strings.TrimSpace(testCase.OutputData) {
			return "Wrong Answer"
		}
	}

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
