package parser

import (
	"fmt"
	"nutshell/frontend/lexer"
	"nutshell/runtime"
	"strconv"
)

type NutParser struct {
	Tokens       *[]*lexer.Token
	Position     int
	CurrentToken *lexer.Token
}

func (n *NutParser) advance() {
	n.Position++
	n.CurrentToken = (*n.Tokens)[n.Position]
}

func (n *NutParser) get_token_type(token *lexer.Token) int {
	return (*token).TokenType
}

func (n *NutParser) get_current_token_type() int {
	return n.get_token_type(n.CurrentToken)
}

func (n *NutParser) get_token_value(token *lexer.Token) string {
	return (*token).Value
}

func (n *NutParser) get_current_token_value() string {
	return n.get_token_value(n.CurrentToken)
}

func (n *NutParser) expect(token_type int) runtime.RuntimeResult[*lexer.Token] {
	if n.get_current_token_type() != token_type {
		var err runtime.Error = runtime.SyntaxError(n.CurrentToken.StartPosition, n.CurrentToken.EndPosition, fmt.Sprintf("Expected %d, got %s", token_type, n.get_current_token_value()))
		return runtime.RuntimeResult[*lexer.Token]{
			Result: nil,
			Error:  &err,
		}
	}

	return runtime.RuntimeResult[*lexer.Token]{
		Result: n.CurrentToken,
		Error:  nil,
	}
}

func (n *NutParser) ParseBlock() runtime.RuntimeResult[*Block] {
	var block Block = *InitBlock()
	for n.Position < len(*n.Tokens)-1 {
		for n.get_current_token_type() == lexer.Semicolon {
			n.advance()
		}

		if n.Position == len(*n.Tokens)-1 {
			break
		}

		var rt runtime.RuntimeResult[*Statement] = n.parse_statement()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Block]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		block.Body = append(block.Body, rt.Result)
	}

	return runtime.RuntimeResult[*Block]{
		Result: &block,
		Error:  nil,
	}
}

func (n *NutParser) parse_statement() runtime.RuntimeResult[*Statement] {
	var returned runtime.RuntimeResult[*Statement]
	switch n.get_current_token_type() {
	case lexer.Let:
		var rt runtime.RuntimeResult[*Statement] = n.parse_let_variable_declaration()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Statement]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		returned = runtime.RuntimeResult[*Statement]{
			Result: rt.Result,
			Error:  nil,
		}
	case lexer.Const:
		var rt runtime.RuntimeResult[*Statement] = n.parse_const_variable_declarartion()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Statement]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		returned = runtime.RuntimeResult[*Statement]{
			Result: rt.Result,
			Error:  nil,
		}
	case lexer.Auto:
		var rt runtime.RuntimeResult[*Statement] = n.parse_type_inference_variable_declaration(false)
		if rt.Error != nil {
			return runtime.RuntimeResult[*Statement]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		returned = runtime.RuntimeResult[*Statement]{
			Result: rt.Result,
			Error:  nil,
		}
	default:
		var rt runtime.RuntimeResult[*Expression] = n.parse_expression()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Statement]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		if n.get_current_token_type() == lexer.Identifier {
			var rt2 runtime.RuntimeResult[*Statement] = n.parse_typed_variable_declaration(rt.Result, false)
			if rt2.Error != nil {
				return runtime.RuntimeResult[*Statement]{
					Result: nil,
					Error:  rt2.Error,
				}
			}

			returned = runtime.RuntimeResult[*Statement]{
				Result: rt2.Result,
				Error:  nil,
			}
		} else {
			var statement Statement = (*rt.Result).(Statement)
			returned = runtime.RuntimeResult[*Statement]{
				Result: &statement,
				Error:  nil,
			}
		}
	}

	if n.get_current_token_type() != lexer.EOF && n.get_current_token_type() != lexer.Semicolon {
		var err runtime.Error = runtime.SyntaxError(n.CurrentToken.StartPosition, n.CurrentToken.EndPosition, fmt.Sprintf("Expected ';' or newline, got %s", n.get_current_token_value()))
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  &err,
		}
	}

	return returned
}

