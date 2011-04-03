package compiler

import "utf8"
import "strconv"
import "util"

import "fmt"

const EOF = -1

type Lexer struct {
    input       []byte
    readOffset  int
    ch          int
    offset      int
}


//
// for testing purpose only
//
func (S *Lexer) GetInput() string {
    return string(S.input)
}

// char offset (utf-8)
func (S *Lexer) GetOffset() int {
    return S.offset
}

// byte offset
func (S *Lexer) GetReadOffset() int {
    return S.readOffset
}

func (S *Lexer) GetCh() int {
    return S.ch
}
//
// end testing
//


func (S *Lexer) Init(input string) *Lexer {
    S.input = []byte(input)
    S.readOffset = 0
    S.offset = 0
    S.advance()
    return S
}

func (S *Lexer) Consume() {
    S.advance()
}

func (S *Lexer) getChar() (ch int, w int) {
    ch,w = int(S.input[S.readOffset]), 1
    switch {
        case ch == 0:
            S.error("illegal 0")
        case ch >= 0x80:
            ch,w = utf8.DecodeRune(S.input[S.readOffset:])
            if ch == utf8.RuneError && w == 1 {
                S.error("illegal utf")
            }
    }
    return
}

func (S *Lexer) advance() {
    if S.readOffset < len(S.input) {
        ch,w := S.getChar()
        S.offset++
        S.readOffset += w
        S.ch = ch
    } else {
        S.ch = EOF
    }
}

func (S *Lexer) Match(x int) {
    if x == S.ch {
        S.Consume()
    } else {
        S.error("expect " + strconv.Itoa(x) + ", found " + strconv.Itoa(S.ch))
    }
}

func (S *Lexer) NextToken() *Token {
    for S.ch != EOF {
        switch S.ch {
            case ' ', '\t': S.WS()
            case '\r','\n': return S.EOL()
            case ';': S.Consume(); return &Token{tokenType: SEMI, text:";"}
            case '.': S.Consume(); return &Token{tokenType: DOT,  text:"."}
            case '{': S.Consume(); return &Token{tokenType: LCURL,text:"{"}
            case '}': S.Consume(); return &Token{tokenType: RCURL,text:"}"}
            case '(': S.Consume(); return &Token{tokenType: LPAR, text:"("}
            case ')': S.Consume(); return &Token{tokenType: RPAR, text:")"}
            case '[': S.Consume(); return &Token{tokenType: LBRAC,text:"["}
            case ']': S.Consume(); return &Token{tokenType: RBRAC,text:"]"}
            case '*': S.Consume(); return &Token{tokenType: STAR, text:"*"}
            case ',': S.Consume(); return &Token{tokenType: COMMA,text:","}
            case ':': S.Consume(); return &Token{tokenType: COLON,text:":"}
            case '?': S.Consume(); return &Token{tokenType: QUESTION,text:"?"}
            case '|': S.Consume(); return &Token{tokenType: OR,   text:"|"}
            case '&': S.Consume(); return &Token{tokenType: AND,  text:"&"}
            case '^': S.Consume(); return &Token{tokenType: XOR,  text:"^"}
            case '!': S.Consume(); return &Token{tokenType: NOT,  text:"!"}
        	case '=': S.Consume(); return &Token{tokenType: EQUAL,text:"="}
            default:
                if S.isLetter() {
                    return S.KeywordOrIdent()
                }
                S.error(fmt.Sprintf("invalid character: '%c' (%d)", S.ch, S.ch))
        }
    }
    return &Token{tokenType:EOF,text:"<EOF>"}
}

func (S *Lexer) isLetter() bool {
    ch := S.ch
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (S *Lexer) LETTER() {
    if S.isLetter() {
        S.Consume()
    } else {
        S.error("expecting LETTER; found " + strconv.Itoa(S.ch))
    }
}

func (S *Lexer) KeywordOrIdent() *Token {
    buf := util.NewStringBuffer()
    for S.isLetter()  {
        buf.Append(S.ch); S.LETTER()
    }
    str := buf.String()
    switch str {
        case "import":
            return &Token{tokenType:IMPORT, text:str}
        case "static":
            return &Token{tokenType:STATIC, text:str}
        case "package":
            return &Token{tokenType:PACKAGE,text:str}
        case "class":
            return &Token{tokenType:CLASS,  text:str}
        case "case":
            return &Token{tokenType:CASE,   text:str}
        case "match":
            return &Token{tokenType:MATCH,  text:str}
        case "return":
            return &Token{tokenType:RETURN, text:str}
        default:
            return &Token{tokenType:IDENT,  text:str}
    }
    return nil
}

func (S *Lexer) WS() {
    for S.ch ==' ' || S.ch =='\t' {
        S.advance()
    }
}

func (S *Lexer) EOL() *Token {
    if S.ch == '\r' {      // CR
        S.Consume()
        if S.ch == '\n' { // CRLF
            S.Consume()
        }
    } else if S.ch == '\n' {  // LF
        S.Consume()
    }
    return &Token{tokenType: EOL, text:"<EOL>"}
}

func (S *Lexer) error(msg string) {
    panic(msg)
}
