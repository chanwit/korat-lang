package compiler_test

import "testing"
import "compiler"

// import "fmt"
import . "util"
import "reflect"

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

func TestTypeWithReflection(t *testing.T) {
    lexer  := new(compiler.Lexer).Init(
    "\n"                      +
    "class B {\n"             +
    "   static main(args){\n" +
    "   }\n"                  +
    "}\n"                     )
    parser := new(compiler.Parser).Init(lexer)
    InvokeByName(parser, "TypeDecl")
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