func (n *NutParser) parse_let_variable_declaration() runtime.RuntimeResult[*Statement] {
	var let_token *lexer.Token = n.CurrentToken
	n.advance()
	var rt runtime.RuntimeResult[*lexer.Token] = n.expect(lexer.Identifier)
	if rt.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var variable_name string = rt.Result.Value
	n.advance()
	if n.get_current_token_type() == lexer.Semicolon || n.get_current_token_type() == lexer.EOF {
		var any_identifier Expression = interface{}(Identifier{
			VariableName: "any",
			IdentifierToken: &lexer.Token{
				TokenType:     lexer.Identifier,
				Value:         "any",
				StartPosition: let_token.StartPosition,
				EndPosition:   let_token.StartPosition,
			},
		}).(Expression)

		var null_identifier Expression = interface{}(Identifier{
			VariableName: "null",
			IdentifierToken: &lexer.Token{
				TokenType:     lexer.Identifier,
				Value:         "null",
				StartPosition: let_token.StartPosition,
				EndPosition:   let_token.EndPosition,
			},
		}).(Expression)

		var variable_declaration_statement Statement = interface{}(VariableDeclaration{
			PositionStart: let_token.StartPosition,
			VariableName:  variable_name,
			DataType:      &any_identifier,
			Value:         &null_identifier,
			IsConstant:    false,
		}).(Statement)

		return runtime.RuntimeResult[*Statement]{
			Result: &variable_declaration_statement,
			Error:  nil,
		}
	}

	rt = n.expect(lexer.Equals)
	if rt.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	n.advance()
	var rt2 runtime.RuntimeResult[*Expression] = n.parse_expression()
	if rt2.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var value *Expression = rt2.Result
	var any_identifier Expression = interface{}(Identifier{
		VariableName: "any",
		IdentifierToken: &lexer.Token{
			TokenType:     lexer.Identifier,
			Value:         "any",
			StartPosition: let_token.StartPosition,
			EndPosition:   let_token.StartPosition,
		},
	}).(Expression)

	var variable_declaration_statement Statement = interface{}(VariableDeclaration{
		PositionStart: let_token.StartPosition,
		VariableName:  variable_name,
		DataType:      &any_identifier,
		Value:         value,
		IsConstant:    false,
	}).(Statement)

	return runtime.RuntimeResult[*Statement]{
		Result: &variable_declaration_statement,
		Error:  nil,
	}
}

func (n *NutParser) parse_const_variable_declarartion() runtime.RuntimeResult[*Statement] {
	var const_token *lexer.Token = n.CurrentToken
	n.advance()

	switch n.get_current_token_type() {
	case lexer.Auto:
		var rt runtime.RuntimeResult[*Statement] = n.parse_type_inference_variable_declaration(true)
		if rt.Error != nil {
			return runtime.RuntimeResult[*Statement]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		variable_declaration_statement, _ := (*rt.Result).(VariableDeclaration)
		variable_declaration_statement.PositionStart = const_token.StartPosition

		var returned Statement = interface{}(variable_declaration_statement).(Statement)
		return runtime.RuntimeResult[*Statement]{
			Result: &returned,
			Error:  nil,
		}
	default:
		var rt runtime.RuntimeResult[*Expression] = n.parse_additive_expression()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Statement]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		if (*rt.Result).Kind() != IdentifierExpr {
			var rt2 runtime.RuntimeResult[*Statement] = n.parse_typed_variable_declaration(rt.Result, true)
			if rt2.Error != nil {
				return runtime.RuntimeResult[*Statement]{
					Result: nil,
					Error:  rt2.Error,
				}
			}

			variable_declaration_statement, _ := (*rt2.Result).(VariableDeclaration)
			variable_declaration_statement.PositionStart = const_token.StartPosition
			var returned Statement = interface{}(variable_declaration_statement).(Statement)
			return runtime.RuntimeResult[*Statement]{
				Result: &returned,
				Error:  nil,
			}
		}

		if n.get_current_token_type() == lexer.Identifier {
			var rt2 runtime.RuntimeResult[*Statement] = n.parse_typed_variable_declaration(rt.Result, true)
			if rt2.Error != nil {
				return runtime.RuntimeResult[*Statement]{
					Result: nil,
					Error:  rt2.Error,
				}
			}

			variable_declaration_statement, _ := (*rt2.Result).(VariableDeclaration)
			variable_declaration_statement.PositionStart = const_token.StartPosition
			var returned Statement = interface{}(variable_declaration_statement).(Statement)
			return runtime.RuntimeResult[*Statement]{
				Result: &returned,
				Error:  nil,
			}
		}

		variable_node, _ := (*rt.Result).(Identifier)
		var variable_name string = variable_node.VariableName

		var rt3 runtime.RuntimeResult[*lexer.Token] = n.expect(lexer.Equals)
		if rt3.Error != nil {
			return runtime.RuntimeResult[*Statement]{
				Result: nil,
				Error:  rt3.Error,
			}
		}

		n.advance()
		var rt2 runtime.RuntimeResult[*Expression] = n.parse_expression()
		if rt2.Error != nil {
			return runtime.RuntimeResult[*Statement]{
				Result: nil,
				Error:  rt2.Error,
			}
		}

		var value *Expression = rt2.Result
		var any_identifier Expression = interface{}(Identifier{
			VariableName: "any",
			IdentifierToken: &lexer.Token{
				TokenType:     lexer.Identifier,
				Value:         "any",
				StartPosition: const_token.StartPosition,
				EndPosition:   const_token.StartPosition,
			},
		}).(Expression)

		var variable_declaration_statement Statement = interface{}(VariableDeclaration{
			PositionStart: const_token.StartPosition,
			VariableName:  variable_name,
			DataType:      &any_identifier,
			Value:         value,
			IsConstant:    true,
		}).(Statement)

		return runtime.RuntimeResult[*Statement]{
			Result: &variable_declaration_statement,
			Error:  nil,
		}
	}
}

