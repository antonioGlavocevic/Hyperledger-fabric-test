package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// MyChaincode example
type MyChaincode struct {
}

// Course struct. Course tags use json
type Course struct {
	CourseID   string `json:"courseID"`
	CourseName string `json:"courseName"`
}

// Student struct. Structure tags use json
type Student struct {
	StudentID string   `json:"studentID"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Courses   []string `json:"course"`
}

// StudentQuery struct. Structure tags use json
type StudentQuery struct {
	StudentID string   `json:"studentID"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Courses   []Course `json:"course"`
}

// Init - Function called on chaincode instantiation
func (t *MyChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("mychaincode Init")
	return shim.Success(nil)
}

// Invoke - Function called on chaincode invoke
func (t *MyChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("mychaincode Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "reg" {
		return t.register(stub, args)
	} else if function == "createCourse" {
		return t.createCourse(stub, args)
	} else if function == "del" {
		return t.delete(stub, args)
	} else if function == "addCourse" {
		return t.addCourse(stub, args)
	} else if function == "changeCourse" {
		return t.changeCourse(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	} else if function == "queryStudent" {
		return t.queryStudent(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"reg\" \"del\" \"addCourse\" \"changeCourse\" \"query\" \"queryStudent\"")
}

// register - register a new student
func (t *MyChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting 4\"}"
		return shim.Error(jsonResp)
	}

	dupe, _ := stub.GetState(args[0])
	if dupe != nil {
		jsonResp := "{\"Error\": " + args[0] + " already exists\"}"
		return shim.Error(jsonResp)
	}

	courseAsBytes, err := stub.GetState(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}
	if courseAsBytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + args[3] + "\"}"
		return shim.Error(jsonResp)
	}

	courses := []string{args[3]}
	student := Student{StudentID: args[0], FirstName: args[1], LastName: args[2], Courses: courses}
	studentAsBytes, err := json.Marshal(student)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(args[0], studentAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// createCourse - create a new course
func (t *MyChaincode) createCourse(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting 2\"}"
		return shim.Error(jsonResp)
	}

	dupe, _ := stub.GetState(args[0])
	if dupe != nil {
		jsonResp := "{\"Error\": " + args[0] + " already exists\"}"
		return shim.Error(jsonResp)
	}

	course := Course{CourseID: args[0], CourseName: args[1]}
	courseAsBytes, err := json.Marshal(course)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(args[0], courseAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// delete - delete a user and return tokens to pool
func (t *MyChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting 1\"}"
		return shim.Error(jsonResp)
	}

	studentAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if studentAsBytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// addCourse - add a course to a student
func (t *MyChaincode) addCourse(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting 2\"}"
		return shim.Error(jsonResp)
	}

	studentAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if studentAsBytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	student := Student{}
	err = json.Unmarshal(studentAsBytes, &student)
	if err != nil {
		return shim.Error(err.Error())
	}

	courseAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	if courseAsBytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + args[1] + "\"}"
		return shim.Error(jsonResp)
	}

	student.Courses = append(student.Courses, args[1])

	studentAsBytes, err = json.Marshal(student)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(args[0], studentAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// changeCourse - change a course name
func (t *MyChaincode) changeCourse(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting 2\"}"
		return shim.Error(jsonResp)
	}

	courseAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if courseAsBytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	course := Course{}
	err = json.Unmarshal(courseAsBytes, &course)
	if err != nil {
		return shim.Error(err.Error())
	}

	course.CourseName = args[1]

	courseAsBytes, err = json.Marshal(course)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(args[0], courseAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// query - get data by state id
func (t *MyChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting student id to query\"}"
		return shim.Error(jsonResp)
	}

	queryAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if queryAsBytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(queryAsBytes)
}

// queryStudent - get student by student id
func (t *MyChaincode) queryStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting student id to query\"}"
		return shim.Error(jsonResp)
	}

	studentAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if studentAsBytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	student := Student{}
	err = json.Unmarshal(studentAsBytes, &student)
	if err != nil {
		return shim.Error(err.Error())
	}

	var courses []Course
	for i := range student.Courses {
		courseAsBytes, err := stub.GetState(student.Courses[i])
		if err != nil {
			return shim.Error(err.Error())
		}
		if studentAsBytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
			return shim.Error(jsonResp)
		}

		course := Course{}
		err = json.Unmarshal(courseAsBytes, &course)
		if err != nil {
			return shim.Error(err.Error())
		}

		courses = append(courses, course)
	}

	studentQuery := StudentQuery{StudentID: student.StudentID, FirstName: student.FirstName, LastName: student.LastName, Courses: courses}
	resp, err := json.Marshal(studentQuery)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(resp)
}

func main() {
	err := shim.Start(new(MyChaincode))
	if err != nil {
		fmt.Printf("Error starting mychaincode: %s", err)
	}
}
