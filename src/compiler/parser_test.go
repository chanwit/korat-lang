package compiler_test

import "testing"
import "compiler"

import "os"
// import "fmt"
import . "util"
import "reflect"
import "strings"
import "sut/govy"
import "ast"

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

    class := CLASS(IDENT("A"),MEMBERS())
    if node.String() != class {
        t.Fatalf("CLASS not parsed")
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

    args   := ARGS(ARG(TYPE("java.lang.Object"),IDENT("args"),NIL))
    mainMethod := METHOD(
            MODIFIERS(STATIC),
            NIL,
            IDENT("main"),
            args,
            "METHOD_BODY")
    class  := CLASS(IDENT("A"),MEMBERS(mainMethod))

    if node.String() != class {
        t.Fatalf("CLASS not parsed")
    }
}

func InvokeByName(obj interface{}, name string, in...reflect.Value) []reflect.Value {
    ov := reflect.NewValue(obj)
    ot := ov.Type()
    var fn *reflect.FuncValue = nil
    for i := 0; i < ot.NumMethod(); i++ {
        if ot.Method(i).Name == name {
            fn = ot.Method(i).Func
            break
        }
    }
    if fn == nil {
        return []reflect.Value{reflect.NewValue(nil)}
    }
    params := make([]reflect.Value, len(in)+1)
    params[0] = ov
    for i := 1; i < len(params); i++ {
        params[i] = in[i-1]
    }
    return fn.Call(params)
}

func RunOneFile(filename string, t *testing.T) {
    f := govy.Open(filename)
    defer f.Close()

    ruleName,_ := f.NextLine() // TypeDecl:
    if ruleName[len(ruleName)-1] != ':' {
        panic("Rule not found")
    }
    ruleName = ruleName[0:len(ruleName)-1]
    // fmt.Printf("%s\n", ruleName)
    src := ""
    for {
        s,_ := f.NextLine()
        if s == "expect:" {
            break
        }
        src = src + s + "\n"
    }
    // fmt.Printf("%s\n", src)

    expect := ""
    for {
        s,e := f.NextLine()
        if e == os.EOF {
            break
        }
        expect = expect + strings.Trim(s, " \n\r\t")
    }
    // fmt.Printf("%s\n", expect)

    lexer  := new(compiler.Lexer).Init(src)
    parser := new(compiler.Parser).Init(lexer)
    v := InvokeByName(parser, ruleName)
    node := v[0].Interface().(*ast.Node)
    if node.String() != expect {
        t.Fatalf("found:  " + node.String() + "\nexpect: " + expect)
    }
}

func TestTypeWithReflection(t *testing.T) {
    lexer  := new(compiler.Lexer).Init(
    "\n"                      +
    "class B {\n"             +
    "   static main(args){\n" +
    "   }\n"                  +
    "}\n"                     )
    parser := new(compiler.Parser).Init(lexer)
    InvokeByName(parser, "TypeDecl")

    f,_ := os.Open("./test", os.O_RDONLY, 0666)
    defer f.Close()
    n,_ := f.Readdirnames(-1)
    for i := 0; i < len(n); i++ {
        if strings.LastIndex(n[i],".kt") != -1 {
            // fmt.Printf("%s\n", n[i])
            RunOneFile("./test/" + n[i], t)
        }
    }
}

//func TestTypeWithMainMethodAndMatch(t *testing.T) {
//    lexer  := new(compiler.Lexer).Init(
//    "\n"                        +
//    "class A {\n"               +
//    "   static main(args){\n"   +
//    "       a := Mul(Const())\n"+
//    "       a := Mul(Const())\n"+
//    "   }\n"                    +
//    "}\n"                       )
//
//
//
//}