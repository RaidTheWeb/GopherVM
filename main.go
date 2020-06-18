package main

import(
    "fmt"
    "flag"
    "os"
    "io/ioutil"
    "strings"
    "regexp"
)

func TokenizeExpression(expr string) string {
    _, REG_TOKEN := regexp.Compile("#[0-9]")
    
    s := strings.Split(expr, " ")
    Token := ""
    switch s[0] {
        case "store":
            if len(s) >= 3 {
                if REG_TOKEN.Match(s[1]) {

                }
            }

    }
    return Token
}

func usage() {
    fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
    flag.PrintDefaults()
    os.Exit(2)
}

func main() {
    flag.Usage = usage
    flag.Parse()

    args := flag.Args()
    if len(args) < 1 {
        fmt.Fprint(os.Stderr, "gophervm: \033[31mpanic:\033[0m input file missing, vm terminated.\n")
        os.Exit(1);
    }
    data, err := ioutil.ReadFile(args[0])
    if err != nil {
        fmt.Fprintf(os.Stderr, "gophervm: \033[31mpanic:\033[0m %s, vm terminated.\n", err)
        os.Exit(2)
    }
    fmt.Println(string(data))
    fmt.Println(TokenizeExpression("store #1 \"hi\""))
}