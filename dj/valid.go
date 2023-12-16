package dj

// Valid returns true if the input is valid json.
// The input can be a string or []byte.
func Valid[DATA string | []byte](json DATA) error {
	_, err := validPayload(json, 0)
	return err
}

func validPayload[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			i, err = validAny(data, i)
			if err != nil {
				return i, err
			}
			for ; i < len(data); i++ {
				switch data[i] {
				default:
					return i, SyntaxError{Index: i, Msg: "invalid character after JSON"}
				case ' ', '\t', '\n', '\r':
					continue
				}
			}
			return i, nil
		case ' ', '\t', '\n', '\r':
			continue
		}
	}
	return i, SyntaxError{Index: i, Msg: "invalid character before JSON"}
}

func validAny[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, SyntaxError{Index: i, Msg: "invalid character beginning JSON"}
		case ' ', '\t', '\n', '\r':
			continue
		case '{':
			return validObject(data, i+1)
		case '[':
			return validArray(data, i+1)
		case '"':
			return validString(data, i+1)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return validNumber(data, i+1)
		case 't':
			return validTrue(data, i+1)
		case 'f':
			return validFalse(data, i+1)
		case 'n':
			return validNull(data, i+1)
		}
	}
	return i, SyntaxError{Index: i, Msg: "no content found"}
}

func validObject[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, SyntaxError{Index: i, Msg: "expected object key or closing brace"}
		case ' ', '\t', '\n', '\r':
			continue
		case '}':
			return i + 1, nil
		case '"':
		key:
			if i, err = validString(data, i+1); err != nil {
				return i, err
			}
			if i, err = validColon(data, i); err != nil {
				return i, err
			}
			if i, err = validAny(data, i); err != nil {
				return i, err
			}
			if i, err = validComma(data, i, '}'); err != nil {
				return i, err
			}
			if data[i] == '}' {
				return i + 1, nil
			}
			i++
			for ; i < len(data); i++ {
				switch data[i] {
				default:
					return i, SyntaxError{Index: i, Msg: "invalid character between object entries"}
				case ' ', '\t', '\n', '\r':
					continue
				case '"':
					goto key
				}
			}
			return i, SyntaxError{Index: i, Msg: "object not closed after entry"}
		}
	}
	return i, SyntaxError{Index: i, Msg: "object not closed"}
}

func validColon[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, SyntaxError{Index: i, Msg: "invalid character for colon"}
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return i + 1, nil
		}
	}
	return i, SyntaxError{Index: i, Msg: "expected colon"}
}

func validComma[DATA string | []byte](data DATA, i int, end byte) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, SyntaxError{Index: i, Msg: "invalid character for comma"}
		case ' ', '\t', '\n', '\r':
			continue
		case ',', end:
			return i, nil
		}
	}
	return i, SyntaxError{Index: i, Msg: "expected comma"}
}

func validArray[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			for ; i < len(data); i++ {
				if i, err = validAny(data, i); err != nil {
					return i, err
				}
				if i, err = validComma(data, i, ']'); err != nil {
					return i, err
				}
				if data[i] == ']' {
					return i + 1, nil
				}
			}
		case ' ', '\t', '\n', '\r':
			continue
		case ']':
			return i + 1, nil
		}
	}
	return i, SyntaxError{Index: i, Msg: "array not closed"}
}

func validString[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		if data[i] < ' ' {
			return i, SyntaxError{Index: i, Msg: "invalid character for string"}
		} else if data[i] == '\\' {
			i++
			if i == len(data) {
				return i, SyntaxError{Index: i, Msg: "escape character at end of data"}
			}
			switch data[i] {
			default:
				return i, SyntaxError{Index: i, Msg: "invalid escape character " + string(data[i:i+1])}
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
			case 'u':
				for j := 0; j < 4; j++ {
					i++
					if i >= len(data) {
						return i, SyntaxError{Index: i, Msg: "too short unicode character"}
					}
					if !((data[i] >= '0' && data[i] <= '9') ||
						(data[i] >= 'a' && data[i] <= 'f') ||
						(data[i] >= 'A' && data[i] <= 'F')) {
						return i, SyntaxError{Index: i, Msg: "invalid unicode character"}
					}
				}
			}
		} else if data[i] == '"' {
			return i + 1, nil
		}
	}
	return i, SyntaxError{Index: i, Msg: "string not closed"}
}

func validNumber[DATA string | []byte](data DATA, i int) (outi int, err error) {
	i--
	// sign
	if data[i] == '-' {
		i++
		if i == len(data) {
			return i, SyntaxError{Index: i, Msg: "sign character at end of data"}
		}
		if data[i] < '0' || data[i] > '9' {
			return i, SyntaxError{Index: i, Msg: "expected digit after sign"}
		}
	}
	// int
	if i == len(data) {
		return i, SyntaxError{Index: i, Msg: "short data for number"}
	}
	if data[i] == '0' {
		i++
	} else {
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	// frac
	if i == len(data) {
		return i, nil
	}
	if data[i] == '.' {
		i++
		if i == len(data) {
			return i, SyntaxError{Index: i, Msg: "expected digit following dot"}
		}
		if data[i] < '0' || data[i] > '9' {
			return i, SyntaxError{Index: i, Msg: "expected digit following dot"}
		}
		i++
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	// exp
	if i == len(data) {
		return i, nil
	}
	if data[i] == 'e' || data[i] == 'E' {
		i++
		if i == len(data) {
			return i, SyntaxError{Index: i, Msg: "expected digit following exponent in exp number"}
		}
		if data[i] == '+' || data[i] == '-' {
			i++
		}
		if i == len(data) {
			return i, SyntaxError{Index: i, Msg: "expected digit following sign in exp number"}
		}
		if data[i] < '0' || data[i] > '9' {
			return i, SyntaxError{Index: i, Msg: "expected valid digit in exp number"}
		}
		i++
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	return i, nil
}

func validTrue[DATA string | []byte](data DATA, i int) (outi int, err error) {
	if i+3 <= len(data) && data[i] == 'r' && data[i+1] == 'u' &&
		data[i+2] == 'e' {
		return i + 3, nil
	}
	return i, SyntaxError{Index: i, Msg: "expected 'true'"}
}

func validFalse[DATA string | []byte](data DATA, i int) (outi int, err error) {
	if i+4 <= len(data) && data[i] == 'a' && data[i+1] == 'l' &&
		data[i+2] == 's' && data[i+3] == 'e' {
		return i + 4, nil
	}
	return i, SyntaxError{Index: i, Msg: "expected 'false'"}
}

func validNull[DATA string | []byte](data DATA, i int) (outi int, err error) {
	if i+3 <= len(data) && data[i] == 'u' && data[i+1] == 'l' &&
		data[i+2] == 'l' {
		return i + 3, nil
	}
	return i, SyntaxError{Index: i, Msg: "expected 'null'"}
}
