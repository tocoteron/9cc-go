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
	NODE_EQ // ==
	NODE_NE // !=
	NODE_LT // <
	NODE_LE // <=
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
	return equality()
}

func equality() *Node {
	node := relational()

	for {
		if tokenizer.Consume("==") {
			node = newNode(NODE_EQ, node, relational())
		} else if tokenizer.Consume("!=") {
			node = newNode(NODE_NE, node, relational())
		} else {
			return node
		}
	}
}

func relational() *Node {
	node := add()

	for {
		if tokenizer.Consume("<") {
			node = newNode(NODE_LT, node, add())
		} else if tokenizer.Consume("<=") {
			node = newNode(NODE_LE, node, add())
		} else if tokenizer.Consume(">") {
			node = newNode(NODE_LT, add(), node)
		} else if tokenizer.Consume(">=") {
			node = newNode(NODE_LE, add(), node)
		} else {
			return node
		}
	}
}

func add() *Node {
	node := mul()

	for {
		if tokenizer.Consume("+") {
			node = newNode(NODE_ADD, node, mul())
		} else if tokenizer.Consume("-") {
			node = newNode(NODE_SUB, node, mul())
		} else {
			return node
		}
	}
}

func mul() *Node {
	node := unary()

	for {
		if tokenizer.Consume("*") {
			node = newNode(NODE_MUL, node, unary())
		} else if tokenizer.Consume("/") {
			node = newNode(NODE_DIV, node, unary())
		} else {
			return node
		}
	}
}

func unary() *Node {
	if tokenizer.Consume("+") {
		return primary()
	}

	if tokenizer.Consume("-") {
		return newNode(NODE_SUB, newNodeNum(0), primary())
	}

	return primary()
}

func primary() *Node {
	if tokenizer.Consume("(") {
		node := expression()
		tokenizer.Expect(")")
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
	case NODE_EQ:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  sete al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	case NODE_NE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setne al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	case NODE_LT:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setl al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	case NODE_LE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setle al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	}

	fmt.Printf("  push rax\n")
}
