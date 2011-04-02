package compiler

import "container/vector"
import "util"
import . "ast"
import "strconv"

// import "fmt"

const FAILED = -1

type Error string
var (
    NoErr = Error("Successful")
    ErrRecognition  = Error("Recognition Error")
    ErrNoViableRule = Error("No viable rule")
    ErrParseFailed  = Error("Parse failed")
)
func (e Error) String() {
    return string(e)
}

type Parser struct {
    input     *Lexer
    lookahead *vector.Vector
    markers   *vector.Vector
    p         int

    listMemo  map[int]int
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
    this.listMemo = map[int]int{};
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

func (this *Parser) AlreadyParsedRule(memoization map[int]int) bool {
    memoI,ok := memoization[this.Index()]
    if ok == false {
        return false
    }
    if memoI == FAILED {
        panic(ErrParseFailed)
    }
    this.Seek(memoI)
    return true
}

func (this *Parser) Index() int { return this.p }

func (this *Parser) Memoize(memoization map[int]int, 
                            startTokenIndex int,
                            failed bool) {
    var stopTokenIndex int
    if failed {
        stopTokenIndex = FAILED
    } else {
        stopTokenIndex = this.Index()
    }
    memoization[startTokenIndex] = stopTokenIndex
}

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
    AT:       true,
    PUBLIC:   true,
    PROTECTED:true,
    STATIC:   true,
    ABSTRACT: true,
    FINAL:    true,
    NATIVE:   true,
    SYNC:     true,
    TRANSIENT:true,
    VOLATILE: true,
    STRICTFP: true,
}

//
// methodBodyDecl: { blockStatement* }
//
func (this *Parser) MethodBodyDecl() *Node {
    // println "methodBodyDecl"
    this.Match(LCURL);  for this.LA(1)==EOL { this.Match(EOL) }
    blockStmts := []*Node{}
    for this.LA(1) != RCURL {
        blockStmts = append(blockStmts, this.BlockStatement())
        this.semiOrEol()
    }
    this.Match(RCURL)
    return NewNode1("METHOD_BODY", blockStmts)
}

// blockStatement
//     :   localVariableDeclarationStatement
//     |   classOrInterfaceDeclaration
//     |   statement
func (this *Parser) BlockStatement() *Node {
    // a,b
    if this.LA(1) == IDENT && this.LA(2) == COMMA {
        return this.MultipleVarDeclStmt()
    // a :=
    } else if this.LA(1) == IDENT && this.LA(2) == COLON && this.LA(3) == EQUAL {
        return this.InferLocalVarDeclStmt()
    // a =
    } else if this.LA(1) == IDENT && this.LA(2) == EQUAL {
        return this.LocalVarDeclStmt()
    }
    //return this.Statement()
    return NewNode0("STMT", this.IDENT())
}

//
// a(,b)+ (:)?= expression
//
func (this *Parser) MultipleVarDeclStmt() *Node {
    return NewNode0("STMT", this.IDENT())
}

func (this *Parser) InferLocalVarDeclStmt() *Node {
    ident := this.IDENT()
    this.Match(COLON)
    this.Match(EQUAL)
    expr := this.Expression()
    return NewNode0("INFER_ASSIGN", ident, expr)
}

func (this *Parser) LocalVarDeclStmt() *Node {
    ident := this.IDENT()
    this.Match(EQUAL)
    expr := this.Expression()
    return NewNode0("ASSIGN", ident, expr)
}

// expression
//     :   conditionalExpression
//         (assignmentOperator expression
//         )?
func (this *Parser) Expression() *Node {
    c := this.ConditionalExpression()
    if this.LA(1) == EQUAL || this.LA(2) == EQUAL {
        a := this.AssignmentOperator()
        e := this.Expression()
        return NewNode0("EXPR", a, e)
    }
    return NewNode0("EXPR", c)
}

// conditionalExpression
//     :   conditionalOrExpression
//         ('?' expression ':' conditionalExpression
//         )?
func (this *Parser) ConditionalExpression() *Node {
    this.ConditionalOrExpression()
    if this.LA(1) == QUESTION {
        this.Match(QUESTION)
        this.Expression()
        this.Match(COLON)
        this.ConditionalExpression()
    }
    return nil
}

