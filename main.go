package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/eyedeekay/onramp"
)

func help() string {
	formattedString := fmt.Sprintf("%s\n", "  -file string")
	formattedString += fmt.Sprintf("%s\n", "        Serve a single file after loading it from the disk")
	formattedString += fmt.Sprintf("%s\n", "  -help")
	formattedString += fmt.Sprintf("%s\n", "        Show help/usage information")
	formattedString += fmt.Sprintf("%s\n", "  -string string")
	formattedString += fmt.Sprintf("%s", "        Serve the string as a file over I2P")
	return formattedString
}

func usage() string {
	formattedString := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		"# SingleFile",
		"",
		"Serve one file over an I2P service, automatically, on all paths",
		"",
		"```sh",
		help(),
		"```",
		"",
	)
	return formattedString
}

func printUsage() {
	formattedString := usage()
	fmt.Println(formattedString)
}

func main() {
	serveString := flag.String("string", "", "Serve the string as a file over I2P")
	serveFile := flag.String("file", "", "Serve a single file after loading it from the disk")
	help := flag.Bool("help", false, "Show help/usage information")
	flag.Parse()
	if *help {
		printUsage()
		os.Exit(0)
	}
	if *serveFile != "" && *serveString != "" {
		printUsage()
		log.Fatal("Serve only a string or a file, not both")
	}
	if *serveFile == "" && *serveString == "" {
		*serveString = usage()
	}
	garlic, err := onramp.NewGarlic("singleFile", "127.0.0.1:7656", []string{})
	if err != nil {
		panic(err)
	}
	ln, err := garlic.Listen()
	if err != nil {
		panic(err)
	}
	if *serveString != "" {
		bytes := []byte(*serveString)
		sb := &serveBytes{
			bytes: bytes,
		}
		if err := http.Serve(ln, sb); err != nil {
			panic(err)
		}
	} else {
		if _, err := os.Stat(*serveFile); err == nil {
			if bytes, err := ioutil.ReadFile(*serveFile); err != nil {
				sb := &serveBytes{
					bytes: bytes,
				}
				if err := http.Serve(ln, sb); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}

		}
	}

}

type serveBytes struct {
	bytes []byte
}

func (s *serveBytes) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	rw.Write(s.bytes)
}
