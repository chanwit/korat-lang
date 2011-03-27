package compiler

import "container/vector"
import "ast"
import "util"

// import "fmt"

const FAILED = -1

type Parser struct {
    input     *Lexer
    lookahead *vector.Vector
    markers   *vector.Vector
    p         int
}

func (this *Parser) Init(input *Lexer) *Parser {
    this.input = input
    this.lookahead = new(vector.Vector)
    this.markers   = new(vector.Vector)
    this.sync(1)
    return this
}

func (this *Parser) Consume() {
    this.p++
    if this.p == this.lookahead.Len() && !this.IsSpeculating() {
        this.p = 0
        this.lookahead = new(vector.Vector)
        // this.clearMemo()
    }
    this.sync(1)
}

func (this *Parser) sync(i int) {
    // fmt.Printf("sync: %d\n", i)
    if this.p+i-1 > (this.lookahead.Len()-1) {
        // fmt.Printf("sync: under if\n")
        n := (this.p+i-1) - (this.lookahead.Len() - 1)
        // fmt.Printf("sync: n := %d\n", n)
        this.fill(n)
    }
}

func (this *Parser) fill(n int) {
    for i:=1; i<=n; i++ {
        // fmt.Printf("fill i := %d\n", i)
        this.lookahead.Push(this.input.NextToken())
    }
    // for i := 0; i<this.lookahead.Len(); i++ {
    //     v,_ := this.lookahead.At(i).(*Token)
    //     fmt.Printf("this.lookahead[%d] = %s\n", i, v)
    // }
}

func (this *Parser) LT(i int) *Token {
    // fmt.Printf("LT(%d)\n", i)
    this.sync(i)
    v,_ := this.lookahead.At(this.p+i-1).(*Token)
    // fmt.Printf("v=%s\n", v)
    return v
}

func (this *Parser) LA(i int) TokenType {
    return this.LT(i).tokenType
}

func (this *Parser) Match(x TokenType) *Token {
    t := this.LT(1)
    if t.tokenType == x {
        this.Consume()
        return t
    }
    panic("Token not match. Found:" + t.tokenType.String() + ", Expect: " + tokens[x])
}

func (this *Parser) Mark() int {
    this.markers.Push(this.p)
    return this.p
}

func (this *Parser) Release() {
    marker,_ := this.markers.At(this.markers.Len()-1).(int)
    this.markers.Delete(this.markers.Len()-1)
    this.Seek(marker)
}

func (this *Parser) Seek(index int) {
    this.p = index
}

func (this *Parser) IsSpeculating() bool {
    return this.markers.Len() > 0
}

func (this *Parser) AlreadyParsedRule(memoization map[int]*int) bool {
    memoI := memoization[this.Index()]
    if memoI == nil {
        return false
    }
    if *memoI == FAILED {
        panic("*memoI == FAILED")
    }
    this.Seek(*memoI)
    return true
}

func (this *Parser) Index() int { return this.p }

func (this *Parser) PackageDecl() *ast.Node {
    this.Match(PACKAGE)
    qname := this.QNAME()
    return &ast.Node{Name:"PACKAGE", Children:[]ast.Node{*qname}}
}

func (this *Parser) QNAME() *ast.Node {
    sb := util.NewStringBuffer()
    sb.AppendStr(this.Match(IDENT).text)
    for this.LA(1) == DOT {
        sb.AppendStr(this.Match(DOT).text)
        sb.AppendStr(this.Match(IDENT).text)
    }    
    return &ast.Node{Name:"QNAME", Text: sb.String()}
}
