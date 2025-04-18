package runtime

type Position struct {
	FileName      string
	FileExtension string
	FileText      string
	Index         int
	Line          int
	Column        int
}

func (p *Position) Advance(current_char *rune) {
	p.Index++
	if current_char == nil {
		p.Column++
		return
	}

	if *current_char == '\n' {
		p.Line++
		p.Column = 0
		return
	}

	p.Column++
}

func (p *Position) GetLine(line int) string {
	var i int = 0
	var returned string = ""
	for _, c := range p.FileText {
		if i < line {
			if c == '\n' {
				i++
			}
		} else {
			if c == '\n' {
				return returned
			}

			returned += string(c)
		}
	}

	return returned
}

func (p Position) Copy() Position {
	return Position{
		FileName:      p.FileName,
		FileExtension: p.FileExtension,
		FileText:      p.FileText,
		Index:         p.Index,
		Line:          p.Line,
		Column:        p.Column,
	}
}