func (n *NutParser) parse_typed_variable_declaration(data_type *Expression, is_constant bool) runtime.RuntimeResult[*Statement] {
	var rt runtime.RuntimeResult[*lexer.Token] = n.expect(lexer.Identifier)
	if rt.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var variable_name string = n.get_current_token_value()
	n.advance()

	if n.get_current_token_type() == lexer.Semicolon || n.get_current_token_type() == lexer.EOF {
		var variable_declaration_statement Statement = interface{}(VariableDeclaration{
			PositionStart: (*data_type).StartPosition(),
			VariableName:  variable_name,
			DataType:      data_type,
			Value:         nil,
			IsConstant:    false,
		}).(Statement)

		return runtime.RuntimeResult[*Statement]{
			Result: &variable_declaration_statement,
			Error:  nil,
		}
	}

	rt = n.expect(lexer.Equals)
	if rt.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	n.advance()
	var rt2 runtime.RuntimeResult[*Expression] = n.parse_expression()
	if rt2.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt2.Error,
		}
	}

	var variable_declaration_statement Statement = interface{}(VariableDeclaration{
		PositionStart: (*data_type).StartPosition(),
		VariableName:  variable_name,
		DataType:      data_type,
		Value:         rt2.Result,
		IsConstant:    is_constant,
	}).(Statement)

	return runtime.RuntimeResult[*Statement]{
		Result: &variable_declaration_statement,
		Error:  nil,
	}
}

func (n *NutParser) parse_type_inference_variable_declaration(is_constant bool) runtime.RuntimeResult[*Statement] {
	var auto_token *lexer.Token = n.CurrentToken
	n.advance()

	var rt runtime.RuntimeResult[*lexer.Token] = n.expect(lexer.Identifier)
	if rt.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var variable_name string = rt.Result.Value
	n.advance()

	rt = n.expect(lexer.Equals)
	if rt.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	n.advance()

	var rt2 runtime.RuntimeResult[*Expression] = n.parse_expression()
	if rt2.Error != nil {
		return runtime.RuntimeResult[*Statement]{
			Result: nil,
			Error:  rt2.Error,
		}
	}

	var variable_declaration_statement Statement = interface{}(VariableDeclaration{
		PositionStart: auto_token.StartPosition,
		VariableName:  variable_name,
		DataType:      nil,
		Value:         rt2.Result,
		IsConstant:    is_constant,
	}).(Statement)

	return runtime.RuntimeResult[*Statement]{
		Result: &variable_declaration_statement,
		Error:  nil,
	}
}

