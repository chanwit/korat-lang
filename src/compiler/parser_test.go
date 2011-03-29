package compiler_test

import "testing"
import "compiler"

import "fmt"

func TestParsingPackage(t *testing.T) {
    lexer  := new(compiler.Lexer).Init("package a.b.c")
    parser := new(compiler.Parser).Init(lexer)
    node   := parser.PackageDecl()
    if node.Name != "PACKAGE" {
        t.Fatalf("Package not parsed")
    }
    if node.Children[0].Name != "QNAME" {
        t.Fatalf("QNAME not parsed")
    }
    if node.Children[0].Text != "a.b.c" {
        t.Fatalf("'a.b.c' not parsed")
    }
}

func TestParsingMultiline(t *testing.T) {
    lexer  := new(compiler.Lexer).Init("package a.b.c\n")
    parser := new(compiler.Parser).Init(lexer)
    node   := parser.PackageDecl()
    if node.Name != "PACKAGE" {
        t.Fatalf("Package not parsed")
    }
    if node.Children[0].Name != "QNAME" {
        t.Fatalf("QNAME not parsed")
    }
    if node.Children[0].Text != "a.b.c" {
        t.Fatalf("'a.b.c' not parsed")
    }
}

func TestBlankType(t *testing.T) {
    lexer  := new(compiler.Lexer).Init("class A { }")
    parser := new(compiler.Parser).Init(lexer)
    node   := parser.TypeDecl()
    if node.Name != "CLASS" {
        t.Fatalf("CLASS not parsed")
    }
    if node.Children[0].Name != "IDENT" {
        t.Fatalf("IDENT not parsed")
    }
    if node.Children[0].Text != "A" {
        t.Fatalf("A not found")
    }
}

func TestTypeWithMainMethod(t *testing.T) {
    lexer  := new(compiler.Lexer).Init(
    "\n"                      +
    "class A {\n"             +
    "   static main(args){\n" +
    "   }\n"                  +
    "}\n"                     )
    parser := new(compiler.Parser).Init(lexer)
    node   := parser.TypeDecl()
    // if node.String() != "CLASS(IDENT('A'),MEMBERS(METHOD(MODIFIERS(STATIC),<nil>,IDENT('main'),ARGS(ARG(TYPE('java.lang.Object'),IDENT('args'),<nil>)),METHOD_BODY)))" {
    //            
    //     t.Fatalf("CLASS not parsed")
    // }
    if node.F("IDENT").String() != "IDENT('A')" {
        fmt.Printf("%s\n", node)
        t.Fatalf("IDENT not parsed")
    }
    if node.F("MEMBERS").Name != "MEMBERS" {
        t.Fatalf("MEMBERS not parsed")
    }
    if node.F("MEMBERS").At(0).String() != "METHOD(MODIFIERS(STATIC),<nil>,IDENT('main'),ARGS(ARG(TYPE('java.lang.Object'),IDENT('args'),<nil>)),METHOD_BODY)" {
        t.Fatalf("METHOD not parsed")
    }
}