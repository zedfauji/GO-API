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

	fileName := fmt.Sprintf("%s.json", id)
	fileLoc := filepath.Join("/tmp/data/secrets", fileName)
	jsonFile, err := os.Open(fileLoc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened ", fileLoc)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var ResponseSecret1 []ResponseSecret
	//fmt.Println(string(byteValue))
	v := json.Unmarshal([]byte(byteValue), &ResponseSecret1)
	fmt.Println(v)
	//fmt.Printf("%#v", ResponseSecret1)
	fmt.Println("This is me", ResponseSecret1[0])
	PrintObject(ResponseSecret1[0])
}

func PrintObject(r ResponseSecret) {
	fmt.Printf("ID is :%s, Token is : %s ,Status : %t", r.ID, r.Token, r.Active)
}
