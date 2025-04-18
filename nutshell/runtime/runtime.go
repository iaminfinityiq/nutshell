package runtime

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type RuntimeResult[T any] struct {
	Result T
	Error  *Error
}

type Error struct {
	StartPosition *Position
	EndPosition   *Position
	ErrorType     string
	Reason        string
}

func (e *Error) DisplayError() {
	c := color.New(color.FgHiRed).Add(color.Bold)
	c.Printf("%s:", e.ErrorType)

	fmt.Printf(" %s\nFile %s.%s, line %d\n", e.Reason, e.StartPosition.FileName, e.StartPosition.FileExtension, e.StartPosition.Line+1)

	var offset int = len(strconv.FormatInt(int64(e.EndPosition.Line+1), 10)) + 4
	if e.StartPosition.Line == e.EndPosition.Line {
		c := color.New(color.FgGreen).Add(color.Bold)
		c.Printf("\n%d || ", e.StartPosition.Line+1)
		fmt.Printf("%s\n", e.StartPosition.GetLine(e.StartPosition.Line))

		for i := 0; i < e.StartPosition.Column+offset; i++ {
			fmt.Print(" ")
		}

		for i := e.StartPosition.Column + offset; i < e.EndPosition.Column+offset; i++ {
			fmt.Print("^")
		}

		fmt.Print("\n")
		return
	}

	for i := e.StartPosition.Line; i <= e.EndPosition.Line; i++ {
		var line string = e.StartPosition.GetLine(i + 1)
		var index_offset int = len(strconv.FormatInt(int64(i+1), 10))

		fmt.Print("\n")
		if strings.TrimSpace(line) == "" {
			c := color.New(color.FgGreen).Add(color.Bold)
			for i := index_offset; i <= offset-4; i++ {
				c.Printf("0")
			}

			c.Printf("%d ||\n", i+1)
			continue
		}

		c := color.New(color.FgGreen).Add(color.Bold)
		for i := index_offset; i <= offset-4; i++ {
			c.Printf("0")
		}

		c.Printf("%d ||\n", i+1)
		fmt.Println(line)

		if i == e.StartPosition.Line {
			for j := 0; j < e.StartPosition.Column+offset; j++ {
				fmt.Print(" ")
			}

			for j := e.StartPosition.Column + offset; j < len(line)+offset; j++ {
				fmt.Print("^")
			}
		} else if i == e.EndPosition.Line {
			for j := offset; j < e.EndPosition.Column+offset; j++ {
				fmt.Print("^")
			}
		} else {
			for j := offset; j < len(line)+offset; j++ {
				fmt.Print("^")
			}
		}

		fmt.Print("\n")
	}
}
