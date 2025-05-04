package runtime

func SyntaxError(position_start *Position, position_end *Position, reason string) Error {
	return Error{
		StartPosition: position_start,
		EndPosition:   position_end,
		ErrorType:     "SyntaxError",
		Reason:        reason,
	}
}

func MathError(position_start *Position, position_end *Position, reason string) Error {
	return Error{
		StartPosition: position_start,
		EndPosition:   position_end,
		ErrorType:     "MathError",
		Reason:        reason,
	}
}

func ArgumentError(position_start *Position, position_end *Position, reason string) Error {
	return Error{
		StartPosition: position_start,
		EndPosition:   position_end,
		ErrorType:     "ArgumentError",
		Reason:        reason,
	}
}

func TypeError(position_start *Position, position_end *Position, reason string) Error {
	return Error{
		StartPosition: position_start,
		EndPosition:   position_end,
		ErrorType:     "TypeError",
		Reason:        reason,
	}
}

func VariableError(position_start *Position, position_end *Position, reason string) Error {
	return Error{
		StartPosition: position_start,
		EndPosition: position_end,
		ErrorType: "VariableError",
		Reason: reason,
	}
}
