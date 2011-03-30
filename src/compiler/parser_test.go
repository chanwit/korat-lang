package compiler_test

import "testing"
import "compiler"

import "os"
import "fmt"
import "reflect"
import "strings"
import "sut/govy"
import "ast"

func invokeByName(obj interface{}, name string, in...reflect.Value) []reflect.Value {
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

func runOneFile(filename string, t *testing.T) {
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
        s,e := f.NextLine()
        if s == "expect:" {
            break
        }
        if e == os.EOF {
            panic("EOF found before 'expect:'")
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

    lexer  := new(compiler.Lexer).Init(src)
    parser := new(compiler.Parser).Init(lexer)
    v := invokeByName(parser, ruleName)
    node := v[0].Interface().(*ast.Node)
    if node.String() != expect {
        t.Fatalf("found:  " + node.String() + "\nexpect: " + expect)
    }
    fmt.Printf("PASSED   : " + filename + "\n")
}

func TestTypeWithReflection(t *testing.T) {
    f,_ := os.Open("./test", os.O_RDONLY, 0666)
    defer f.Close()
    n,_ := f.Readdirnames(-1)
    fmt.Printf("\n")
    for i := 0; i < len(n); i++ {
        if strings.LastIndex(n[i],".kt") != -1 {
            // fmt.Printf("%s\n", n[i])
            runOneFile("./test/" + n[i], t)
        }
    }
}
