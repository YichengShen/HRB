package Server

import (
	"HRB/HRBAlgorithm"
	"fmt"
)


var serverList []string
var MyId string //basically the IP of individual server
var isFault bool //check whether I should be the faulty based on configuration

var faultyCount int
var trustedCount int

type messageChan chan TcpMessage

var SendChans map[string] messageChan
var ReadChans chan TcpMessage

var isLocalMode bool //indicate whether this is a local mode
var source bool //A flag to indicate whether I am the sender

/*
Only used in Local Mode for debugging purpose
 */
var localId int


//Start up the peer
func peerStartup(local bool) {
	trustedPath := "./Configuration/trusted"
	faultedPath := "./Configuration/faulty"
	if local {
		serverList, MyId, isFault = readServerListLocal(trustedPath, faultedPath, localId)
	} else {
		serverList, MyId, isFault = readServerListNetwork(trustedPath, faultedPath)
	}
	if isFault {
		fmt.Println(MyId + " is faulty")
	}
}

func LocalModeStartup(id int) {
	isLocalMode = true
	localId = id

	peerStartup(isLocalMode)
	if serverList[0] == MyId {
		source = true
	} else {
		source = false
	}
	setUpRead()
	setUpWrite()
	HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)
	simpleTest()
}

func NetworkModeStartup() {
	isLocalMode = false
	peerStartup(isLocalMode)
	setUpRead()
	setUpWrite()
	HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)
}



/*
Reading from the network
*/

func setUpRead() {
	ReadChans = make (chan TcpMessage)
	//Start listening data
	go TcpReader(ReadChans, MyId)
	//Channel that filters the data based on the message type
	go simpleFilterRecData(ReadChans)
	//go filterRecData(readChans)
}

func simpleFilterRecData(ch chan TcpMessage) {
	for {
		message := <- ch
		data := message.Message

		switch v := data.(type) {
		case HRBAlgorithm.MSGStruct:
			fmt.Println("Msg")
			HRBAlgorithm.SimpleMsgHandler(data)
		case HRBAlgorithm.ECHOStruct:
			fmt.Println("Echo")
			HRBAlgorithm.SimpleEchoHandler(data)
		case HRBAlgorithm.REQStruct:
			fmt.Println("Req")
			HRBAlgorithm.SimpleReqHandler(data)
		case HRBAlgorithm.FWDStruct:
			fmt.Println("FWD")
			HRBAlgorithm.SimpleFwdHandler(data)
		default:
			fmt.Printf("Sending : %+v\n", v)
			fmt.Println("I do ot understand what you send")
		}
	}
}

func filterRecData (ch chan TcpMessage) {
	for {
		message := <- ch
		data := message.Message

		switch v := data.(type) {
		case HRBAlgorithm.MSGStruct:
			fmt.Println("Msg")
			//receiveMsg(data)
		case HRBAlgorithm.ECHOStruct:
			fmt.Println("Echo")
			//receiveEcho(data)
		case HRBAlgorithm.ACCStruct:
			fmt.Println("Acc")
			//receiveAcc(data)
		case HRBAlgorithm.REQStruct:
			fmt.Println("Req")
			//receiveReq(data)
		case HRBAlgorithm.FWDStruct:
			fmt.Print("FWD")
		default:
			fmt.Printf("Sending : %+v\n", v)
			fmt.Println("I do ot understand what you send")
		}
	}
}

func receiveMsg (data HRBAlgorithm.Message) {
	HRBAlgorithm.Msghandler(data)
}

func receiveEcho (data HRBAlgorithm.Message) {
	HRBAlgorithm.EchoHandler(data)
}

func receiveAcc (data HRBAlgorithm.Message) {
	HRBAlgorithm.AccHandler(data)
}

func receiveReq(data HRBAlgorithm.Message) {
	HRBAlgorithm.ReqHandler(data)
}

func receiveFwd (data HRBAlgorithm.Message) {
	HRBAlgorithm.FwdHandler(data)
}

/*
Writing to the Network
*/
//Setting up writting Channels for individual sever

func setUpWrite() {
	SendChans = make (map[string] messageChan)

	for _, serverId := range serverList {
		SendChans[serverId] = make(chan TcpMessage)
		go deliver(serverId, SendChans[serverId])
	}
	go reqSendListener()
}

func reqSendListener() {
	//fmt.Println("Inside reqSendListener")
	for {
		req := <- HRBAlgorithm.SendReqChan
		//fmt.Printf("Sending Msg: %+v\n",req.M)
		if req.SendTo == "all" {
			for _, sendChan := range SendChans {
				tcpMessage := TcpMessage{Message:req.M}
				sendChan <- tcpMessage
			}
		} else {
			sendChan := SendChans[req.SendTo]
			tcpMessage := TcpMessage{Message:req.M}
			sendChan <- tcpMessage
		}
	}

}

func deliver(ipPort string, ch chan TcpMessage) {
	TcpWriter(ipPort, ch)
}


func simpleTest() {
	//test
	if source {
		//fmt.Println("I am the source")
		m := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, Data:"abc", Header:0, Round:0}
		for _ , server := range serverList {
			tcpMessage := TcpMessage{Message:m}
			SendChans[server] <- tcpMessage
		}

		//m = HRBAlgorithm.FWDStruct{Id:MyId, SenderId:MyId}
		//message := TcpMessage{Message: m, ID:MyId}
		//for _,ch := range SendChans {
		//	ch <- message
		//}

		//m = HRBAlgorithm.ACCStruct{Id:MyId, SenderId:MyId}
		//message = TcpMessage{Message: m, ID:MyId}
		//for _,ch := range SendChans {
		//	ch <- message
		//}
		//
		//m = HRBAlgorithm.REQStruct{Id:MyId, SenderId:MyId}
		//message = TcpMessage{Message: m, ID:MyId}
		//
		//for _,ch := range SendChans {
		//	ch <- message
		//}
		//
		//m = HRBAlgorithm.ECHOStruct{Id:MyId, SenderId:MyId}
		//message = TcpMessage{Message: m, ID:MyId}
		//
		//for _,ch := range SendChans {
		//	ch <- message
		//}
		//
		//m = HRBAlgorithm.MSGStruct{Id:MyId, SenderId:MyId}
		//message = TcpMessage{Message: m, ID:MyId}
		//for _,ch := range SendChans {
		//	ch <- message
		//}
	}
}


