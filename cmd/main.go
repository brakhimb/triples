package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"tripleS/pkg"
	"tripleS/pkg/handler"
	"tripleS/pkg/utils"
)

func main() {
	help := flag.Bool("help", false, "help")
	port := flag.String("port", "8000", "port")
	directory := flag.String("directory", "data", "directory")

	flag.Parse()
	// if _, err := os.Stat(*directory); os.IsNotExist(err) {
	// 	log.Fatalf("Directory does not exist: %s", *directory)
	// }
	if *help {
		fmt.Print("Simple Storage Service.\n\n" +
			"**Usage:**\n" +
			"    triple-s [-port <N>] [-dir <S>]  \n" +
			"    triple-s --help\n\n" +
			"**Options:**\n" +
			"- --help     Show this screen.\n" +
			"- --port N   Port number\n" +
			"- --dir S    Path to the directory")
	}

	if *directory == "cmd" || *directory == "pkg" {
		log.Fatal("Conflict with folder name")
		return
	}

	if _, err := os.Stat(*directory); os.IsNotExist(err) {
		err := os.MkdirAll(*directory, 0o775)
		if err != nil {
			log.Fatal(err)
			return
		}
		path := *directory + "/buckets.csv"
		err = utils.HandlerCSV(path, "PUT buckets.csv", "", "", "", "")
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			handler.CreatBucketHandler(w, r, directory)
		case "GET":
			handler.ListBucketHandler(w, r, directory)
		case "DELETE":
			handler.DeleteBucketHandler(w, r, directory)
		default:
			utils.HandlerXML(w, pkg.Response{Status: http.StatusMethodNotAllowed, Message: "Method Not Allowed"})
		}
	})

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalf("ERROR with server %s", err)
	}
}
