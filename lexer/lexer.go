package lexer

func TokenizeExpression(expr string) string {
    s = strings.Split(expr, " ")
    switch s[0] {
        case "help":
            return "help"
    }
}