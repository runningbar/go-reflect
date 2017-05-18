package main

import mirror "github.com/hyperledger/fabric/common/mirror/server"

func main() {

	type test struct {
		a string
		b string
		c int
	}
	var t2 = test{"a", "b", 1}
	var t3 = test{"a", "d", 2}
	var t4 = []test{{"c", "d", 1}, {"c", "d", 2}}
	var t1 = make(map[string]test)
	t1["t21"] = t2
	t1["t31"] = t3
	mirror.PutInMirror("t11", t1)
	mirror.PutInMirror("t41", t4)
	mirror.PutInMirror("test", nil)
	go mirror.StartMirrorServer(61011)
	mirror.StartMirrorServer(61012)
	//fmt.Println(err1)
	//fmt.Println(err2)
}
