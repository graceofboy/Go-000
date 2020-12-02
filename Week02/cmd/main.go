package main

import "fmt"

func main() {
	service := service.NewService()
	data, err := service.FindPersion()
	if err != nil {
		fmt.Println("find persion error: %+v", err)
	}
	fmt.Println("persion info: %+v", data)
}
