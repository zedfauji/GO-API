# GO-API
Only Containerized API written in GOLANG, No DOCKERFILES 
if response == "broken" {
		fmt.Println("Its broken")
	} else {
		w.WriteHeader(http.StatusOK)
		//TODO: Return ID with Accesstoken
		w.Write([]byte(s.AccessToken))
	}
