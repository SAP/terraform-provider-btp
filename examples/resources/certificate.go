// This file contains a function to set the tf variable "certificate", to be used in the API credential resources
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)


func main() {
	err := tfutils.GenerateCertificate()

	if err!=nil{
		fmt.Printf("Error generating a certificate : %s", err)
		return 
	}

	data, err := os.ReadFile("cert.pem")
	if err != nil {
		fmt.Printf("Error reading the certificate : %s", err)
		return 
	}

	pemString := string(data)

	output := map[string]string{
        "certificate": pemString,
	}

	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		fmt.Printf("Error encoding JSON : %s", err)
		return
	}

	err = os.Remove("cert.pem")
	if err != nil {
		fmt.Println("Unable to delete PEM file")
		return
	}

}

