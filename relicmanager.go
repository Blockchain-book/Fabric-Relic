package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//作接收器
type RelicChaincode struct{}

type RelicOrder struct {
	OrderID        string `json:"orderID"`
	OrderValue     string `json:"orderValue"`
	OrderDate      string `json:"orderDate"`
	OrderStatus    string `json:"orderStatus"`
	ProvideID      string `json:"provideID"`
	BuyerID        string `json:"buyerID"`
	SellerID       string `json:"sellerID"`
	RelicID        string `json:"relicID"`
	GovNum         string `json:"govNum"`
	RelicName      string `json:"relicName"`
	RelicDescribe  string `json:"relicDescribe"`
	RelicDataURL   string `json:"relicDataURL"`
	ImageURL       string `json:"imageURL"`
	InputDate      string `json:"inputDate"`
	JudgeName      string `json:"judgeName"`
	JudgeNum       string `json:"judgeNum"`
	JudgeOrgID     string `json:"judgeOrgID"`
	Evaluation     string `json:"evaluation"`
	EvaluationName string `json:"evaluationName"`
	EvaluationNum  string `json:"evaluationNum"`
	NewValue       string `json:"newValue"`
	NewValueDate   string `json:"newValueDate"`
	OwnerID        string `json:"ownerID"`
	RelicStatus    string `json:"relicStatus"`
}

func (t *RelicChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("Success Init"))
}

func (t *RelicChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "addneworder":
		return t.AddNewOrder(stub, args)
	case "getorder":
		return t.GetOrder(stub, args)
	case "getorderbyrelicid":
		return t.GetOrderByrelicID(stub, args)
	default:
		fmt.Println("function error")
		return shim.Error("function" + function + "is not exist")
	}
}

func (t *RelicChaincode) AddNewOrder(stub shim.ChaincodeStubInterface, orderData []string) peer.Response {
	if len(orderData) != 24 {
		//return fmt.Errorf("the number of args is %d, not 24", len(orderData))
		return shim.Error("the number of args is %d, not 24")
	}
	relicOrder := new(RelicOrder)
	relicOrder.OrderID = orderData[0]
	relicOrder.OrderValue = orderData[1]
	relicOrder.OrderDate = orderData[2]
	relicOrder.OrderStatus = orderData[3]
	relicOrder.ProvideID = orderData[4]
	relicOrder.BuyerID = orderData[5]
	relicOrder.SellerID = orderData[6]
	relicOrder.RelicID = orderData[7]
	relicOrder.GovNum = orderData[8]
	relicOrder.RelicName = orderData[9]
	relicOrder.RelicDescribe = orderData[10]
	relicOrder.RelicDataURL = orderData[11]
	relicOrder.ImageURL = orderData[12]
	relicOrder.InputDate = orderData[13]
	relicOrder.JudgeName = orderData[14]
	relicOrder.JudgeNum = orderData[15]
	relicOrder.JudgeOrgID = orderData[16]
	relicOrder.Evaluation = orderData[17]
	relicOrder.EvaluationName = orderData[18]
	relicOrder.EvaluationNum = orderData[19]
	relicOrder.NewValue = orderData[20]
	relicOrder.NewValueDate = orderData[21]
	relicOrder.OwnerID = orderData[22]
	relicOrder.RelicStatus = orderData[23]

	data1, err := stub.GetState(relicOrder.OrderID)
	if data1 != nil && err == nil {
		return shim.Error("relic has been existed")
	}
	data, err := json.Marshal(relicOrder)
	if err != nil {
		return shim.Error("relic order marshal is failed for " + err.Error())
	}
	err = stub.PutState(relicOrder.OrderID, data)
	if err != nil {
		return shim.Error("put State for relic order failed")
	}
	return shim.Success([]byte("add relic order success"))
}

func (t *RelicChaincode) GetOrder(stub shim.ChaincodeStubInterface, orderID []string) peer.Response {
	if len(orderID) != 1 {
		return shim.Error("the number of args is not 1")
	}
	relicData, err := stub.GetState(orderID[0])
	if err != nil {
		return shim.Error("getting relic order error")
	}
	if relicData == nil {
		return shim.Error("relic order is not exist")
	}
	return shim.Success(relicData)
}

func (t *RelicChaincode) GetOrderByrelicID(stub shim.ChaincodeStubInterface, relicID []string) peer.Response {
	if len(relicID) != 1 {
		return shim.Error("the number of args is not 1")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"relicID\":\"%s\"}}", relicID[0])
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error("failed to query by relicID")
	}
	return shim.Success(queryResult)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
	resultIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		fmt.Println(resultIterator)
		fmt.Println(err)
		return nil, err
	}
	defer resultIterator.Close()
	buffer, err := constructQueryResponseFromIterator(resultIterator)
	if err != nil {
		return nil, err
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return &buffer, nil
}

func main() {
	fmt.Println("hello, this is a project to manager relic order info")
	err := shim.Start(new(RelicChaincode))
	if err != nil {
		fmt.Println("Error when starting Relic Chain code:" + err.Error())
	}
}
