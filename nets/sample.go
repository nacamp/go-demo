package nets

import (
	"fmt"
	"log"
	"net/http"
)

//https://golang.org/pkg/net/http/#ListenAndServeTLS
//https://golang.org/pkg/net/http/
//http://golang.site/go/article/111-간단한-웹-서버-HTTP-서버
//http://jeonghwan-kim.github.io/dev/2019/02/18/go-todo-web-application.html

func test1() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":5000", nil)
}

func test2() {
	http.Handle("/", new(testHandler1))
	http.Handle("/t2", new(testHandler2))

	http.ListenAndServe(":5000", nil)
}

type testHandler1 struct {
	http.Handler
}

func (h *testHandler1) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	str := "Your Request Path is " + req.URL.Path
	w.Write([]byte(str))
}

type testHandler2 struct {
	http.Handler
}

func (h *testHandler2) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	str := "Your Request Path2 is " + req.URL.Path
	w.Write([]byte(str))
}

//https://gist.github.com/denji/12b3a568f092ab951456
func tls() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello World"))
	})

	err := http.ListenAndServeTLS(":443", "/opt/server.crt", "/opt/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	/*
		genrsa -out server.key 2048
		openssl ecparam -genkey -name secp384r1 -out server.key
		openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
		openssl x509 -in server.crt -out server.pem -outform PEM

		requests.get('https://127.0.0.1:443/hello', verify=False).text
		#만들때 Common Name 을 임의의 도메인을 기록해서 실행하면 에러가 난다.
		requests.get('https://127.0.0.1:443/hello', verify = os.path.join(os.getcwd(),'server.pem').text
	*/
}

func SampleMain() {
	fmt.Println("nets...")
	// test1()
	// test2()
	tls()
}
