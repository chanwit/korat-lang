package korat_test

import "testing"
import "korat"

func TestStringBuffer(t *testing.T) {
    b := korat.NewStringBuffer()
    s := "abcกขค" 
    for _,ch := range s {
        b.Append(ch)
    }
    if b.String() != s {
        t.Fatalf("string buffer error")
    }
}