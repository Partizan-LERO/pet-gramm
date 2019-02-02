package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
)

type ViewData struct{
	Title string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := ViewData{
		Title : "Anagram generator",
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data)
}

func anagramHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["word"]
	m := getAnagrams(id)

	var anagrams[]string

	for _, a := range m {
		if len(a) > 1 {
			for _, v := range a {
				anagrams = append(anagrams, string(v))
			}
		}
	}

	js, _ := json.Marshal(anagrams)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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

func errorHandler(w http.ResponseWriter, r *http.Request)  {
	js, _ := json.Marshal(nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(js)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/word/{word:[a-z]+}", anagramHandler)
	router.HandleFunc("/word/{word:[1-9]+}", errorHandler)
	router.HandleFunc("/word/}", errorHandler)
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	http.Handle("/",router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}

type byteSlice []byte

func (b byteSlice) Len() int           { return len(b) }
func (b byteSlice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byteSlice) Less(i, j int) bool { return b[i] < b[j] }