package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Person struct {
	Name string
	Age  int
}
type Package struct {
	Person        *Person
	PersonFriends []*Person
}

//Формирование пакета данных
func NewPackage() Package {

	rand.Seed(time.Now().UTC().UnixNano())
	var rndAge int
	rndAge = rand.Intn(100)

	pkg := Package{
		Person: &Person{
			Name: "Mikhail", Age: 23},
	}

	pkg.PersonFriends = append(pkg.PersonFriends, &Person{Name: "Andrey", Age: 32}, &Person{Name: "Ivan", Age: rndAge})

	return pkg
}

//Обработчик запроса от клиента
func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	data, err := json.Marshal(NewPackage())
	if err != nil {
		log.Fatal("json.Marshal: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(data))
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
