package dj

// Valid returns true if the input is valid json.
// The input can be a string or []byte.
func Valid[DATA string | []byte](data DATA) error {
	// Use validAny directly to avoid allocating the Result.Raw field.
	_, err := validPayload(data, 0)
	if err != nil {
		return err
	}
	return nil
}

// valid* functions are slower than parse* equivalents because they check for invalid JSON.
// Once the JSON is validated, the parse* functions can be used to unmarshal the JSON faster.

func validPayload[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			i, err = validAny(data, i)
			if err != nil {
				return 0, err
			}
			for ; i < len(data); i++ {
				switch data[i] {
				default:
					return 0, NewSyntaxError(i, "invalid character after JSON")
				case ' ', '\t', '\n', '\r':
					continue
				}
			}
			return i, nil
		case ' ', '\t', '\n', '\r':
			continue
		}
	}
	return 0, NewSyntaxError(i, "empty content")
}

func validAny[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return 0, NewSyntaxError(i, "invalid character beginning JSON")
		case ' ', '\t', '\n', '\r':
			continue
		case '{':
			return validObject(data, i+1)
		case '[':
			return validArray(data, i+1)
		case '"':
			return validString(data, i+1)
		case 't':
			return validTrue(data, i+1)
		case 'f':
			return validFalse(data, i+1)
		case 'n':
			return validNull(data, i+1)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return validNumber(data, i+1)
		}
	}
	return 0, NewSyntaxError(i, "empty content")
}

func validObject[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, NewSyntaxError(i, "expected object key or closing brace")
		case ' ', '\t', '\n', '\r':
			continue
		case '}':
			return i + 1, nil
		case '"':
		key:
			if i, err = validString(data, i+1); err != nil {
				return 0, err
			}
			if i, err = validColon(data, i); err != nil {
				return 0, err
			}
			if i, err = validAny(data, i); err != nil {
				return 0, err
			}
			if i, err = validComma(data, i, '}'); err != nil {
				return 0, err
			}
			if data[i] == '}' {
				return i + 1, nil
			}
			i++
			for ; i < len(data); i++ {
				switch data[i] {
				default:
					return i, NewSyntaxError(i, "invalid character between object entries")
				case ' ', '\t', '\n', '\r':
					continue
				case '"':
					goto key
				}
			}
			return i, NewSyntaxError(i, "object not closed after entry")
		}
	}
	return i, NewSyntaxError(i, "object not closed")
}

func validColon[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, NewSyntaxError(i, "invalid character for colon")
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return i + 1, nil
		}
	}
	return i, NewSyntaxError(i, "expected colon")
}

func validComma[DATA string | []byte](data DATA, i int, end byte) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, NewSyntaxError(i, "invalid character for comma")
		case ' ', '\t', '\n', '\r':
			continue
		case ',', end:
			return i, nil
		}
	}
	return i, NewSyntaxError(i, "expected comma")
}

func validArray[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			for ; i < len(data); i++ {
				if i, err = validAny(data, i); err != nil {
					return 0, err
				}
				if i, err = validComma(data, i, ']'); err != nil {
					return 0, err
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
	return i, NewSyntaxError(i, "array not closed")
}

func validString[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		if data[i] < ' ' {
			return i, NewSyntaxError(i, "invalid character for string")
		} else if data[i] == '\\' {
			i++
			if i == len(data) {
				return i, NewSyntaxError(i, "escape character at end of data")
			}
			switch data[i] {
			default:
				return i, NewSyntaxError(i, "invalid escape character "+string(data[i:i+1]))
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
			case 'u':
				for j := 0; j < 4; j++ {
					i++
					if i >= len(data) {
						return i, NewSyntaxError(i, "too short unicode character")
					}
					if !((data[i] >= '0' && data[i] <= '9') ||
						(data[i] >= 'a' && data[i] <= 'f') ||
						(data[i] >= 'A' && data[i] <= 'F')) {
						return i, NewSyntaxError(i, "invalid unicode character")
					}
				}
			}
		} else if data[i] == '"' {
			return i + 1, nil
		}
	}
	return i, NewSyntaxError(i, "string not closed")
}

func validNumber[DATA string | []byte](data DATA, i int) (outi int, err error) {
	i--
	// sign
	if data[i] == '-' {
		i++
		if i == len(data) {
			return i, NewSyntaxError(i, "sign character at end of data")
		}
		if data[i] < '0' || data[i] > '9' {
			return i, NewSyntaxError(i, "expected digit after sign")
		}
	}
	// int
	if i == len(data) {
		return i, NewSyntaxError(i, "short data for number")
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
			return i, NewSyntaxError(i, "expected digit following dot")
		}
		if data[i] < '0' || data[i] > '9' {
			return i, NewSyntaxError(i, "expected digit following dot")
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
			return i, NewSyntaxError(i, "expected digit following exponent in exp number")
		}
		if data[i] == '+' || data[i] == '-' {
			i++
		}
		if i == len(data) {
			return i, NewSyntaxError(i, "expected digit following sign in exp number")
		}
		if data[i] < '0' || data[i] > '9' {
			return i, NewSyntaxError(i, "expected valid digit in exp number")
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
	return 0, NewSyntaxError(i, "expected 'true'")
}

func validFalse[DATA string | []byte](data DATA, i int) (outi int, err error) {
	if i+4 <= len(data) && data[i] == 'a' && data[i+1] == 'l' &&
		data[i+2] == 's' && data[i+3] == 'e' {
		return i + 4, nil
	}
	return 0, NewSyntaxError(i, "expected 'false'")
}

func validNull[DATA string | []byte](data DATA, i int) (outi int, err error) {
	if i+3 <= len(data) && data[i] == 'u' && data[i+1] == 'l' && data[i+2] == 'l' {
		return i + 3, nil
	}
	return 0, NewSyntaxError(i, "expected 'null'")
}
