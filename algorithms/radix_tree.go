package algorithm

import (
	"context"
)

type (
	nodeType      uint8
	HandlerFunc   func(context.Context)
	HandlersChain []HandlerFunc

	node struct {
		path      string
		indicies  string // 字符数组，用于对child 做索引
		children  []*node
		handlers  HandlersChain
		nType     nodeType
		maxParams uint8
		wildChild bool
	}
)

const (
	static nodeType = iota
	root
	param    // /:param
	catchAll // /*
)

func (n *node) AddRoute(path string, handlers HandlersChain) {

	if n.path == path {
		n.handlers = append(n.handlers, handlers)
		return
	}

	if len(n.children) == 0 && n.path == "" {
		n.InsertPath(path, handlers)
		return
	}

	publicPath := calcPublicPath(n.path, path)
	if publicPath == n.path {
		trimPath := path[len(publicPath):]

	}

	maxParams := calcParams(path)
}