// conditionalOrExpression
//     :   conditionalAndExpression
//         ('||' conditionalAndExpression
//         )*
func (this *Parser) ConditionalOrExpression() *Node {
    this.ConditionalAndExpression()
    for this.LA(1) == OR && this.LA(2) == OR {
        this.ConditionalAndExpression()
    }
    return nil
}

// conditionalAndExpression
//     :   inclusiveOrExpression
//         ('&&' inclusiveOrExpression
//         )*
func (this *Parser) ConditionalAndExpression() *Node {
    this.InclusiveOrExpression()
    for this.LA(1) == AND && this.LA(2) == AND {
        this.InclusiveOrExpression()
    }
    return nil
}

// inclusiveOrExpression
//     :   exclusiveOrExpression
//         ('|' exclusiveOrExpression
//         )*
func (this *Parser) InclusiveOrExpression() *Node {
    this.ExclusiveOrExpression()
    for this.LA(1) == OR {
        this.ExclusiveOrExpression()
    }
    return nil
}

// exclusiveOrExpression
//     :   andExpression
//         ('^' andExpression
//         )*
func (this *Parser) ExclusiveOrExpression() *Node {
    this.AndExpression()
    for this.LA(1) == XOR {
        this.AndExpression()
    }
    return nil
}

// andExpression
//     :   equalityExpression
//         ('&' equalityExpression
//         )*
func (this *Parser) AndExpression() *Node {
    this.EqualityExpression()
    for this.LA(1) == AND {
        this.EqualityExpression()
    }
    return nil
}

// equalityExpression
//     :   instanceOfExpression
//         (
//             (   '=='
//             |   '!='
//             )
//             instanceOfExpression
//         )*
func (this *Parser) EqualityExpression() *Node {
    this.InstanceOfExpression()
    tok1 := this.LA(1)
    for tok1 == EQUAL || tok1 == NOT {
        if(tok1 == EQUAL) {
            this.Match(EQUAL)
        } else {
            this.Match(NOT)
        }
        this.Match(EQUAL)
        this.InstanceOfExpression()
    }
    return nil
}

// instanceOfExpression
//     :   relationalExpression
//         ('instanceof' type
//         )?
func (this *Parser) InstanceOfExpression() *Node {
    this.RelationalExpression()
    if this.LA(1) == INSTANCE_OF {
        this.Match(INSTANCE_OF)
        this.Type()
    }
    return nil
}

//
// relationalExpression
//     :   shiftExpression (relationalOp shiftExpression)*
//
func (this *Parser) RelationalExpression() *Node {
    this.ShiftExpression()
    for this.LA(1) == LANGLE || this.LA(1)== RANGLE {
        this.RelationalOp()
        this.ShiftExpression()
    }
    return nil
}

//
// relationalOp
//     :    '<' '='
//     |    '>' '='
//     |    '<'
//     |    '>'
//
func (this *Parser) RelationalOp() *Node {
    tok1,tok2 := this.LA(1),this.LA(2)
    if tok1 == LANGLE {
        this.Match(LANGLE)
        if tok2 == EQUAL {
            this.Match(EQUAL)
            return NewNode0("LESS_THAN_OR_EQUAL")
        }
        return NewNode0("LESS_THAN")
    } else if tok1 == RANGLE {
        this.Match(RANGLE)
        if tok2 == EQUAL {
            this.Match(EQUAL)
            return NewNode0("GREATER_THAN_OR_EQUAL")
        }
        return NewNode0("GREATER_THAN")
    }
    panic("Unreachable code")
}

// shiftExpression 
//     :   additiveExpression
//         (shiftOp additiveExpression
//         )*
func (this *Parser) ShiftExpression() *Node {
    this.AdditiveExpression()
    for this.LA(1) == LANGLE || this.LA(1)== RANGLE {
        this.ShiftOp()
        this.AdditiveExpression()
    }
    return nil
}

// shiftOp 
//     :    '<' '<'
//     |    '>' '>' '>'
//     |    '>' '>'
func (this *Parser) ShiftOp() *Node {
    if this.LA(1) == LANGLE {
        this.Match(LANGLE)
        this.Match(LANGLE)
        return NewNode0("SHL")
    } else {
        this.Match(RANGLE)
        this.Match(RANGLE)    
        return NewNode0("SHR")   
    }
    panic("Unreachable code")
}

