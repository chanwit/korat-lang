package util

import "fmt"

const (
    STATIC = "STATIC"
    NIL    = "<nil>"
)

func format(tag string, s...string) string {
    if len(s) == 0 {return tag}

    r := s[0]
    for i := 1; i < len(s); i++ {
        r = r + "," + s[i]
    }
    return fmt.Sprintf("%s(%s)", tag, r)    
}

func CLASS(s...string) string {
    return format("CLASS", s...)
}

func IDENT(s string) string {
    return fmt.Sprintf("IDENT('%s')", s)
}

func METHOD(s...string) string {
    return format("METHOD", s...)
}

func MEMBERS(s...string) string {
    return format("MEMBERS", s...)
}

func MODIFIERS(s...string) string {
    return format("MODIFIERS", s...)
}

func ARGS(s...string) string {
    return format("ARGS", s...)
}

func ARG(s...string) string {
    return format("ARG", s...)
}

func TYPE(s string) string {
    return fmt.Sprintf("TYPE('%s')", s)
}
