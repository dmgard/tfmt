# Typed String Formatter

`tfmt` is like `strings.Builder` with typed methods for formatting builtin datatypes, making it much faster than `fmt.Printf` (benchmark results pending).

Usage:
```go
import "fmt"
import "github.com/dmgard/tfmt"

func main() {
	f := tfmt.Fmt()
	
	f.Str("Let's print some integers: ")
	
	for i := range 10 {
		f.Int(i).Str(", ")
	}
	
	f.TrimRight(", ")
	
	fmt.Println(f.String())
}
```

Also has some helpful features for nested indentation:

```go

f.Str("type someType struct {").Indent().Ln().
	Str("field1 type1").Ln().
	Str("field2 type2").Ln().
	Str("field3 type3")
f.OutdentLn().Str("}")

fmt.Println(f.String())
```

prints

```go
type someType struct {
	field1 type1
	field2 type2
	field3 type3
}
```