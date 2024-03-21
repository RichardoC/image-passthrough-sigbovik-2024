package main

import (
	// "io/ioutil"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// // data := map[string]string{
		// // 	"Region": os.Getenv("FLY_REGION"),
		// // }

		// // t.ExecuteTemplate(w, "index.html.tmpl", data)
		// client := &http.Client{}
		// requ, err := http.NewRequest(r.Method, "http://localhost.airdns.org:32707/", r.Body)
		// if err != nil {
		// 	http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		// 	return
		// }
		// log.Printf("Original headers %+v \n", r.Header)
		// origBody, err := io.ReadAll(r.Body) // TODO, have this streaming
		// if err != nil {
		// 	http.Error(w, "Error with request", http.StatusBadRequest)
		// 	return
		// }
		// log.Printf("Original body %+v something\n", origBody)

		// requ.Header.Add("content-type", r.Header.Get("content-type"))
		// // requ.Header.Add("content-length", r.Header.Get("content-length"))
		// requ.Body = r.Body

		// res, err := client.Do(requ)
		// if err != nil {
		// 	http.Error(w, "Error with request", http.StatusInternalServerError)
		// 	return
		// }
		// defer res.Body.Close()
		// log.Printf("Got response %+v \n", res)

		// body, err := io.ReadAll(res.Body) // TODO, have this streaming
		// log.Printf("Got response body: %+v \n", body)

		// This code is largely based on https://stackoverflow.com/a/66632056/16365314

		body, err := io.ReadAll(r.Body)
		inputBodyStr := string(body)
		log.Printf("Input body bytes: %s \n", inputBodyStr)
		if err != nil {
			log.Print(err)
			return
		}

		// r.Body = io.NopCloser(bytes.NewReader(body))
		// If Server A and B are separate docker images, you may need to use their docker subnet IP, like below.
		proxyReq, err := http.NewRequest(r.Method, fmt.Sprintf("http://localhost.airdns.org:32707/%s", r.RequestURI), bytes.NewReader(body))
		if err != nil {
			log.Print(err)
			return
		}

		for header, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(header, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			log.Print(err)
			return
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return
		}

		if err != nil {
			msg := fmt.Sprintf("Error with reading body: %+v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.Write(respBody)
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
