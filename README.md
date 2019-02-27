# GO-API
Only Containerized API written in GOLANG, No DOCKERFILES 
if response == "broken" {
		fmt.Println("Its broken")
	} else {
		w.WriteHeader(http.StatusOK)
		//TODO: Return ID with Accesstoken
		w.Write([]byte(s.AccessToken))
	}


-- Spec Page

1. All the diagrams and information about the current architecutre :- same format with documentation

2. Fix the diagrams 
3. Implement user/password , for initial set up as primary credentials 
4. implement admin user which can see the secrets 
5. If we want to support multi user
6. Edge1 user credentials generation technique 
7. we want to join edge1 user with the organization 

Make assumption on everything 
user mode and local mode 

1. Current Architecutre :- Diagram about PS and secret Storage talking, generating token and storing 
2. Secret Storage :- Current arch and diag of secret storage, how its storing token and how its validating incoming jwt token from PS 
3. Datadog implementation :- 
4. Initial provisioning :- Implement Username and password for Secret Storage incoming with Edge1 User
5. Provisioning Primary Credentials during IFQA Process :- Describing what we discussed about creating an api in skycatch cloud to check edge1 user credentials 	

2 users 
1. admin permissiong level 
2. normal permission level 
Assuming edge is not connected to cloud. 

end of provisioning is to be connected with the org 
