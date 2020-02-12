package Server

import (
	"HRB/HRBAlgorithm"
	"math/rand"
	"time"
)

/*
Note id is always the Ip + Port for Local
 */

//var protocalSendChan map[string] messageChan
var protocalSendChan chan TcpMessage
var protocalReadChan chan TcpMessage


/*
Only used in Local Mode for debugging purpose
 */


//For One round
func ProtocalStart() {
	//fmt.Println("Local Setup")

	//HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)
	if isSourceFault {
		trustedCount = trustedCount - 1;
		faultyCount = faultyCount + 1;
		//fmt.Println("FaultyCount ", faultyCount)
		if MyId == serverList[0] {
			isFault = true
			isSourceFault = true
		}
	}

	//General setup

	HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)


	if algorithm == 1 {
		hashSimpleSetup()
		if source {
			simpleBroadcast(dataSize, round)
		}
	} else if algorithm == 2 {
		hashComplexSetup()
		if source {
			simpleBroadcast(dataSize, round)
		}
	} else if algorithm == 3 {
		hashECSimpleSetup()
		if isSourceFault {

		} else {
			if source {
				HRBAlgorithm.SimpleECBroadCast(dataSize, round)
			}
		}
	} else if algorithm == 4 {
		hashECComplexSetup()
		if source {
			HRBAlgorithm.ComplexECBroadCast(dataSize,round)
		}
	} else if algorithm == 5 {
		digestSetup()
		if source {
			//fmt.Println("Digest Broadcast")
			HRBAlgorithm.BroadcastPrepare(dataSize, round)
		}

	} else if algorithm == 6 {
		codedSetup()
		if source {
			//fmt.Println("NCBA")
			HRBAlgorithm.ECByzBroadCast(dataSize, round)
		}

	} else if algorithm == 7 {
		codedCrashSetup()
		if source {
			//fmt.Println("Crash Ccoded Broadcast")
			HRBAlgorithm.CrashECBroadCast(dataSize, round)
		}
	} else if algorithm == 8 {
		optimalSetup()
		if source {
			HRBAlgorithm.OptimalBroadcast(dataSize,round)
		}
	} else {
		//fmt.Println("Do not understand what you give")
	}

}


/*
Seven different algorithms to choose
 */

func hashSimpleSetup() {
	protocalReadChan = setUpRead()
	go filterSimple(protocalReadChan)
	go setUpWrite()
}

func hashComplexSetup() {
	protocalReadChan = setUpRead()
	go filter(protocalReadChan)
	go setUpWrite()
}


func hashECSimpleSetup() {
	protocalReadChan = setUpRead()
	go filterSimpleEC(protocalReadChan)
	go setUpWrite()
}


func hashECComplexSetup() {
	protocalReadChan = setUpRead()
	go filterComplexEc(protocalReadChan)
	go setUpWrite()
}

func digestSetup() {
	HRBAlgorithm.InitDigest()
	protocalReadChan = setUpRead()
	go filterDigest(protocalReadChan)
	go setUpWrite()
}

func codedSetup() {
	HRBAlgorithm.InitByzCode()
	protocalReadChan= setUpRead()
	go filterByzCode(protocalReadChan)
	go setUpWrite()
}

func codedCrashSetup() {
	HRBAlgorithm.InitCrash()
	protocalReadChan= setUpRead()
	go filterCrashCoded(protocalReadChan)
	go setUpWrite()
}

func optimalSetup() {
	HRBAlgorithm.InitOptimal()
	protocalReadChan= setUpRead()
	go filterOptimal(protocalReadChan)
	go setUpWrite()
}


func NetworkModeStartup(id int) {
	isLocalMode = false
	setUpRead()
	HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)
}


/*
Reading from the network
*/

func setUpRead() chan TcpMessage{
	protocalReadChan = make (chan TcpMessage)
	//Start listening data
	go TcpReader(protocalReadChan, MyId)
	return protocalReadChan
}

