package main

import(
    "fmt"
    "flag"
    "os"
    "io/ioutil"
    "strings"
    "regexp"
    "strconv"
    "bufio"
    "io"
)

import (
    "sync"

    "github.com/cheekybits/genny/generic"
)


// Dict/Env Implementation
// Key the key of the dictionary
type Key generic.Type

// Value the content of the dictionary
type Value generic.Type

// ValueDictionary the set of Items
type ValueDictionary struct {
    items map[Key]Value
    lock  sync.RWMutex
}

// Set adds a new item to the dictionary
func (d *ValueDictionary) Set(k Key, v Value) {
    d.lock.Lock()
    defer d.lock.Unlock()
    if d.items == nil {
        d.items = make(map[Key]Value)
    }
    d.items[k] = v
}

// Delete removes a value from the dictionary, given its key
func (d *ValueDictionary) Delete(k Key) bool {
    d.lock.Lock()
    defer d.lock.Unlock()
    _, ok := d.items[k]
    if ok {
        delete(d.items, k)
    }
    return ok
}

// Has returns true if the key exists in the dictionary
func (d *ValueDictionary) Has(k Key) bool {
    d.lock.RLock()
    defer d.lock.RUnlock()
    _, ok := d.items[k]
    return ok
}

// Get returns the value associated with the key
func (d *ValueDictionary) Get(k Key) Value {
    d.lock.RLock()
    defer d.lock.RUnlock()
    return d.items[k]
}

// Clear removes all the items from the dictionary
func (d *ValueDictionary) Clear() {
    d.lock.Lock()
    defer d.lock.Unlock()
    d.items = make(map[Key]Value)
}

// Size returns the amount of elements in the dictionary
func (d *ValueDictionary) Size() int {
    d.lock.RLock()
    defer d.lock.RUnlock()
    return len(d.items)
}

// Keys returns a slice of all the keys present
func (d *ValueDictionary) Keys() []Key {
    d.lock.RLock()
    defer d.lock.RUnlock()
    keys := []Key{}
    for i := range d.items {
        keys = append(keys, i)
    }
    return keys
}

// Values returns a slice of all the values present
func (d *ValueDictionary) Values() []Value {
    d.lock.RLock()
    defer d.lock.RUnlock()
    values := []Value{}
    for i := range d.items {
        values = append(values, d.items[i])
    }
    return values
}








// Errors

// LexerError
func ThrowLexerError(err string, lineno int) {
    fmt.Fprintf(os.Stderr, "gophervm: \033[31mLexerError:\033[0m lineno:%d %s, vm terminated.\n", lineno, err)
    os.Exit(2)
}

// ParserError
func ThrowParserError(err string, lineno int) {
    fmt.Fprintf(os.Stderr, "gophervm: \033[31mParserError:\033[0m lineno:%d %s, vm terminated.\n", lineno, err)
    os.Exit(2)
}



