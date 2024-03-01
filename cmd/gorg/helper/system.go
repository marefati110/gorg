package helper

import "os"

func GetWD() string {
	_, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return "/home/ali/gorg/test"
}
