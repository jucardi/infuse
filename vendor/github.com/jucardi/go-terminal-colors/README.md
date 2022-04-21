# go-terminal-colors

Helpful utility to print colored output in the terminal.

#### Quick Start

To keep up to date with the most recent version:

```bash
go get github.com/jucardi/go-terminal-colors
```

To get a specific version:

```bash
go get gopkg.in/jucardi/go-terminal-colors.v1
```

#### Using the library

Supports two modalities, using the `fmt` function signatures with `fmtc` or creating a new `IColorFlow` which allows printing
messages with chaining functions.

#### The `fmt` way

To print a message with colors using the `fmt` way, start by using `fmtc.WithColors(colors...)` which allows to pass multiple predefined terminal colors
in the `fmtc` package, followed by the `fmt` function would use.

**Example**

```go
fmtc.WithColors(fmtc.White, fmtc.Bold, fmtc.BgRed).Println("An error has ocurred!")
```

#### Using the color flow

The color flow allows using a custom `io.Writer`, or if not provided, it will generate the resulting string in memory which can then be retrieve by
invoking the `String()` function.

**Example using Stdout as the `io.Writer`**

```go
fmtc.New(os.Stdout).
	Print(" ERROR ", fmtc.White, fmtc.Bold, fmtc.BgRed).
	Print(" ").
	Print(timestamp, fmtc.Cyan).
	Print(" ").
	Println(err.Error)
```

**Same example but using an in-memory writer to obtain the resulting string**

```go
result := fmtc.New().
	Print(" ERROR ", fmtc.White, fmtc.Bold, fmtc.BgRed).
	Print(" ").
	Print(timestamp, fmtc.Cyan).
	Print(" ").
	Println(err.Error).
	String()
```