func filter (ch chan TcpMessage) {
	for {
		message := <-ch
		HRBAlgorithm.FilterRecData(message.Message)
	}
	//Call the filter method
}
func filterSimple(ch chan TcpMessage) {
	for {
		message := <-ch
		HRBAlgorithm.SimpleFilterRecData(message.Message)
	}
	//Call the filter method
}

func filterSimpleEC(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterSimpleErasureCodeRecData(message.Message)
	}
}

func filterComplexEc(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterComplexErasureRecData(message.Message)
	}
}

func filterDigest(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterDigest(message.Message)
	}
}

func filterByzCode(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterByzCode(message.Message)
	}
}

func filterCrashCoded(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterCrashCode(message.Message)
	}
}

func filterOptimal(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterOptimal(message.Message)
	}
}

/*
Writing to the Network
*/
//Setting up writting Channels for individual sever



func setUpWrite() {
	protocalSendChan = make(chan TcpMessage)
	go deliver(MyId, protocalSendChan)
	for {
		req := <- HRBAlgorithm.SendReqChan
		//fmt.Printf("Sending Msg: %+v\n",req.M)
		if req.SendTo == "all" ||  req.SendTo == ""{
			//fmt.Println("Protocal send to all")
			for _, id := range serverList {
				//fmt.Println("Protocal send to ", id)
				tcpMessage := TcpMessage{Message:req.M, ID:id}
				protocalSendChan <- tcpMessage
			}
		} else {
			//fmt.Println("Protocal send to ", req.SendTo)
			tcpMessage := TcpMessage{Message:req.M, ID:req.SendTo}
			protocalSendChan <- tcpMessage
		}
	}
}



func deliver(ipPort string, ch chan TcpMessage) {
	TcpWriter(ch)
}


/*
Simple Testing
 */

//func testEcSourceFault(s string) {
//	//test
//	if source {
//		shards := HRBAlgorithm.Encode(s, faultyCount + 1, trustedCount - 1)
//		//Get the string version of the string
//		hashStr := HRBAlgorithm.ConvertBytesToString(HRBAlgorithm.Hash([]byte(s)))
//
//		for id , server := range serverList {
//			if id == 2 || id == 3{
//				wrong := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, HashData: hashStr, Data: "", Header:0, Round:0}
//				tcpMessage := TcpMessage{Message: wrong}
//				protocalSendChan[server] <- tcpMessage
//			} else if id == 4 || id ==5 {
//
//			} else {
//				codeString := HRBAlgorithm.ConvertBytesToString(shards[id])
//				correct := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, HashData: hashStr, Data: codeString, Header:0, Round:0}
//				tcpMessage := TcpMessage{Message:correct}
//				protocalSendChan[server] <- tcpMessage
//			}
//		}
//	}
//}
//
//func testSourceFault(s string) {
//	if source {
//		m := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, Data:s, Header:0, Round:0}
//		faultym :=  HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, Data:"", Header:0, Round:0}
//		for id , server := range serverList {
//			if id == 2 || id == 3 {
//				tcpMessage := TcpMessage{Message:faultym}
//				protocalSendChan[server] <- tcpMessage
//			} else if id == 4 || id == 5 {
//
//			} else {
//				tcpMessage := TcpMessage{Message:m}
//				protocalSendChan[server] <- tcpMessage
//			}
//		}
//	}
//}
//

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}


/*
Used in Algorithm 1 and 2
 */
func simpleBroadcast(byte_length, round int) {
	if source {
		for i := 0; i < round; i++ {
			time.Sleep(1*time.Second)
			s := RandStringBytes(byte_length)
			m := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, Data: s, Header:0, Round:i}
			for _, id := range serverList {
				//fmt.Println("Protocal send to ", id)
				tcpMessage := TcpMessage{Message:m, ID:id}
				protocalSendChan <- tcpMessage
			}
		}
	}
}





