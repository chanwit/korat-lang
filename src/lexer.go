package korat

import "utf8"
import "strconv"
// import "fmt"

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
func (lex *Lexer) GetInput() string {
    return string(lex.input)
}

// char offset (utf-8)
func (lex *Lexer) GetOffset() int {
    return lex.offset
}

// byte offset
func (lex *Lexer) GetReadOffset() int {
    return lex.readOffset
}

func (lex *Lexer) GetCh() int {
    return lex.ch
}
//
// end testing
//


func (lex *Lexer) Init(input string) *Lexer {
    lex.input = []byte(input)
    lex.readOffset = 0
    lex.offset = 0    
    lex.advance()
    return lex
}

func (lex *Lexer) Consume() {
    lex.advance()
}

func (lex *Lexer) GetChar() (ch int, w int) {
    ch,w = int(lex.input[lex.readOffset]), 1
    switch {
        case ch == 0:
            lex.error("illegal 0")
        case ch >= 0x80:
            ch,w = utf8.DecodeRune(lex.input[lex.readOffset:])
            if ch == utf8.RuneError && w == 1 {
                lex.error("illegal utf")
            }
    }
    return
}

func (lex *Lexer) advance() {
    if lex.readOffset < len(lex.input) {        
        ch,w := lex.GetChar()
        lex.offset++
        lex.readOffset += w
        lex.ch = ch
    } else {
        lex.ch = EOF
    }
}

func (lex *Lexer) Match(x int) {
    if x == lex.ch {
        lex.Consume()
    } else {
        lex.error("expect " + strconv.Itoa(x) + ", found " + strconv.Itoa(lex.ch))
    }
}

func (lex *Lexer) NextToken() *Token {    
    for lex.ch != EOF {
        switch lex.ch {
            case ' ', '\t': lex.ws()
            case '\r','\n': return lex.eol()
            case ';': return &Token{tokenType: SEMI, text: ";"}
            default:
                if lex.isLetter() {
                    return lex.identifier()
                }
                lex.error("invalid character " + strconv.Itoa(lex.ch))
        }
    }
    return nil
}

func (lex *Lexer) isLetter() bool {
    // fmt.Printf(">>>> isLetter\n")
    ch := lex.ch
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (lex *Lexer) letter() {
    if lex.isLetter() {
        lex.Consume()
    } else {
        lex.error("expecting LETTER; found " + strconv.Itoa(lex.ch))
    }
}

func (lex *Lexer) identifier() *Token {
    buf := NewStringBuffer()
    for lex.isLetter()  {
        buf.Append(lex.ch); lex.letter()
    }
    str := buf.String()
    switch str {
        case "import":
            return &Token{tokenType: IMPORT, text: str}
        case "static":
            return &Token{tokenType: STATIC, text: str}
        case "package":
            return &Token{tokenType: PACKAGE,text: str}
        case "class":
            return &Token{tokenType: CLASS,  text: str}
        case "case":
            return &Token{tokenType: CASE,   text: str}
        case "match":
            return &Token{tokenType: MATCH,  text: str}
        case "return":
            return &Token{tokenType: RETURN, text: str}
        default:
            return &Token{tokenType: NAME,   text: str}
    }
    return nil
}

func (lex *Lexer) ws() {

}

func (lex *Lexer) eol() *Token {
    return nil
}

func (lex *Lexer) error(msg string) {
}