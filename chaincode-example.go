/*
Copyright IBM Corp. 2016 All Rights Reserved.

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
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	var myValue int
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	myValue, err = strconv.Atoi(args[0])
        if err != nil {
                return shim.Error("Expecting integer value for asset holding")
        }

	stub.PutState("defaultValue", []byte(strconv.Itoa(myValue)))
	// Initialize the chaincode
	fmt.Printf("Hello world, this is bupttest chaincode\n")
	fmt.Printf("default value = %d\n", myValue)

	return shim.Success(nil)
}

func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        // Transaction makes payment of X units from A to B
	var A, B string    // Entities
        var Aval, Bval int // Asset holdings
        var X int          // Transaction value
        var err error

        if len(args) != 3 {
                return shim.Error("Incorrect number of arguments. Expecting 3")
        }

        A = args[0]
        B = args[1]

        // Get the state from the ledger
        // TODO: will be nice to have a GetAllState call to ledger
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

		fmt.Printf("This is bupt test chaincode, invoke sucess!\n")
        return shim.Success(nil)
}

func (t *SimpleChaincode) addComKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	Type := args[0]
	Key1 := args[1]
	Key2 := args[2]
	Value := args[3]
	ComKey , err := stub.CreateCompositeKey(Type, []string {Key1, Key2})
    if err != nil {
		return shim.Error(err.Error())
    }
    err = stub.PutState(ComKey, []byte(Value))
    if err != nil {
        return shim.Error(err.Error())
    }

	fmt.Printf("add ComKey Sucess Type:%s; Key1: %s; Key2: %s; Value: %s;\n", Type, Key1, Key2, Value)
    return shim.Success(nil)
}

func (t *SimpleChaincode) queryComKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	searchKey1 := args[0]
	searchKey2 := args[1]
	rs, err := stub.GetStateByPartialCompositeKey(searchKey1,[]string{searchKey2})
	if err != nil {
		fmt.Println("find error", err.Error())
		return shim.Error(err.Error())
	}
	defer rs.Close()
	for rs.HasNext() {
		item, _ := rs.Next()
		objectType, otherPart, err := stub.SplitCompositeKey(item.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Type:" + objectType)
		fmt.Println("Key1:" + otherPart[0])
		fmt.Println("Key2:" + otherPart[1])
		fmt.Println("Value:" + string(item.Value))
	}

    return shim.Success(nil)
}

func (t *SimpleChaincode) queryRangeKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	startstr := args[0]
	endstr := args[1]
	resultIterator, _ := stub.GetStateByRange(startstr, endstr)

	defer resultIterator.Close()
	for resultIterator.HasNext() {
		item, _ := resultIterator.Next()
		fmt.Println("find Key:%s; Value:%s;", string(item.Key), string(item.Value))
	}
    return shim.Success(nil)
}
func (t *SimpleChaincode) addNewKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var myValue int

	me := args[0]
        myValue, err := strconv.Atoi(args[1])
        if err != nil {
                fmt.Printf("Error convert %s to integer: %s", args[1], err)
                return shim.Error(fmt.Sprintf("Error convert %s to integer: %s", args[1], err))
        }

	err = stub.PutState(me, []byte(strconv.Itoa(myValue)))
        if err != nil {
                return shim.Error(err.Error())
        }
	fmt.Printf("add new key %s with %d sucess!\n", me, myValue)
	return shim.Success(nil)
}
func (t *SimpleChaincode) delKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	me := args[0]

	err := stub.DelState(me)
        if err != nil {
                return shim.Error("delete error:" + me)
        }
	fmt.Printf("delete %s \n", me)
	return shim.Success(nil)
}
func (t *SimpleChaincode) queryKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var myValue int

	me := args[0]

	myValueBytes, err := stub.GetState(me)
        if err != nil {
                fmt.Printf("find %s error: %s", args[0], err)
                return shim.Error(fmt.Sprintf("Error find %s : %s", args[0], err))
        }
	if myValueBytes == nil {
                fmt.Printf("find %s null!", me)
                return shim.Error("Entity not found:" + me)
        }
	myValue, _ = strconv.Atoi(string(myValueBytes))


	fmt.Printf("the value of %s is %d \n", me, myValue)
	return shim.Success(nil)
}
func (t *SimpleChaincode) queryKeyHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var me string

	me = args[0]

	keysIter, err := stub.GetHistoryForKey(me)
        if err != nil {
                return shim.Error(fmt.Sprintf("Error getting %s history.Error accessing state:%s", args[0], err))
        }

	defer keysIter.Close()

	for keysIter.HasNext() {
		response, iterErr := keysIter.Next()
		if iterErr != nil {
			return shim.Error(fmt.Sprintf("GetHistory failed, error access state: %s", err))
		}
		txid := response.TxId
		txvalue := response.Value
		txstatus := response.IsDelete
		txtimesamp := response.Timestamp

		fmt.Printf(" Tx info - txid: %s   value: %s  isdelete: %t  datetime: %s \n", txid, string(txvalue), txstatus, txtimesamp)
	}
	fmt.Printf("query history for %s success\n", me)
	return shim.Success(nil)
}

func (t *SimpleChaincode) invokeOtherChaincode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	me := args[0]
	params1 := []string{"query", "a"}
	queryArgs := make([][]byte, len(params1))
	for i, arg := range params1 {
		queryArgs[i] = []byte(arg)
	}

	response := stub.InvokeChaincode(me, queryArgs, "mychannel")

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("failed to invoke chaincode %s, got error %s", me, response.Payload)
	fmt.Printf(errStr)
	return shim.Error(errStr)
	}
	result := string(response.Payload)

	fmt.Printf("invoke chaincode %s success, result: %s", me, result)

	return shim.Success([]byte("success InvokeChaincode"))

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		return t.invoke(stub, args)
	}else if function == "invokeOtherChaincode" {
		return t.invokeOtherChaincode(stub, args)
	}else if function == "queryKeyHistory" {
		return t.queryKeyHistory(stub, args)
	}else if function == "queryKey" {
		return t.queryKey(stub, args)
	}else if function == "delKey" {
		return t.delKey(stub, args)
	}else if function == "addNewKey" {
		return t.addNewKey(stub, args)
	}else if function == "addComKey" {
		return t.addComKey(stub, args)
	}else if function == "queryComKey" {
		return t.queryComKey(stub, args)
	}else if function == "queryRangeKey" {
		return t.queryRangeKey(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"invoke\"")
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
