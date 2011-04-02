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
    AND
    AT
    CASE
    CLASS
    COLON
    COMMA
    DEFAULT
    DIV
    DOT
    EOL
    EQUAL
    FINAL
    IDENT    
    IMPORT
    INSTANCE_OF
    LANGLE
    LCURL
    LPAR
    LBRAC    
    MATCH
    MINUS
    NATIVE
    NOT
    OR
    PACKAGE
    PERCENT
    PLUS
    PROTECTED
    PUBLIC
    QNAME
    QUESTION
    RANGLE
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
    TILD
    VOLATILE
    XOR
)

var tokens = map[TokenType]string{
    NA:  "N/A",
    EOF: "<EOF>",
    EOL: "<EOL>",

    ABSTRACT: "abstract",
    AND:      "&",
    AT:       "@",
    CASE:     "case",
    CLASS:    "class",
    
    COLON:    ":",
    COMMA:    ",",    
    DEFAULT:  "default",
    DIV:      "/",
    DOT:      ".",
    FINAL:    "final",
    
    IDENT:    "<IDENT>",
    INSTANCE_OF: "instanceof",
    IMPORT:   "import",
    LCURL:    "{",
    LPAR:     "(",
    LBRAC:    "[",
    
    MATCH:    "match",
    MINUS:    "-",
    NATIVE:   "native",
    NOT:      "!",

    OR:        "|",
    PACKAGE:   "package",
    PERCENT:   "%",
    PLUS:      "+",
    PROTECTED: "protected",
    PUBLIC:    "public",
    QNAME:     "<QNAME>",
    QUESTION:  "?",
    RETURN:    "return",
    
    RCURL:     "}",
    RPAR:   ")",
    RBRAC:  "]",
    SEMI:   ";",
    STAR:   "*",
    STATIC: "static",

    STRICTFP:  "strictfp",
    SYNC:      "synchronized",
    TILD:      "~",
    TRANSIENT: "transient",
    VOLATILE:  "volatile",
    
    XOR:    "^",
}

func (t TokenType) String() string {
    if str, exists := tokens[t]; exists {
        return "(" + strconv.Itoa(int(t)) + "," + str + ")"
    }
    return "<" + strconv.Itoa(int(t)) + ">"
}