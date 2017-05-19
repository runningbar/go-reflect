package main

import (
	"time"

	mirror "github.com/runningbar/go-reflect/server"
)

func main() {
	type test struct {
		a string
		b string
		c int
	}
	var t2 = test{"ad", "be", 1}
	var t3 = test{"a3", "d3", 2}
	var t4 = []test{{"cc", "dc", 1}, {"cc2", "dc2", 2}}
	var t1 = make(map[string]*test)
	t1["t21"] = &t2
	t1["t31"] = &t3
	var ti7 = 7
	var testp = 6
	var ptr = &testp

	var n = test{"test1", "test2", 1}
	var cn = make(map[string]*test)
	cn["n"] = &n
	//fmt.Println("cn = ", cn)
	//var nr = reflect.ValueOf(&cn).Elem()
	//var d = cn["n"]
	//var nra = reflect.ValueOf(&d).Elem().FieldByName("a")
	//var nra = nr.MapIndex(nr.MapKeys()[0]).FieldByName("a")
	//fmt.Println("nra = ", nra.String())

	mirror.PutInMirror("t11", &t1)
	mirror.PutInMirror("t41", &t4)
	mirror.PutInMirror("test", &t2)
	mirror.PutInMirror("ti7", &ti7)
	mirror.PutInMirror("ptr", &ptr)
	serve := make(chan error)
	go mirror.StartMirrorServer(8001)
	time.Sleep(5 * time.Second)
	testp = 69
	t1["t21"].a = "dflkj"
	ti7 = 8
	//fmt.Println("main, [t2, ti7] = ", t2, ti7)
	//mirror.GetByKey_test("t11")

	n.a = "fsalfkjsl"
	cn["n"].a = "fasl"
	//fmt.Println("correct:", n, *cn["n"])
	//fmt.Println("now:", nra.String(), nra.Kind(), nra.CanAddr(), nra.CanSet())
	//fmt.Println("conclusion:", nr)

	<-serve
	//mirror.StartMirrorServer(61012)
	//fmt.Println(err1)
	//fmt.Println(err2)
}
