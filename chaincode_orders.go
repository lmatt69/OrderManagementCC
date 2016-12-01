/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Order struct {
	OrderNumber			int 				`json:"orderNumber"`
	OrderDate				string 			`json:"orderDate"`
	OrderSubmitter  string 			`json:"orderSubmitter"`
	OrderValue 			float64 		`json:"orderValue"`
	OrderLines      []OrderLine `json:"orderLine"`
	OrderTotal			float64			`json:"orderTotal"`
	OrderStatus			string			`json:"orderStatus"`
}

type OrderLine struct {
	PartNumber			string 				`json:"partNumber"`
	PartDescription	string				`json:"partDescription"`
	Quantity				int 					`json:"quantity"`
	UnitPrice				float64				`json:"unitPrice"`
	LineTotal				float64				`json:"lineTotal"`
	LineStatus			string				`json:"lineStatus"`
}
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	/*from cp_web: adding in an array of Order Numbers for easier searching*/

	var order []string
	orderBytes, _ := json.Marshal(&order)

	err := stub.PutState("OrderKeys", orderBytes)
	if err != nil {
		return nil, err
	}

	fmt.Println("finished initalizing chaincode state")

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}else if function == "createOrder" {
		return t.createOrder(stub, args)
	}else if function == "addOrderLines" {
		return t.addOrderLines(stub,args)
	}else if function == "deleteOrderLines" {
		return t.deleteOrderLines(stub, args)
	}else if function == "deleteOrder" {
		return t.deleteOrder(stub, args)
	}else if function == "updateLineStatus" {
		return t.updateLineStatus( stub, args)
	}else if function == "updateOrderStatus"{
		return t.updateOrderStatus(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

func(t *SimpleChaincode) createOrder( stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. Order to be created")
	}
	var order Order
	var err error

	fmt.Println("retrieving Order")
  /*Unmarshal the Order object from the request*/
	err = json.Unmarshal([]byte(args[0]), &order)
	/*throw error if marshaling failed */
	if err != nil {
		fmt.Println("error invalid order received")
		return nil, errors.New("Invalid order received")
	}

	/*We have an order, check it has its component parts*/
	fmt.Println("running createOrder()")
	if(order.OrderNumber  == nil){
		fmt.Println("Order received with no order number")
		return nil, errors.News("Order received with no order number")
	}

	/*retrieve order numbers*/
	orderBytes, err := stub.GetState("OrderKeys")
	if err != nil {
		return nil, err
	}
	var orders []Order
	err = json.Unmarshal(orderBytes, &orders)
	if err != nil {
		fmt.Println("Unable to retrieve order number list")
		return nil, errors.new("Unable to retrieve orders")
	}

/*find order and add if missing (leave alone if not) */
 alreadyOrder := false
 for _, orderKey := range orders {
	 if orderKey == orderNumber {
		 alreadyOrder = truew
	 }
}
if alreadyOrder == false {
	orders append(orders, orderNumber)
	orderBytesWrite, err = json.Marshal(&orders)
	if err != nil {
		fmt.Println("Cannot create order numbers")
		return nil, errors.News("Error creating order numbers")
	}
	fmt.Println("Storing order numbers")
	err = stub.PutState("OrderKeys", orderBytesWrite)
	if err != nil {
		fmt.Println("Cannot write Order keys")
		return nil, errors.New("Error writing keys back")
	}
}
	fmt.Println("Finished writing ")
	return nil, nil }

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {											//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) read (stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of ket to query")
	}

	key = args[0]
valAsBytes, err := stub.GetState(key)
if err != nil {
	jsonResp = "{\"Error\":\"Failed to get state for "+ key +"\"}"
	return nil, errors.New(jsonResp)
}

return valAsBytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s ", err)
	}
}
