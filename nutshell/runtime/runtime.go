package runtime

import (
	"fmt"
	"strconv"

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

	c = color.New(color.FgHiGreen).Add(color.Bold)
	var offset int = len(strconv.FormatInt(int64(e.EndPosition.Line+1), 10)) + 4
	if e.StartPosition.Line == e.EndPosition.Line {
		c.Printf("%d || ", e.StartPosition.Line+1)
		fmt.Println(e.StartPosition.GetLine(e.StartPosition.Line))
		for j := 0; j < e.StartPosition.Column+offset; j++ {
			fmt.Print(" ")
		}

		for j := e.StartPosition.Column + offset; j < e.EndPosition.Column+offset; j++ {
			fmt.Print("^")
		}

		fmt.Print("\n")
		return
	}

	for i := e.StartPosition.Line; i <= e.EndPosition.Line; i++ {
		var line_repr string = strconv.FormatInt(int64(i+1), 10)
		for j := 0; j < offset-4-len(line_repr); j++ {
			c.Printf("0")
		}

		c.Printf("%d || ", i+1)

		var line string = e.StartPosition.GetLine(i)
		fmt.Print(line)
		if line != "" {
			fmt.Print("\n")
			if i == e.StartPosition.Line {
				for j := 0; j < e.StartPosition.Column+offset; j++ {
					fmt.Print(" ")
				}

				for j := e.StartPosition.Column + offset; j < len(line)+offset; j++ {
					fmt.Print("^")
				}
			} else if i == e.EndPosition.Line {
				for j := 0; j < offset; j++ {
					fmt.Print(" ")
				}

				for j := offset; j < e.EndPosition.Column+offset; j++ {
					fmt.Print("^")
				}
			} else {
				for j := 0; j < offset; j++ {
					fmt.Print(" ")
				}

				for j := offset; j < len(line)+offset; j++ {
					fmt.Print("^")
				}
			}
		}
		
		fmt.Print("\n")
	}
}
