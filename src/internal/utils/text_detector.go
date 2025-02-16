package utils

import (
	"strings"
)

// BannedWords 脏话库
var BannedWords = []string{"fuck", "shit", "asshole", "dick", "pussy", "puta", "bitch", "nmsl", "cnm", "草你妈", "屌", "妈逼", "妈的", "操", "傻逼", "你妈死了", "狗逼", "xjp", "习近平"}

// TrieNode 表示 Trie 树节点结构
// children：存储当前节点的所有子节点，key 为字符（rune），value 为对应的 TrieNode 节点
// isEnd：标识该节点是否为一个完整词汇的结尾
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// ProfanityFilter 脏话过滤器，内部使用 Trie 树存储词库
type ProfanityFilter struct {
	root *TrieNode
}

// NewProfanityFilter 根据词库初始化一个新的脏话过滤器
func NewProfanityFilter(words []string) *ProfanityFilter {

	// 初始化根节点，并创建过滤器实例
	filter := &ProfanityFilter{root: &TrieNode{children: make(map[rune]*TrieNode)}}

	// 将每个词加入 Trie 树中
	for _, word := range words {
		filter.AddWord(word)
	}
	return filter
}

// AddWord 将一个词汇添加到 Trie 树中
func (pf *ProfanityFilter) AddWord(word string) {
	node := pf.root

	// 转换为小写，确保匹配不区分大小写
	for _, ch := range strings.ToLower(word) {

		// 如果子节点不存在则创建
		if node.children == nil {
			node.children = make(map[rune]*TrieNode)
		}
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}

		// 移动到下一个子节点
		node = node.children[ch]
	}

	// 标记该节点为词汇结束
	node.isEnd = true
}

// ContainsProfanity 检测输入文本中是否包含脏话
func (pf *ProfanityFilter) ContainsProfanity(text string) bool {
	// 将文本转换为小写
	text = strings.ToLower(text)
	// 将字符串转换为 rune 切片，正确处理多字节字符（如中文）
	runes := []rune(text)

	// 遍历文本中每个字符，将其作为可能的匹配起点
	for i := 0; i < len(runes); i++ {
		node := pf.root

		// 逐字符尝试在 Trie 树中匹配脏话
		for j := i; j < len(runes); j++ {
			ch := runes[j]
			if next, ok := node.children[ch]; ok {

				// 如果存在，则更新当前节点为子节点
				node = next

				// 如果当前节点标记为词尾，说明找到了完整的脏话词汇
				if node.isEnd {
					return true
				}
			} else {
				// 当前路径不存在，结束本轮匹配
				break
			}
		}
	}
	return false
}

// defaultFilter 是默认的脏话过滤器，基于全局的 BannedWords 初始化
var defaultFilter = NewProfanityFilter(BannedWords)

// ContainsProfanity 检测输入文本中是否包含脏话，直接调用默认过滤器
func ContainsProfanity(text string) bool {
	return defaultFilter.ContainsProfanity(text)
}
