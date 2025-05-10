package utils

import "fmt"

func PrettyPrint(msg map[string][]string) {
	for uid, content := range msg {
		fmt.Println(uid, content)
	}
}
