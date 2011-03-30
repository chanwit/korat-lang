package compiler

import "container/vector"
import "util"
import . "ast"
import "strconv"

// import "fmt"

const FAILED = -1

type Parser struct {
    input     *Lexer
    lookahead *vector.Vector
    markers   *vector.Vector
    listMemo  map[int]*int
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
        this.clearMemo()
    }
    this.sync(1)
}

func (this *Parser) clearMemo() {
    this.listMemo = map[int]*int{};
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

func (this *Parser) Match(x TokenType) (t *Token) {
    t = this.LT(1)
    if t.tokenType == x {
        this.Consume()
        // fmt.Printf("matched %s\n", t)
        return
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

//
// Production Rules
//
//
func (this *Parser) CompilationUnit() *Node {
    for this.LA(1) == EOL { this.Match(EOL) }
    if  this.LA(1) == AT && this.LA(2) == IDENT {
        // TODO: annotations()
    }
    if  this.LA(1) == PACKAGE {
        packageDecl := this.PackageDecl()
        return NewNode0("UNIT", packageDecl)
    }
    if  this.LA(1) == IMPORT {
        this.ImportDecls()
    }
    this.TypeDecls()
    return nil
}

func (this *Parser) PackageDecl() *Node {
    this.Match(PACKAGE)
    qname := this.QNAME()
    return NewNode0("PACKAGE", qname)
}

func (this *Parser) ImportDecls() *Node {
    for this.LA(1)==SEMI || this.LA(1)==EOL {
        // this.semiOrEol()
    }
    imports := []*Node{}
    imports = append(imports, this.ImportDecl())
    for this.LA(1) == IMPORT {
        imports = append(imports, this.ImportDecl())
    }
    return NewNode1("IMPORTS", imports)
}

func (this *Parser) ImportDecl() *Node {
    this.Match(IMPORT)
    foundStatic := false
    if this.LA(1) == STATIC {
        foundStatic = true
        this.Match(STATIC)
    }
    qname := this.QnameForImport()
    this.semiOrEol()
    if foundStatic {
        return NewNode0("IMPORT_STATIC", qname)
    }
    return NewNode0("IMPORT", qname)
}

func (this *Parser) TypeDecls() *Node {
    return nil
}

func (this *Parser) IDENT() *Node {
    return NewNode2("IDENT", this.Match(IDENT).text)
}

// typeDecl: 'class' name '{' '}'
func (this *Parser) TypeDecl() *Node {
    for this.LA(1)==EOL { this.Match(EOL) }

    caseClass := false

    if this.LA(1) == CASE {
        caseClass = true
        this.Match(CASE)
    }
    this.Match(CLASS)
    name := this.IDENT(); for this.LA(1)==EOL { this.Match(EOL) }
    this.Match(LCURL)
    members := this.Members()
    this.Match(RCURL)

    if(caseClass) {
        return NewNode0("CASE_CLASS", name, members)
    }
    return NewNode0("CLASS", name, members)
}

func (this *Parser) semiOrEol() {
    switch this.LA(1)  {
        case SEMI: this.Match(SEMI); return
        case EOL:  this.Match(EOL);  return
        case EOF:  this.Match(EOF);  return
    }
    panic("Expect semi-colon or EOL")
}

func (this *Parser) Members() *Node {
    members := []*Node{}
    for this.LA(1)==SEMI || this.LA(1)==EOL {
        this.semiOrEol()
    }
    for this.LA(1) != RCURL {
        members = append(members, this.MemberDecl())
    }

    for this.LA(1)==EOL { this.Match(EOL) }

    return NewNode1("MEMBERS", members)
}

func (this *Parser) MemberDecl() *Node {
    //if(LA(1) in MODIFIER_LIST)
    //fieldDecl()
    return this.MethodDecl()
}

func (this *Parser) MethodDecl() *Node  {
    for this.LA(1)==EOL { this.Match(EOL) }

    modifiers := this.Modifiers()

    var returnType *Node = nil
    if this.LA(2) == IDENT {
        returnType = this.Type()
    }
    methodName := this.IDENT()

    this.Match(LPAR)
    argDecls := this.ArgumentDecls()
    this.Match(RPAR);  for this.LA(1)==EOL { this.Match(EOL) }

    var body *Node = nil
    if this.LA(1) == LCURL {
        body = this.MethodBodyDecl()
    }

    for this.LA(1)==SEMI || this.LA(1)==EOL { this.semiOrEol() }

    if body == nil {
        return NewNode0("INTERFACE_METHOD", modifiers, returnType, methodName, argDecls)
    }
    return NewNode0("METHOD", modifiers, returnType, methodName, argDecls, body)
}

//
// type: qname []*
//
func (this *Parser) Type() *Node {
    qname := this.QNAME()
    dim := 0
    if this.LA(1) == LBRAC {        
        this.Match(LBRAC)
        this.Match(RBRAC)
        dim++
    }
    if dim > 0 {
        return NewNode3("TYPE", qname.Text, NewNode2("DIM", strconv.Itoa(dim)))
    } 
    return NewNode2("TYPE", qname.Text)
}

var modifiers = map[TokenType]bool {
    AT:true,
    PUBLIC:true,
    PROTECTED: true,
    STATIC: true,
    ABSTRACT: true,
    FINAL: true,
    NATIVE: true,
    SYNC: true,
    TRANSIENT: true,
    VOLATILE: true,
    STRICTFP: true,
}

func (this *Parser) MethodBodyDecl() *Node {
    // println "methodBodyDecl"
    this.Match(LCURL);  for this.LA(1)==EOL { this.Match(EOL) }
    this.Match(RCURL)
    return NewNode0("METHOD_BODY") // #TODO
}

// 
// modifiers: modifier*
//
func (this *Parser) Modifiers() *Node {
    m := []*Node{}
    for modifiers[this.LA(1)] == true {
        m = append(m, this.Modifier())
    }
    return NewNode1("MODIFIERS", m)
}

func (this *Parser) Modifier() *Node {
    switch this.LA(1) {
        case AT:        return this.Annotation()
        case PUBLIC:    this.Match(PUBLIC)    ; return NewNode0("PUBLIC")
        case PROTECTED: this.Match(PROTECTED) ; return NewNode0("PROTECTED")
        case STATIC:    this.Match(STATIC)    ; return NewNode0("STATIC") 
        case ABSTRACT:  this.Match(ABSTRACT)  ; return NewNode0("ABSTRACT")
        case FINAL:     this.Match(FINAL)     ; return NewNode0("FINAL")  
        case NATIVE:    this.Match(NATIVE)    ; return NewNode0("NATIVE")
        case SYNC:      this.Match(SYNC)      ; return NewNode0("SYNC")    
        case TRANSIENT: this.Match(TRANSIENT) ; return NewNode0("TRANSIENT")
        case VOLATILE:  this.Match(VOLATILE)  ; return NewNode0("VOLATILE") 
        case STRICTFP:  this.Match(STRICTFP)  ; return NewNode0("STRICTFP")

        default:
            panic("expecting a modifier, found " + this.LA(1).String())
    }
    return nil
}

func (this *Parser) ArgumentDecls() *Node {
    if this.LA(1) == RPAR { return NewNode0("ARGS") }

    a := []*Node{}
    a = append(a, this.ArgumentDecl())
    for this.LA(1) == COMMA {
        this.Match(COMMA)
        a = append(a, this.ArgumentDecl())
    }
    return NewNode1("ARGS", a)
}


var DEFAULT_TYPE = &Node{Name:"TYPE", Text:"java.lang.Object"}

func (this *Parser) ArgumentDecl() *Node {
    var annotations *Node = nil
    if this.LA(1) == AT {
        annotations = this.Annotations()
    }
    argType := DEFAULT_TYPE
    if this.LA(2) == LBRAC || this.LA(2) == IDENT {
        argType = this.Type()
    }
    name := this.IDENT()
    return NewNode0("ARG", argType, name, annotations)
}

func (this *Parser) QNAME() *Node {
    sb := util.NewStringBuffer()
    sb.AppendStr(this.Match(IDENT).text)
    for this.LA(1) == DOT {
        sb.AppendStr(this.Match(DOT).text)
        sb.AppendStr(this.Match(IDENT).text)
    }
    return NewNode2("QNAME", sb.String())
}

func (this *Parser) QnameForImport() *Node {
    sb := util.NewStringBuffer()
    sb.AppendStr(this.Match(IDENT).text)
    for this.LA(1) == DOT {
        sb.AppendStr(this.Match(DOT).text)
        if this.LA(1) == STAR {
            sb.AppendStr(this.Match(STAR).text)
        } else {
            sb.AppendStr(this.Match(IDENT).text)
        }
    }
    return NewNode2("QNAME", sb.String())
}


//
// annotations: annotation*
//
func (this *Parser) Annotations() *Node {
    anns := []*Node{}
    anns = append(anns, this.Annotation())
    for this.LA(1) == AT {
        anns = append(anns, this.Annotation())
    }
    return NewNode1("ANNOTATIONS", anns)
}

//
// annotation: '@' ident '(' args ')'
//
func (this *Parser) Annotation() *Node {
    this.Match(AT)
    this.Match(IDENT)
    if this.LA(1) == LPAR {
        this.Match(LPAR)
        // annotationArgs()
        this.Match(RPAR)
    }
    return NewNode0("ANNOTATION" /* #TODO */)
}
