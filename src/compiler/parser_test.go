package compiler_test

import "testing"
import "compiler"

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