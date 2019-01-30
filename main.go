package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
)

type ViewData struct{
	Title string
	Anagrams []string
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	m := getAnagrams(id)
	var anagrams []string
	for _, a := range m {
		if len(a) > 1 {
			for _, v := range a {
				anagrams = append(anagrams, string(v))
			}
		}
	}

	data := ViewData{
		Title : id,
		Anagrams : anagrams,
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data)
}

func getAnagrams(arg string) map[string][][]byte  {
	r, err := http.Get("https://raw.githubusercontent.com/first20hours/google-10000-english/master/20k.txt")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var bs byteSlice
	m := make(map[string][][]byte)

	barg := byteSlice([]byte(arg))

	sort.Sort(barg)

	for _, word := range bytes.Fields(b) {
		if len(barg) == len(word) {
			bs = append(bs[:0], byteSlice(word)...)
			sort.Sort(bs)

			if bytes.Equal(bs, barg) {
				k := string(bs)
				a := append(m[k], word)
				m[k] = a
			}
		}
	}

	return m
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{id:[a-z]+}", productsHandler)
	http.Handle("/",router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}

type byteSlice []byte

func (b byteSlice) Len() int           { return len(b) }
func (b byteSlice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byteSlice) Less(i, j int) bool { return b[i] < b[j] }