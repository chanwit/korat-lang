package ast

type Node struct {
    Name     string
    Children []*Node
    Text     string
}

func NewNode0(name string, nodes...*Node) *Node {
    return &Node{Name: name, Text: "", Children: []*Node(nodes)}
}

func NewNode1(name string, nodes []*Node) *Node {
    return &Node{Name: name, Text: "", Children: nodes}
}

func NewNode2(name string, text string) *Node {
    return &Node{Name: name, Text: text}
}

func f(nodes []*Node) (r string) {
    if(len(nodes) == 0) { r = ""; return }
    r = nodes[0].String()
    for i := 1; i < len(nodes); i++ {
        r = r + "," + nodes[i].String()
    }
    return
}

func (this *Node) String() string {    
    if(this == nil) { return "<nil>" }
    if len(this.Text) == 0 {
        if this.Children == nil || len(this.Children) == 0{
            return this.Name
        }
        return this.Name + "(" + f(this.Children) + ")"
    }
    if this.Children == nil || len(this.Children) == 0 {
        return this.Name + "('" + this.Text + "')"
    }
    return this.Name + "(" + this.Text + "," + f(this.Children) + ")"
}