/*
GopherVM Lexer
Regex/Switch
*/
func TokenizeExpression(expr string, lineno int) string {
    REG_TOKEN, _ := regexp.Compile("#[0-9]")
    STRING_TOKEN, _ := regexp.Compile("\".*\"")
    INT_TOKEN, _ := regexp.Compile("[0-9]")
    START_STRING, _ := regexp.Compile("\".*")
    
    s := strings.Split(expr, " ")
    Token := ""
    ValueType := ""
    switch s[0] {
        case "store":
            Type := "STORE"
            if len(s) >= 3 {
                if REG_TOKEN.MatchString(s[1]) {
                    Register := s[1]
                    StoreVal := s[2]
                    if START_STRING.MatchString(s[2]) {
                        StoreVal = strings.Join(s[2:], " ")
                    }
                    if STRING_TOKEN.MatchString(StoreVal) {
                        ValueType = "STRING"
                    } else if INT_TOKEN.MatchString(s[2]) {
                        ValueType = "INT"
                        StoreVal = s[2]
                    } else {
                        ThrowLexerError("Invalid Type For STORE Call", lineno)
                    }
                    Token = fmt.Sprintf("%s,%s,%s,%s", Type, Register, ValueType, StoreVal)
                } else {
                    ThrowLexerError("Invalid Register", lineno)
                }

            } else {
                ThrowLexerError("Invalid Number Of Arguments For STORE Call", lineno)
            }


        case "print_str":
            Type := "PRINT_STR"
            if len(s) == 2 {
                if REG_TOKEN.MatchString(s[1]) {
                    Register := s[1]
                    Token = fmt.Sprintf("%s,%s", Type, Register)
                } else {
                    ThrowLexerError("Invalid Register", lineno)
                }
            } else {
                ThrowLexerError("Invalid Number Of Arguments For PRINT_STR Call", lineno)
            }
        
        case "print_int":
            Type := "PRINT_INT"
            if len(s) == 2 {
                if REG_TOKEN.MatchString(s[1]) {
                    Register := s[1]
                    Token = fmt.Sprintf("%s,%s", Type, Register)
                } else {
                    ThrowLexerError("Invalid Register", lineno)
                }
            } else {
                ThrowLexerError("Invalid Number Of Arguments For PRINT_INT Call", lineno)
            }

        case "exit":
            if len(s) == 1 {
                Type := "EXIT"
                Token = fmt.Sprintf("%s", Type)
            } else {
                ThrowLexerError("Invalid Number Of Arguments For EXIT Call", lineno)
            }

        case "add":
            if len(s) == 4 {
                Type := "ADD"
                if REG_TOKEN.MatchString(s[1]) {
                    StoreReg := s[1]
                    if REG_TOKEN.MatchString(s[2]) {
                        FirstReg := s[2]
                        if REG_TOKEN.MatchString(s[3]) {
                            SecondReg := s[3]
                            Token = fmt.Sprintf("%s,%s,%s,%s", Type, StoreReg, FirstReg, SecondReg)


                        } else {
                            ThrowLexerError("Invalid Register", lineno)
                        }
                    } else {
                        ThrowLexerError("Invalid Register", lineno)
                    }
                } else {
                    ThrowLexerError("Invalid Register", lineno)
                }
            } else {
                ThrowLexerError("Invalid Number Of Arguments For ADD Call", lineno)
            }

        case "sub":
            if len(s) == 4 {
                Type := "SUB"
                if REG_TOKEN.MatchString(s[1]) {
                    StoreReg := s[1]
                    if REG_TOKEN.MatchString(s[2]) {
                        FirstReg := s[2]
                        if REG_TOKEN.MatchString(s[3]) {
                            SecondReg := s[3]
                            Token = fmt.Sprintf("%s,%s,%s,%s", Type, StoreReg, FirstReg, SecondReg)


                        } else {
                            ThrowLexerError("Invalid Register", lineno)
                        }
                    } else {
                        ThrowLexerError("Invalid Register", lineno)
                    }
                } else {
                    ThrowLexerError("Invalid Register", lineno)
                }
            } else {
                ThrowLexerError("Invalid Number Of Arguments For SUB Call", lineno)
            }

        case "mul":
            if len(s) == 4 {
                Type := "MUL"
                if REG_TOKEN.MatchString(s[1]) {
                    StoreReg := s[1]
                    if REG_TOKEN.MatchString(s[2]) {
                        FirstReg := s[2]
                        if REG_TOKEN.MatchString(s[3]) {
                            SecondReg := s[3]
                            Token = fmt.Sprintf("%s,%s,%s,%s", Type, StoreReg, FirstReg, SecondReg)


                        } else {
                            ThrowLexerError("Invalid Register", lineno)
                        }
                    } else {
                        ThrowLexerError("Invalid Register", lineno)
                    }
                } else {
                    ThrowLexerError("Invalid Register", lineno)
                }
            } else {
                ThrowLexerError("Invalid Number Of Arguments For MUL Call", lineno)
            }

        case "div":
            if len(s) == 4 {
                Type := "DIV"
                if REG_TOKEN.MatchString(s[1]) {
                    StoreReg := s[1]
                    if REG_TOKEN.MatchString(s[2]) {
                        FirstReg := s[2]
                        if REG_TOKEN.MatchString(s[3]) {
                            SecondReg := s[3]
                            Token = fmt.Sprintf("%s,%s,%s,%s", Type, StoreReg, FirstReg, SecondReg)


                        } else {
                            ThrowLexerError("Invalid Register", lineno)
                        }
                    } else {
                        ThrowLexerError("Invalid Register", lineno)
                    }
                } else {
                    ThrowLexerError("Invalid Register", lineno)
                }
            } else {
                ThrowLexerError("Invalid Number Of Arguments For DIV Call", lineno)
            }

        
        case "":
            if true == true {
                ;
            }

        case "define":
            Type := "DEFINE"
            Register := s[1]
            StoreVal := strings.Join(s[2:], " ")
            Token = fmt.Sprintf("%s,%s,%s", Type, Register, StoreVal)

        case "call":
            Type := "CALL"
            Register := s[1]
            Token = fmt.Sprintf("%s,%s", Type, Register)
        
        default:
            COMMENT, _ := regexp.Compile("//.*")
            if COMMENT.MatchString(expr) {
                if true == true {
                    ;
                }    
            } else {
                ThrowLexerError("Unknown Call", lineno)
            }

    }
    return Token
}




