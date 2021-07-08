package client

import (
	"fmt"
	"net/rpc/jsonrpc"
	"strconv"

	"example.com/common"
)

func Start() {

	// get JSON-RPC client by dialing TCP connection
	client, _ := jsonrpc.Dial("tcp", "127.0.0.1:9002")

	fmt.Print("Enter A to get all vegetables\n")
	fmt.Print("Enter B to search a vegetable\n")
	fmt.Print("Enter C to add a vegetable\n")
	fmt.Print("Enter D to update a vegetable\n")
	fmt.Print("=============================\n")

	for i := 0; ; i++ {
		var command string
		fmt.Print("Enter your command: ")
		fmt.Scanln(&command)
		switch command {

		case "A":
			fmt.Printf("All vegetables")
			fmt.Printf("--------------\n")
			var allVeges []*common.Vegetable
			if err := client.Call("Store.GetAll", "1", &allVeges); err != nil {
				fmt.Println("Error:1 Store.GetAll()", err)
			} else {
				for _, veg := range allVeges {
					fmt.Println(veg.Details())
				}
			}
			break

		case "B":
			var carrot common.Vegetable
			fmt.Printf("Get vegetable details")
			fmt.Printf("---------------------\n")
			var name string
			fmt.Print("Enter vegetable name: ")
			fmt.Scanln(&name)
			if err := client.Call("Store.Get", name, &carrot); err != nil {
				fmt.Println("Error while Store.Get()", err)
			} else {
				fmt.Println(name+" details: ", carrot.Details())
			}
			break

		case "C":
			var allVeges []*common.Vegetable
			fmt.Printf("Add vegetable Leeks")
			fmt.Printf("-------------------\n")
			var addVeg common.Vegetable
			addVeg.Name = "Leeks"
			addVeg.Price = 100.5
			addVeg.QTY = 10

			var input string
			fmt.Print("Enter vegetable name: ")
			fmt.Scanln(&input)
			addVeg.Name = input

			fmt.Print("Enter vegetable price: ")
			fmt.Scanln(&input)
			f, err := strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Print("Bad value, bye")
			}
			addVeg.Price = f

			fmt.Print("Enter vegetable QTY: ")
			fmt.Scanln(&input)
			f, err = strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Print("Bad value, bye")
			}
			addVeg.QTY = f

			fmt.Println("beFore adding vegetable details: ", addVeg.Details())

			if err := client.Call("Store.AddVeg", &addVeg, &allVeges); err != nil {
				fmt.Println("Error while Store.AddVeg()", err)
			} else {
				for _, veg := range allVeges {
					if veg.Name == addVeg.Name {
						fmt.Println(veg.Details())
					}
				}
			}
			break

		case "D":
			var allVeges []*common.Vegetable
			fmt.Printf("Update vegetable")
			fmt.Printf("----------------\n")
			var addVeg common.Vegetable
			addVeg.Name = "Leeks"
			addVeg.Price = 100.5
			addVeg.QTY = 10

			var input string
			fmt.Print("Enter vegetable name: ")
			fmt.Scanln(&input)
			addVeg.Name = input

			fmt.Print("Enter vegetable price: ")
			fmt.Scanln(&input)
			f, err := strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Print("Bad value, bye")
			}
			addVeg.Price = f

			fmt.Print("Enter vegetable QTY: ")
			fmt.Scanln(&input)
			f, err = strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Print("Bad value, bye")
			}
			addVeg.QTY = f

			fmt.Println("beFore adding vegetable details: ", addVeg.Details())

			if err := client.Call("Store.UpdateVeg", &addVeg, &allVeges); err != nil {
				fmt.Println("Error while Store.AddVeg()", err)
			} else {
				for _, veg := range allVeges {
					if veg.Name == addVeg.Name {
						fmt.Println(veg.Details())
					}
				}
			}
			break

		default:
			fmt.Print("K. Bye\n")
			fmt.Print("======\n")
			return
		}
	}
}
