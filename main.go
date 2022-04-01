package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/nebhale/client-go/bindings"
	"github.com/urfave/negroni"

	//"time"
	"embed"
	"io/fs"
	"io/ioutil"
)

var db *sql.DB

//go:embed web/dist/shopbasket
var webStaticContent embed.FS

func HandleGetInventory(w http.ResponseWriter, r *http.Request) {
	datastore := Datastore{db}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	inventory, err := datastore.GetInventory(id)
	fmt.Println(inventory)
	response, err := json.Marshal(inventory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func HandleCreateInventory(w http.ResponseWriter, r *http.Request) {
	datastore := Datastore{db}
	var req Inventory
	json.NewDecoder(r.Body).Decode(&req)
	inventory, err := datastore.CreateInventory(req)
	response, err := json.Marshal(inventory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func HandleListInventory(w http.ResponseWriter, r *http.Request) {
	datastore := Datastore{db}
	inventory, err := datastore.ListInventory()
	fmt.Println(inventory)
	response, err := json.Marshal(inventory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func HandleDeleteInventory(w http.ResponseWriter, r *http.Request) {
	datastore := Datastore{db}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := datastore.DeleteInventory(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"build": "1"}`))
}

func main() {
	var err error
	db, err = InitializeDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	if db != nil {
		defer db.Close()
		var greeting string

		err = db.QueryRow("select 1").Scan(&greeting)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		}

		fmt.Println(greeting)
	}

	n, router := Route()
	n.UseHandler(router)
	n.Run(":8080")

}
func InitializeDB() (*sql.DB, error) {
	var err error
	b := bindings.FromServiceBindingRoot()
	b = bindings.Filter(b, "postgresql")
	if len(b) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Incorrect number of PostgreSQL drivers: %d\n", len(b))
		os.Exit(1)
	}
	connectionString, ok := bindings.Get(b[0], "pgbouncer-uri")
	if !ok {
		_, _ = fmt.Fprintln(os.Stderr, "No pgbouncer-uri in binding")
		//os.Exit(1)
	}
	db, err = sql.Open("pgx", connectionString)
	return db, err
}
func Route() (n *negroni.Negroni, rt *mux.Router) {
	router := mux.NewRouter()

	router.HandleFunc("/status", StatusHandler).Methods("GET")
	router.HandleFunc("/api/inventory/{id}", HandleGetInventory).Methods("GET")
	router.HandleFunc("/api/inventory", HandleCreateInventory).Methods("POST")
	router.HandleFunc("/api/inventory", HandleListInventory).Methods("GET")
	router.HandleFunc("/api/inventory/{id}", HandleDeleteInventory).Methods("DELETE")
	webStaticContentRoot, _ := fs.Sub(webStaticContent, "web/dist/shopbasket")
	n = negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.FS(webStaticContentRoot)))
	return n, router
}
