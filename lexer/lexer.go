package lexer

import "fmt"

func Test() {
    fmt.Println("testing, testing... 1, 2, 3.")
}

func TokenizeExpression(expr string) string {
    s = strings.Split(expr, " ")
    switch s[0] {
        case "help":
            return "help"
    }
}