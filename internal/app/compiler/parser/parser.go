package parser

import (
	"fmt"

	"github.com/tocoteron/9cc-go/internal/app/compiler/tokenizer"
)

type NodeKind int

const (
	NODE_ADD NodeKind = iota
	NODE_SUB
	NODE_MUL
	NODE_DIV
	NODE_NUM
)

type Node struct {
	kind NodeKind
	lhs  *Node
	rhs  *Node
	val  int
}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	node := &Node{}
	node.kind = kind
	node.lhs = lhs
	node.rhs = rhs

	return node
}

func newNodeNum(val int) *Node {
	node := &Node{}
	node.kind = NODE_NUM
	node.val = val

	return node
}

func Parse() *Node {
	return expression()
}

func expression() *Node {
	node := mul()

	for {
		if tokenizer.Consume('+') {
			node = newNode(NODE_ADD, node, mul())
		} else if tokenizer.Consume('-') {
			node = newNode(NODE_SUB, node, mul())
		} else {
			return node
		}
	}
}

func mul() *Node {
	node := unary()

	for {
		if tokenizer.Consume('*') {
			node = newNode(NODE_MUL, node, unary())
		} else if tokenizer.Consume('/') {
			node = newNode(NODE_DIV, node, unary())
		} else {
			return node
		}
	}
}

func unary() *Node {
	if tokenizer.Consume('+') {
		return primary()
	}

	if tokenizer.Consume('-') {
		return newNode(NODE_SUB, newNodeNum(0), primary())
	}

	return primary()
}

func primary() *Node {
	if tokenizer.Consume('(') {
		node := expression()
		tokenizer.Expect(')')
		return node
	}

	return newNodeNum(tokenizer.ExpectNumber())
}

func Generate(node *Node) {
	if node.kind == NODE_NUM {
		fmt.Printf("  push %d\n", node.val)
		return
	}

	Generate(node.lhs)
	Generate(node.rhs)

	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")

	switch node.kind {
	case NODE_ADD:
		fmt.Printf("  add rax, rdi\n")
		break
	case NODE_SUB:
		fmt.Printf("  sub rax, rdi\n")
		break
	case NODE_MUL:
		fmt.Printf("  imul rax, rdi\n")
		break
	case NODE_DIV:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
		break
	}

	fmt.Printf("  push rax\n")
}
