package parser

import (
	"errors"
	"strings"
)

type Parser interface {
	Parse(line string) (string, map[string]string, error)
}

type ArgsParser struct{}

// GetArgs принимает список аргументов и возвращает словарь, где ключ - это название флага, а значение - это значение данного флага
func (ArgsParser) GetArgs(args []string) (map[string]string, error) {
	result := make(map[string]string)
	ok := true
	i := 0
	for i < len(args) {
		arg := args[i]
		var key, value string

		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) == 2 {
				key = parts[0]
				value = parts[1]
			} else if i+1 < len(args) {
				key = parts[0]
				value = args[i+1]
				i++
			} else {
				ok = false
			}
		} else if strings.HasPrefix(arg, "-") {
			parts := strings.SplitN(arg[1:], "=", 2)
			if len(parts) == 2 {
				key = parts[0]
				value = parts[1]
			} else if i+1 < len(args) {
				key = parts[0]
				value = args[i+1]
				i++
			}
		} else {
			ok = false
		}

		if key != "" && value != "" {
			result[key] = value
		} else {
			ok = false
		}

		i++
	}
	if ok {
		return result, nil
	} else {
		return nil, errors.New("invalid flag")
	}
}

func (p ArgsParser) Parse(input string) (string, map[string]string, error) {
	parts := strings.Fields(input)
	var cmd string
	var argList []string
	if len(parts) == 0 {
		return "", nil, errors.New("empty line")
	} else if len(parts) == 1 {
		cmd, argList = parts[0], make([]string, 0)
	} else {
		cmd, argList = parts[0], parts[1:]
	}
	args, err := p.GetArgs(argList)
	if err != nil {
		return "", nil, err
	}

	return cmd, args, nil
}
