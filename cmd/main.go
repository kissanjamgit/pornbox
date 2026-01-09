package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kissanjamgit/pornbox"
)

func main() {
	ID := flag.String("id", "", "Content ID")
	flag.Parse()
	if *ID == "" {
		fmt.Printf("ERROR: `*ID == ''` e: %v", *ID)
		os.Exit(1)
	}
	reg := regexp.MustCompile(` +`)
	str := reg.ReplaceAllString(*ID, " ")
	list := strings.Split(str, " ")
	if len(list) == 1 {
		pb := pornbox.New(list[0])
		cr, err := pb.Video()
		if err != nil {
			fmt.Printf("ERROR: `err != nil` e: %v", err)
			os.Exit(1)
		}
		fmt.Printf("#EXTM3U\n#EXTINF:-1,%s\n%s\n", cr.Name, cr.URL)
		return
	}
	res, err := pornbox.Queue(list)
	code := 0
	if len(err) == 0 || len(res) == 0 {
		fmt.Printf("ERROR: `len(err) == 0  || len(res) == 0` e: %v %v", err, res)
		code = 1
	}
	var buff strings.Builder
	buff.WriteString("#EXTM3U\n")
	for _, cr := range res {
		fmt.Fprintf(&buff, "#EXTINF:-1,%s\n%s\n", cr.Name, cr.URL)
	}
	fmt.Print(buff.String())

	os.Exit(code)
}
