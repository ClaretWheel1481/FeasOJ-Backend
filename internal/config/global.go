package config

import "github.com/docker/docker/api/types/container"

/////////////////////////////////////////// 服务器配置 //////////////////////////////////////////////

// ServerAddress 服务器Address
const ServerAddress = "127.0.0.1:37882"

// RabbitMQAddress RabbitMQ地址
const RabbitMQAddress = "amqp://rabbitmq:Sing5200@cloud.claret.space:5672/"

// ConsulAddress Consul地址
const ConsulAddress = "http://localhost:8500"

// EnableHTTPS 启用https(当该选项启用时，请确保下方证书与私钥路径已经填写，服务器会创建/certificate目录，请将证书与私钥放入其中)
const EnableHTTPS = false

// ServerCertPath 服务器证书路径(./certificate/fullchain.pem)
const ServerCertPath = "./certificate/fullchain.pem"

// ServerKeyPath 服务器私钥路径(./certificate/privkey.key)
const ServerKeyPath = "./certificate/privkey.key"

// ImageGuardEnabled ImageGuard是否启用
const ImageGuardEnabled = true

// ProfanityDetectorEnabled ProfanityDetector是否启用
const ProfanityDetectorEnabled = true

/////////////////////////////////////////// MySQL配置 //////////////////////////////////////////////

// MaxOpenConns MySQL最大连接数
const MaxOpenConns = 240

// MaxIdleConns MySQL最大空闲连接数
const MaxIdleConns = 100

// MaxLifeTime MySQL连接最大生命周期（单位：秒）
const MaxLifeTime = 32

/////////////////////////////////////////// Docker配置 //////////////////////////////////////////////

// SandBoxConfig SandBox配置（Docker）
var SandBoxConfig = container.Resources{
	Memory:    512 * 1024 * 1024, // 512MB限制，1G：1024 * 1024 * 1024，2G：2 * 1024 * 1024 * 1024
	NanoCPUs:  0.5 * 1e9,         // 50%的一个CPU核心限制，1个CPU核心：1 * 1e9，2个CPU核心：2 * 1e9
	CPUShares: 1024,              // CPU权重，默认为1024，越高优先级越高
}

// MaxSandbox SandBox最大并发数
const MaxSandbox = 5
