package lexer

import (
	"fmt"
	"nutshell/runtime"
)

func is_digit(r rune) bool {
	return '0' <= r && r <= '9'
}

func is_alpha(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') || r == '_'
}

func is_legal(r rune) bool {
	return is_digit(r) || is_alpha(r)
}

func is_whitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

type Lexer struct {
	CurrentChar *rune
	Position    runtime.Position
}

func (l *Lexer) Advance() {
	l.Position.Advance(l.CurrentChar)
	if l.Position.Index < len(l.Position.FileText) {
		var r rune = rune(l.Position.FileText[l.Position.Index])
		l.CurrentChar = &r
	} else {
		l.CurrentChar = nil
	}
}

func (l *Lexer) Tokenize() runtime.RuntimeResult[*[]*Token] {
	var tokens []*Token = []*Token{}

	var keywords map[string]int = make(map[string]int)
	keywords["let"] = Let
	keywords["var"] = Var
	keywords["const"] = Const
	keywords["auto"] = Auto

	for l.CurrentChar != nil {
		switch *l.CurrentChar {
		case '\n':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Semicolon, "\\n"))
		case ';':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Semicolon, ";"))
		case '+':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Plus, "+"))
		case '-':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Minus, "-"))
		case '*':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Multiply, "*"))
		case '/':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Divide, "/"))
		case '^':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Power, "^"))
		case '%':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Modulo, "%"))
		case '(':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, LeftParenthese, "("))
		case ')':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, RightParenthese, ")"))
		case '=':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Equals, "="))
		case ',':
			var position_start runtime.Position = l.Position.Copy()
			var position_end runtime.Position = position_start.Copy()
			position_end.Advance(nil)
			tokens = append(tokens, CreateToken(&position_start, &position_end, Comma, ","))
		default:
			if is_whitespace(*l.CurrentChar) {
				l.Advance()
				continue
			}

			if *l.CurrentChar == '"' {
				var position_start runtime.Position = l.Position.Copy()
				l.Advance()
				var value string = ""

				for l.CurrentChar != nil && *l.CurrentChar != '"' {
					value += string(*l.CurrentChar)
					l.Advance()

					if l.CurrentChar != nil {
						if *l.CurrentChar == '\\' {
							value += string(*l.CurrentChar)
							l.Advance()
						}
					}
				}

				if l.CurrentChar == nil {
					var position_end runtime.Position = position_start.Copy()
					position_end.Advance(nil)

					var err runtime.Error = runtime.SyntaxError(&position_start, &position_end, "Expected '\"' to be closed, got EOF")
					return runtime.RuntimeResult[*[]*Token]{
						Result: nil,
						Error:  &err,
					}
				}

				var position_end = l.Position.Copy()
				tokens = append(tokens, CreateToken(&position_start, &position_end, String, fmt.Sprintf("\"%s", value)))

				l.Advance()
				continue
			}

			if *l.CurrentChar == '\'' {
				var position_start runtime.Position = l.Position.Copy()
				l.Advance()
				var value string = ""

				for l.CurrentChar != nil && *l.CurrentChar != '\'' {
					value += string(*l.CurrentChar)
					l.Advance()

					if l.CurrentChar != nil {
						if *l.CurrentChar == '\\' {
							value += string(*l.CurrentChar)
							l.Advance()
						}
					}
				}

				if l.CurrentChar == nil {
					var position_end runtime.Position = position_start.Copy()
					position_end.Advance(nil)

					var err runtime.Error = runtime.SyntaxError(&position_start, &position_end, "Expected ''' to be closed, got EOF")
					return runtime.RuntimeResult[*[]*Token]{
						Result: nil,
						Error:  &err,
					}
				}

				var position_end = l.Position.Copy()
				tokens = append(tokens, CreateToken(&position_start, &position_end, String, fmt.Sprintf("'%s", value)))

				l.Advance()
				continue
			}

			if *l.CurrentChar == '.' {
				l.Advance()
				if is_digit(*l.CurrentChar) {
					var number string = ""
					var dot_count int = 1

					var position_start runtime.Position = l.Position.Copy()
					for is_digit(*l.CurrentChar) || *l.CurrentChar == '.' {
						if *l.CurrentChar == '.' {
							dot_count++
						}

						number += string(*l.CurrentChar)
						l.Advance()

						if l.CurrentChar == nil {
							break
						}
					}

					var position_end runtime.Position = l.Position.Copy()
					switch dot_count {
					case 1:
						tokens = append(tokens, CreateToken(&position_start, &position_end, Double, number))
					default:
						l.Advance()

						var err runtime.Error = runtime.SyntaxError(&position_start, &position_end, fmt.Sprintf("Expected 0 or 1 '.' in a number, got %d", dot_count))
						return runtime.RuntimeResult[*[]*Token]{
							Result: nil,
							Error:  &err,
						}
					}

					continue
				}

				var position_start runtime.Position = l.Position.Copy()
				var position_end runtime.Position = position_start.Copy()
				position_end.Advance(nil)
				tokens = append(tokens, CreateToken(&position_start, &position_end, Dot, "."))
				continue
			}

			if is_digit(*l.CurrentChar) {
				var number string = ""
				var dot_count int = 0

				var position_start runtime.Position = l.Position.Copy()
				for is_digit(*l.CurrentChar) || *l.CurrentChar == '.' {
					if *l.CurrentChar == '.' {
						dot_count++
					}

					number += string(*l.CurrentChar)
					l.Advance()

					if l.CurrentChar == nil {
						break
					}
				}

				var position_end runtime.Position = l.Position.Copy()
				switch dot_count {
				case 0:
					tokens = append(tokens, CreateToken(&position_start, &position_end, Int, number))
				case 1:
					tokens = append(tokens, CreateToken(&position_start, &position_end, Double, number))
				default:
					l.Advance()
					var err runtime.Error = runtime.SyntaxError(&position_start, &position_end, fmt.Sprintf("Expected 0 or 1 '.' in a number, got %d", dot_count))
					return runtime.RuntimeResult[*[]*Token]{
						Result: nil,
						Error:  &err,
					}
				}

				continue
			}

			if is_alpha(*l.CurrentChar) {
				var identifier string = ""

				var position_start runtime.Position = l.Position.Copy()
				for is_legal(*l.CurrentChar) {
					identifier += string(*l.CurrentChar)
					l.Advance()

					if l.CurrentChar == nil {
						break
					}
				}

				var position_end runtime.Position = l.Position.Copy()
				tt, ok := keywords[identifier]
				if ok {
					tokens = append(tokens, CreateToken(&position_start, &position_end, tt, identifier))
				} else {
					tokens = append(tokens, CreateToken(&position_start, &position_end, Identifier, identifier))
				}

				continue
			}

			var position_start runtime.Position = l.Position.Copy()
			var character rune = *l.CurrentChar

			l.Advance()
			var position_end runtime.Position = l.Position.Copy()

			var err runtime.Error = runtime.SyntaxError(&position_start, &position_end, fmt.Sprintf("Invalid character: '%s'", string(character)))
			return runtime.RuntimeResult[*[]*Token]{
				Result: nil,
				Error:  &err,
			}
		}

		l.Advance()
	}

	var position_start runtime.Position = l.Position.Copy()

	l.Advance()
	var position_end runtime.Position = l.Position.Copy()

	tokens = append(tokens, &Token{
		StartPosition: &position_start,
		EndPosition:   &position_end,
		TokenType:     EOF,
		Value:         "EOF",
	})

	return runtime.RuntimeResult[*[]*Token]{
		Result: &tokens,
		Error:  nil,
	}
}

func InitLexer(file_name string, file_extension string, code string) *Lexer {
	var lexer *Lexer = &Lexer{
		CurrentChar: nil,
		Position: runtime.Position{
			FileName:      file_name,
			FileExtension: file_extension,
			FileText:      code,
			Index:         -1,
			Line:          0,
			Column:        -1,
		},
	}

	lexer.Advance()
	return lexer
}
