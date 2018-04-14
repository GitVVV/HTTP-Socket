package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

var (
	addressOut  string
	portOut     string
	protocolOut string
	portIn      string
)

func init() {
	flag.StringVar(&addressOut, "a", "logs.svc.aku.com", "Адрес службы логирования")
	flag.StringVar(&portOut, "po", "3338", "Порт службы логирования")
	flag.StringVar(&protocolOut, "p", "tcp", "Протокол службы логирования (tcp или udp)")
	flag.StringVar(&portIn, "pi", "8800", "Порт входящих сообщений")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", indexHandleFunc)
	http.HandleFunc("/ping/", pingHandleFunc)

	http.ListenAndServe(":"+portIn, nil)
}

func connectBody(mess string) error {
	conn, err := net.Dial(protocolOut, addressOut+":"+portOut)
	if err != nil {
		return err
	}

	fmt.Fprintf(conn, mess)
	conn.Close()

	return nil
}

func indexHandleFunc(w http.ResponseWriter, r *http.Request) {
	resp := []byte("")

	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	err := connectBody(string(body))
	if err != nil {
		resp = []byte(err.Error())
	} else {
		resp = []byte("OK")
	}
	w.Write(resp)
}

func pingHandleFunc(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}
