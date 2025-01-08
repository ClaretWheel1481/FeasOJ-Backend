package config

import "github.com/docker/docker/api/types/container"

/////////////////////////////////////////// 服务器配置 //////////////////////////////////////////////

// 后端模式
const DebugMode = true

// 服务器Address
const ServerAddress = "127.0.0.1:37882"

// 启用https(当该选项启用时，请确保下方证书与私钥路径已经填写，服务器会创建/certificate目录，请将证书与私钥放入其中)
const EnableHTTPS = true

// 服务器证书路径(./certificate/fullchain.pem)
const ServerCertPath = "./certificate/fullchain.pem"

// 服务器私钥路径(./certificate/privkey.key)
const ServerKeyPath = "./certificate/privkey.key"

/////////////////////////////////////////// MySQL配置 //////////////////////////////////////////////

// MySQL最大连接数
const MaxOpenConns = 240

// MySQL最大空闲连接数
const MaxIdleConns = 100

// MySQL连接最大生命周期（单位：秒）
const MaxLifeTime = 32

/////////////////////////////////////////// Docker配置 //////////////////////////////////////////////

// SandBox配置（Docker）
var SandBoxConfig = container.Resources{
	Memory:    512 * 1024 * 1024, // 512MB限制，1G：1024 * 1024 * 1024，2G：2 * 1024 * 1024 * 1024
	NanoCPUs:  0.5 * 1e9,         // 50%的一个CPU核心限制，1个CPU核心：1 * 1e9，2个CPU核心：2 * 1e9
	CPUShares: 1024,              // CPU权重，默认为1024，越高优先级越高
}

// SandBox最大并发数
const MaxWorkers = 5
