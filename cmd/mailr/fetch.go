package main

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
)

func main() {

	fmt.Println("Fetching", os.Args[1])

	res, _ := http.Get(os.Args[1])
	defer res.Body.Close()

	xb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("RESPONSE BODY -----------------------------------------")
	fmt.Println(string(xb))
	fmt.Println("END RESPONSE BODY -----------------------------------------")
}

