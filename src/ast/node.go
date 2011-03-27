package ast

type Node struct {
    Name     string
    Children []Node
    Text     string
}

func (this *Node) String() string {
    return "(" + this.Name + ")"
}