// additiveExpression 
//     :   multiplicativeExpression
//         (   
//             (   '+'
//             |   '-'
//             )
//             multiplicativeExpression
//          )*
func (this *Parser) AdditiveExpression() *Node {
    this.MultiplicativeExpression()
    tok1 := this.LA(1)
    for tok1 == PLUS || tok1 == MINUS {
        if tok1 == PLUS {
            this.Match(PLUS)  
        } else {
            this.Match(MINUS)
        }
        this.MultiplicativeExpression()
        tok1 = this.LA(1)
    }
    return nil
}

// multiplicativeExpression 
//     :
//         unaryExpression
//         (   
//             (   '*'
//             |   '/'
//             |   '%'
//             )
//             unaryExpression
//         )*
func (this *Parser) MultiplicativeExpression() *Node {
    this.UnaryExpression()
    tok := this.LA(1)
    for tok == STAR || tok == DIV || tok == PERCENT {
        if tok == STAR {
            this.Match(STAR)
        } else if tok == DIV {
            this.Match(DIV)
        } else {
            this.Match(PERCENT)
        }
        this.UnaryExpression()
        tok = this.LA(1)
    }
    return nil
}

// unaryExpression 
//     :   '+'  unaryExpression
//     |   '-' unaryExpression
//     |   '++' unaryExpression
//     |   '--' unaryExpression
//     |   unaryExpressionNotPlusMinus
func (this *Parser) UnaryExpression() *Node {
    tok1, tok2 := this.LA(1), this.LA(2)
    if tok1 == PLUS {
        this.Match(PLUS)
        if tok2 == PLUS {
            this.Match(PLUS)
            return NewNode0("INC", this.UnaryExpression())
        }
        return NewNode0("U_PLUS", this.UnaryExpression())
    } else if tok1 == MINUS {
        this.Match(MINUS)
        if tok2 == MINUS {
            this.Match(MINUS)
            return NewNode0("DEC", this.UnaryExpression())
        }
        return NewNode0("U_MINUS", this.UnaryExpression())
    } else {
        return this.UnaryExpressionNotPlusMinus()        
    }
    panic("Unreachable code")    
}

func (this *Parser) _UnaryExpressionNotPlusMinus() *Node {
}

// unaryExpressionNotPlusMinus 
//     :   '~' unaryExpression
//     |   '!' unaryExpression
//     |   castExpression
//     |   primary (selector)* ('++' | '--')?
func (this *Parser) UnaryExpressionNotPlusMinus() *Node {
    tok1 := this.LA(1)
    //,tok2,tok3 := this.LA(1),this.LA(2),this.LA(3)
    if tok1 == TILD {
        this.Match(TILD)        
        return NewNode0("TILD", this.UnaryExpression())
    } else if tok1 == NOT {
        this.Match(NOT)
        return NewNode0("NOT", this.UnaryExpression())                
    } else if tok1 == LPAR {
        return this.Primary()        
    } else {
        return this.CastExpression()
    }
    panic("Unreachable code")    
}



// assignmentOperator
//     :   '='
//     |   '+='
//     |   '-='
//     |   '*='
//     |   '/='
//     |   '&='
//     |   '|='
//     |   '^='
//     |   '%='
//     |    '<' '<' '='
//     |    '>' '>' '>' '='
//     |    '>' '>' '='
//     ;
func (this *Parser) AssignmentOperator() *Node {
    t1 := this.LA(1)
    t2 := this.LA(2)
    // #TODO
    switch {
        case t1 == EQUAL:
            this.Match(EQUAL)
            return NewNode0("ASSIGN_OP")
        case t1 == PLUS && t2 == EQUAL:
            this.Match(PLUS); this.Match(EQUAL)
            return NewNode0("PLUS_ASSIGN_OP")
        case t1 == COLON && t2 == EQUAL:
            this.Match(COLON); this.Match(EQUAL)
            return NewNode0("INFER_ASSIGN_OP")
    }
    panic("assign op")
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
