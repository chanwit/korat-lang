package ast

type Node struct {
    Name     string
    Children []*Node
    Text     string
}

func NewNode0(name string, nodes...*Node) *Node {
    return &Node{Name: name, Children: []*Node(nodes)}
}

func NewNode1(name string, nodes []*Node) *Node {
    return &Node{Name: name, Children: nodes}
}

func NewNode2(name string, text string) *Node {
    return &Node{Name: name, Text: text}
}

func (this *Node) String() string {
    return "(" + this.Name + ")"
}