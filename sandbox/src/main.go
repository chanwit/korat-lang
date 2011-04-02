package main

import "os"
import "fmt"


type Error string
type Node  string

func (node Node) String() string {
	return string(node)
}

func (e Error) String() string {
	return string(e)
}

var (
    ListNode = Node("list")
    Success = Error("success")
	ErrRecognition = Error("Recognition Error")
)

func _list() Node {
    panic(ErrInternal)
    return ListNode
}

func list() (n Node) {
    failed := Success
    defer func() { if e := recover(); e != nil {
        failed = ErrInternal
        error = failed
    }}()
    int startTokenIndex = index()
    if isSpeculating() && memN := alreadyParsedRule(listMemo) {
        n = memN
        return
    }
    n = _list()
    if isSpeculating() {
        memoize(listMemo, startTokenIndex, failed)
    }
    return
}

func main() {
    fmt.Printf("%s\n", list())
}
