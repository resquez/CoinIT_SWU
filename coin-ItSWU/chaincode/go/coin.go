package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

// 구조체
type Wallet struct {
	Name string `json: "name"`
	ID string `json: "id"`
	Mileage string `json: "mileage"`
}

type TransferLog struct {
	Index string `json: "index"`
	Time string `json: "time"`
	Source string `json: "source"`
	Destination string `json: "destination"`
	Value string `json: "value"`
}

type LogIndex struct {
	ID string
	LastIdx int
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
		// init 시 인수 입력X 초깃값X
		return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	function, args := APIstub.GetFunctionAndParameters()
	
	if function == "initWallet" {
		return s.initWallet(APIstub, args)
	} else if function == "getWallet" {
		return s.getWallet(APIstub, args)
	} else if function == "testGetLog" {
		return s.testGetLog(APIstub, args)
	} else if function == "issueMileage" {
		return s.issueMileage(APIstub, args)
	} else if function == "purchaseWithMileage" {
		return s.purchaseWithMileage(APIstub, args)
	} else if function == "resetMileage" {
		return s.resetMileage(APIstub, args)
	} else {
		return s.getLog(APIstub, args)
	}
	fmt.Println("function error : " + function)
	return shim.Error("Invalid function")
}

func (s *SmartContract) initWallet(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Declare wallets
	var wallet = Wallet{Name: args[0], ID: args[1], Mileage: "0"}
	// Convert my_wallet to []byte
	walletasJSONBytes, _ := json.Marshal(wallet)
	err := APIstub.PutState(wallet.ID, walletasJSONBytes)
	if err != nil {
		return shim.Error("Failed to create asset " + wallet.Name)
	}

	// Declare logIndex
	var logIndex = LogIndex{ID: args[1], LastIdx: -1}
	var key = args[1] + "Log"
	logIndexAsBytes, _ := json.Marshal(logIndex)
	APIstub.PutState(key, logIndexAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getWallet(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	walletAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		fmt.Println(err.Error())
	}

	wallet := Wallet{}
	json.Unmarshal(walletAsBytes, &wallet)

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	if bArrayMemberAlreadyWritten == true {
		buffer.WriteString(",")
	}
	buffer.WriteString("{\"Name\":")
	buffer.WriteString("\"")
	buffer.WriteString(wallet.Name)
	buffer.WriteString("\"")

	buffer.WriteString(", \"ID\":")
	buffer.WriteString("\"")
	buffer.WriteString(wallet.ID)
	buffer.WriteString("\"")

	buffer.WriteString(", \"Mileage\":")
	buffer.WriteString("\"")
	buffer.WriteString(wallet.Mileage)
	buffer.WriteString("\"")

	buffer.WriteString("}")
	bArrayMemberAlreadyWritten = true
	buffer.WriteString("]\n")

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) issueMileage(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1. 변수
	var Aval, X int
	var err error
	var key, now string

	// 2. 인수 유효성 체크
	if len(args) != 2 {
		return shim.Error("Error Incorrect arguments.")
	}

	t := time.Now().UTC()
  	t = t.In(time.FixedZone("KST", 9*60*60))
  	now = t.Format("2006-01-02 15:04:05")

 	key = args[0] + "Log"
	logIndexAsBytes, _ := APIstub.GetState(key)
	logIndex := LogIndex{}
	json.Unmarshal(logIndexAsBytes, &logIndex)
	logIndex.LastIdx = logIndex.LastIdx + 1
	idx := strconv.Itoa(logIndex.LastIdx)

	updatedIdxBytes, _ := json.Marshal(logIndex)
	APIstub.PutState(args[0] + "Log", updatedIdxBytes)

	// 3. 인수에서 정보 가져와 바이트 변환 
	walletAsBytes, err := APIstub.GetState(args[0])  //(a노드 정보)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if walletAsBytes == nil {
		return shim.Error("Entity not found")
	}
	wallet := Wallet{}
	json.Unmarshal(walletAsBytes, &wallet)

	var log = TransferLog{Index: idx, Time: now, Source: "SWU", Destination: args[0], Value: "+" + args[1]}	

	// 4. 마일리지 누적
	Aval, _ = strconv.Atoi(string(wallet.Mileage))
	X, err = strconv.Atoi(args[1])

	wallet.Mileage = strconv.Itoa(Aval + X)

	// 5. 원장에 저장
	updatedWalletAsBytes, _ := json.Marshal(wallet)
	APIstub.PutState(args[0], updatedWalletAsBytes)

	logasJSONBytes, _ := json.Marshal(log)
	key = args[0] + "LogNo." + idx
	APIstub.PutState(key, logasJSONBytes)

	// 6. 결과 출력
	var buffer bytes.Buffer

	buffer.WriteString("Current mileage\":")
	buffer.WriteString("\"")
	buffer.WriteString(wallet.Mileage)
	buffer.WriteString("\"")
	buffer.WriteString("\n")

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) purchaseWithMileage(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Aval, Bval int
	var X int
	var key_A, key_B, now string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	t := time.Now().UTC()
  	t = t.In(time.FixedZone("KST", 9*60*60))
  	now = t.Format("2006-01-02 15:04:05") 

  	key_A = args[0] + "Log"
	logIndexAsBytes_A, _ := APIstub.GetState(key_A)
	logIndex_A := LogIndex{}
	json.Unmarshal(logIndexAsBytes_A, &logIndex_A)
	logIndex_A.LastIdx = logIndex_A.LastIdx + 1
	idx_A := strconv.Itoa(logIndex_A.LastIdx)

	updatedIdxBytes_A, _ := json.Marshal(logIndex_A)
	APIstub.PutState(args[0] + "Log", updatedIdxBytes_A)

	key_B = args[1] + "Log"
	logIndexAsBytes_B, _ := APIstub.GetState(key_B)
	logIndex_B := LogIndex{}
	json.Unmarshal(logIndexAsBytes_B, &logIndex_B)
	logIndex_B.LastIdx = logIndex_B.LastIdx + 1
	idx_B := strconv.Itoa(logIndex_B.LastIdx)

	updatedIdxBytes_B, _ := json.Marshal(logIndex_B)
	APIstub.PutState(args[1] + "Log", updatedIdxBytes_B)

	walletAAsbytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if walletAAsbytes == nil {
		return shim.Error("Entity not found")
	}

	wallet_A := Wallet{}
	json.Unmarshal(walletAAsbytes, &wallet_A)
	Aval, _ = strconv.Atoi(wallet_A.Mileage)

	walletBAsbytes, err := APIstub.GetState(args[1])
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if walletBAsbytes == nil {
		return shim.Error("Entity not found")
	}

	wallet_B := Wallet{}
	json.Unmarshal(walletBAsbytes, &wallet_B)
	Bval, _ = strconv.Atoi(wallet_B.Mileage)

	X, err = strconv.Atoi(args[2])

	wallet_A.Mileage = strconv.Itoa(Aval - X)
	wallet_B.Mileage = strconv.Itoa(Bval + X)
	updatedAAsBytes, _ := json.Marshal(wallet_A)
	updatedBAsBytes, _ := json.Marshal(wallet_B)
	APIstub.PutState(args[0], updatedAAsBytes)
	APIstub.PutState(args[1], updatedBAsBytes)

	var log_A = TransferLog{Index: idx_A, Time: now, Source: args[0], Destination: args[1], Value: "-" + args[2]}
	key_A = args[0] + "LogNo." + idx_A
	logasJSONBytes_A, _ := json.Marshal(log_A)
	APIstub.PutState(key_A, logasJSONBytes_A)

	var log_B = TransferLog{Index: idx_B, Time: now, Source: args[0], Destination: args[1], Value: "+" + args[2]}
	key_B = args[1] + "LogNo." + idx_B
	logasJSONBytes_B, _ := json.Marshal(log_B)
	APIstub.PutState(key_B, logasJSONBytes_B)

	return shim.Success(nil)
}