func (n *NutParser) parse_expression() runtime.RuntimeResult[*Expression] {
	var rt runtime.RuntimeResult[*Expression] = n.parse_assignment_expression()
	if rt.Error != nil {
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	return runtime.RuntimeResult[*Expression]{
		Result: rt.Result,
		Error:  nil,
	}
}

func (n *NutParser) parse_assignment_expression() runtime.RuntimeResult[*Expression] {
	var rt runtime.RuntimeResult[*Expression] = n.parse_additive_expression()
	if rt.Error != nil {
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	if n.get_current_token_type() != lexer.Equals {
		return runtime.RuntimeResult[*Expression]{
			Result: rt.Result,
			Error:  nil,
		}
	}

	n.advance()
	var left *Expression = rt.Result

	rt = n.parse_expression()
	if rt.Error != nil {
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var assignment_expression Expression = interface{}(AssignmentExpression{
		Callee:        left,
		Value:         rt.Result,
		PositionStart: (*left).StartPosition(),
		PositionEnd:   (*rt.Result).EndPosition(),
	}).(Expression)

	return runtime.RuntimeResult[*Expression]{
		Result: &assignment_expression,
		Error:  nil,
	}
}

func (n *NutParser) parse_additive_expression() runtime.RuntimeResult[*Expression] {
	var rt runtime.RuntimeResult[*Expression] = n.parse_multiplicative_expression()
	if rt.Error != nil {
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var left *Expression = rt.Result
	for n.get_current_token_type() == lexer.Plus || n.get_current_token_type() == lexer.Minus {
		var operator int = n.get_current_token_type()
		n.advance()

		rt = n.parse_multiplicative_expression()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Expression]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		var binary_expression Expression = interface{}(BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    rt.Result,
		}).(Expression)

		left = &binary_expression
	}

	return runtime.RuntimeResult[*Expression]{
		Result: left,
		Error:  nil,
	}
}

func (n *NutParser) parse_multiplicative_expression() runtime.RuntimeResult[*Expression] {
	var rt runtime.RuntimeResult[*Expression] = n.parse_unary_expression()
	if rt.Error != nil {
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var left *Expression = rt.Result
	for n.get_current_token_type() == lexer.Multiply || n.get_current_token_type() == lexer.Divide || n.get_current_token_type() == lexer.Modulo {
		var operator int = n.get_current_token_type()
		n.advance()

		rt = n.parse_unary_expression()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Expression]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		var binary_expression Expression = interface{}(BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    rt.Result,
		}).(Expression)

		left = &binary_expression
	}

	return runtime.RuntimeResult[*Expression]{
		Result: left,
		Error:  nil,
	}
}

func (n *NutParser) parse_unary_expression() runtime.RuntimeResult[*Expression] {
	if n.get_current_token_type() != lexer.Plus && n.get_current_token_type() != lexer.Minus {
		var rt runtime.RuntimeResult[*Expression] = n.parse_call_expression()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Expression]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		return runtime.RuntimeResult[*Expression]{
			Result: rt.Result,
			Error:  nil,
		}
	}

	var sign int = lexer.Plus
	const sum int = lexer.Plus + lexer.Minus
	var start_token *lexer.Token = n.CurrentToken
	for n.get_current_token_type() == lexer.Plus || n.get_current_token_type() == lexer.Minus {
		if n.get_current_token_type() == lexer.Plus {
			n.advance()
			continue
		}

		sign = sum - sign
		n.advance()
	}

	var rt runtime.RuntimeResult[*Expression] = n.parse_call_expression()
	if rt.Error != nil {
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var unary_expression Expression = interface{}(UnaryExpression{
		Sign:           sign,
		StartSignToken: start_token,
		Value:          rt.Result,
	}).(Expression)

	return runtime.RuntimeResult[*Expression]{
		Result: &unary_expression,
		Error:  nil,
	}
}

func (n *NutParser) parse_call_expression() runtime.RuntimeResult[*Expression] {
	var rt runtime.RuntimeResult[*Expression] = n.parse_member_expression()
	if rt.Error != nil {
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	if n.get_current_token_type() != lexer.LeftParenthese {
		return runtime.RuntimeResult[*Expression]{
			Result: rt.Result,
			Error:  nil,
		}
	}

	var callee *Expression = rt.Result

	var left_parenthese_token *lexer.Token = n.CurrentToken
	n.advance()

	var arguments []*Expression = []*Expression{}
	for n.get_current_token_type() != lexer.RightParenthese && n.get_current_token_type() != lexer.EOF {
		rt = n.parse_expression()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Expression]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		arguments = append(arguments, rt.Result)
		if n.get_current_token_type() == lexer.RightParenthese || n.get_current_token_type() == lexer.EOF {
			break
		}

		var rt2 runtime.RuntimeResult[*lexer.Token] = n.expect(lexer.Comma)
		if rt2.Error != nil {
			return runtime.RuntimeResult[*Expression]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		n.advance()
	}

	if n.get_current_token_type() == lexer.EOF {
		var err runtime.Error = runtime.SyntaxError(left_parenthese_token.StartPosition, left_parenthese_token.EndPosition, "Expected ')', got EOF")
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  &err,
		}
	}

	var right_parenthese_token *lexer.Token = n.CurrentToken
	n.advance()

	var call_expression Expression = interface{}(CallExpression{
		Callee:               callee,
		Arguments:            arguments,
		LeftParentheseToken:  left_parenthese_token,
		RightParentheseToken: right_parenthese_token,
	}).(Expression)

	return runtime.RuntimeResult[*Expression]{
		Result: &call_expression,
		Error:  nil,
	}
}

func (n *NutParser) parse_member_expression() runtime.RuntimeResult[*Expression] {
	var rt runtime.RuntimeResult[*Expression] = n.parse_primary_expression()
	if rt.Error != nil {
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var left *Expression = rt.Result

	for n.get_current_token_type() == lexer.Dot {
		n.advance()
		var rt2 runtime.RuntimeResult[*lexer.Token] = n.expect(lexer.Identifier)
		if rt2.Error != nil {
			return runtime.RuntimeResult[*Expression]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		n.advance()
		var member_expression Expression = interface{}(MemberExpression{
			Parent: left,
			Member: rt2.Result,
		}).(Expression)

		left = &member_expression
	}

	return runtime.RuntimeResult[*Expression]{
		Result: left,
		Error:  nil,
	}
}

func (n *NutParser) parse_primary_expression() runtime.RuntimeResult[*Expression] {
	var returned runtime.RuntimeResult[*Expression]
	switch n.get_current_token_type() {
	case lexer.Int:
		value, _ := strconv.ParseInt(n.get_current_token_value(), 10, 64)
		var int_expression Expression = interface{}(Int{
			Value:    value,
			IntToken: n.CurrentToken,
		}).(Expression)

		returned = runtime.RuntimeResult[*Expression]{
			Result: &int_expression,
			Error:  nil,
		}
	case lexer.Double:
		value, _ := strconv.ParseFloat(n.get_current_token_value(), 64)
		var double_expression Expression = interface{}(Double{
			Value:       value,
			DoubleToken: n.CurrentToken,
		}).(Expression)

		returned = runtime.RuntimeResult[*Expression]{
			Result: &double_expression,
			Error:  nil,
		}
	case lexer.String:
		//var quote rune = rune(n.get_current_token_value()[0])
		value, _ := strconv.Unquote(fmt.Sprintf("\"%s\"", n.get_current_token_value()[1:]))
		var string_expression Expression = interface{}(String{
			Value:       value,
			StringToken: n.CurrentToken,
		}).(Expression)

		returned = runtime.RuntimeResult[*Expression]{
			Result: &string_expression,
			Error:  nil,
		}
	case lexer.Identifier:
		var identifier_expression Expression = interface{}(Identifier{
			VariableName:    n.CurrentToken.Value,
			IdentifierToken: n.CurrentToken,
		}).(Expression)

		returned = runtime.RuntimeResult[*Expression]{
			Result: &identifier_expression,
			Error:  nil,
		}
	case lexer.LeftParenthese:
		var left_parenthese *lexer.Token = n.CurrentToken
		n.advance()

		var rt runtime.RuntimeResult[*Expression] = n.parse_expression()
		if rt.Error != nil {
			return runtime.RuntimeResult[*Expression]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		var expression *Expression = rt.Result

		var rt2 runtime.RuntimeResult[*lexer.Token] = n.expect(lexer.RightParenthese)
		if rt2.Error != nil {
			return runtime.RuntimeResult[*Expression]{
				Result: nil,
				Error:  rt2.Error,
			}
		}

		var right_parenthese *lexer.Token = rt2.Result

		var bracket_expression Expression = interface{}(BracketExpression{
			Value:                expression,
			LeftParentheseToken:  left_parenthese,
			RightParentheseToken: right_parenthese,
		}).(Expression)

		returned = runtime.RuntimeResult[*Expression]{
			Result: &bracket_expression,
		}
	default:
		var err runtime.Error = runtime.SyntaxError(n.CurrentToken.StartPosition, n.CurrentToken.EndPosition, "Invalid syntax!")
		return runtime.RuntimeResult[*Expression]{
			Result: nil,
			Error:  &err,
		}
	}

	n.advance()
	return returned
}

func InitNutParser(tokens *[]*lexer.Token) *NutParser {
	var parser *NutParser = &NutParser{
		Tokens:       tokens,
		Position:     -1,
		CurrentToken: nil,
	}

	parser.advance()
	return parser
}
