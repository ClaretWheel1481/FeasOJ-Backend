package global

import "sync"

// 全局变量 - 本地配置文件路径
var ParentDir string
var ConfigDir string
var LogDir string
var AvatarsDir string
var CurrentDir string
var CertDir string
var DocsDir string

// 全局变量 - 容器ID
var ContainerIDs sync.Map

// 全局变量 - 微服务地址
var ImageGuardAddr string
var ProfanityDetectorAddr string
var JudgeCoreAddr string
