package compiler

import "strconv"

type TokenType int

type Token struct {
    tokenType TokenType
    text      string   
}

//
// for testing purpose
//
func (t *Token) GetTokenType() TokenType {
    return t.tokenType
}

func (t *Token) GetText() string {
    return t.text
}
//
// for testing purpose
//

const (
    NA        TokenType = iota

    ABSTRACT
    AT
    CASE
    CLASS
    COLON
    COMMA
    DEFAULT
    DOT
    EOL
    EQUAL
    FINAL
    IDENT    
    IMPORT
    LCURL
    LPAR
    LBRAC    
    MATCH
    NATIVE
    PACKAGE
    PROTECTED
    PUBLIC
    QNAME
    RCURL
    RPAR
    RBRAC
    RETURN
    SEMI
    STAR
    STATIC
    STRICTFP
    SYNC
    TRANSIENT
    VOLATILE
)

var tokens = map[TokenType]string{
    NA:  "N/A",
    EOF: "<EOF>",
    EOL: "<EOL>",

    ABSTRACT: "abstract",
    AT:       "@",
    CASE:     "case",
    CLASS:    "class",
    COLON:    ":",

    COMMA:    ",",    
    DEFAULT:  "default",
    DOT:      ".",
    FINAL:    "final",
    IDENT:    "<IDENT>",
        
    IMPORT:   "import",
    LCURL:    "{",
    LPAR:     "(",
    LBRAC:    "[",
    MATCH:    "match",
    NATIVE:   "native",

    PACKAGE:   "package",
    PROTECTED: "protected",
    PUBLIC:    "public",
    QNAME:     "<QNAME>",
    RETURN: "return",
    
    RCURL:     "}",
    RPAR:   ")",
    RBRAC:  "]",
    SEMI:   ";",
    STAR:   "*",
    STATIC: "static",

    STRICTFP:  "strictfp",
    SYNC:      "synchronized",
    TRANSIENT: "transient",
    VOLATILE:  "volatile",
}

func (t TokenType) String() string {
    if str, exists := tokens[t]; exists {
        return "(" + strconv.Itoa(int(t)) + "," + str + ")"
    }
    return "<" + strconv.Itoa(int(t)) + ">"
}