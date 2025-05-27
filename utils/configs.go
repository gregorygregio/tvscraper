package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var arguments = make(map[string]string)
var isLoaded = false

func LoadArguments() {
	if len(os.Args) <= 1 {
		fmt.Println("Informar usuário e senha do qbittorrent")
	}

	for i := 0; i < len(os.Args); i++ {
		arg := os.Args[i]
		fmt.Println(arg)
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		arg = strings.Replace(arg, "--", "", 1)
		occuranceOfEqualSigns := strings.Count(arg, "=")
		if occuranceOfEqualSigns != 1 {
			fmt.Printf("Invalid format on argument '%s' \n", arg)
			panic("Invalid format on argument")
		}
		argSplit := strings.Split(arg, "=")
		if len(argSplit) != 2 {
			fmt.Printf("Invalid format on argument '%s' \n", arg)
			panic("Invalid format on argument")
		}

		arguments[argSplit[0]] = argSplit[1]
	}
	fmt.Println(arguments)
	isLoaded = true
}

func GetConfig(configName string) (string, error) {
	if !isLoaded {
		return "", errors.New("Config not found")
	}

	val, ok := arguments[configName]
	if !ok {
		return "", errors.New("Config not found")
	}

	return val, nil
}
