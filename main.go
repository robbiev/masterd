package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/info/refs", func(ow http.ResponseWriter, r *http.Request) {
		fmt.Println("incoming")
		service := r.URL.Query()["service"][0]
		fmt.Printf("/info/refs?service=%s\n", service)
		ow.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")

		buf := bytes.NewBuffer(nil)
		w := io.MultiWriter(ow, buf)

		writePacketLine(w, []byte("# service=git-upload-pack"))
		writePacketEnd(w)
		caps := "multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done symref=HEAD:refs/heads/master agent=git/2:2.1.1+github-607-gfba4028"
		writePacketLine(w, []byte(string(append([]byte("fa2fa601958e6ac7f9ccbb9a7052e4eed33eefe0 HEAD"), 0))+caps))
		writePacketLine(w, []byte("fa2fa601958e6ac7f9ccbb9a7052e4eed33eefe0 refs/heads/master"))
		writePacketEnd(w)

		fmt.Println(buf)
		//w.WriteHeader(200)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("unknown location")
		fmt.Println(r.URL)
		w.WriteHeader(404)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func writePacketLine(w io.Writer, bytes []byte) {
	fmt.Fprintf(w, "%04x", len(bytes)+5)
	w.Write(bytes)
	w.Write([]byte("\n"))
}
func writePacketEnd(w io.Writer) {
	w.Write([]byte("0000"))
}
