// Пакет для инициализации .env файла
// и получения слайса из последовательности значений переменных окружения

package env

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadEnv(filePath string) error {
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("no .env file found, using system environment variables: %s", err)
	}
	defer func(file *os.File) {
		errClose := file.Close()
		if errClose != nil {
			fmt.Printf("failed to close config file %s: %v\n", filePath, errClose)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}

		line = strings.ReplaceAll(line, "\"", "")

		keyValue := strings.SplitN(line, "=", 2)
		if len(keyValue) != 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		setEnvErr := os.Setenv(key, value)
		if setEnvErr != nil {
			return setEnvErr
		}
	}

	if scannerErr := scanner.Err(); scannerErr != nil {
		return fmt.Errorf("error reading .env file: %s", scannerErr)
	}

	return nil
}

func GetIntFromEnv(key string) (int, error) {
	str := os.Getenv(key)

	str = strings.TrimSpace(str)

	if str == "" {
		return 0, fmt.Errorf("env %s is empty", key)
	}

	return strconv.Atoi(str)
}

func GetSliceIntFromEnv(key string) ([]int, error) {
	str := os.Getenv(key)

	str = strings.TrimSpace(str)

	if str == "" {
		return nil, fmt.Errorf("env %s is empty", key)
	}

	strSlice := strings.Split(str, ",")

	intSlice := make([]int, len(strSlice))

	for i, s := range strSlice {
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed convert element %s to int in env %s", s, key)
		}

		intSlice[i] = num
	}

	return intSlice, nil
}
