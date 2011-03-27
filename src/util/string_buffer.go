package util

import "utf8"
//import "fmt"

type StringBuffer struct {
    bytes []byte
    index int
}

func NewStringBuffer() *StringBuffer {
    r := new(StringBuffer)
    r.bytes = make([]byte, 2048)
    r.index = 0
    return r
}

func (S *StringBuffer) Append(ch int) *StringBuffer {
    // fmt.Printf("append: %c", ch)
    w := utf8.EncodeRune(S.bytes[S.index:], ch)
    S.index +=w
    return S
}

func (S *StringBuffer) AppendStr(s string) *StringBuffer {
    // fmt.Printf("append: %c", ch)
    for _,ch := range s {
        w := utf8.EncodeRune(S.bytes[S.index:], ch)
        S.index +=w
    }
    return S
}

func (S *StringBuffer) String() string {
    S.bytes[S.index] = 0
    return string(S.bytes[0:S.index])
}