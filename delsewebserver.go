package main

import (
	"io"
	"net/http"
	"bufio"
	"os"
	"fmt"
	"strconv"
	"regexp"
	"strings"
)

const defaultport int = 80
const portprefix string = ":"

var listeningAddr string = ""
var mux map[string]func(http.ResponseWriter, *http.Request)

// Stampa il messaggio di benvenuto e la lista dei comandi
func welcomeMessage() {
	println("[delsewebserver]")
	println("")
	println("Questo software avvia un server-web sulla porta desiderata e rimane in attesa di comandi.")
	println("")
	println("COMANDI DISPONIBILI: ")
	println("/quitserver - Arresta il webserver ed esce.")
	println("help - Stampa la guida. DA IMPLEMENTARE!")
	println("")
}

func requestListeningPort() string {
	reader := bufio.NewReader(os.Stdin)
	inputstring, _ := reader.ReadString('\n')
	return inputstring
}

func testListeningPort(p string) int {
	var resultport int
	var re *regexp.Regexp

	if p != "" {
		//fmt.Println("Porta da testare: ", p)

		re = regexp.MustCompile(`\r?\n`)
		p = re.ReplaceAllString(p, " ")		
		//fmt.Println("Porta da testare dopo regexp.ReplaceAllString: ", p)
		//provo a convertire p (string) in port (int)
		p = strings.TrimSpace(p)
		//fmt.Println("Porta da testare dopo TrimSpace: ", p)
		port, err := strconv.Atoi(p)
		//fmt.Println("Porta da testare dopo Atoi: ", port)
		if (err != nil) {
			//fmt.Println("Numero di porta consentito da 0 a 65535 - <0 ", "Verrà utilizzata la porta di default: ", defaultport, err)
		 	resultport = defaultport
		} else {
			resultport = port
		}

		if resultport < 0 {
			//fmt.Println("Numero di porta consentito da 0 a 65535 - <0 ", "Verrà utilizzata la porta di default: ", defaultport)
			resultport = defaultport
		}
		if resultport > 65535 {
			//fmt.Println("Numero di porta consentito da 0 a 65535 - >65535 ", "Verrà utilizzata la porta di default: ", defaultport)
			resultport = defaultport	
		}
	}
	return resultport
}

func askAndTestListeningPort() {
	fmt.Println("***********ATTENZIONE!***********")
	fmt.Println("Verrà utilizzata la porta ", defaultport)
	fmt.Println("Premere INVIO per continuare o digitare il numero di porta desiderato seguito da INVIO.")
	
	setListeningPort(testListeningPort(requestListeningPort()))
	
	fmt.Println("")
	fmt.Println("Verrà utilizzata la porta :", listeningAddr)
	fmt.Println("")
}

func setListeningPort(port int) {
	listeningAddr = ":" + strconv.Itoa(port)
	//fmt.Println("Listening Address ", listeningAddr)
}

func startServer() {
	fmt.Println("")
	fmt.Println("Imposto il server sulla porta ", listeningAddr, "... OK")
	fmt.Println("")

	server := http.Server{
		Addr:    listeningAddr,
		Handler: &myHandler{},
	}

	fmt.Println("")
	fmt.Println("Creazione funzionalità disponibili...")

	mux = make(map[string]func(http.ResponseWriter, *http.Request))

	mux["/"] = hello
	//Per ogni coppia di valori Stringa nomefunzione creare una funzione che gestisca w http.ResponseWriter, r *http.Request
	//
	// Es con nomefunzione = "hello":
	// func hello(w http.ResponseWriter, r *http.Request) {
	// 		io.WriteString(w, "Hello world!")
	//}
	
	mux["/mario"] = mario
	mux["/getter"] = get
	mux["/quitserver"] = quitserver
	
	fmt.Println(".....................................OK")
	fmt.Println("")
	
	fmt.Println("Server in ascolto ai seguenti indirizzi")
	indirizziDisponibili()

	server.ListenAndServe()
	
	//listenForCommands()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "My server: "+r.URL.String())
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func mario(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello Mario!")
}

func get(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Chiamata GET!")
}

func quitserver(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Programma terminato da remoto.")
	os.Exit(0)
}

func listenForCommands() {
	fmt.Println("listenForCommands")
	reader := bufio.NewReader(os.Stdin)
	commandstring, _ := reader.ReadString('\n')
	
	fmt.Println("Comando inserito: '", commandstring, "'")
}

func indirizziDisponibili() {
	//Si dovrebbero listare tutti gli indirizzi della macchina e aggiungere la porta.
	fmt.Println("Server in ascolto all'indirizzo:")
	if listeningAddr == ":80" {
		fmt.Println("http://localhost/")
	} else {
		fmt.Println("http://localhost" + strings.TrimSpace(listeningAddr) + "/")
	}	
}

func inputparameters() {
	//programmaEParametri := os.Args
    parametri := os.Args[1:]
	
	if len(parametri) != 0 {
		//Controllare ed impostare i parametri
		fmt.Println(parametri)
	} else {
		fmt.Println("Nessun parametro di input!")
	}
}

func main() {
	inputparameters()
	welcomeMessage()
	askAndTestListeningPort()
	startServer()
	
	//dopo lo startServer non fa niente...
	//fmt.Println("listenForCommands")
	//listenForCommands()
	//fmt.Println("listenForCommands end")
}