func (s *SmartContract) resetMileage(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Aval string
	var key,now string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	t := time.Now().UTC()
 	t = t.In(time.FixedZone("KST", 9*60*60))
  	now = t.Format("2006-01-02 15:04:05") 

 	key = args[0] + "Log"
	logIndexAsBytes, _ := APIstub.GetState(key)
	logIndex := LogIndex{}
	json.Unmarshal(logIndexAsBytes, &logIndex)
	logIndex.LastIdx = logIndex.LastIdx + 1
	idx := strconv.Itoa(logIndex.LastIdx)

	updatedIdxBytes, _ := json.Marshal(logIndex)
	APIstub.PutState(args[0] + "Log", updatedIdxBytes)

	walletAsbytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if walletAsbytes == nil {
		return shim.Error("Entity not found")
	}

	wallet := Wallet{}
	json.Unmarshal(walletAsbytes, &wallet)
	Aval = wallet.Mileage
	
	wallet.Mileage = strconv.Itoa(0)
	updatedAsBytes, _ := json.Marshal(wallet)
	APIstub.PutState(args[0], updatedAsBytes)

	var log = TransferLog{Index: idx, Time: now, Source: "SWU", Destination: args[0], Value: "reset -" + Aval}
	key = args[0] + "LogNo." + idx
	logasJSONBytes, _ := json.Marshal(log)
	APIstub.PutState(key, logasJSONBytes)

	return shim.Success(nil) // args[0] : A   string(Aval) : 0
}

func (s *SmartContract) getLog(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	var idxkey = args[0] + "Log"
	logIdxAsBytes, _ := APIstub.GetState(idxkey )
	logIdx := LogIndex{}
	json.Unmarshal(logIdxAsBytes, &logIdx)
	idxStr := strconv.Itoa(logIdx.LastIdx + 1)

	var startKey = args[0] + "LogNo.0"
	var endKey = args[0] + "LogNo." + idxStr
	fmt.Println(startKey)
	fmt.Println(endKey)

	resultsIter, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIter.Close()
	
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIter.HasNext() {
		queryResponse, err := resultsIter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"ID\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		
		buffer.WriteString(", \"Log\":")
		
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]\n")
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) testGetLog(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	key := args[0]+"LogNo."+args[1]
	logAsBytes, err := APIstub.GetState(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	transferLog:= TransferLog{}
	json.Unmarshal(logAsBytes, &transferLog)

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	if bArrayMemberAlreadyWritten == true {
		buffer.WriteString(",")
	}
	buffer.WriteString("{\"Index\":")
	buffer.WriteString("\"")
	buffer.WriteString(transferLog.Index)
	buffer.WriteString("\"")

	buffer.WriteString(", \"Time \":")
	buffer.WriteString("\"")
	buffer.WriteString(transferLog.Time)
	buffer.WriteString("\"")

	buffer.WriteString(", \"Source\":")
	buffer.WriteString("\"")
	buffer.WriteString(transferLog.Source)
	buffer.WriteString("\"")

	buffer.WriteString(", \"Destination\":")
	buffer.WriteString("\"")
	buffer.WriteString(transferLog.Destination)
	buffer.WriteString("\"")

	buffer.WriteString(", \"Value\":")
	buffer.WriteString("\"")
	buffer.WriteString(transferLog.Value)
	buffer.WriteString("\"")

	buffer.WriteString("}")
	bArrayMemberAlreadyWritten = true
	buffer.WriteString("]\n")

	return shim.Success(buffer.Bytes())

}

func main() {
	if err := shim.Start(new(SmartContract)); err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}