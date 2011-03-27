package korat_test

import "testing"

import "utf8"
import "strconv"
// import "fmt"

import "korat"

func TestInitLexer(t *testing.T) {
    l := new(korat.Lexer).Init("package a")

    if l.GetInput() != "package a" {
        t.Fatalf("l.input failed")
    }
    if l.GetOffset() != 1 {
        t.Fatalf("l.readOffset failed != 1")
    }
    if l.GetCh() != 'p' {
        t.Fatalf("l.readOffset failed != p")
    }
}

func TestConsume(t *testing.T) {
    
    s := "package นี้"
    l := new(korat.Lexer).Init(s)
    if l.GetInput() != "package นี้" {
        t.Fatalf("l.input failed")
    }    

    str := utf8.NewString(s)    
    for i := 0; i < str.RuneCount(); i++ {
        if l.GetOffset() != i+1 {
            t.Fatalf("l.readOffset failed != " + strconv.Itoa(i+1) )
        }
        if l.GetCh() != str.At(i) {
            t.Fatalf("ch " + strconv.Itoa(l.GetCh()) + " != at " + strconv.Itoa(str.At(i)))
        }        
        l.Consume()        
    }
    if l.GetCh() != korat.EOF {
        t.Fatalf("EOF not found")
    }       
}

func TestToken(t *testing.T) {
    l := new(korat.Lexer).Init("package")
    tok := l.NextToken()
    // fmt.Printf("%s\n", tok)
    if tok.GetTokenType() != korat.PACKAGE {
        t.Fatalf("Fail : PACKAGE not found")
    }
}