package getoperation

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"../GenerateToken"
)

func GetIds(w http.ResponseWriter, r *http.Request) {
	validToken, err := generatetoken.GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8000/deviceSecret", nil)
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(body))
}
