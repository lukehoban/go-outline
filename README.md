# Go Outline

Simple utility for extracting a JSON representation of the declarations in a 
Go source file.

## Installing

```bash
go get -u github.com/lukehoban/go-outline
```

## Using
```bash
> go-outline -f file.go
[{"label":"proc","type":"package",<...>}]
```

### Schema
```go
type Declaration struct {
	Label       string          `json:"label"`
	Type        string          `json:"type"`
	Start       token.Pos       `json:"start"`
	End         token.Pos       `json:"end"`
	Children    []Declaration   `json:"children,omitempty"`
}
```