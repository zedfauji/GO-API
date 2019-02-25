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

type ResponseSecret []struct {
	ID     string `json:"id,omitempty"`
	Token  string `json:"token,omitempty"`
	Active bool
}

func ReadID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	fileName := fmt.Sprintf("%s.json", id)
	fileLoc := filepath.Join("/tmp/data/secrets", fileName)
	jsonFile, err := os.Open(fileLoc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened ", fileLoc)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var ResponseSecret1 ResponseSecret
	json.Unmarshal(byteValue, &ResponseSecret1)
	fmt.Printf("%#v", ResponseSecret1)
	//fmt.Println("This is me", ResponseSecret1.ID)
}
