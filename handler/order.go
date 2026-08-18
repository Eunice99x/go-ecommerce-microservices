package handler

import (
	"fmt"
	"net/http"
)

type Order struct {}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create and order")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("list all orders")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get one order")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update one order")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete one order")
}