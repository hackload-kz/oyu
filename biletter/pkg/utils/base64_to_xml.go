package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
)

func Base64ToXmlString(encodedData string) (string, error) {
	// Декодируем base64-данные
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования base64: %v", err)
	}

	xmlString := string(decodedData)

	// Проверяем, что строка начинается с XML-заголовка
	if len(xmlString) == 0 || xmlString[0] != '<' {
		return "", errors.New("некорректные данные: не является XML")
	}

	// Дополнительно можно распарсить XML, если это нужно
	//var xmlData interface{}
	//err = xml.Unmarshal(decodedData, &xmlData)
	//if err != nil {
	//	return "", fmt.Errorf("ошибка разбора XML: %v", err)
	//}

	return xmlString, nil
}
