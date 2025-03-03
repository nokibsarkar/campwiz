package routes

import "fmt"

func HandleError(endpoint string) {
	// Handle error here
	if r := recover(); r != nil {
		// log the error
		fmt.Println("Error at endpoint", endpoint, ":", r)
	}
}
