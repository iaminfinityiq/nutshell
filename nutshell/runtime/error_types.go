package runtime

func SyntaxError(position_start *Position, position_end *Position, reason string) Error {
	return Error{
		StartPosition: position_start,
		EndPosition:   position_end,
		ErrorType:     "SyntaxError",
		Reason:        reason,
	}
}
