package main

import ("log"
"os"
"fmt"
 "net/http"
 "github.com/gorilla/mux"
 "github.com/100yo/go-sofia/internal/diagnostics"
 )

type serverConf struct {
		port string
		router http.handler
		name string
}
	
func main() {

	log.Print("Hello, World")

	blPort := os.Getenv("PORT")
	if len(blPort) == 0 {
		log.Fatal("The application port should be set")
	}	

	diagPort := os.Getenv("DIAG_PORT")
	if len(diagPort) == 0 {
		log.Fatal("The diagnostics port should be set")
	}	

	router := mux.NewRouter()
	router.HandleFunc("/", hello)

	possibleErrors := make(chan error, 2)

	configurations := []serverConf {
		{
			port:blPort,
			router: router,
			name: "application server",
		},
		{
			port: diagPort,
			router diagnostics,
			name: "diagnostics server"
		}
	}

	servers := make(http.Server, 2)

	for _, c := range configurations {
		go func (conf serverConf) {
			server := &http.Server{
				Addr: ":" + conf.port,
				handler: conf.router,
			}
			err := server.ListenAndServe()
			// server.Shutdown()
			if err != nil {
				possibleErrors <- err	
			}
			
		}(c)
	}

	

	go func() {
		diagnostics := diagnostics.NewDiagnostics()
		server := &http.Server{
				Addr: ":" + diagPort,
				handler: diagnostics,
			}
		err := server.ListenAndServe()
		if err != nil {
			possibleErrors <- err	
		}
	}()

	select {
	case err := <-possibleErrors:
		log.Fatal(err)
	}
}
 func hello(w http.ResponseWriter, r *http.Request) {
 	log.Print("The hello handler was called")
 	fmt.Fprint(w, http.StatusText(http.StatusOK))
 }