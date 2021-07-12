package snowflake

import "github.com/gogf/gf/frame/g"

var node *Node

func Service() *Node {
	if node == nil {
		node, _ = NewNode(g.Cfg().GetInt64("api.node", 1))
	}
	return node
}
