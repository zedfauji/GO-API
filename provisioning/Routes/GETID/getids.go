package getoperation

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"../../Utilities"
	"../GenerateToken"
)

func GetId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	validToken, err := generatetoken.GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	} else {
		client := &http.Client{}

		urlStr := Utility.MakeURL(params["id"])
		req, _ := http.NewRequest("GET", urlStr, nil)
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
}
