package parser

import (
	"errors"
	"strings"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
)

type Parser interface {
	Parse(line string) (commands.Command, error)
}

type MyParser struct {
}

// getArgs принимает список аргументов и возвращает словарь, где ключ - это название флага, а значение - это значение данного флага
func getArgs(args []string) (map[string]string, error) {
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
		return nil, errors.New("Invalid flag")
	}
}

func (MyParser) Parse(input string) (commands.Command, error) {
	parts := strings.Fields(input)
	var cmd string
	var argList []string
	if len(parts) == 0 {
		return commands.AcceptOrd{}, errors.New("Empty line")
	} else if len(parts) == 1 {
		cmd, argList = parts[0], make([]string, 0)
	} else {
		cmd, argList = parts[0], parts[1:]
	}
	args, err := getArgs(argList)
	if err != nil {
		return commands.NewAcceptOrd(), err
	}
	var ans commands.Command
	switch cmd {
	case "acceptOrd":
		ans = commands.NewAcceptOrd()
	case "acceptReturn":
		ans = commands.NewAcceptReturn()
	case "deliverOrd":
		ans = commands.NewDeliverOrd()
	case "getOrds":
		ans = commands.NewGetOrds()
	case "getReturns":
		ans = commands.NewGetReturns()
	case "returnOrd":
		ans = commands.NewReturnOrd()
	case "help":
		ans = commands.NewHelp()
	default:
		return commands.NewAcceptOrd(), errors.New("unknown error")
	}
	return ans.Validate(args)
}
