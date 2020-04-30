package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)


type User2 struct {
	Id      int    `json:"id"`
	Name    string `json:"username"`
	Age     int    `json:"age,omitempty"`
	Address string `json:"-"`
}

func main() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	u := User2{
		Id:      12,
		Name:    "wendell",
		Age:     1,
		Address: "成都高新区",
	}

	data, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))

	u2 := &User2{}
	data = []byte(`{"User":{"M":"123","N":"456"},"id":12,"username":"wendell","age":1}`)
	err = json.Unmarshal(data, u2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", u2)
}