// GopherVM Parser/Evaluator
func ParseExpression(tokens string, lineno int, env ValueDictionary) ValueDictionary {

    SplitTokens := strings.Split(tokens, ",")
    switch SplitTokens[0] {
        case "STORE":
            RegisterLoc := SplitTokens[1]
            RegisterVal := SplitTokens[3]
            env.Set(RegisterLoc, RegisterVal)
        case "PRINT_STR":
            RegisterLoc := SplitTokens[1]
            RegisterVal := env.Get(RegisterLoc)
            Val := fmt.Sprintf("%s", RegisterVal)
            STRING, _ := regexp.Compile("\".*\"")
            if STRING.MatchString(Val) {
                OutVal := Val[1:len(Val) - 1]
                fmt.Fprint(os.Stdout, OutVal + "\n")
            } else {
                ThrowParserError("Invalid Argument Type To PRINT_STR Call", lineno)
            }
        
        case "DEFINE":
            RegisterLoc := SplitTokens[1]
            RegisterCode := SplitTokens[2]
            x := TokenizeExpression(RegisterCode, lineno)
            env.Set(RegisterLoc, x)

        case "CALL":
            RegisterLoc := SplitTokens[1]
            LexedCode := fmt.Sprintf("%s", env.Get(RegisterLoc))
            fmt.Println(LexedCode)
            env = ParseExpression(LexedCode, lineno, env)

        case "ADD":
            RegisterLoc := SplitTokens[1]
            FirstReg := SplitTokens[2]
            SecondReg := SplitTokens[3]
            FirstNum := fmt.Sprintf("%s", env.Get(FirstReg))
            SecondNum := fmt.Sprintf("%s", env.Get(SecondReg))
            INT, _ := regexp.Compile("[0-9]")
            if INT.MatchString(FirstNum) && INT.MatchString(SecondNum) {
                num1, _ := strconv.Atoi(FirstNum)
                num2, _ := strconv.Atoi(SecondNum)
                OutVal := num1 + num2
                result := strconv.Itoa(OutVal)
                env.Set(RegisterLoc, result)
            }

        case "SUB":
            RegisterLoc := SplitTokens[1]
            FirstReg := SplitTokens[2]
            SecondReg := SplitTokens[3]
            FirstNum := fmt.Sprintf("%s", env.Get(FirstReg))
            SecondNum := fmt.Sprintf("%s", env.Get(SecondReg))
            INT, _ := regexp.Compile("[0-9]")
            if INT.MatchString(FirstNum) && INT.MatchString(SecondNum) {
                num1, _ := strconv.Atoi(FirstNum)
                num2, _ := strconv.Atoi(SecondNum)
                OutVal := num1 - num2
                result := strconv.Itoa(OutVal)
                env.Set(RegisterLoc, result)
            }

        case "MUL":
            RegisterLoc := SplitTokens[1]
            FirstReg := SplitTokens[2]
            SecondReg := SplitTokens[3]
            FirstNum := fmt.Sprintf("%s", env.Get(FirstReg))
            SecondNum := fmt.Sprintf("%s", env.Get(SecondReg))
            INT, _ := regexp.Compile("[0-9]")
            if INT.MatchString(FirstNum) && INT.MatchString(SecondNum) {
                num1, _ := strconv.Atoi(FirstNum)
                num2, _ := strconv.Atoi(SecondNum)
                OutVal := num1 * num2
                result := strconv.Itoa(OutVal)
                env.Set(RegisterLoc, result)
            }

        case "DIV":
            RegisterLoc := SplitTokens[1]
            FirstReg := SplitTokens[2]
            SecondReg := SplitTokens[3]
            FirstNum := fmt.Sprintf("%s", env.Get(FirstReg))
            SecondNum := fmt.Sprintf("%s", env.Get(SecondReg))
            INT, _ := regexp.Compile("[0-9]")
            if INT.MatchString(FirstNum) && INT.MatchString(SecondNum) {
                num1, _ := strconv.Atoi(FirstNum)
                num2, _ := strconv.Atoi(SecondNum)
                OutVal := num1 / num2
                result := strconv.Itoa(OutVal)
                env.Set(RegisterLoc, result)
            }


        case "PRINT_INT":
            RegisterLoc := SplitTokens[1]
            RegisterVal := env.Get(RegisterLoc)
            Val := fmt.Sprintf("%s", RegisterVal)
            INT, _ := regexp.Compile("[0-9]")
            if INT.MatchString(Val) {
                fmt.Fprint(os.Stdout, Val + "\n")
            } else {
                ThrowParserError("Invalid Argument Type To PRINT_STR Call", lineno)
            }
        case "EXIT":
            os.Exit(0)
    }
    return env
} 


// Repl/Shell
func StartRepl(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    env := ValueDictionary{}
    x := ""
    lineno := 1
    for {
		fmt.Printf("gvm >> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
        x = TokenizeExpression(line, lineno)
        env = ParseExpression(x, lineno, env)
        lineno++
	}
}







// usage
func usage() {
    fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
    flag.PrintDefaults()
    os.Exit(2)
}


// main
func main() {
    flag.Usage = usage
    flag.Parse()

    args := flag.Args()
    if len(args) < 1 {
        StartRepl(os.Stdin, os.Stdout)
        os.Exit(1);
    }
    data, err := ioutil.ReadFile(args[0])
    if err != nil {
        fmt.Fprintf(os.Stderr, "gophervm: \033[31mpanic:\033[0m %s, vm terminated.\n", err)
        os.Exit(2)
    }
    lines := strings.Split(string(data), "\n")
    env := ValueDictionary{}
    for i, l := range lines {
        x := TokenizeExpression(l, i + 1)
        env = ParseExpression(x, i + 1, env)
    }
    
}