package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func compress(data []byte) string {
	var result strings.Builder
	for _, b := range data {
		binary := fmt.Sprintf("%08b", b)
		count := 1
		current := binary[0]

		for i := 1; i < len(binary); i++ {
			if binary[i] == current {
				count++
			} else {
				result.WriteString(strconv.Itoa(count))
				result.WriteByte(current)
				count = 1
				current = binary[i]
			}
		}
		result.WriteString(strconv.Itoa(count))
		result.WriteByte(current)
		result.WriteByte(' ') // разделитель между байтами
	}
	return result.String()
}

func decompress(compressed string) []byte {
	var result []byte
	bytes := strings.Split(compressed, " ")

	for _, byteStr := range bytes {
		if byteStr == "" {
			continue
		}

		var binary strings.Builder
		i := 0
		for i < len(byteStr) {
			count := 0
			// Читаем число
			for i < len(byteStr) && byteStr[i] >= '0' && byteStr[i] <= '9' {
				digit, _ := strconv.Atoi(string(byteStr[i]))
				count = count*10 + digit
				i++
			}
			// Добавляем соответствующее количество бит
			bit := byteStr[i]
			for j := 0; j < count; j++ {
				binary.WriteByte(bit)
			}
			i++
		}

		// Преобразуем бинарную строку в байт
		if binary.Len() == 8 {
			val, _ := strconv.ParseUint(binary.String(), 2, 8)
			result = append(result, byte(val))
		}
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: program <filename>")
		return
	}

	filename := os.Args[1]
	ext := filepath.Ext(filename)

	if ext == ".z" {
		// Распаковка
		compressed, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Ошибка при чтении файла: %v\n", err)
			return
		}

		decompressed := decompress(string(compressed))
		outFile := strings.TrimSuffix(filename, ".z")
		err = ioutil.WriteFile(outFile, decompressed, 0644)
		if err != nil {
			fmt.Printf("Ошибка при записи файла: %v\n", err)
			return
		}
		fmt.Printf("Файл успешно распакован: %s\n", outFile)
	} else {
		// Сжатие
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Ошибка при чтении файла: %v\n", err)
			return
		}

		compressed := compress(data)
		outFile := filename + ".z"
		err = ioutil.WriteFile(outFile, []byte(compressed), 0644)
		if err != nil {
			fmt.Printf("Ошибка при записи файла: %v\n", err)
			return
		}
		fmt.Printf("Файл успешно сжат: %s\n", outFile)
	}
}
