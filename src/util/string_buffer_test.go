package util_test

import "testing"
import "util"

func TestStringBuffer(t *testing.T) {
    b := util.NewStringBuffer()
    s := "abcกขค" 
    for _,ch := range s {
        b.Append(ch)
    }
    if b.String() != s {
        t.Fatalf("string buffer error")
    }
}