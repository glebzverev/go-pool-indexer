package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/glebzverev/go-pool-indexer/arb"
)

type Request struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func New(ARB *arb.Arb) {
	go func() {
		serverPort := 3000
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			fmt.Fprintf(w, "Hello from router\n")
			tokens := make([]string, 0)
			for symbol := range arb.TokenAddresses {
				tokens = append(tokens, symbol)
			}
			fmt.Fprintf(w, "Tokens available for trade: %+v\n", tokens)
			fmt.Fprintf(w, "Data config {from:<token>, to:<token>, amount:<float>}\n")
		})
		mux.HandleFunc("/way/", func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			data := new(Request)
			fmt.Println(string(body))
			json.Unmarshal([]byte(body), data)
			if data == nil {
				fmt.Fprintf(w, "Invalid request")
				return
			}
			from, ok := arb.TokenAddresses[data.From]
			if !ok {
				fmt.Fprintf(w, "Uknown token")
				return
			}
			to, ok := arb.TokenAddresses[data.To]
			if !ok {
				fmt.Fprintf(w, "Uknown token")
				return
			}
			way, err := ARB.FindOptimal(from, to, data.Amount)
			fmt.Fprintf(w, "%+v\n", *way)
		})

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", serverPort),
			Handler: mux,
		}
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("error running http server: %s\n", err)
			}
		}
	}()
}
