package main

import (
	"errors"
	"flag"
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
	flag.StringVar(&protocolOut, "p", "tcp", "Протокол службы логирования (tcp, tcp4, tcp6 или udp, udp4, udp6)")
	flag.StringVar(&portIn, "pi", "8800", "Порт входящих сообщений")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", indexHandleFunc)
	http.HandleFunc("/ping/", pingHandleFunc)
	http.HandleFunc("/data/", dataHandleFunc)

	http.ListenAndServe(":"+portIn, nil)
}

func connectBody(mess []byte) error {

	if protocolOut == "tcp" || protocolOut == "tcp4" || protocolOut == "tcp6" {

		remAddr, err := net.ResolveTCPAddr(protocolOut, addressOut+":"+portOut)
		if err != nil {
			return err
		}

		conn, err := net.DialTCP(protocolOut, nil, remAddr)
		if err != nil {
			return err
		}

		conn.Write(mess)
		conn.Close()

	} else if protocolOut == "udp" || protocolOut == "udp4" || protocolOut == "udp6" {

		remAddr, err := net.ResolveUDPAddr(protocolOut, addressOut+":"+portOut)
		if err != nil {
			return err
		}

		conn, err := net.DialUDP(protocolOut, nil, remAddr)
		if err != nil {
			return err
		}

		conn.Write(mess)
		conn.Close()

	} else {
		return errors.New("неверный протокол логирования, должен быть tcp, tcp4, tcp6 или udp, udp4, udp6")
	}

	return nil
}

func dataHandleFunc(w http.ResponseWriter, r *http.Request) {
	resp := []byte("")

	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp = []byte(err.Error())
	}

	errCon := connectBody(body)
	if errCon != nil {
		resp = []byte(err.Error())
	} else {
		resp = []byte("OK")
	}
	w.Write(resp)
}

func pingHandleFunc(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}

func indexHandleFunc(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hi There!"))
}
