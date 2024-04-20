package main

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func hashPassword(password, randString string) string {
	arr1 := strings.Split(password, "")
	arr2 := strings.Split(randString, "")
	var result strings.Builder
	index1 := 0
	index2 := 0
	for (index1 < len(arr1)) || (index2 < len(arr2)) {
		if index1 < len(arr1) {
			result.WriteString(arr1[index1])
			index1++
		}
		if index2 < len(arr2) {
			result.WriteString(arr2[index2])
			index2++
		}
	}
	hash := md5.Sum([]byte(result.String()))
	return fmt.Sprintf("%x", hash)
}
