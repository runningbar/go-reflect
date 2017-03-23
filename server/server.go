package mirror

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func createJSON(key string) []byte {
	var js []byte
	var r []Mirror
	if key == "all" {
		for k, v := range valueMap {
			if strings.Index(k, ".") == -1 {
				r = append(r, startReflect(v, k, "", 1))
			}
		}
	} else {
		r = append(r, reflectByKey(key))
	}
	js, _ = json.Marshal(r)
	return js
}

/*func index(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./static/index.html")
}*/

func query(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	req.ParseForm()
	var key = req.FormValue("key")
	var js = createJSON(key)
	w.Write(js)
}

func createTestData() {
	type test struct {
		a string
		b string
		c int
	}
	var t2 = test{"a", "b", 1}
	var t3 = test{"a", "d", 2}
	var t4 = []test{{"c", "d", 1}, {"c", "d", 2}}
	var t1 = make(map[string]test)
	t1["t2"] = t2
	t1["t3"] = t3
	PutInMirror("t1", t1)
	PutInMirror("t4", t4)
}

// StartMirrorServer that's all
func StartMirrorServer() error {
	createTestData()
	const PORT = 12345
	//http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/query", query)
	err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
	return err
}
