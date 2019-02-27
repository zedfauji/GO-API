package readoperation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type ResponseSecret struct {
	ID     string `json:"id"`
	Token  string `json:"token"`
	Active bool   `json:"Active"`
}

func ReadID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	m := map[string]string{}
	var ResponseSecret1 []ResponseSecret
	fileName := fmt.Sprintf("%s.json", id)
	fileLoc := filepath.Join("/tmp/data/secrets", fileName)
	jsonFile, err := os.Open(fileLoc)
	if err != nil {
		status := fmt.Sprintf("%s doesn't exist", id)
		m["status"] = status
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Opened ", fileLoc)
		byteValue, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal([]byte(byteValue), &ResponseSecret1)
		fmt.Println(err)
		if err != nil {
			m["status"] = "Error While reading file "
		} else {
			m["token"] = ResponseSecret1[0].Token
		}
		//fmt.Printf("%#v", ResponseSecret1)
	}

	_ = json.NewEncoder(w).Encode(m)

}
