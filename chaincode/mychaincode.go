package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// MyChaincode example
type MyChaincode struct {
}

// Init - Function called on chaincode instantiation
func (t *MyChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("mychaincode Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	total, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for total")
	}

	err = stub.PutState("TOTAL", []byte(strconv.Itoa(total)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Init complete. TOTAL=%s\n", total)
	return shim.Success(nil)
}

// Invoke - Function called on chaincode invoke
func (t *MyChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("mychaincode Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "reg" {
		return t.register(stub, args)
	} else if function == "del" {
		return t.delete(stub, args)
	} else if function == "move" {
		return t.move(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"reg\" \"del\" \"move\" \"query\"")
}

// register - register a new user with 10 tokens
func (t *MyChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	tBytes, err := stub.GetState("TOTAL")
	if err != nil {
		return shim.Error("Failed to get state for TOTAL")
	}
	total, err := strconv.Atoi(string(tBytes))
	if total < 10 {
		return shim.Error("Cannot register any new users")
	}

	name := args[0]
	err = stub.PutState(name, []byte(strconv.Itoa(10)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("TOTAL", []byte(strconv.Itoa(total-10)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// delete - delete a user and return tokens to pool
func (t *MyChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	name := args[0]
	aBytes, err := stub.GetState(name)
	if err != nil {
		resp := "Failed to get the state of " + name
		return shim.Error(resp)
	}
	amount, err := strconv.Atoi(string(aBytes))

	err = stub.DelState(name)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	tBytes, err := stub.GetState("TOTAL")
	if err != nil {
		return shim.Error("Failed to get state for TOTAL")
	}
	total, err := strconv.Atoi(string(tBytes))

	err = stub.PutState("TOTAL", []byte(strconv.Itoa(total+amount)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// move - move amount between users
func (t *MyChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// query - check value of user
func (t *MyChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(MyChaincode))
	if err != nil {
		fmt.Printf("Error starting mychaincode: %s", err)
	}
}
