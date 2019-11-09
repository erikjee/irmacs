package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mdp/qrterminal"
	irma "github.com/privacybydesign/irmago"
	"github.com/privacybydesign/irmago/server"
	"github.com/privacybydesign/irmago/server/irmaserver"
)

func printQr(qr *irma.Qr, noqr bool) error {
	qrBts, err := json.Marshal(qr)
	if err != nil {
		return err
	}
	if noqr {
		fmt.Println(string(qrBts))
	} else {
		qrterminal.GenerateWithConfig(string(qrBts), qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout,
			BlackChar: qrterminal.BLACK,
			WhiteChar: qrterminal.WHITE,
		})
	}
	return nil
}

func main() {
	configuration := &server.Configuration{
		// Replace with address that IRMA apps can reach
		URL: "http://172.20.10.4:1234/irma",
	}

	err := irmaserver.Initialize(configuration)
	if err != nil {
		// ...
	}

	http.Handle("/irma/", irmaserver.HandlerFunc())
	http.HandleFunc("/createrequest", createFullnameRequest)

	// Start the server
	fmt.Println("Going to listen on :1234")
	err = http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println("Failed to listen on :1234")
	}
}

func createFullnameRequest(w http.ResponseWriter, r *http.Request) {
	request := `{
        "type": "disclosing",
        "content": [{ "label": "Full name", "attributes": [ "pbdf.pbdf.irmatube.type" ]}]
    }`

	sessionPointer, token, err := irmaserver.StartSession(request, func(r *server.SessionResult) {
		fmt.Println("Session done, result: ", server.ToJson(r))
	})
	if err != nil {
		// ...
	}

	fmt.Println("Created session with token ", token)

	// Print QR in terminal for testing purposes
	printQr(sessionPointer, false)

	// Send session pointer to frontend, which can render it as a QR
	w.Header().Add("Content-Type", "text/json")

	jsonSessionPointer, _ := json.Marshal(sessionPointer)
	w.Write(jsonSessionPointer)
}
