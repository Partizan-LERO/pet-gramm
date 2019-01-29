package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		return
	}

	arg := os.Args[1]

	r, err := http.Get("https://raw.githubusercontent.com/first20hours/google-10000-english/master/20k.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
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

	for _, a := range m {
		if len(a) > 1 {
			fmt.Printf("%s\n", a)
		}
	}
}

type byteSlice []byte

func (b byteSlice) Len() int           { return len(b) }
func (b byteSlice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byteSlice) Less(i, j int) bool { return b[i] < b[j] }