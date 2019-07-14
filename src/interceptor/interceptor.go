/*	Step 1 : Setting Path
*	system path/twt_vidaa_8/src/assets/_req/scripts/go  --> current dir path
*	export GOPATH=system path/twt_vidaa_8/src/assets/_req/scripts/go
*	export GOBIN=$GOPATH/bin
*	export PATH=$PATH:$GOBIN
*/

/*	Step 2 : External Dependencies
*	go get "github.com/graphql-go/graphql"
*	go get "github.com/graphql-go/handler"
*	go get "golang.org/x/crypto/bcrypt"
*	go get "github.com/gorilla/sessions"
*	go get "gopkg.in/mgo.v2"
*	go get "gopkg.in/mgo.v2/bson"
* go get "github.com/SparkPost/gosparkpost"
*/

/*	Step 3 : Execute script
*	go install graphOnGo (Compiling program)
*	graphOnGo		(to run type & hit enter)
*/

/*	Step 4 : Mongodb
*	To start mongo `sudo mongod`
*	Open new tab `mongo`
*	create database	`use vidaa_db`
*	Create db user with `username : root, password : root` with `dbAdmin`
*/

/*	Step 5 : Check go Mongo connect session
*	file loc 		(twt_vidaa_8/src/assets/_req/scripts/go/src/pkgs/Mongoose/Mongoose.go)
*	function Name `GetConnected()`
*	change to app `mgo.ParseURL("mongodb://localhost:27017")`
/Users/sohailshaikh/AndroidStudioProjects/XOR/app/src/main/graphql/com.takshan.rdmuniversal.xor
*
/* Strp 6 : Install API
* `go install interceptor`
* `interceptor`
*/

package main

import (
	"encoding/json"
	"fmt"
	//"time"
	//"os"
	"net/http"
	"packages/Mongoose"
	"packages/StructConfig"
	"packages/TokenManager"
	"strings"
	"gopkg.in/mgo.v2"

)
var sess *mgo.Session
var collection *mgo.Collection
func setOrigin(w http.ResponseWriter, r *http.Request){
	//fmt.Println("In set origin")
	if origin := r.Header.Get("Origin"); origin != "" {
		//fmt.Println("Origin : ",origin)
		//fmt.Println("Auth : ",r.Header.Get("Authorization"))
		w.Header().Set("Access-Control-Allow-Origin", /*"http://localhost:7888"*/origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,UPDATE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",/*"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"*/"Content-Type, Authorization")
		w.Header().Set("Content-Type", "x-www-form-urlencoded, application/json, text/*")

		//for file uploading
		r.ParseMultipartForm(32 << 20)
		r.ParseForm()
		/*for key, values := range r.Form {   // range over map
  		for _, value := range values {    // range over []string
     		fmt.Println("Key : ", key, "Value : ",value)
  		}
		}*/
		if r.FormValue("function") != "Login" && r.FormValue("function") != "GetSponsorList" && r.FormValue("function") != "AddUserDetails" && r.FormValue("function") != "LoginAdmin"{
			if len(strings.TrimSpace(r.FormValue("token"))) > 0 {
				if tokenIs, isTokenValidErr := TokenManager.IsTokenValid(r.FormValue("token")); isTokenValidErr != nil {
					fmt.Println("Invalid token : ",isTokenValidErr)
					responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Invalid token ... Try again"}})
				}else if tokenIs {
					fmt.Println("Token is valid")
					argumentList := StructConfig.ArgumentList{Function:r.FormValue("function"),Arguments:r.FormValue("args")}
			    Mongoose.Caller(argumentList,w,r)
				}else{
					fmt.Println("Invalid token")
					responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Invalid token ... Try again"}})
				}
			}else{
				fmt.Println("Token is empty... Invalid ",r.FormValue("token"))
				responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Token is empty ... Try again"}})
			}
		}else{
			argumentList := StructConfig.ArgumentList{Function:r.FormValue("function"),Arguments:r.FormValue("args")}
	    Mongoose.Caller(argumentList,w,r)
		}
		//fmt.Fprintln(w,Mongoose.Caller(r.FormValue("function")));

	}else{
		fmt.Println("Origin empty hai")
		responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Origin is empty ... Try again"}})
	}

}

func responder(w http.ResponseWriter,result interface{}){
  //fmt.Println("Response : ",result)
  w.WriteHeader(http.StatusOK)
  buff, _ := json.Marshal(result)
  w.Write(buff)
}

func main() {
	/*ms :=time.Now().UnixNano() / (int64(time.Millisecond))
	fmt.Println(ms)*/
	fmt.Println("Environment set, now u can communicate with Rozgar")
	http.HandleFunc("/",setOrigin)
	http.ListenAndServe(":7825",/*context.ClearHandler(http.DefaultServeMux*/nil)
}

/*

db.userDetails_collection.insert({"_id":"5cc5a7850953133b24cde4e7","user_id" : "5cc5a7850953133b24cde4e7","user_name" : "mazhar","user_role" : "admin","identity_code" : "hrp,02","firebase_token" : "ct1Mg9os19I:APA91bGBBsJkxtOU_VR1qCHNt_nGF3T-R8prGZ1UM8JK--AogyWjl0x4hWi48svME3YNtiUt6mzV9EBYl2WdCzOsV661NwP8ypmYv_hkDIFiNja-p1H8mksjSXAPJKBP_7Hn45aXwOsH","sponsor_uname" : "rozgar","hrp" : 10000,"account_status" : "Active","user_added_on" : "1554289344754","personal_info" : {"full_name" : "Mazhar Shaikh","mobile_number" : "9975147441","dob" : "12/11/1990","gender" : "Male"},"bank_details" : {"account_number" : "1234567890","ifs_code" : "ABCD0123456"},"transaction_history" : [],"direct_child_count" : 0})



db.userDetails_collection.insert({"_id":"5cc5a7260953133b24cde4e6","user_id" : "5c90b95eecd0ebeadd424fd6","user_name" : "sallu","user_role" : "admin","identity_code" : "hrp,01","firebase_token" : "fiJkCqlHvCQ:APA91bFlG2kB2WxD4PKukfqUC86ke0bDAou6Q8cerJmdMh6BgoKbHYhIRiIQU9zaJvXjmzsDLd_J42NzljjTWvyHU5LXSaZ_dY0OHoNMMPsbVDN5ZiW0dM4fmY2SUsNx2oHfKcOcacpV","sponsor_uname" : "rozgar","hrp" : 10000,"account_status" : "Active","user_added_on" : "1554289344754","personal_info" : {"full_name" : "Salahuddin Shaikh","mobile_number" : "8888687577","dob" : "12/11/1990","gender" : "Male"},"bank_details" : {"account_number" : "1234567890","ifs_code" : "ABCD0123456"},"transaction_history" : [],"direct_child_count" : 0})

*/
