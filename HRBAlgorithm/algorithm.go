package HRBAlgorithm

import (
	"encoding/gob"
	"fmt"
)

//(SenderID + h, bool)
var MessageReceiveSet map[string] bool
//var MessageSentSet map[string] bool

//(SenderId + h, bool)
var EchoReceiveSet map[string] bool
//(sendToID + h, bool)
var EchoSentSet map[string] bool
//Used in Not Simple Version
var EchoRecCountSet map[ECHOStruct] int
//Used in SimpleVersion because of the place where they get accept
var simpleEchoRecCountSet map[ECHOStruct] []string

//(SenderId + h, bool)
var AccReceiveSet map[string] bool
var AccSentSet map[string] bool
//(HashStr, list of ids that send Acc)
var AccRecCountSet map[ACCStruct] []string

var ReqReceiveSet map[string] bool
//(HashStr, list of ids that you send request to)
var ReqSentSet map[REQStruct] []string


var FwdReceiveSet map[string] bool
//(SendToId, bool)
//var FwdSentSet map[string] bool

//Key: value , Value: Hash(value)
var DataSet map[string] string

/*
Used in the Erasure Coding
 */
var ecDataSet map[string] [][]byte


// Send Phase to the TCPWriter
var sendChan chan Message

var faulty int
var trusted int
var total int
var MyID string

var acceptData []interface{}

var SendReqChan chan PrepareSend

//key: IP_ID, Value: index in the serverList
var serverMap map[string] int
var serverList []string



func AlgorithmSetUp(myID string, servers []string, trustedCount, faultyCount int) {
	serverMap = make(map[string] int)
	for index, server := range servers {
		serverMap[server] = index
	}
	serverList = servers

	fmt.Println("These are the servers", serverMap)
	MessageReceiveSet = make(map[string] bool)
	//MessageSentSet = make(map[string] bool)

	EchoReceiveSet = make(map[string] bool)
	EchoSentSet = make(map[string] bool)
	//Used in Acc version
	EchoRecCountSet = make (map[ECHOStruct] int)
	//used in Simple
	simpleEchoRecCountSet = make (map[ECHOStruct] []string)

	AccReceiveSet = make(map[string] bool)
	AccSentSet = make(map[string] bool)
	AccRecCountSet = make(map[ACCStruct] []string)

	ReqReceiveSet = make(map[string] bool)
	ReqSentSet = make(map[REQStruct] []string)

	FwdReceiveSet = make(map[string] bool)
	//FwdSentSet = make(map[string] bool)

	DataSet = make (map[string] string)

	ecDataSet = make(map[string] [][]byte)

	sendChan = make(chan Message)

	SendReqChan = make (chan PrepareSend)

	//change later based on config
	trusted = trustedCount
	faulty = faultyCount
	total = trusted + faulty
	MyID = myID

	//Register the concrete type for interface
	gob.Register(ACCStruct{})
	gob.Register(FWDStruct{})
	gob.Register(REQStruct{})
	gob.Register(MSGStruct{})
	gob.Register(ECHOStruct{})
}


func hasSent(l []string, val string) bool{
	for _, v := range l {
		if v == val {
			return true
		}
	}
	return false
}

func checkDataExist(expectedHash string) (bool, string) {
	for k,v := range DataSet {
		if v == expectedHash {
			//fmt.Println("Check exist" + expectedHash)
			return true, k
		}
	}
	return false,""
}


func SimpleFilterRecData(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		SimpleMsgHandler(message)
	case ECHOStruct:
		SimpleEchoHandler(message)
	case REQStruct:
		SimpleReqHandler(message)
	case FWDStruct:
		SimpleFwdHandler(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I do ot understand what you send")
	}
}

func FilterRecData (message Message) {
	switch v := message.(type) {
	case MSGStruct:
		fmt.Println("Msg")
		Msghandler(message)
	case ECHOStruct:
		fmt.Println("Echo")
		EchoHandler(message)
	case ACCStruct:
		fmt.Println("Acc")
		AccHandler(message)
	case REQStruct:
		fmt.Println("Req")
		ReqHandler(message)
	case FWDStruct:
		fmt.Print("FWD")
		FwdHandler(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I do ot understand what you send")
	}
}

func FilterSimpleErasureCodeRecData(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		fmt.Println("Msg")
		SimpleECMessageHandler(message)
	case ECHOStruct:
		fmt.Println("Echo")
		SimpleECEchoHandler(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I do ot understand what you send")
	}
}