package Mongoose

import(
  "fmt"
  "net/http"
  "time"
  "strconv"
  "strings"
  "gopkg.in/mgo.v2"
  "encoding/json"
  "packages/StructConfig"
  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2/txn"
  "packages/Underdog"
  "packages/TokenManager"
  "reflect"
  "math"
  "github.com/maddevsio/fcm"
  //"errors"
  // "regexp"
  //"log"
  //jwt "github.com/dgrijalva/jwt-go"
)
const (
	 serverKey = "AAAATDP8Kjk:APA91bE3rFM1XQGEwz381Grx5SxcZqIDhgK46lvUpGtyO56siPOGUl03otvVCDXzNka5VB1MpWHpX2n87sdsYKElAxvxwHyY8i5Y7uS5KEiH0Gu4FzLUegHcwqFeuhMj1VHhqOvWz7nc"
)
var sess *mgo.Session
var collection *mgo.Collection
var sessUValName string
// time.Now().UnixNano() / (int64(time.Millisecond)).string():=time.Now().UnixNano() / (int64(time.Millisecond)).string()
var funcMap = map[string]interface{}{
	"Hello": Hello,
  "Login": Login,
  "LoginAdmin": LoginAdmin,
  "AddUserDetails": AddUserDetails,
  "GetSponsorList": GetSponsorList,
  "DisplayAllFollowers": DisplayAllFollowers,
  "UpdateBankDetails": UpdateBankDetails,
  "GetFollowerListByLevel": GetFollowerListByLevel,
  "GetPendingUserListBySponsor": GetPendingUserListBySponsor,
  "TopupAccount": TopupAccount,
  "GetUserListToTransfer": GetUserListToTransfer,
  "TransferHRP": TransferHRP,
  "GetTransactionHistory": GetTransactionHistory,
  "UpdateUserProfile": UpdateUserProfile,
  "GetBankDetails": GetBankDetails,
  "WithdrawalRequest": WithdrawalRequest,
  "DonateToCompany": DonateToCompany,
  "RefreshFeed": RefreshFeed,
  "Logout": Logout,
  "SetFeedback": SetFeedback,
  //Admin related
  "UsersList": UsersList,
  "WithdrawalRequestList": WithdrawalRequestList,
  "GetSpecificUserDetails": GetSpecificUserDetails,
  "GenerateHRP": GenerateHRP,
  "ChangeUserPassword": ChangeUserPassword,
  "UpdateUserStatus": UpdateUserStatus,
  "ValidateUser": ValidateUser,
  "TransferHRPAdmin": TransferHRPAdmin,
  "GetAdminList": GetAdminList,
  "DonateToCompanyAdmin": DonateToCompanyAdmin,
  "ActionOnWithdrawalRequest": ActionOnWithdrawalRequest,
  "RequestHRP": RequestHRP,
  "GetAdminActivityLogs": GetAdminActivityLogs,
  "SetBroadcast": SetBroadcast,
  "UnsetBroadcast": UnsetBroadcast,
  "RefreshAdminFeed": RefreshAdminFeed,
  "GetFeedback": GetFeedback,
  "GetCompanyEarningRecords": GetCompanyEarningRecords,
  "GetReports": GetReports,
  "GetNewUserReports": GetNewUserReports,
}

func Caller(argumentList StructConfig.ArgumentList,w http.ResponseWriter, r *http.Request){
  fmt.Println("\n====================================================== "+argumentList.Function+" :- ",time.Now().Format("Jan _2 15:04:05"))
  funcMap[argumentList.Function].(func(http.ResponseWriter,*http.Request,string))(w,r,argumentList.Arguments)
}

func responder(w http.ResponseWriter,result interface{}){
  //fmt.Println("Response : ",result)
  w.WriteHeader(http.StatusOK)
  buff, _ := json.Marshal(result)
  w.Write(buff)
}

func setCollection(dbName string, collectionName string) *mgo.Collection {
	if sess == nil {
		fmt.Println("Not connected... Connecting to Mongo")
		sess = GetConnected()
	}
	collection = sess.DB(dbName).C(collectionName)
	return collection
}

func GetConnected() *mgo.Session {
	dialInfo, err := mgo.ParseURL("mongodb://127.0.0.1:27017")
	dialInfo.Direct = true
	dialInfo.FailFast = true
	dialInfo.Database = "rozgar_db"
	dialInfo.Username = "root"
	dialInfo.Password = "tiger"
	sess, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect to mongo, go error %v\n", err)
		panic(err)
	} else {
		return sess
		defer sess.Close()
	}
	return sess
}


func Hello(w http.ResponseWriter, r *http.Request,interfaceName string){
  collection = setCollection("rozgar_db", "table_collection")
  fmt.Println(collection)
  fmt.Println("In hello")
  icv := collection.Find(bson.M{"table_status":"available"})
  fmt.Println(icv)
}

//  Employee Managment staarted
func AddUserDetails(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    if isUserExist,isUserExistErr := isUserValid(credMap["userName"].(string),credMap["mobileNumber"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if !isUserExist {
      fmt.Println("User already exist ...")
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"User already exist... try something else"}})
    }else{
      if sponsorId,generatedIdCode,generatedIdCodeErr := generateIdentityCode(credMap["sponsorUName"].(string)); generatedIdCodeErr != nil {
        fmt.Println("Error while generateIdentityCode : ",generatedIdCodeErr)
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        userDetailsStruct := StructConfig.UserDetails{}
        personalInfoStruct := StructConfig.PersonalInfo{FullName:credMap["fullName"].(string),/*Address:credMap["address"].(string),*/MobileNumber:credMap["mobileNumber"].(string),Gender:credMap["gender"].(string),DOB:credMap["dob"].(string)}
        bankDetailsStruct := StructConfig.BankDetails{AccountNumber:credMap["accountNumber"].(string),/*BankName:credMap["bankName"].(string),*/IFSCode:credMap["ifsCode"].(string)/*,BankAddress:credMap["bankAddress"].(string)*/}
        userId := bson.ObjectId(bson.NewObjectId()).Hex()
        //userDetailsStruct = StructConfig.UserDetails{UserId:userId,Username:credMap["userName"].(string),IdentityCode:credMap["sponsorIdentityCode"].(string)+","+strconv.Itoa(int(generatedIdCode)),FirebaseToken:"",SponsorUname:credMap["sponsorUName"].(string),HRP:0.00,AccountStatus:"Pending",UserAddedOn:strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)),10),PersonalInfo:personalInfoStruct,BankDetails:bankDetailsStruct}
        userDetailsStruct = StructConfig.UserDetails{UserId:userId,Username:credMap["userName"].(string),IdentityCode:credMap["sponsorIdentityCode"].(string)+","+strconv.Itoa(int(generatedIdCode)),FirebaseToken:"",SponsorUname:credMap["sponsorUName"].(string),HRP:0.00,AccountStatus:"Pending",UserAddedOn:time.Now().UnixNano()/(int64(time.Millisecond)),PersonalInfo:personalInfoStruct,BankDetails:bankDetailsStruct}
        change := bson.M{"$inc": bson.M{"direct_child_count": 1}}
        runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))

        ops := []txn.Op{{
            //Adding new user instance
            C:      "userInstance_collection",
            Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
            //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
            Insert: &StructConfig.UserInstance{Username:credMap["userName"].(string),Password:credMap["password"].(string),MobileNo:credMap["mobileNumber"].(string)},
          },
          {
            //Adding new user details
            C:      "userDetails_collection",
            Id:     userId,
            //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
            Insert: userDetailsStruct,
          },
          {
            //incrementing count of direct_child_count of sponsor
            C:      "userDetails_collection",
            Id:     sponsorId,
            Update: change,
          },
        }
        id := bson.NewObjectId() // Optional
        runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"AddUserDetails"))
        if runnerErr != nil {
          fmt.Println("Error while runner : ",runnerErr)
          if resumeErr := runner.Resume(id); resumeErr != nil {
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
          }else{

            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }else{
          fmt.Println("Executed")
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
        }
      }

    }
  }
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    if isUserExist,isUserExistErr := isUserExist(credMap["adminUserName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      personalInfoStruct := StructConfig.PersonalInfo{FullName:credMap["fullName"].(string),/*Address:credMap["address"].(string),*/MobileNumber:credMap["mobileNumber"].(string),Gender:credMap["gender"].(string),DOB:credMap["dob"].(string)}
      actDetailsStruct := StructConfig.ActivityDetails{To:"",From:"",Amount:0.0}
      adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["adminUserName"].(string),ActivityFor:"Change Password",ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:credMap["adminUserName"].(string),ActivityDetails:actDetailsStruct}
      runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
      ops := []txn.Op{{
          //Adding new user instance
          C:      "userDetails_collection",
          Id:     credMap["userId"].(string),
          //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
          Update: bson.M{"$set":bson.M{"personal_info":personalInfoStruct}},
        },
        {
          C:      "adminActivityLog_collection",
          Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
          //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
          Insert: adminLogStruct,
        },
      }
      id := bson.NewObjectId() // Optional
      runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"UpdateUserProfile"))
      if runnerErr != nil {
        fmt.Println("Error while runner : ",runnerErr)
        if resumeErr := runner.Resume(id); resumeErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{

          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
        }
      }else{
        fmt.Println("Executed")
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
      }
    }else{
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Incorrect Password."}})
    }
  }
}

func generateIdentityCode(userName string)(string,float64,error){
  collection = setCollection("rozgar_db", "userDetails_collection")
  m := bson.M{}
  err := collection.Find(bson.M{"user_name":userName}).Select(bson.M{"_id":1,"direct_child_count":1}).One(&m)
  if err != nil {
    fmt.Println("Error while getting direct child count : ",err)
    return "",0,err
  }else{
    fmt.Println("gen id count : ",m["direct_child_count"])
    floatVal,errConv := getFloat(m["direct_child_count"])
    if errConv != nil {
      fmt.Println("Error while getting direct child count converting float : ",errConv)
      return "",0,errConv
    }else{
      return m["_id"].(string),floatVal + 1,nil
    }
  }
}

func getFloat(unk interface{}) (float64, error) {
  var floatType = reflect.TypeOf(float64(0))
    v := reflect.ValueOf(unk)
    v = reflect.Indirect(v)
    if !v.Type().ConvertibleTo(floatType) {
        return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
    }
    fv := v.Convert(floatType)
    return fv.Float(), nil
}

func setInfoStruct(TransBy,TransFor string)StructConfig.InfoStruct{
  return StructConfig.InfoStruct{TransBy:TransBy,TransFor:TransFor,TransOn:time.Now().UnixNano()/(int64(time.Millisecond))}
}

func UpdateBankDetails(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      bankDetailsStruct := StructConfig.BankDetails{AccountNumber:credMap["accountNumber"].(string),/*BankName:credMap["bankName"].(string),*/IFSCode:credMap["ifsCode"].(string)/*,BankAddress:credMap["bankAddress"].(string)*/}
      ops := []txn.Op{}
      updt := bson.M{"$inc": bson.M{"hrp": -50.0},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:"company",From:credMap["userName"].( string),Units:-50.0,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Bank Details Update",BankTransId:"",Status:"Done",Level:""}}}
      updtBank := bson.M{"$inc": bson.M{"hrp": 50.0},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:"company",From:credMap["userName"].( string),Units:50.0,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Bank Details Update",BankTransId:"",Status:"Done",Level:""}}}
      query := bson.M{"$set":bson.M{"bank_details":bankDetailsStruct}}
      runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
      //if user_name != "" && level != "" && id != "" {
      if credMap["extraCharge"].(string) == "true" {
        if companyId,comIdErr := getCompanyId(); comIdErr != nil {
          fmt.Println("Error while getting companyId : ",comIdErr)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          fmt.Println("Extra charge applied")
          /*ops = append(ops,txn.Op{
            C:      "userDetails_collection",
            Id:     credMap["userId"].(string),
            Assert: bson.M{"hrp": bson.M{"$gte": 50.0}},
            Update: updt,
          })*/
          ops = []txn.Op{{
              //Adding new user instance
              C:      "userDetails_collection",
              Id:     credMap["userId"].(string),
              Assert: bson.M{"hrp": bson.M{"$gte": 50.0}},
              Update: updt,
            },
            {
                //Adding new user instance
                C:      "userDetails_collection",
                Id:     companyId,
                //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
                Update: updtBank,
            },
          }
        }
      }

      ops = append(ops,txn.Op{
        C:      "userDetails_collection",
        Id:     credMap["userId"].(string),
        //Assert: bson.M{"hrp": bson.M{"$gte": totalAmt}},
        Update: query,
      })
      id := bson.NewObjectId() // Optional
      runnerErr := runner.Run(ops, id, nil)
      if runnerErr != nil {
        if resumeErr := runner.Resume(id); resumeErr != nil {
          fmt.Println("Error while runner @ UpdateBankDetails : ",runnerErr)
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }
        //fmt.Println("Error while runner : ",runnerErr)
        //responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        fmt.Println("Bank details updated")
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
      }
    }else{
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Incorrect Password."}})
    }

  }

}

func WithdrawalRequest(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      if isValidForWithdrawal(credMap["userName"].(string)) {
        if companyId,comIdErr := getCompanyId(); comIdErr != nil {
          fmt.Println("Error while getting companyId : ",comIdErr)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          //deduct respective amount and create a transacation
          //Add respective amount to company acc
          withdrawalAmt,floatErr := strconv.ParseFloat(credMap["withdrawalAmt"].(string),64)
          if floatErr != nil {
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
          }else{
            transId := bson.ObjectId(bson.NewObjectId()).Hex()
            updt := bson.M{"$inc": bson.M{"hrp": -withdrawalAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:transId,To:"company",From:credMap["userName"].( string),Units:-withdrawalAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Withdrawal",BankTransId:"",Status:"Processing",Level:""}}}
            updtCmpny := bson.M{"$inc": bson.M{"hrp": withdrawalAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:transId,To:"company",From:credMap["userName"].( string),Units:withdrawalAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Withdrawal",BankTransId:"",Status:"Done",Level:""}}}
            runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
            ops := []txn.Op{{
                //Adding new user instance
                C:      "userDetails_collection",
                Id:     credMap["userId"].(string),
                Assert: bson.M{"hrp": bson.M{"$gte": withdrawalAmt}},
                Update: updt,
              },
              {
                  //Adding new user instance
                  C:      "userDetails_collection",
                  Id:     companyId,
                  //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
                  Update: updtCmpny,
              },
            }
            id := bson.NewObjectId() // Optional
            runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"UpdateUserProfile"))
            if runnerErr != nil {
              fmt.Println("Error while runner : ",runnerErr)
              if resumeErr := runner.Resume(id); resumeErr != nil {
                  responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
              }else{
                data := make(map[string]string)
                data["title"] = "Happy Rozgar"
                data["body"] = "Withdrawal request from "+credMap["userName"].(string)
                _,fbTokens:= getAdminsFBToken();
                //var ids = []string{fbTokens}
                sendNotificationWithFCM(fbTokens,data)
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
              }
            }else{
              fmt.Println("Executed")
              data := make(map[string]string)
              data["title"] = "Happy Rozgar"
              data["body"] = "Withdrawal request from "+credMap["userName"].(string)
              _,fbTokens:= getAdminsFBToken();
              //var ids = []string{fbTokens}
              sendNotificationWithFCM(fbTokens,data)
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
            }
          }
        }
      }else{
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Withdrawal can be done once in a week on Sunday only."}})
      }

    }else{
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"User does not exist "}})
    }
  }
}

func isValidForWithdrawal(userName string)bool{
  collection = setCollection("rozgar_db","userDetails_collection")
  if int(time.Now().Weekday()) == 0 {
    m := []bson.M{}
    yesterday := time.Now().AddDate(0,0,-1).UnixNano()/(int64(time.Millisecond))
    //fmt.Println("Yesterday : ",yesterday)
    //pipeErr := collection.Pipe([]bson.M{bson.M{"$match":bson.M{"$and":[]bson.M{bson.M{"transaction_history.transaction_for":"Withdrawal"},bson.M{"user_name":userName}}}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$gt":[]interface{}{"$$t.transaction_on",yesterday}}}}}}}).All(&m)
    pipeErr := collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":userName}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$gte":[]interface{}{"$$t.transaction_on",yesterday}}}}}}}}}).All(&m)
    if pipeErr != nil {
      fmt.Println("Error while isValidForWithdrawal : ",pipeErr)
      fmt.Println("Invalid for Withdrawal err ")
      return false
    }else{
      //fmt.Println("isValid Withdrawal m : ",m)
      if len(m) > 0 {
        if retArr,_:= Underdog.InterfaceArrToMap(m[0]["transaction_history"]); len(retArr) > 0 {
          fmt.Println("Invalid for Withdrawal already requested ")
          return false
        }else{
          fmt.Println("valid for Withdrawal exist")
          return true
        }
      }else {
        fmt.Println("Invalid for Withdrawal not exist")
        return true
      }
    }
  }else{
    fmt.Println("Today is : ",time.Now().Weekday())
    return false
  }
}

func ActionOnWithdrawalRequest(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    /*withdrawalAmt,floatErr := strconv.ParseFloat(credMap["withdrawalAmt"].(string),64)
    if floatErr != nil {
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else{
      fmt.Println("credMap : ",credMap)
      updt := bson.M{"$set":bson.M{"transaction_history.$.status":credMap["status"].(string),"transaction_history.$.bank_trans_id":credMap["paidId"].(string)}}
      //updtCmpny := bson.M{"$inc": bson.M{"hrp": withdrawalAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:"company",From:credMap["userName"].( string),Units:withdrawalAmt,TransactionOn:strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)),10),TransactionFor:"Withdrawal",BankTransId:"",Status:"Done",Level:""}}}
      actDetailsStruct := StructConfig.ActivityDetails{To:"",From:"",Amount:0.0}
      adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["adminUserName"].(string),ActivityFor:"Withdrawal request set to "+credMap["status"].(string),ActivityOn:strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)),10),ActivityStatus:"Done",ActivityPerformedOn:credMap["userName"].(string),ActivityDetails:actDetailsStruct}
      fmt.Println("Updt query : ",updt)
      runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
      ops := []txn.Op{{
          //Adding new user instance
          C:      "userDetails_collection",
          Id:     credMap["userId"].(string),
          Assert: bson.M{"user_id":credMap["userId"].(string),"transaction_history.trans_id": credMap["transactionId"].(string)},
          Update: updt,
        },
        {
          C:      "adminActivityLog_collection",
          Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
          //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
          Insert: adminLogStruct,
        },
      }
      fmt.Println("ops : ",ops)
      id := bson.NewObjectId() // Optional
      runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"ActionOnWithdrawalRequest"))
      if runnerErr != nil {
        fmt.Println("Error while runner : ",runnerErr)
        if resumeErr := runner.Resume(id); resumeErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
        }
      }else{
        fmt.Println("Executed")
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
      }
    }*/
    collection = setCollection("rozgar_db", "userDetails_collection")
    var updtStatusErr error
    withdrawalAmt,floatErr := strconv.ParseFloat(credMap["withdrawalAmt"].(string),64)
    if floatErr != nil {
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else{
      if credMap["status"].(string) == "Declined" {
        if isUpdated,isUpdtErr := updateCompanyOnDeclined(credMap,withdrawalAmt); isUpdtErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Failed to take action on withdrawal"}})
        }else if isUpdated {
          //- is to make withdrawal amt positive
          updtStatusErr = collection.Update(bson.M{"_id":credMap["userId"].(string),"transaction_history.trans_id":credMap["transactionId"].(string)},bson.M{"$inc":bson.M{"hrp":-withdrawalAmt},"$set":bson.M{"transaction_history.$.status":credMap["status"].(string),"transaction_history.$.bank_trans_id":credMap["paidId"].(string)}})
        }else{
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Failed to take action on withdrawal"}})
        }
      }else{
        updtStatusErr = collection.Update(bson.M{"_id":credMap["userId"].(string),"transaction_history.trans_id":credMap["transactionId"].(string)},bson.M{"$set":bson.M{"transaction_history.$.status":credMap["status"].(string),"transaction_history.$.bank_trans_id":credMap["paidId"].(string)}})
      }
      if updtStatusErr != nil {
        fmt.Println("Error while updating withdrawal request status : ",updtStatusErr)
        //Rollback for company
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Failed to take action on withdrawal"}})
      }else{
        collection = setCollection("rozgar_db", "adminActivityLog_collection")
        actDetailsStruct := StructConfig.ActivityDetails{To:"",From:"",Amount:withdrawalAmt}
        adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["adminUserName"].(string),ActivityFor:"Withdrawal request set to "+credMap["status"].(string),ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:credMap["userName"].(string),ActivityDetails:actDetailsStruct}
        adminLogStructErr := collection.Insert(adminLogStruct)
        if adminLogStructErr != nil {
          fmt.Println("Error while inserting admin log : ",adminLogStructErr)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Failed to take action on withdrawal"}})
        }else{
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
        }
      }
    }
  }
}

func updateCompanyOnDeclined(credMap map[string]interface{},withdrawalAmt float64)(bool,error){
  if companyId,comIdErr := getCompanyId(); comIdErr != nil {
    fmt.Println("Error while getting companyId : ",comIdErr)
    //responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    return false,comIdErr
  }else{
    collection = setCollection("rozgar_db", "userDetails_collection")
    updtStatusErr := collection.Update(bson.M{"_id":companyId,"transaction_history.trans_id":credMap["transactionId"].(string)},bson.M{"$inc":bson.M{"hrp":withdrawalAmt},"$set":bson.M{"transaction_history.$.units":0,"transaction_history.$.status":credMap["status"].(string),"transaction_history.$.bank_trans_id":credMap["paidId"].(string)}})
    if updtStatusErr != nil {
      fmt.Println("Error while updating company trans on withdrawal declined : ",updtStatusErr)
      return false,updtStatusErr
    }
    return true,nil
  }
}

func DonateToCompany(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      if companyId,comIdErr := getCompanyId(); comIdErr != nil {
        fmt.Println("Error while getting companyId : ",comIdErr)
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        //deduct respective amount and create a transacation
        //Add respective amount to company acc
        donateAmt,floatErr := strconv.ParseFloat(credMap["donateAmt"].(string),64)
        if floatErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          updt := bson.M{"$inc": bson.M{"hrp": -donateAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:"company",From:credMap["userName"].( string),Units:-donateAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Donation",BankTransId:"",Status:"Done",Level:""}}}
          updtCmpny := bson.M{"$inc": bson.M{"hrp": donateAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:"company",From:credMap["userName"].( string),Units:donateAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Donation",BankTransId:"",Status:"Done",Level:""}}}
          runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
          ops := []txn.Op{{
              //Adding new user instance
              C:      "userDetails_collection",
              Id:     credMap["userId"].(string),
              Assert: bson.M{"hrp": bson.M{"$gte": donateAmt}},
              Update: updt,
            },
            {
                //Adding new user instance
                C:      "userDetails_collection",
                Id:     companyId,
                //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
                Update: updtCmpny,
            },
          }
          id := bson.NewObjectId() // Optional
          runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"UpdateUserProfile"))
          if runnerErr != nil {
            fmt.Println("Error while runner : ",runnerErr)
            if resumeErr := runner.Resume(id); resumeErr != nil {
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
            }else{
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
            }
          }else{
            fmt.Println("Executed")
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }
      }
    }else{
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"User does not exist "}})
    }
  }
}

func DonateToCompanyAdmin(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    if isUserExist,isUserExistErr := isUserExist(credMap["adminUserName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      if companyId,comIdErr := getCompanyId(); comIdErr != nil {
        fmt.Println("Error while getting companyId : ",comIdErr)
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        //deduct respective amount and create a transacation
        //Add respective amount to company acc
        donateAmt,floatErr := strconv.ParseFloat(credMap["donateAmt"].(string),64)
        if floatErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          updt := bson.M{"$inc": bson.M{"hrp": -donateAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:"company",From:credMap["userName"].( string),Units:-donateAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Donation",BankTransId:"",Status:"Done",Level:""}}}
          updtCmpny := bson.M{"$inc": bson.M{"hrp": donateAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:"company",From:credMap["userName"].( string),Units:donateAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Donation",BankTransId:"",Status:"Done",Level:""}}}
          actDetailsStruct := StructConfig.ActivityDetails{To:"Company",From:credMap["userName"].(string),Amount:donateAmt}
          adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["adminUserName"].(string),ActivityFor:"Donate HRP",ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:credMap["userName"].(string),ActivityDetails:actDetailsStruct}

          runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
          ops := []txn.Op{{
              //Adding new user instance
              C:      "userDetails_collection",
              Id:     credMap["userId"].(string),
              Assert: bson.M{"hrp": bson.M{"$gte": donateAmt}},
              Update: updt,
            },
            {
                //Adding new user instance
                C:      "userDetails_collection",
                Id:     companyId,
                //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
                Update: updtCmpny,
            },
            {
              C:      "adminActivityLog_collection",
              Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
              //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
              Insert: adminLogStruct,
            },
          }
          id := bson.NewObjectId() // Optional
          runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"UpdateUserProfile"))
          if runnerErr != nil {
            fmt.Println("Error while runner : ",runnerErr)
            if resumeErr := runner.Resume(id); resumeErr != nil {
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
            }else{
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
            }
          }else{
            fmt.Println("Executed")
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }
      }
    }else{
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"User does not exist "}})
    }
  }
}

func getUserId(userName string)string{
  collection = setCollection("rozgar_db","userDetails_collection")
  //levelDocStruct := StructConfig.LevelDoc{}
  m := bson.M{}
  err := collection.Find(bson.M{"user_name":userName}).Select(bson.M{"_id":1}).One(&m); if err != nil {
    fmt.Println("Error while get userId @ uname : ",userName," err is : ",err)
    return ""
  }else{
    //return levelDocStruct.Id
    return m["_id"].(string)
  }
}

func AddUserInstance(email string,password string,mobileNo string)(bool,error){
  collection = setCollection("rozgar_db", "userInstance_collection")
  UserInstanceStruct := &StructConfig.UserInstance{Username:email,Password:password,MobileNo:mobileNo}
  insertInstanceErr := collection.Insert(UserInstanceStruct)
  if insertInstanceErr != nil {
    return false,insertInstanceErr
  } else {
    return true,nil
  }
}

func Login(w http.ResponseWriter, r *http.Request,interfaceName string) {
  collection = setCollection("rozgar_db","userInstance_collection")
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    //fmt.Println(credMap)
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      fmt.Println("User exist ...",isUserExist)
      userDetailsStruct := StructConfig.UserDetails{}
      collection = setCollection("rozgar_db","userDetails_collection")
      detailsErr := collection.Find(bson.M{"user_name":credMap["userName"].(string),"user_role":""}).Select(bson.M{"transaction_history":0}).One(&userDetailsStruct)
      if detailsErr != nil {
        fmt.Println("Error while fetching user details : ",detailsErr)
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        tokenString,tokenErr := TokenManager.GenerateToken(userDetailsStruct)
        if tokenErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Token not generated"}})
        }else{
          fmt.Println("Login Successfully !!! ")
          UpdateFBToken(credMap["userName"].(string),credMap["fbToken"].(string))
          responder(w,[]StructConfig.LoginResponse{StructConfig.LoginResponse{Response:"true",UserDetails:userDetailsStruct,TokenString:tokenString,FollowerCounts:getFollowerCounts(userDetailsStruct.IdentityCode),TeamCounts:getTeamCounts(userDetailsStruct.IdentityCode),Earned:totalEarnedAmount(userDetailsStruct.IdentityCode),BroadcastDetails:getActiveBroadcast(),ErrInResponse:""}})
        }
      }
    }else{
      fmt.Println("Username or password not exist.")
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Check username or Password."}})
    }
  }
}

func RefreshFeed(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    userDetailsStruct := StructConfig.UserDetails{}
    collection = setCollection("rozgar_db","userDetails_collection")
    detailsErr := collection.Find(bson.M{"user_name":credMap["userName"].(string)}).Select(bson.M{"transaction_history":0}).One(&userDetailsStruct)
    if detailsErr != nil {
      fmt.Println("Error while fetching user details : ",detailsErr)
    }else{
      fmt.Println("Latest data sent !!! ")
      /*data := make(map[string]string)
      var fbTokenArr = []string{"daMY85gg6nY:APA91bFjw8GmqNk7Cv1HrrG5wlQ4vayBD9DSqTNj3nPkBlmK6wqgjg8TjqzrCl7pSoxlCOdasdA14owdGdsGfgdYlLTK5DJM5tSaGEeRFs0LsqUR8HJaf-_i7R4SJHkylukyCr8qqSxn","ef9nH7biwog:APA91bETicEK2cYdI5z4Nlt5OqEkkHXpDrl60VfDHfZelYriScaNPNCTghFTjTajODeJFrVhGpJK-xgpJuABL1k7rhw_rQCrf-mkkgOOLZM1OiM6J8OfszjUuVsbWnzmvc6wEDpy3X-K","dDdSTii9LQo:APA91bHaDaoAPReWkDF0xHSUOVzMnTxx0dD6f7JBSHcqIaN5MTY7pgGHo0w_BVqb0PJiCbRV9sWlAjBBm0-9MAbTFzyV-JvXHmhFas0hCRuWyEPPgWJ1duOpAVYoQt86feDcPpa0Z02Q","e5N_qfg7ll8:APA91bFQhik4MeHvQ7sRbks-WbPUMwYHpC30XA0J38SOBClf2pDO1Le59hzU3KK-nICmZ8fi-pqXW_7GEI70_E3wLlGsYUhEeLiURCXlVwfWUEHBxMo-6a6OxZlb_DDrV8mEoaScYjvh"}
      data["title"] = "Happy Rozgar"
      data["body"] = "Feed Refreshed."
      sendNotification(fbTokenArr,data)
      sendNotificationWithFCM()*/
      responder(w,[]StructConfig.RefreshResponse{StructConfig.RefreshResponse{Response:"true",UserDetails:userDetailsStruct,FollowerCounts:getFollowerCounts(userDetailsStruct.IdentityCode),TeamCounts:getTeamCounts(userDetailsStruct.IdentityCode),Earned:totalEarnedAmount(userDetailsStruct.IdentityCode),BroadcastDetails:getActiveBroadcast(),ErrInResponse:""}})
    }
  }
}

func RefreshAdminFeed(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    userDetailsStruct := StructConfig.UserDetails{}
    collection = setCollection("rozgar_db","userDetails_collection")
    detailsErr := collection.Find(bson.M{"user_name":credMap["userName"].(string)}).Select(bson.M{"transaction_history":0}).One(&userDetailsStruct)
    if detailsErr != nil {
      fmt.Println("Error while fetching user details : ",detailsErr)
    }else{
      fmt.Println("Latest data sent !!! ")
      /*data := make(map[string]string)
      var fbTokenArr = []string{"daMY85gg6nY:APA91bFjw8GmqNk7Cv1HrrG5wlQ4vayBD9DSqTNj3nPkBlmK6wqgjg8TjqzrCl7pSoxlCOdasdA14owdGdsGfgdYlLTK5DJM5tSaGEeRFs0LsqUR8HJaf-_i7R4SJHkylukyCr8qqSxn","ef9nH7biwog:APA91bETicEK2cYdI5z4Nlt5OqEkkHXpDrl60VfDHfZelYriScaNPNCTghFTjTajODeJFrVhGpJK-xgpJuABL1k7rhw_rQCrf-mkkgOOLZM1OiM6J8OfszjUuVsbWnzmvc6wEDpy3X-K","dDdSTii9LQo:APA91bHaDaoAPReWkDF0xHSUOVzMnTxx0dD6f7JBSHcqIaN5MTY7pgGHo0w_BVqb0PJiCbRV9sWlAjBBm0-9MAbTFzyV-JvXHmhFas0hCRuWyEPPgWJ1duOpAVYoQt86feDcPpa0Z02Q","e5N_qfg7ll8:APA91bFQhik4MeHvQ7sRbks-WbPUMwYHpC30XA0J38SOBClf2pDO1Le59hzU3KK-nICmZ8fi-pqXW_7GEI70_E3wLlGsYUhEeLiURCXlVwfWUEHBxMo-6a6OxZlb_DDrV8mEoaScYjvh"}
      data["title"] = "Happy Rozgar"
      data["body"] = "Feed Refreshed."
      sendNotification(fbTokenArr,data)
      sendNotificationWithFCM()*/
      responder(w,[]StructConfig.RefreshAdminResponse{StructConfig.RefreshAdminResponse{Response:"true",UserDetails:userDetailsStruct,AdminCounts:adminCounts(userDetailsStruct.IdentityCode),BroadcastDetails:getActiveBroadcast(),ErrInResponse:""}})
    }
  }
}

func UpdateFBToken(userName string,fbToken string){
  collection = setCollection("rozgar_db","userDetails_collection")
  if existErr := collection.Update(bson.M{"firebase_token":fbToken},bson.M{"$set":bson.M{"firebase_token":""}}); existErr != nil{
    if existErr.Error() != "not found" {
      fmt.Println("Error while removing existing fbToken with same Handset... ",existErr)
    }
  }

  err := collection.Update(bson.M{"user_name":userName},bson.M{"$set":bson.M{"firebase_token":fbToken}})
  if err != nil {
    fmt.Println("Error while updating firebase token : ",err)
  }else{
    fmt.Println("Firebase token updated.")
  }

}

func getFBToken(userId string)(error,string){
  m := bson.M{}
  collection = setCollection("rozgar_db","userDetails_collection")
  err := collection.Find(bson.M{"_id":userId}).Select(bson.M{"firebase_token":1,"_id":0}).One(&m)
  if err != nil {
    fmt.Println("Error while getFBToken : ",err)
    return err,""
  }else{
    //fmt.Println("Firebase token updated.")
    return nil,m["firebase_token"].(string)
  }
}

func getAdminsFBToken()(error,[]string){
  var fbTokens []string
  m := []bson.M{}
  collection = setCollection("rozgar_db","userDetails_collection")
  err := collection.Find(bson.M{"user_role":"admin"}).Select(bson.M{"firebase_token":1,"_id":0}).All(&m)
  if err != nil {
    fmt.Println("Error while getAdminsFBToken : ",err)
    return err,fbTokens
  }else{
    for i := 0; i < len(m); i++ {
      fbTokens = append(fbTokens,m[i]["firebase_token"].(string))
    }
    return nil,fbTokens
  }
}

func getUsersFBToken(userId string)(error,[]string){
  var fbTokens []string
  m := []bson.M{}
  collection = setCollection("rozgar_db","userDetails_collection")
  err := collection.Find(bson.M{"user_role":""}).Select(bson.M{"firebase_token":1,"_id":0}).All(&m)
  if err != nil {
    fmt.Println("Error while updating firebase token : ",err)
    return err,fbTokens
  }else{
    for i := 0; i < len(m); i++ {
      fbTokens = append(fbTokens,m[i]["firebase_token"].(string))
    }
    return nil,fbTokens
  }
}

/*func sendNotification(ids []string,data map[string]string){
  var sendIds []string
  for _,val := range ids {
    if val != "" {
      sendIds = append(sendIds,val)
    }
  }
  var xds []string/*{
      "token5",
      "token6",
      "token7",
  }
	c := fcm.NewFcmClient(serverKey)
    c.NewFcmRegIdsMsg(sendIds, data)
    c.AppendDevices(xds)
	status, err := c.Send()

	if err == nil {
    status.PrintResults()
	} else {
		fmt.Println(err)
	}
}*/

func sendNotificationWithFCM(tokenArr []string,data map[string]string){
  /*data := map[string]string{
		"msg": "Hello World1",
		"sum": "Happy Day",
	}*/
	c := fcm.NewFCM(serverKey)
	//var tokenArr = []string{"daMY85gg6nY:APA91bFjw8GmqNk7Cv1HrrG5wlQ4vayBD9DSqTNj3nPkBlmK6wqgjg8TjqzrCl7pSoxlCOdasdA14owdGdsGfgdYlLTK5DJM5tSaGEeRFs0LsqUR8HJaf-_i7R4SJHkylukyCr8qqSxn","e5N_qfg7ll8:APA91bFQhik4MeHvQ7sRbks-WbPUMwYHpC30XA0J38SOBClf2pDO1Le59hzU3KK-nICmZ8fi-pqXW_7GEI70_E3wLlGsYUhEeLiURCXlVwfWUEHBxMo-6a6OxZlb_DDrV8mEoaScYjvh","f7fU0JjYClA:APA91bE1YmmKovXy_31ZpiOuhGXJ2kgYpYo84rh2yrhQ3ANVF5a6Z18050Ndtjvte_FhDjhPBwikXgSeh-w-rXruE_WTgB184_5f91bOiuNSyMvYgPerr_9nqpuQsKfGo3KAWdP6cC7C"}
	response, err := c.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  tokenArr,
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
		Notification: fcm.Notification{
			Title: data["title"],
			Body:  data["body"],
		},
	})
	if err != nil {
		//log.Fatal(err)
    fmt.Println("Error while sending with fcm : ",err)
	}
	fmt.Println("Status Code   :", response.StatusCode)
	fmt.Println("Success       :", response.Success)
	fmt.Println("Fail          :", response.Fail)
	fmt.Println("Canonical_ids :", response.CanonicalIDs)
	fmt.Println("Topic MsgId   :", response.MsgID)
}

func Logout(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    collection = setCollection("rozgar_db","userDetails_collection")
    err := collection.Update(bson.M{"user_name":credMap["userName"].(string)},bson.M{"$set":bson.M{"firebase_token":""}})
    if err != nil {
      fmt.Println("Error while updating firebase token Logout : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Logout failed ... "}})
    }else{
      fmt.Println("Firebase token updated.")
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
    }
  }
}

func isUserExist(email,password string) (bool,error)  {
  collection = setCollection("rozgar_db","userInstance_collection")
  usr,err := collection.Find(bson.M{"user_name":email,"password":password}).Count()
  if err != nil {
    if err.Error() == "not found" {
    	fmt.Println("Incorrect username or does not exist")
			return false, nil
		}else{
			fmt.Println("Something went wrong isUserExist me : ", err)
			return false, err
		}
  }else if usr < 1 {
    //User does not exist
    return false,nil
  }else{
    return true,nil
  }
}

func isUserValid(userName string,mobileNo string) (bool,error)  {
  collection = setCollection("rozgar_db","userInstance_collection")
  usr,err := collection.Find(bson.M{"$or":[]bson.M{bson.M{"user_name":userName},bson.M{"mobile_no":mobileNo}}}).Count()
  fmt.Println("usr : ",usr)
  if err != nil {
    if err.Error() == "not found" {
    //if errors.IsNotFound(err) {
			fmt.Println("Incorrect username or mobile does not exist")
			return false, nil
		}else{
			fmt.Println("Something went wrong isUserValid me : ", err)
			return false, err
		}
  }else if usr >= 1 {
    //User does not exist
    fmt.Println("User already exist")
    return false,nil
  }else{
    return true,nil
  }
}

func getFollowerCounts(identityCode string)StructConfig.FollowerCounts{
  collection = setCollection("rozgar_db","userDetails_collection")

  cnt1,err1 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,\d+(?!,)$`,""}}).Count()
  if err1 != nil {
    fmt.Println("Error while getting level1 follower counts : ",err1)
    cnt1 = 0
  }

  cnt2,err2 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,\d+,\d+(?!,)$`,""}}).Count()
  if err2 != nil {
    fmt.Println("Error while getting level2 follower counts : ",err2)
    cnt2 = 0
  }

  cnt3,err3 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,\d+,\d+,\d+(?!,)$`,""}}).Count()
  if err3 != nil {
    fmt.Println("Error while getting level3 follower counts : ",err3)
    cnt3 = 0
  }

  cnt4,err4 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,\d+,\d+,\d+,\d+(?!,)$`,""}}).Count()
  if err4 != nil {
    fmt.Println("Error while getting level4 follower counts : ",err4)
    cnt4 = 0
  }

  cnt5,err5 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,\d+,\d+,\d+,\d+,\d+(?!,)$`,""}}).Count()
  if err5 != nil {
    fmt.Println("Error while getting level5 follower counts : ",err5)
    cnt5 = 0
  }
  fmt.Println("Level 1 : ",cnt1," Level 2 : ",cnt2," Level 3 : ",cnt3," Level 4 : ",cnt4,"Level 5 : ",cnt5)
  return StructConfig.FollowerCounts{Level1Count:int(cnt1),Level2Count:int(cnt2),Level3Count:int(cnt3),Level4Count:int(cnt4),Level5Count:int(cnt5)}
}

func getTeamCounts(identityCode string)StructConfig.TeamCounts {
  collection = setCollection("rozgar_db","userDetails_collection")

  active,err1 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,`,""},"account_status":"Active"}).Count()
  if err1 != nil {
    fmt.Println("Error while getting active team counts : ",err1)
    active = 0
  }

  non_active,err2 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,`,""},"$or":[]bson.M{bson.M{"account_status":"Pending"},bson.M{"account_status":"Suspended"}}}).Count()
  if err1 != nil {
    fmt.Println("Error while getting non_active team counts : ",err2)
    non_active = 0
  }

  total,err3 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,`,""}}).Count()
  if err3 != nil {
    fmt.Println("Error while getting total team counts : ",err3)
    total = 0
  }
  return StructConfig.TeamCounts{ActiveMembersCount:int(active),NonActiveMembersCount:int(non_active),TotalMembersCount:int(total)}
}

func totalEarnedAmount(identityCode string)float64{
  m := []bson.M{}
  totalEarned := 0.0
  err := collection.Pipe([]bson.M{bson.M{"$match":bson.M{"identity_code":identityCode}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Level Income"}}}}}}}).All(&m)
  if err != nil {
    fmt.Println("Error while total Earned ammount : ",err)
  }else{
      if convertMap,convertErr := Underdog.InterfaceArrToMap(m[0]["transaction_history"]); convertErr != nil {
        fmt.Println("Error while InterfaceToArrMap totalEarnedAmount : ",convertErr)
      } else {
        for i := 0; i < len(convertMap); i++ {
          totalEarned = totalEarned + convertMap[i]["units"].(float64)
        }
      }
  }
  return totalEarned
}

func totalEarnedCompany(identityCode string)float64{
  m := []bson.M{}
  totalEarned := 0.0
  err := collection.Pipe([]bson.M{bson.M{"$match":bson.M{"identity_code":identityCode}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$or":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Topup"}},bson.M{"$eq":[]interface{}{"$$t.transaction_for","Bank Details Update"}}}}}}}}}).All(&m)
  if err != nil {
    fmt.Println("Error while total Earned ammount : ",err)
  }else{
      if convertMap,convertErr := Underdog.InterfaceArrToMap(m[0]["transaction_history"]); convertErr != nil {
        fmt.Println("Error while InterfaceToArrMap totalEarnedAmount : ",convertErr)
      } else {
        for i := 0; i < len(convertMap); i++ {
          totalEarned = totalEarned + convertMap[i]["units"].(float64)
        }
      }
  }
  return totalEarned
}

func LoginAdmin(w http.ResponseWriter, r *http.Request,interfaceName string) {
  collection = setCollection("rozgar_db","userInstance_collection")
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    //fmt.Println(credMap)
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      //fmt.Println("User already exist ...")
      userDetailsStruct := StructConfig.UserDetails{}
      collection = setCollection("rozgar_db","userDetails_collection")
      detailsErr := collection.Find(bson.M{"user_name":credMap["userName"].(string),"user_role":"admin"}).Select(bson.M{"transaction_history":0}).One(&userDetailsStruct)
      if detailsErr != nil {
        fmt.Println("Error while fetching user details : ",detailsErr)
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        tokenString,tokenErr := TokenManager.GenerateToken(userDetailsStruct)
        if tokenErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Token not generated"}})
        }else{
          fmt.Println("Login Successfully !!! ")
          UpdateFBToken(credMap["userName"].(string),credMap["fbToken"].(string))
          //responder(w,[]StructConfig.LoginResponse{StructConfig.LoginResponse{Response:"true",UserDetails:userDetailsStruct,TokenString:tokenString,FollowerCounts:getFollowerCounts(userDetailsStruct.IdentityCode),TeamCounts:getTeamCounts(userDetailsStruct.IdentityCode),Earned:totalEarnedAmount(userDetailsStruct.IdentityCode),ErrInResponse:""}})
          responder(w,[]StructConfig.LoginAdminResponse{StructConfig.LoginAdminResponse{Response:"true",UserDetails:userDetailsStruct,TokenString:tokenString,AdminCounts:adminCounts(userDetailsStruct.IdentityCode),BroadcastDetails:getActiveBroadcast(),ErrInResponse:""}})
        }
      }
    }else{
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"User does not exist "}})
    }
  }
}

func adminCounts(identityCode string)StructConfig.AdminCounts {
  //https://stackoverflow.com/questions/38449780/transact-sql-equivalent-in-golang-mongodb-aggregate
  //All admin active users
  //All admin Inactive users
  //All admin activities
  //admin total earnings
  //admin balance
  //Broadcast

  //All company active users
  //All company Inactive users
  //All withdrawal requests
  //New joinee this month
  //company total earnings
  //company balance


  /*err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":credMap["userName"].(string)}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Bank Details Update"}}}}}}}).All(&m)
  db.userDetails_collection.aggregate([{"$facet":{"ActiveUsersCount":[{"$match":{"account_status":"Active"}},{"$count":"ActiveUsersCount"},],"InactiveUsersCount":[{"$match":{"account_status":"Inactive"}},{"$count":"InactiveUsersCount"},],"WithdrawalRequestsCount":[{"$match":{"transaction_history":{"$elemMatch":{"$and":[{"transaction_for":"Withdrawal"},{"status":"Processing"}]}}}},{"$count":"WithdrawalRequestsCount"},]}},{"$project":{"ActiveUsersCount":{$arrayElemAt:["$ActiveUsersCount.ActiveUsersCount",0]},"InactiveUsersCount":{$arrayElemAt:["$InactiveUsersCount.InactiveUsersCount",0]},"WithdrawalRequestsCount":{$arrayElemAt:["$WithdrawalRequestsCount.WithdrawalRequestsCount",0]}}}])*/
  /*m := []bson.M{}
  //OG
  err:=collection.Pipe([]bson.M{bson.M{"$facet":bson.M{"ActiveUsersCount":[]bson.M{bson.M{"$match":bson.M{"account_status":"Active"}},bson.M{"$count":"ActiveUsersCount"},},"InactiveUsersCount":[]bson.M{bson.M{"$match":bson.M{"account_status":"Inactive"}},bson.M{"$count":"InactiveUsersCount"},},"WithdrawalRequestsCount":[]bson.M{bson.M{"$match":{"transaction_history":bson.M{"$elemMatch":bson.M{"$and":[]bson.M{bson.M{"transaction_for":"Withdrawal"},{"status":"Processing"}}}}}},bson.M{"$count":"WithdrawalRequestsCount"},}}},bson.M{"$project":bson.M{"ActiveUsersCount":bson.M{"$arrayElemAt":[]interface{}{"$ActiveUsersCount.ActiveUsersCount",0}},"InactiveUsersCount":bson.M{"$arrayElemAt":[]interface{}{"$InactiveUsersCount.InactiveUsersCount",0}},"WithdrawalRequestsCount":bson.M{"$arrayElemAt":[]interface{}{"$WithdrawalRequestsCount.WithdrawalRequestsCount",0}}}}}).All(&m)
  fmt.Println("m @ adminCounts : ",m)*/
  /*err:=collection.Pipe([]bson.M{bson.M{"$facet":bson.M{"ActiveUsersCount":[]bson.M{bson.M{"$match":bson.M{"account_status":"Active"}},bson.M{"$count":"ActiveUsersCount"}},"InactiveUsersCount":[]bson.M{bson.M{"$match":bson.M{"account_status":"Inactive"}},bson.M{"$count":"InactiveUsersCount"},},"WithdrawalRequestsCount":[]bson.M{bson.M{"$match":{"transaction_history":bson.M{"$elemMatch":bson.M{"$and":[]bson.M{bson.M{"transaction_for":"Withdrawal"},{"status":"Processing"}}}}}}},bson.M{"$count":"WithdrawalRequestsCount"},}}},bson.M{"$project":bson.M{"ActiveUsersCount":bson.M{"$arrayElemAt":[]interface{} {"$ActiveUsersCount.ActiveUsersCount",0}},"InactiveUsersCount":bson.M{"$arrayElemAt":[]interface{} {"$InactiveUsersCount.InactiveUsersCount",0}},"WithdrawalRequestsCount":bson.M{"$arrayElemAt":[]interface{} {"$WithdrawalRequestsCount.WithdrawalRequestsCount",0}}}}).All(&m)*/
  collection = setCollection("rozgar_db","adminActivityLog_collection")
  companyBalance := 0.0
  activities,actErr := collection.Find(bson.M{}).Count()
  if actErr != nil {
    fmt.Println("Error while getting activities count : ",actErr)
    activities = 0
  }
  collection = setCollection("rozgar_db","userDetails_collection")
  active,err1 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,`,""},"account_status":"Active"}).Count()
  if err1 != nil {
    fmt.Println("Error while getting active team counts : ",err1)
    active = 0
  }

  non_active,err2 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,`,""},"account_status":"Pending"}).Count()
  if err1 != nil {
    fmt.Println("Error while getting non_active team counts : ",err2)
    non_active = 0
  }

  /*total,err3 := collection.Find(bson.M{"identity_code":bson.RegEx{``+identityCode+`,`,""}}).Count()
  if err3 != nil {
    fmt.Println("Error while getting total team counts : ",err3)
    total = 0
  }*/

  companyActive,err4 := collection.Find(bson.M{"identity_code":bson.RegEx{`hrp`,""},"account_status":"Active"}).Count()
  if err4 != nil {
    fmt.Println("Error while getting company active team counts : ",err4)
    companyActive = 0
  }

  companyInactive,err5 := collection.Find(bson.M{"identity_code":bson.RegEx{`hrp`,""},"account_status":"Pending"}).Count()
  if err5 != nil {
    fmt.Println("Error while getting company non_active team counts : ",err5)
    companyInactive = 0
  }
  userDetailsStruct := StructConfig.UserDetails{}
  balErr := collection.Find(bson.M{"identity_code":bson.RegEx{`hrp`,""}}).Select(bson.M{"hrp":1}).One(&userDetailsStruct)
  if balErr != nil {
    fmt.Println("Error while fetching company bal : ",balErr)
  }else{
    companyBalance = userDetailsStruct.HRP
  }
  m := []bson.M{}
  withdrawalCount := 0
  //withdrawalCount,err6 := collection.Find(bson.M{"transaction_history.$.transaction_for":"Withdrawal","transaction_history.$.status":"Processing"}).Count()
  err6 := collection.Pipe([]bson.M{bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Processing"}}}}}},"_id":0}}}).All(&m)
  if err6 != nil {
    fmt.Println("Error while withdrawal count : ",err6)
    withdrawalCount = 0
  }else{
    grandTransactionHistoryStruct := []StructConfig.TransactionHistory{}
    transactionHistoryStruct := []StructConfig.TransactionHistory{}
    for i := 0; i < len(m); i++ {
      b, errMar := json.Marshal(m[i]["transaction_history"])
      if errMar != nil {
        fmt.Println("error while marshal : ", errMar)
        //fmt.Fprintf(w, "%s", b)
      } else {
        errUnm := json.Unmarshal([]byte(b), &transactionHistoryStruct)
        if errUnm != nil {
          fmt.Println("Error while unmarshal Grand : ", errUnm)
        }else{
          for ind := 0; ind < len(transactionHistoryStruct); ind++ {
            grandTransactionHistoryStruct = append(grandTransactionHistoryStruct,transactionHistoryStruct[ind])
          }
        }
      }
    }
    withdrawalCount = len(grandTransactionHistoryStruct)
  }
  return StructConfig.AdminCounts{AdminActiveUsers:int(active),AdminInactiveUsers:int(non_active),AdminActivitiesCount:int(activities),AdminTotalEarnings:totalEarnedAmount(identityCode),CompanyActiveUsers:int(companyActive),CompanyInactiveUsers:int(companyInactive),AllWithdrawalRequest:withdrawalCount,CompanyNewJoineeMonth:0,CompanyTotalEarnings:totalEarnedCompany("hrp"),CompanyBalance:companyBalance,CompanyFollowerCounts:getFollowerCounts("hrp")}
}

func getActiveBroadcast() StructConfig.BroadcastDetails {
  collection = setCollection("rozgar_db","broadcast_collection")
  broadcastDetailsStruct := StructConfig.BroadcastDetails{}
  err := collection.Find(bson.M{"broad_status":"Active"}).One(&broadcastDetailsStruct)
  if err != nil {
    if err.Error() != "not found" {
      fmt.Println("Error while getting broadcast details : ",err)
    }
  }
  return broadcastDetailsStruct
}

func GetSponsorList(w http.ResponseWriter, r *http.Request,interfaceName string){
  /*if credMap,credErr:=Underdog.StringToMap(interfaceName); credErr != nil {
    fmt.Println("Error while interface to map @GetSponsorList : ",credErr)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  } else if len(credMap) > 0 {*/
  //fmt.Println("in GetSponsorList")
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := []StructConfig.UserDetails{}
    err := collection.Find(bson.M{}).Select(bson.M{"_id":1,"user_name":1,"personal_info.full_name":1,"identity_code":1}).All(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching user details list : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //fmt.Println(userDetailsStruct)
      responder(w,[]StructConfig.SponsorListResponse{StructConfig.SponsorListResponse{Response:"true",SponsorList:userDetailsStruct,ErrInResponse:""}})
    }
  /*}else{
    fmt.Println("Creadential map is empty...")
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  }*/
}

func GetPendingUserListBySponsor(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := []StructConfig.UserDetails{}
    err := collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,`,""},"account_status":"Pending"}).Select(bson.M{"_id":1,"user_name":1,"personal_info.full_name":1,"identity_code":1}).All(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching user details list : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //fmt.Println(userDetailsStruct)
      responder(w,[]StructConfig.SponsorListResponse{StructConfig.SponsorListResponse{Response:"true",SponsorList:userDetailsStruct,ErrInResponse:""}})
    }
  }
}

func TopupAccount(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    if isUserExist,isUserExistErr := isUserExist(credMap["sponsorUName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      var amt []float64
      hrp,floatErr := strconv.ParseFloat(credMap["hrp"].(string),64)
      if floatErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        if userId,userIdErr := getIdByIdCode(credMap["accountUserIdentityCode"].(string),"Pending"); userIdErr != nil {
          fmt.Println("Error while getting pending userId : ",userIdErr)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong, try again"}})
        }else{
          remainingAmt := hrp
          var userIdArr []string
          var userNameArr []string
          var fbTokenArr []string
          //levelMap := distributeHRPs(credMap["sponsorIdentityCode"].(string))
          levelMap := distributeHRPs(getSponsorIdByIdCode(credMap["accountUserIdentityCode"].(string)))

          if val, ok := levelMap["level1"]; ok {
            //500 (33.33)  accountUserIdentityCode
            //fmt.Println("idCode @ level1 : ",val)
            if strings.Contains(val,",") {
              //Transfer to val
              if valId,valIdErr := getIdByIdCode(val,""); valIdErr != nil {
                fmt.Println("Unable to fetch id by id @ level1 : ",valIdErr)
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong, try again"}})
              }else{
                //fmt.Println("ValId : ",valId)
                userIdArr = append(userIdArr,valId["id"])
                userNameArr = append(userNameArr,valId["user_name"])
                fbTokenArr = append(fbTokenArr,valId["firebase_token"])
                amt = append(amt,math.Round((hrp*33.33/100)/0.5*0.5))
                remainingAmt = remainingAmt - amt[0]
              }
            }else{
              //Transfer remaining to company
            }
          }
          if val, ok := levelMap["level2"]; ok {
            //200 (13.33)
            //fmt.Println("idCode @ level2 : ",val)
            if strings.Contains(val,",") {
              //Transfer to val
              if valId,valIdErr := getIdByIdCode(val,""); valIdErr != nil {
                fmt.Println("Unable to fetch id by id @ level2 : ",valIdErr)
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong, try again"}})
              }else{
                userIdArr = append(userIdArr,valId["id"])
                userNameArr = append(userNameArr,valId["user_name"])
                fbTokenArr = append(fbTokenArr,valId["firebase_token"])
                amt = append(amt,math.Round((hrp*13.33/100)/0.5*0.5))
                remainingAmt = remainingAmt - amt[1]
              }
            }else{
              //Transfer remaining to company
            }
          }
          if val, ok := levelMap["level3"]; ok {
            //100 (6.66)
            //fmt.Println("idCode @ level3 : ",val)
            if strings.Contains(val,",") {
              //Transfer to val
              if valId,valIdErr := getIdByIdCode(val,""); valIdErr != nil {
                fmt.Println("Unable to fetch id by id @ level3 : ",valIdErr)
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong, try again"}})
              }else{
                userIdArr = append(userIdArr,valId["id"])
                userNameArr = append(userNameArr,valId["user_name"])
                fbTokenArr = append(fbTokenArr,valId["firebase_token"])
                amt = append(amt,math.Round((hrp*6.66/100)/0.5*0.5))
                remainingAmt = remainingAmt - amt[2]
              }
            }else{
              //Transfer remaining to company
            }
          }
          if val, ok := levelMap["level4"]; ok {
            //100 (6.66)
            //fmt.Println("idCode @ level4 : ",val)
            if strings.Contains(val,",") {
              //Transfer to val
              if valId,valIdErr := getIdByIdCode(val,""); valIdErr != nil {
                fmt.Println("Unable to fetch id by id @ level4 : ",valIdErr)
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong, try again"}})
              }else{
                userIdArr = append(userIdArr,valId["id"])
                userNameArr = append(userNameArr,valId["user_name"])
                fbTokenArr = append(fbTokenArr,valId["firebase_token"])
                amt = append(amt,math.Round((hrp*6.66/100)/0.5*0.5))
                remainingAmt = remainingAmt - amt[3]
              }
            }else{
              //Transfer remaining to company
            }
          }
          if val, ok := levelMap["level5"]; ok {
            //100 (6.66)
            //fmt.Println("idCode @ level5 : ",val)
            if strings.Contains(val,",") {
              //Transfer to val
              if valId,valIdErr := getIdByIdCode(val,""); valIdErr != nil {
                fmt.Println("Unable to fetch id by id @ level5 : ",valIdErr)
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong, try again"}})
              }else{
                userIdArr = append(userIdArr,valId["id"])
                userNameArr = append(userNameArr,valId["user_name"])
                fbTokenArr = append(fbTokenArr,valId["firebase_token"])
                amt = append(amt,math.Round((hrp*6.66/100)/0.5*0.5))
                remainingAmt = remainingAmt - amt[4]
              }
            }else{
              //Transfer remaining to company
            }
          }

          isTrans,transErr := TransactionQueries(credMap["userId"].(string),credMap["userName"].(string),userIdArr,userNameArr,fbTokenArr,amt,remainingAmt,userId["id"],userId["user_name"])
          if transErr != nil {
            fmt.Println("Error in TransactionQueries : ",transErr)
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Transaction Failed, try again !!!"}})
          }else if isTrans {
            //fmt.Println("Transaction successfully done")
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }else{
            fmt.Println("Transaction failed  ")
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Transaction Failed, try again !!!"}})
          }
        }
      }
    }else{
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"User does not exist "}})
    }
  }
}

func TransactionQueries(sponsorId string,sponsorUserName string,userId []string,userNames []string,fbTokenArr[]string,amt []float64,remainingAmt float64,pendingUserId string,pendingUserName string)(bool,error){
  if companyId,comIdErr := getCompanyId(); comIdErr != nil {
    fmt.Println("Error while getting companyId : ",comIdErr)
    return false,nil
  }else{
    //fmt.Println("Len of userId : ",len(userId))
    //fmt.Println("Sponsor id: ",sponsorId," remainingAmt : ",remainingAmt," userId : ",userId," amt : ",amt)
    level := 0
    totalAmt := 1500.00
    updtDec := bson.M{"$inc": bson.M{"hrp": -totalAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:pendingUserName,From:sponsorUserName,Units:-totalAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Topup",BankTransId:"",Status:"Done",Level:""}}}
    //fmt.Println("debiting query : ",updtDec)
    runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
    ops := []txn.Op{}

    ops = append(ops,txn.Op{
      C:      "userDetails_collection",
      Id:     sponsorId,
      Assert: bson.M{"hrp": bson.M{"$gte": totalAmt}},
      Update: updtDec,
    })

    for index := 0;  index < len(userId); index++ {
      level++
      //if index != atPos {
        updt := bson.M{"$inc": bson.M{"hrp": amt[index]},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:userNames[index],From:sponsorUserName,Units:amt[index],TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Level Income",BankTransId:"",Status:"Done",Level:"Level "+strconv.Itoa(level)}}}
        //fmt.Println("@ index : ",index," query is : ",updt," user id : ",userId[index])

        /*data := make(map[string]string)

        data["title"] = "Happy Rozgar"
        data["body"] = "Login Successfully done."
        if fbTokenErr,fbToken := getFBToken("sallu"); fbTokenErr != nil {
          fmt.Println("Error while getting fbToken : ",fbTokenErr)
        }else{
          var ids = []string{fbToken}
          sendNotification(ids,data)
        }*/
        ops = append(ops,txn.Op{
          C:      "userDetails_collection",
          Id:     userId[index],
          //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
          // Assert: txn.DocMissing,
          Update: updt,
        })

    }
    updtCmpn := bson.M{"$inc": bson.M{"hrp": remainingAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:pendingUserName,From:sponsorUserName,Units:remainingAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Topup",BankTransId:"",Status:"Done",Level:""}}}

    //fmt.Println("Credit in company query : ",updtCmpn," Company id : ",companyId)
    ops = append(ops,txn.Op{
      C:      "userDetails_collection",
      Id:     companyId,
      //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},*/
      Update: updtCmpn,
    })
    //fmt.Println("Pending user id : ",pendingUserId)
    ops = append(ops,txn.Op{
      C:      "userDetails_collection",
      Id:     pendingUserId,
      Assert: bson.M{"account_status":"Pending"},
      Update: bson.M{"$set":bson.M{"account_status":"Active"}},
    })

    //fmt.Println("Len of ops : ",len(ops))
    //fmt.Println("ops : ",ops)
    id := bson.NewObjectId() // Optional
    runnerErr := runner.Run(ops,id,setInfoStruct(sponsorId,"TopupAccount"))
    if runnerErr != nil {
      fmt.Println("Error while runner : ",runnerErr)
      if resumeErr := runner.Resume(id); resumeErr != nil {
          return false,resumeErr
      }else{
        data := make(map[string]string)
        data["title"] = "Happy Rozgar"
        data["body"] = "Level Income credited."
        //sendNotification(fbTokenArr,data)
        sendNotificationWithFCM(fbTokenArr,data)
        if fbTokenErr,fbToken := getFBToken(pendingUserId); fbTokenErr != nil {
          fmt.Println("Error while getting fbToken : ",fbTokenErr)
        }else{
          var ids = []string{fbToken}
          data["body"] = "Your account is now Active."
          //sendNotification(ids,data)
          sendNotificationWithFCM(ids,data)
        }
        return true,nil
      }
     }else{
       data := make(map[string]string)
       data["title"] = "Happy Rozgar"
       data["body"] = "Level Income credited."
       sendNotificationWithFCM(fbTokenArr,data)
       if fbTokenErr,fbToken := getFBToken(pendingUserId); fbTokenErr != nil {
         fmt.Println("Error while getting fbToken : ",fbTokenErr)
       }else{
         data["body"] = "Your account is now Active."
         var ids = []string{fbToken}
         sendNotificationWithFCM(ids,data)
       }
       return true,nil
     }
     return false,nil
  }

   /**


   **/
}

func getCompanyId()(string,error){
  collection = setCollection("rozgar_db","userDetails_collection")
  m := bson.M{}
  err := collection.Find(bson.M{"user_role": "company"}).Select(bson.M{"_id":1}).One(&m)
  if err != nil {
    fmt.Println("Error while getCompanyId : ",err)
    return "",err
  }else{
    return m["_id"].(string),nil
  }
}

func getIdByIdCode(identityCode string,status string)(map[string]string,error){
  collection = setCollection("rozgar_db","userDetails_collection")
  m := bson.M{}
  var err error
  idUsernameMap := make(map[string]string)
  if status == "Pending" {
    err = collection.Find(bson.M{"identity_code": identityCode}).Select(bson.M{"_id":1,"user_name":1,"firebase_token":1}).One(&m)
  }else if status == "" {
    err = collection.Find(bson.M{"identity_code": identityCode,"account_status":"Active"}).Select(bson.M{"_id":1,"user_name":1,"firebase_token":1}).One(&m)
  }else{
    err = collection.Find(bson.M{"identity_code": identityCode}).Select(bson.M{"_id":1,"user_name":1,"firebase_token":1}).One(&m)
  }
  if err != nil {
    if err.Error() == "not found" {
      if status == "" {
        compId,_:= getCompanyId()
        idUsernameMap["id"] = compId
        idUsernameMap["user_name"] = "rozgar"
      }
      return idUsernameMap,nil
    }
    fmt.Println("Error while getIdByIdCode : ",err)
    return idUsernameMap,err
  }else{
    idUsernameMap["id"] = m["_id"].(string)
    idUsernameMap["user_name"] = m["user_name"].(string)
    return idUsernameMap,nil
  }
}

func getSponsorIdByIdCode(identityCode string)(string){
  sponsorId := ""
  result := strings.Split(identityCode,",")
	if len(result) > 1 {
    for i := 0; i < len(result)-1; i++ {
      if i == len(result)-2{
		      sponsorId = sponsorId+result[i]
	    }else{
	       sponsorId = sponsorId+result[i]+","
	    }
    }
  }
  return sponsorId
}

func distributeHRPs(sponsorIdentityCode string)map[string]string{
  levelMap := make(map[string]string)
  result := strings.Split(sponsorIdentityCode, ",")
  if(len(result) > 5){
	  levelMap["level1"] = sponsorIdentityCode
	  for i := 0; i < len(result)-1; i++ {
	    if i == len(result)-2{
		      levelMap["level2"] = levelMap["level2"]+result[i]
	    }else{
	        levelMap["level2"] = levelMap["level2"]+result[i]+","
	    }
	  }
	  for i := 0; i < len(result)-2; i++ {
	    if i == len(result)-3{
		      levelMap["level3"] = levelMap["level3"]+result[i]
	    }else{
	        levelMap["level3"] = levelMap["level3"]+result[i]+","
	    }
	  }
	  for i := 0; i < len(result)-3; i++ {
	     if i == len(result)-4{
		    levelMap["level4"] = levelMap["level4"]+result[i]
	    }else{
        levelMap["level4"] = levelMap["level4"]+result[i]+","
	    }
	  }
	  for i := 0; i < len(result)-4; i++ {
	    if i == len(result)-5{
		    levelMap["level5"] = levelMap["level5"]+result[i]
	    }else{
	    	levelMap["level5"] = levelMap["level5"]+result[i]+","
	    }
	  }

	}else if(len(result) > 4) {
	  if(len(result) == 5) {
	    //"hrp,01,2,1,1"
	    levelMap["level1"] = sponsorIdentityCode
	    //fmt.Println(result[0:4])
	    levelMap["level2"] = result[0]+","+result[1]+","+result[2]+","+result[3]
	    levelMap["level3"] = result[0]+","+result[1]+","+result[2]
	    levelMap["level4"] = result[0]+","+result[1]
	    levelMap["level5"] = result[0]
	  }
	}else if(len(result) > 3) {
	  if(len(result) == 4) {
	    //"hrp,01,2,1"
	    levelMap["level1"] = sponsorIdentityCode
	    levelMap["level2"] = result[0]+","+result[1]+","+result[2]
	    levelMap["level3"] = result[0]+","+result[1]
	    levelMap["level4"] = result[0]
	  }
	}else if(len(result) > 2) {
	  if(len(result) == 3) {
	    //"hrp,01,2"
	    levelMap["level1"] = sponsorIdentityCode
	    levelMap["level2"] = result[0]+","+result[1]
	    levelMap["level3"] = result[0]
	  }
	}else if(len(result) > 1) {
	  if(len(result) == 2) {
	    //"hrp,01"
	    levelMap["level1"] = sponsorIdentityCode
	    levelMap["level2"] = result[0]
	  }
	}else if(len(result) > 0) {
	  if(len(result) == 1) {
	    //"hrp"
	    levelMap["level1"] = sponsorIdentityCode
	  }
	}
  return levelMap
}

func GetUserListToTransfer(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := []StructConfig.UserDetails{}
    err := collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,\d+`,""},"account_status":"Active"}).Select(bson.M{"_id":1,"user_name":1,"personal_info.full_name":1,"identity_code":1}).All(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching user details list : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //fmt.Println(userDetailsStruct)
      responder(w,[]StructConfig.SponsorListResponse{StructConfig.SponsorListResponse{Response:"true",SponsorList:userDetailsStruct,ErrInResponse:""}})
    }
  }
}

func GetBankDetails(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := StructConfig.UserDetails{}
    err := collection.Find(bson.M{"user_name":credMap["userName"].(string)}).Select(bson.M{"bank_details":1,"_id":0}).One(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching bank details : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      responder(w,[]StructConfig.BankDetailsResponse{StructConfig.BankDetailsResponse{Response:"true",BankDetails:userDetailsStruct.BankDetails,ErrInResponse:""}})
    }
  }
}

func TransferHRP(w http.ResponseWriter, r *http.Request,interfaceName string){
  if credMap,credErr := Underdog.StringToMap(interfaceName); credErr != nil {
    fmt.Println("Error while interface to map @TransferHRP : ",credErr)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  } else if len(credMap) > 0 {
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      amtToTransfer,floatErr := strconv.ParseFloat(credMap["hrp"].(string),64)
      if floatErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        if toUserId,toUserIdErr := getIdByIdCode(credMap["toIdentityCode"].(string),"transfer"); toUserIdErr != nil {
          fmt.Println("Error while getting toUserId : ",toUserIdErr)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          creditAmt := bson.M{"$inc": bson.M{"hrp": amtToTransfer},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:toUserId["user_name"],From:credMap["userName"].( string),Units:amtToTransfer,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Fund Receive",BankTransId:"",Status:"Done",Level:""}}}
          debitAmt := bson.M{"$inc": bson.M{"hrp": -amtToTransfer},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:toUserId["user_name"],From:credMap["userName"].( string),Units:-amtToTransfer,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Fund Transfer",BankTransId:"",Status:"Done",Level:""}}}
          runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
          ops := []txn.Op{{
              //Transfer from
              C:      "userDetails_collection",
              Id:     credMap["userId"].(string),
              Assert: bson.M{"hrp": bson.M{"$gte": amtToTransfer}},
              Update: debitAmt,
            },
            {
              //Transfer to
              C:      "userDetails_collection",
              Id:     toUserId["id"],
              //Assert: bson.M{"identityCode": credMap["toIdentityCode"].(string)},
              Update: creditAmt,
            },
          }
          id := bson.NewObjectId()
          runnerErr := runner.Run(ops,id,setInfoStruct(credMap["userName"].(string),"TransferHRP"))
          if runnerErr != nil {
            fmt.Println("Error while runner : ",runnerErr)
            if resumeErr := runner.Resume(id); resumeErr != nil {
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
            }else{
              if fbTokenErr,fbToken := getFBToken(toUserId["id"]); fbTokenErr != nil {
                fmt.Println("Error while getting fbToken : ",fbTokenErr)
              }else{
                data := make(map[string]string)
                data["title"] = "Happy Rozgar"
                data["body"] = "HRP received from "+credMap["userName"].(string)+"."
                var ids = []string{fbToken}
                sendNotificationWithFCM(ids,data)
              }
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
            }
          }else{
            fmt.Println("Executed")
            if fbTokenErr,fbToken := getFBToken(toUserId["id"]); fbTokenErr != nil {
              fmt.Println("Error while getting fbToken : ",fbTokenErr)
            }else{
              data := make(map[string]string)
              data["title"] = "Happy Rozgar"
              data["body"] = "HRP received from "+credMap["userName"].(string)+"."
              var ids = []string{fbToken}
              sendNotificationWithFCM(ids,data)
            }
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }
      }
    }else{
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Incorrect Password."}})
    }
  }else{
    fmt.Println("Credential map is empty...")
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  }
}

func TransferHRPAdmin(w http.ResponseWriter, r *http.Request,interfaceName string){
  if credMap,credErr := Underdog.StringToMap(interfaceName); credErr != nil {
    fmt.Println("Error while interface to map @TransferHRPAdmin : ",credErr)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  } else if len(credMap) > 0 {
    if isUserExist,isUserExistErr := isUserExist(credMap["adminUserName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      amtToTransfer,floatErr := strconv.ParseFloat(credMap["hrp"].(string),64)
      if floatErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        if toUserId,toUserIdErr := getIdByIdCode(credMap["toIdentityCode"].(string),"transfer"); toUserIdErr != nil {
          fmt.Println("Error while getting toUserId : ",toUserIdErr)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          //creditAmt := bson.M{"$inc": bson.M{"hrp": amtToTransfer},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:toUserId["user_name"],From:credMap["userName"].(string),Units:amtToTransfer,TransactionOn:strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)),10),TransactionFor:"Fund Receive",BankTransId:"",Status:"Done",Level:""}}}
          creditAmt := bson.M{"$inc": bson.M{"hrp": amtToTransfer},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:toUserId["user_name"],From:credMap["userName"].(string),Units:amtToTransfer,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Fund Receive",BankTransId:"",Status:"Done",Level:""}}}
          //debitAmt := bson.M{"$inc": bson.M{"hrp": -amtToTransfer},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:toUserId["user_name"],From:credMap["userName"].(string),Units:-amtToTransfer,TransactionOn:strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)),10),TransactionFor:"Fund Transfer",BankTransId:"",Status:"Done",Level:""}}}
          debitAmt := bson.M{"$inc": bson.M{"hrp": -amtToTransfer},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:toUserId["user_name"],From:credMap["userName"].(string),Units:-amtToTransfer,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Fund Transfer",BankTransId:"",Status:"Done",Level:""}}}
          actDetailsStruct := StructConfig.ActivityDetails{To:toUserId["user_name"],From:credMap["userName"].(string),Amount:amtToTransfer}
          adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["adminUserName"].(string),ActivityFor:"Transfer HRP",ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:credMap["userName"].(string),ActivityDetails:actDetailsStruct}
          runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
          ops := []txn.Op{{
              //Transfer from
              C:      "userDetails_collection",
              Id:     credMap["userId"].(string),
              Assert: bson.M{"hrp": bson.M{"$gte": amtToTransfer}},
              Update: debitAmt,
            },
            {
              //Transfer to
              C:      "userDetails_collection",
              Id:     toUserId["id"],
              //Assert: bson.M{"identityCode": credMap["toIdentityCode"].(string)},
              Update: creditAmt,
            },
            {
              C:      "adminActivityLog_collection",
              Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
              //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
              Insert: adminLogStruct,
            },
          }
          id := bson.NewObjectId()
          runnerErr := runner.Run(ops,id,setInfoStruct(credMap["userName"].(string),"TransferHRP"))
          if runnerErr != nil {
            fmt.Println("Error while runner : ",runnerErr)
            if resumeErr := runner.Resume(id); resumeErr != nil {
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
            }else{
              if fbTokenErr,fbToken := getFBToken(toUserId["id"]); fbTokenErr != nil {
                fmt.Println("Error while getting fbToken : ",fbTokenErr)
              }else{
                data := make(map[string]string)
                data["title"] = "Happy Rozgar"
                data["body"] = "HRP received from "+credMap["userName"].(string)+"."
                var ids = []string{fbToken}
                sendNotificationWithFCM(ids,data)
              }
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
            }
          }else{
            fmt.Println("Executed")
            if fbTokenErr,fbToken := getFBToken(toUserId["id"]); fbTokenErr != nil {
              fmt.Println("Error while getting fbToken : ",fbTokenErr)
            }else{
              data := make(map[string]string)
              data["title"] = "Happy Rozgar"
              data["body"] = "HRP received from "+credMap["userName"].(string)+"."
              var ids = []string{fbToken}
              sendNotificationWithFCM(ids,data)
            }
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }
      }
    }else{
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Incorrect Password."}})
    }
  }else{
    fmt.Println("Credential map is empty...")
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  }
}

func GetFollowerListByLevel(w http.ResponseWriter, r *http.Request,interfaceName string){
  /*level 1 : hrp,01,\d(?!,)
  level 2 : hrp,01,\d,\d(?!,)
  level 3 : hrp,01,\d,\d,\d(?!,)
  level 4 : hrp,01,\d,\d,\d,\d(?!,)
  level 5 : hrp,01,\d,\d,\d,\d,\d(?!,)*/

  if credMap,credErr := Underdog.StringToMap(interfaceName); credErr != nil {
    fmt.Println("Error while interface to map @GetFollowerListByLevel : ",credErr)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  } else if len(credMap) > 0 {
    var err error
    userDetailsStruct := []StructConfig.UserDetails{}
    collection = setCollection("rozgar_db","userDetails_collection")
    if credMap["level"].(string) == "level1" {
      err = collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,\d+(?!,)$`,""}}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
    }else if credMap["level"].(string) == "level2" {
      err = collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,\d+,\d+(?!,)$`,""}}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
    }else if credMap["level"].(string) == "level3" {
      err = collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,\d+,\d+,\d+(?!,)$`,""}}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
    }else if credMap["level"].(string) == "level4" {
      err = collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,\d+,\d+,\d+,\d+(?!,)$`,""}}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
    }else if credMap["level"].(string) == "level5" {
      err = collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,\d+,\d+,\d+,\d+,\d+(?!,)$`,""}}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
    }else if credMap["level"].(string) == "active" {
      //err = collection.Find(bson.M{"sponsor_uname":credMap["userName"].(string),"account_status":"Active"}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
      err = collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,`,""},"account_status":"Active"}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
    }else if credMap["level"].(string) == "inactive" {
      err = collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,`,""},"$or":[]bson.M{bson.M{"account_status":"Pending"},bson.M{"account_status":"Suspended"}}}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
    }else if credMap["level"].(string) == "total" {
      err = collection.Find(bson.M{"identity_code":bson.RegEx{``+credMap["identityCode"].(string)+`,`,""}}).Select(bson.M{"_id":0,"firebase_token":0,"bank_details":0,"direct_child_count":0,"transaction_history":0}).All(&userDetailsStruct)
    }else{
      fmt.Println("Level not provided or undefined")
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }
    if err != nil {
      fmt.Println("Error while retrieving followers by level : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      responder(w,[]StructConfig.FollowersListResponse{StructConfig.FollowersListResponse{Response:"true",FollowerList:userDetailsStruct,ErrInResponse:""}})
    }
  }else{
    fmt.Println("Creadential map is empty...")
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  }
}

func DisplayAllFollowers(w http.ResponseWriter, r *http.Request,interfaceName string){
  if credMap,credErr := Underdog.StringToMap(interfaceName); credErr != nil {
    fmt.Println("Error while interface to map @DisplayFollowers : ",credErr)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  } else if len(credMap)  > 0 {
    userDetailsStruct := []StructConfig.UserDetails{}
    //err := collection.Find(bson.M{"user_name":credMap["userName"].(string)}).One(&levelDocStruct)
    err := collection.Find(bson.M{"identity_code":bson.RegEx{"/"+credMap["identityCode"].(string)+",/",""}}).All(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while retrieving followers : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      responder(w,[]StructConfig.FollowersListResponse{StructConfig.FollowersListResponse{Response:"true",FollowerList:userDetailsStruct,ErrInResponse:""}})
    }
  }else{
    fmt.Println("Creadential map is empty...")
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  }
}

func GetTransactionHistory(w http.ResponseWriter, r *http.Request,interfaceName string){
  if credMap,credErr := Underdog.StringToMap(interfaceName); credErr != nil {
    fmt.Println("Error while interface to map @GetTransactionHistory : ",credErr)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  } else if len(credMap) > 0 {
    collection = setCollection("rozgar_db","userDetails_collection")
    //userDetailsStruct := StructConfig.UserDetails{}
    transactionHistoryStruct := []StructConfig.TransactionHistory{}
    var err error
    /*if credMap["transactionType"].(string) == "All" {
      err = collection.Find(bson.M{"user_name":credMap["userName"].(string)}).One(&userDetailsStruct)
      if len(userDetailsStruct.TransactionHistory) > 0 {
        transactionHistoryStruct = userDetailsStruct.TransactionHistory
      }
    }else if credMap["transactionType"].(string) == "Incoming" {
      err = collection.Find(bson.M{"user_name":credMap["userName"].(string)}).Select(bson.M{"transaction_history":1}).One(&userDetailsStruct)
      for i := 0; i < len(userDetailsStruct.TransactionHistory); i++ {
        if userDetailsStruct.TransactionHistory[i].Units > 0 {
          fmt.Println("Positive value : ",userDetailsStruct.TransactionHistory[i].Units)
          transactionHistoryStruct = append(transactionHistoryStruct,userDetailsStruct.TransactionHistory[i])
        }
      }
    }else if credMap["transactionType"].(string) == "Outgoing" {
      err = collection.Find(bson.M{"user_name":credMap["userName"].(string)}).Select(bson.M{"transaction_history":1}).One(&userDetailsStruct)
      for i := 0; i < len(userDetailsStruct.TransactionHistory); i++ {
        if userDetailsStruct.TransactionHistory[i].Units < 0 {
          transactionHistoryStruct = append(transactionHistoryStruct,userDetailsStruct.TransactionHistory[i])
        }
      }
    }else if credMap["transactionType"].(string) == "Fund Transfer" {

    }*/
    fmt.Println("credMap : ",credMap)
    m := []bson.M{}
    if credMap["transactionType"].(string) == "Others" {
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":credMap["userName"].(string)}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":[]interface{}{bson.M{"$or":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Donation"}},bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.transaction_for","Bank Details Update"}}}}}}}}}}).All(&m)
      err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":credMap["userName"].(string)}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Bank Details Update"}}}}}}}).All(&m)
    }else if credMap["transactionType"].(string) == "All" {
      err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":credMap["userName"].(string)}},bson.M{"$project":bson.M{"transaction_history":1}}}).All(&m)
    }else{
      err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":credMap["userName"].(string)}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for",credMap["transactionType"].(string)}}}}}}}).All(&m)
    }
    b, errMar := json.Marshal(m[0]["transaction_history"])
    if errMar != nil {
      fmt.Println("error while marshal : ", errMar)
      //fmt.Fprintf(w, "%s", b)
    } else {
      errUnm := json.Unmarshal([]byte(b), &transactionHistoryStruct)
      if errUnm != nil {
        fmt.Println("Error while unmarshal Grand : ", errUnm)
      }
    }
    /*for i := 0; i < len(m["transaction_history"]); i++ {
      transactionHistoryStruct = append(transactionHistoryStruct,StructConfig.TransactionHistory{To:m["transaction_history"][i]["to"],From:m["transaction_history"][i]["from"],Units:m["transaction_history"][i]["units"],TransactionOn:m["transaction_history"][i]["transaction_on"],TransactionFor:m["transaction_history"][i]["transaction_for"],Level:m["transaction_history"][i]["level"]})
    }*/
    //fmt.Println(credMap["transactionType"].(string)+" struct : ",transactionHistoryStruct)

    //db.userDetails_collection.aggregate([{$match:{"user_name":"mazhar",}},{$project:{"transaction_history":{$filter:{input:"$transaction_history",as:"t",cond:{$eq:["$$t.transaction_for","Fund Transfer"]}}}}}]).pretty()

    //err = collection.Find(bson.M{"user_name":credMap["userName"].(string)}).All(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching transaction history : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //fmt.Println("userDetailsStruct : ",transactionHistoryStruct)
      responder(w,[]StructConfig.TransactionHistoryResponse{StructConfig.TransactionHistoryResponse{Response:"true",TransactionHistory:transactionHistoryStruct,ErrInResponse:""}})
    }
  }else{
    fmt.Println("Creadential map is empty...")
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  }
}

// Admin related operations
func UsersList(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    //fmt.Println("cred map : ",credMap)
    var err error
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := []StructConfig.UserDetails{}
    if credMap["filter"].(string) == "all" {
      if credMap["isByHRP"].(bool) {
        hrp,floatErr := strconv.ParseFloat(credMap["keyword"].(string),64)
        if floatErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          if credMap["expression"].(string) == "isGreaterThan" {
            err = collection.Find(bson.M{"user_role":"","hrp":bson.M{"$gt":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isLessThan" {
            err = collection.Find(bson.M{"user_role":"","hrp":bson.M{"$lt":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isEqualsTo" {
            err = collection.Find(bson.M{"user_role":"","hrp":hrp}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isNotEqual" {
            err = collection.Find(bson.M{"user_role":"","hrp":bson.M{"$ne":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          }
        }
      }else if credMap["isByUsername"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          err = collection.Find(bson.M{"user_role":"","user_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        } else if credMap["expression"].(string) == "isNotEqual" {
          err = collection.Find(bson.M{"user_role":"","user_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        }
      }else if credMap["isByName"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          err = collection.Find(bson.M{"user_role":"","personal_info.full_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        } else if credMap["expression"].(string) == "isNotEqual" {
          err = collection.Find(bson.M{"user_role":"","personal_info.full_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        }
      }else{
        err = collection.Find(bson.M{"user_role":""}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
      }

    }else if credMap["filter"].(string) == "Active"{
      if credMap["isByHRP"].(bool) {
        hrp,floatErr := strconv.ParseFloat(credMap["keyword"].(string),64)
        if floatErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          if credMap["expression"].(string) == "isGreaterThan"{
            err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string),"hrp":bson.M{"$gt":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isLessThan"{
            err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string),"hrp":bson.M{"$lt":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isEqualsTo"{
            err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string),"hrp":hrp}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isNotEqual"{
            err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string),"hrp":bson.M{"$ne":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          }
        }
      }else if credMap["isByUsername"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string),"user_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        } else if credMap["expression"].(string) == "isNotEqual" {
          err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string),"user_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        }
      }else if credMap["isByName"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo" {
          err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string),"personal_info.full_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        } else if credMap["expression"].(string) == "isNotEqual" {
          err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string),"personal_info.full_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        }
      }else{
        err = collection.Find(bson.M{"user_role":"","account_status":credMap["filter"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
      }

    }else if credMap["filter"].(string) == "Inactive"{
      if credMap["isByHRP"].(bool) {
        hrp,floatErr := strconv.ParseFloat(credMap["keyword"].(string),64)
        if floatErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          if credMap["expression"].(string) == "isGreaterThan"{
            err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"},"hrp":bson.M{"$gt":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isLessThan"{
            err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"},"hrp":bson.M{"$lt":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isEqualsTo"{
            err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"},"hrp":hrp}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          } else if credMap["expression"].(string) == "isNotEqual"{
            err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"},"hrp":bson.M{"$ne":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
          }
        }
      }else if credMap["isByUsername"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"},"user_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        } else if credMap["expression"].(string) == "isNotEqual" {
          err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"},"user_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        }
      }else if credMap["isByName"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo" {
          err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"},"personal_info.full_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        } else if credMap["expression"].(string) == "isNotEqual" {
          err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"},"personal_info.full_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
        }
      }else{
        err = collection.Find(bson.M{"user_role":"","account_status":bson.M{"$ne":"Active"}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
      }
    }
    if err != nil {
      fmt.Println("Error while fetching user list : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      responder(w,[]StructConfig.UsersListResponse{StructConfig.UsersListResponse{Response:"true",UsersList:userDetailsStruct,ErrInResponse:""}})
    }
  }
}

func WithdrawalRequestList(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    collection = setCollection("rozgar_db","userDetails_collection")
    //fmt.Println("cred map : ",credMap)
    var pipeErr error
    userDetailsStruct := []StructConfig.UserDetails{}
    m := []bson.M{}
    //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"transaction_history.transaction_for":"Withdrawal","transaction_history.status":"Processing"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":[]interface{}{bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Processing"}}}}}}}}}}).All(&m)
    if credMap["filter"].(string) == "all" {
      if credMap["isByHRP"].(bool) {
        hrp,floatErr := strconv.ParseFloat(credMap["keyword"].(string),64)
        if floatErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          if credMap["expression"].(string) == "isLessThan" {
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$gt":[]interface{}{"$$t.units",-hrp}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isGreaterThan" {
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$lt":[]interface{}{"$$t.units",-hrp}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isEqualsTo" {
            //err = collection.Find(bson.M{"user_role":"","hrp":hrp}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.units",-hrp}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isNotEqual" {
            //err = collection.Find(bson.M{"user_role":"","hrp":bson.M{"$ne":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$ne":[]interface{}{"$$t.units",-hrp}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          }
        }
      }else if credMap["isByUsername"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          //err = collection.Find(bson.M{"user_role":"","user_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          //pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"$and":[]bson.M{bson.M{"transaction_history.transaction_for":"Withdrawal"},bson.M{"user_name":credMap["keyword"].(string)}}}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":credMap["keyword"].(string),"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        } else if credMap["expression"].(string) == "isNotEqual" {
          //err = collection.Find(bson.M{"user_role":"","user_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":bson.M{"$ne":credMap["keyword"].(string)},"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        }
      }else if credMap["isByName"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          //err = collection.Find(bson.M{"user_role":"","personal_info.full_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          //pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"$and":[]bson.M{bson.M{"transaction_history.transaction_for":"Withdrawal"},bson.M{"personal_info.full_name":credMap["keyword"].(string)}}}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"personal_info.full_name":credMap["keyword"].(string)}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        } else if credMap["expression"].(string) == "isNotEqual" {
          //err = collection.Find(bson.M{"user_role":"","personal_info.full_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          //pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"$and":[]bson.M{bson.M{"transaction_history.transaction_for":"Withdrawal"},bson.M{"personal_info.full_name":bson.M{"$ne":credMap["keyword"].(string)}}}}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"personal_info.full_name":bson.M{"$ne":credMap["keyword"].(string)},"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        }
      }else{
        pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
      }
    }else if credMap["filter"].(string) == "accepted" {
      if credMap["isByHRP"].(bool) {
        hrp,floatErr := strconv.ParseFloat(credMap["keyword"].(string),64)
        if floatErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          if credMap["expression"].(string) == "isLessThan" {
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$gt":[]interface{}{"$$t.units",-hrp}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isGreaterThan" {
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$lt":[]interface{}{"$$t.units",-hrp}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isEqualsTo" {
            //err = collection.Find(bson.M{"user_role":"","hrp":hrp}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.units",-hrp}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isNotEqual" {
            //err = collection.Find(bson.M{"user_role":"","hrp":bson.M{"$ne":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$ne":[]interface{}{"$$t.units",-hrp}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          }
        }
      }else if credMap["isByUsername"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          //err = collection.Find(bson.M{"user_role":"","user_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":credMap["keyword"].(string),"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        } else if credMap["expression"].(string) == "isNotEqual" {
          //err = collection.Find(bson.M{"user_role":"","user_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":bson.M{"$ne":credMap["keyword"].(string)},"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        }
      }else if credMap["isByName"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          //err = collection.Find(bson.M{"user_role":"","personal_info.full_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"personal_info.full_name":credMap["keyword"].(string),"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        } else if credMap["expression"].(string) == "isNotEqual" {
          //err = collection.Find(bson.M{"user_role":"","personal_info.full_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"personal_info.full_name":bson.M{"$ne":credMap["keyword"].(string)},"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        }
      }else{
        pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Paid"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        //err = collection.Pipe([]bson.M{bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for",credMap["transactionType"].(string)}},bson.M{"$gte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["startDate"].(string))}},bson.M{"$lte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["endDate"].(string))}}}}}}}}}).All(&m)
      }
    }else if credMap["filter"].(string) == "declined" {
      if credMap["isByHRP"].(bool) {
        hrp,floatErr := strconv.ParseFloat(credMap["keyword"].(string),64)
        if floatErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          if credMap["expression"].(string) == "isLessThan" {
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$gt":[]interface{}{"$$t.units",-hrp}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isGreaterThan" {
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$lt":[]interface{}{"$$t.units",-hrp}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isEqualsTo" {
            //err = collection.Find(bson.M{"user_role":"","hrp":hrp}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.units",-hrp}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          } else if credMap["expression"].(string) == "isNotEqual" {
            //err = collection.Find(bson.M{"user_role":"","hrp":bson.M{"$ne":hrp}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
            condition := bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$ne":[]interface{}{"$$t.units",-hrp}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}
            pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":condition}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
          }
        }
      }else if credMap["isByUsername"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          //err = collection.Find(bson.M{"user_role":"","user_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":credMap["keyword"].(string),"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        } else if credMap["expression"].(string) == "isNotEqual" {
          //err = collection.Find(bson.M{"user_role":"","user_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_name":bson.M{"$ne":credMap["keyword"].(string)},"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        }
      }else if credMap["isByName"].(bool) {
        if credMap["expression"].(string) == "isEqualsTo"{
          //err = collection.Find(bson.M{"user_role":"","personal_info.full_name":credMap["keyword"].(string)}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"personal_info.full_name":credMap["keyword"].(string),"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        } else if credMap["expression"].(string) == "isNotEqual" {
          //err = collection.Find(bson.M{"user_role":"","personal_info.full_name":bson.M{"$ne":credMap["keyword"].(string)}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1}).All(&userDetailsStruct)
          pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"personal_info.full_name":bson.M{"$ne":credMap["keyword"].(string)},"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        }
      }else{
        pipeErr = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":""}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Withdrawal"}},bson.M{"$eq":[]interface{}{"$$t.status","Declined"}}}}}},"user_name":1,"personal_info.full_name":1,"user_id":1}}}).All(&m)
        //err = collection.Pipe([]bson.M{bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for",credMap["transactionType"].(string)}},bson.M{"$gte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["startDate"].(string))}},bson.M{"$lte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["endDate"].(string))}}}}}}}}}).All(&m)
      }
    }
    if pipeErr != nil {
      fmt.Println("Error while fetching withdrawal request list : ",pipeErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //fmt.Println("Map returns : ",m)
      b, errMar := json.Marshal(m)
      if errMar != nil {
        fmt.Println("error while marshal : ", errMar)
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
        //fmt.Fprintf(w, "%s", b)
      } else {
        errUnm := json.Unmarshal([]byte(b), &userDetailsStruct)
        if errUnm != nil {
          fmt.Println("Error while unmarshal Grand : ", errUnm)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
        }
      }
      responder(w,[]StructConfig.UsersListResponse{StructConfig.UsersListResponse{Response:"true",UsersList:userDetailsStruct,ErrInResponse:""}})
    }
  }
}

func GetSpecificUserDetails(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := StructConfig.UserDetails{}
    err := collection.Find(bson.M{"user_id":credMap["userId"].(string)}).Select(bson.M{"transaction_history":0}).One(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching Specific User Details : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //responder(w,[]StructConfig.UserDetailsResponse{StructConfig.UserDetailsResponse{Response:"true",UserDetails:userDetailsStruct,ErrInResponse:""}})
      responder(w,[]StructConfig.SpecificUserResponse{StructConfig.SpecificUserResponse{Response:"true",UserDetails:userDetailsStruct,SponsorFullName:getUserFullname(userDetailsStruct.SponsorUname),FollowerCounts:getFollowerCounts(userDetailsStruct.IdentityCode),TeamCounts:getTeamCounts(userDetailsStruct.IdentityCode),Earnings:getSpecificUsersEarnings(userDetailsStruct.IdentityCode),ErrInResponse:""}})
    }
  }
}

func getUserFullname(userName string)string{
  m := bson.M{}
  var userMap map[string]interface{}
  err := collection.Find(bson.M{"user_name":userName}).Select(bson.M{"personal_info.full_name":1}).One(&m)
  if err != nil {
    fmt.Println("Error while fetching full name : ",err)
    return ""
  }else{
    fmt.Println("M : ",m)
    b, errMar := json.Marshal(m["personal_info"])
    if errMar != nil {
      fmt.Println("error while marshal : ", errMar)
      //fmt.Fprintf(w, "%s", b)
    } else {
      errUnm := json.Unmarshal([]byte(b), &userMap)
      if errUnm != nil {
        fmt.Println("Error while unmarshal Grand : ", errUnm)
      }
    }
    return userMap["full_name"].(string)
  }
  return ""
}

func getSpecificUsersEarnings(identityCode string)StructConfig.Earnings{
  //TeamCounts:getTeamCounts(userDetailsStruct.IdentityCode)
  //totalEarnedAmount(userDetailsStruct.IdentityCode)
  //getFollowerCounts(userDetailsStruct.IdentityCode)
  m := []bson.M{}
  totalEarned := 0.0
  levelisedEarnings := []float64{0.0, 0.0, 0.0, 0.0, 0.0}
  err := collection.Pipe([]bson.M{bson.M{"$match":bson.M{"identity_code":identityCode}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Level Income"}}}}}}}).All(&m)
  if err != nil {
    fmt.Println("Error while getSpecificUsersEarnings amount : ",err)
  }else{
      if convertMap,convertErr := Underdog.InterfaceArrToMap(m[0]["transaction_history"]); convertErr != nil {
        fmt.Println("Error while InterfaceToArrMap totalEarnedAmount : ",convertErr)
      } else {
        for i := 0; i < len(convertMap); i++ {
          totalEarned = totalEarned + convertMap[i]["units"].(float64)
          levelisedEarnings[i] = convertMap[i]["units"].(float64)
        }
      }
  }
  return StructConfig.Earnings{Level1Earnings:levelisedEarnings[0],Level2Earnings:levelisedEarnings[1],Level3Earnings:levelisedEarnings[2],Level4Earnings:levelisedEarnings[3],Level5Earnings:levelisedEarnings[4],TotalEarnings:totalEarned}
}

func GenerateHRP(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      HrpAmt,floatErr := strconv.ParseFloat(credMap["hrp"].(string),64)
      if floatErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        if companyId,comIdErr := getCompanyId(); comIdErr != nil {
          fmt.Println("Error while getting companyId : ",comIdErr)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          //change := bson.M{"$inc": bson.M{"hrp": HrpAmt}}
          creditAmt := bson.M{"$inc": bson.M{"hrp": HrpAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:"company",From:credMap["userName"].(string),Units:HrpAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Generate HRP",BankTransId:"",Status:"Done",Level:""}}}
          actDetailsStruct := StructConfig.ActivityDetails{To:"",From:"",Amount:HrpAmt}
          adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["userName"].(string),ActivityFor:"Generate HRP",ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:"Company",ActivityDetails:actDetailsStruct}
          runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
          ops := []txn.Op{{
              //Adding new user instance
              C:      "userDetails_collection",
              Id:     companyId,
              //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
              Update: creditAmt,
            },
            {
              C:      "adminActivityLog_collection",
              Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
              //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
              Insert: adminLogStruct,
            },
          }
          id := bson.NewObjectId() // Optional
          runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"AddUserDetails"))
          if runnerErr != nil {
            fmt.Println("Error while runner : ",runnerErr)
            if resumeErr := runner.Resume(id); resumeErr != nil {
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
            }else{
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
            }
          }else{
            fmt.Println("Executed")
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }
      }
    }
  }
}

func RequestHRP(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      HrpAmt,floatErr := strconv.ParseFloat(credMap["hrp"].(string),64)
      if floatErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        if companyId,comIdErr := getCompanyId(); comIdErr != nil {
          fmt.Println("Error while getting companyId : ",comIdErr)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{
          //change := bson.M{"$inc": bson.M{"hrp": HrpAmt}}
          debitAmt := bson.M{"$inc": bson.M{"hrp": -HrpAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:credMap["userName"].(string),From:"company",Units:-HrpAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Request HRP",BankTransId:"",Status:"Done",Level:""}}}
          creditAmt := bson.M{"$inc": bson.M{"hrp": HrpAmt},"$push":bson.M{"transaction_history":&StructConfig.TransactionHistory{TransId:bson.ObjectId(bson.NewObjectId()).Hex(),To:credMap["userName"].(string),From:"company",Units:HrpAmt,TransactionOn:time.Now().UnixNano()/(int64(time.Millisecond)),TransactionFor:"Request HRP",BankTransId:"",Status:"Done",Level:""}}}
          actDetailsStruct := StructConfig.ActivityDetails{To:credMap["userName"].(string),From:"Company",Amount:HrpAmt}
          adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["userName"].(string),ActivityFor:"Request HRP",ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:"Company",ActivityDetails:actDetailsStruct}
          runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
          ops := []txn.Op{{
              //Adding new user instance
              C:      "userDetails_collection",
              Id:     companyId,
              Assert: bson.M{"hrp": bson.M{"$gte": HrpAmt}},
              Update: debitAmt,
            },
            {
              //Adding new user instance
              C:      "userDetails_collection",
              Id:     credMap["userId"].(string),
              //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
              Update: creditAmt,
            },
            {
              C:      "adminActivityLog_collection",
              Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
              //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
              Insert: adminLogStruct,
            },
          }
          id := bson.NewObjectId() // Optional
          runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"AddUserDetails"))
          if runnerErr != nil {
            fmt.Println("Error while runner : ",runnerErr)
            if resumeErr := runner.Resume(id); resumeErr != nil {
                responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
            }else{
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
            }
          }else{
            fmt.Println("Executed")
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }
      }
    }
  }
}

func ChangeUserPassword(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    if userId,uInstanceErr := getUserInstanceId(credMap["userName"].(string)); uInstanceErr != nil {
      fmt.Println("Error while fetching user instance id : ",uInstanceErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else{
      if isUserExist,isUserExistErr := isUserExist(credMap["adminUserName"].(string),credMap["adminPassword"].(string)); isUserExistErr != nil {
        fmt.Println("Error while checking user exist or not : ",isUserExistErr)
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else if isUserExist {
        change := bson.M{"$set":bson.M{"password": credMap["password"].(string)}}
        actDetailsStruct := StructConfig.ActivityDetails{To:"",From:"",Amount:0.0}
        adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["adminUserName"].(string),ActivityFor:"Change Password",ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:credMap["userName"].(string),ActivityDetails:actDetailsStruct}
        runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
        ops := []txn.Op{{
            //Adding new user instance
            C:      "userInstance_collection",
            Id:     userId,
            //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
            Update: change,
          },
          {
            C:      "adminActivityLog_collection",
            Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
            //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
            Insert: adminLogStruct,
          },
        }
        id := bson.NewObjectId() // Optional
        runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"AddUserDetails"))
        if runnerErr != nil {
          fmt.Println("Error while runner : ",runnerErr)
          if resumeErr := runner.Resume(id); resumeErr != nil {
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
          }else{
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }else{
          fmt.Println("Executed")
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
        }
      }else{
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Incorrect admin password."}})
      }
    }
  }
}

func getUserInstanceId(userName string)(string,error){
  collection = setCollection("rozgar_db","userInstance_collection")
  m := bson.M{}
  err := collection.Find(bson.M{"user_name":userName}).Select(bson.M{"_id":1}).One(&m)
  if err != nil {
    fmt.Println("Error while fetching user instance id : ",err)
    return "",err
  }else{
    return m["_id"].(string),nil
  }
}

func UpdateUserStatus(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    if isUserExist,isUserExistErr := isUserExist(credMap["adminUserName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      change := bson.M{"$set":bson.M{"account_status": credMap["status"].(string)}}
      actDetailsStruct := StructConfig.ActivityDetails{To:"",From:"",Amount:0.0}
      adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["adminUserName"].(string),ActivityFor:"User status "+credMap["status"].(string),ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:credMap["userName"].(string),ActivityDetails:actDetailsStruct}
      runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
      ops := []txn.Op{{
          //Adding new user instance
          C:      "userDetails_collection",
          Id:     credMap["userId"].(string),
          //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
          Update: change,
        },
        {
          C:      "adminActivityLog_collection",
          Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
          //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
          Insert: adminLogStruct,
        },
      }
      id := bson.NewObjectId() // Optional
      runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"UpdateUserProfile"))
      if runnerErr != nil {
        fmt.Println("Error while runner : ",runnerErr)
        if resumeErr := runner.Resume(id); resumeErr != nil {
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
        }else{

          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
        }
      }else{
        fmt.Println("Executed")
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
      }
    }else{
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Incorrect Password."}})
    }
  }
}

func ValidateUser(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  }else{
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := StructConfig.UserDetails{}
    err := collection.Find(bson.M{"user_name":credMap["userName"].(string)}).Select(bson.M{"user_id":1,"identity_code":1,"personal_info":1,"account_status":1,"hrp":1}).One(&userDetailsStruct); if err != nil {
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Incorrect Username."}})
    }else{
      responder(w,[]StructConfig.UserDetailsResponse{StructConfig.UserDetailsResponse{Response:"true",UserDetails:userDetailsStruct,ErrInResponse:""}})
    }
  }
}

func GetAdminList(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    fmt.Println("creadMap : ",credMap)
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := []StructConfig.UserDetails{}
    err := collection.Find(bson.M{"user_role":"admin"}).Select(bson.M{"user_id":1,"user_name":1,"personal_info.full_name":1,"identity_code":1}).All(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching user details list : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //fmt.Println(userDetailsStruct)
      responder(w,[]StructConfig.UsersListResponse{StructConfig.UsersListResponse{Response:"true",UsersList:userDetailsStruct,ErrInResponse:""}})
    }
  }
}

func GetAdminActivityLogs(w http.ResponseWriter, r *http.Request,interfaceName string){
  /*credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    fmt.Println("creadMap : ",credMap)*/
    collection = setCollection("rozgar_db","adminActivityLog_collection")
    adminLogStruct := []StructConfig.AdminActivityLogs{}
    err := collection.Find(bson.M{}).All(&adminLogStruct)
    if err != nil {
      fmt.Println("Error while fetching user details list : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //fmt.Println(userDetailsStruct)
      responder(w,[]StructConfig.AdminLogListResponse{StructConfig.AdminLogListResponse{Response:"true",AdminActivityLogs:adminLogStruct,ErrInResponse:""}})
    }
  //}
}

func SetBroadcast(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    if isUserExist,isUserExistErr := isUserExist(credMap["userName"].(string),credMap["password"].(string)); isUserExistErr != nil {
      fmt.Println("Error while checking user exist or not : ",isUserExistErr)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
    }else if isUserExist {
      if disableBroadcast() {
        //change := bson.M{"$inc": bson.M{"hrp": HrpAmt}}
        broadcastDetailsStruct := StructConfig.BroadcastDetails{BroadMsg:credMap["msg"].(string),BroadBy:credMap["userName"].(string),BroadOn:time.Now().UnixNano()/(int64(time.Millisecond)),BroadReason:credMap["reason"].(string),BroadStatus:"Active"}
        actDetailsStruct := StructConfig.ActivityDetails{To:"",From:"",Amount:0}
        adminLogStruct := StructConfig.AdminActivityLogs{ActivityBy:credMap["userName"].(string),ActivityFor:"Message Broadcasted",ActivityOn:time.Now().UnixNano()/(int64(time.Millisecond)),ActivityStatus:"Done",ActivityPerformedOn:"Company",ActivityDetails:actDetailsStruct}
        runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
        ops := []txn.Op{{
            C:      "broadcast_collection",
            Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
            //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
            Insert: broadcastDetailsStruct,
          },
          {
            C:      "adminActivityLog_collection",
            Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
            //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
            Insert: adminLogStruct,
          },
        }
        id := bson.NewObjectId() // Optional
        runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"AddUserDetails"))
        if runnerErr != nil {
          fmt.Println("Error while runner : ",runnerErr)
          if resumeErr := runner.Resume(id); resumeErr != nil {
              responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
          }else{
            data := make(map[string]string)
            data["title"] = "Happy Rozgar"
            data["body"] = "New Broadcast Message."
            _,fbTokens:= getAdminsFBToken();
            //var ids = []string{fbTokens}
            sendNotificationWithFCM(fbTokens,data)
            responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
          }
        }else{
          fmt.Println("Executed")
          data := make(map[string]string)
          data["title"] = "Happy Rozgar"
          data["body"] = "New Broadcast Message."
          _,fbTokens:= getAdminsFBToken();
          //var ids = []string{fbTokens}
          sendNotificationWithFCM(fbTokens,data)
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
        }
      }
    }
  }
}

func UnsetBroadcast(w http.ResponseWriter, r *http.Request,interfaceName string){
  if disableBroadcast() {
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
  }else{
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
  }
}


func disableBroadcast()bool{
  collection = setCollection("rozgar_db","broadcast_collection")
  err := collection.Update(bson.M{"broad_status":"Active"},bson.M{"$set":bson.M{"broad_status":"Inactive"}})
  if err != nil {
    if err.Error() == "not found" {
      return true
    }else{
      fmt.Println("Failed to disable broadcast : ",err)
      return false
    }
  }else{
    return true
  }
}

func SetFeedback(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    feedbackDetailsStruct := StructConfig.FeedbackDetails{FeedMsg:credMap["feedbackMsg"].(string),FeedBy:credMap["userName"].(string),FeedOn:time.Now().UnixNano()/(int64(time.Millisecond)),FeedReason:credMap["reason"].(string),FeedType:credMap["from"].(string)}
    runner := txn.NewRunner(setCollection("rozgar_db","transaction_collection"))
    ops := []txn.Op{{
        C:      "feedback_collection",
        Id:     bson.ObjectId(bson.NewObjectId()).Hex(),
        //Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
        Insert: feedbackDetailsStruct,
      },
    }
    id := bson.NewObjectId() // Optional
    runnerErr := runner.Run(ops, id,setInfoStruct(credMap["userName"].(string),"SetFeedback"))
    if runnerErr != nil {
      fmt.Println("Error while runner : ",runnerErr)
      if resumeErr := runner.Resume(id); resumeErr != nil {
          responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong"}})
      }else{
        data := make(map[string]string)
        data["title"] = "Happy Rozgar"
        data["body"] = "New Feedback/Ticket Message."
        _,fbTokens:= getAdminsFBToken();
        //var ids = []string{fbTokens}
        sendNotificationWithFCM(fbTokens,data)
        responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
      }
    }else{
      fmt.Println("Executed")
      data := make(map[string]string)
      data["title"] = "Happy Rozgar"
      data["body"] = "New Feedback/Ticket Message."
      _,fbTokens:= getAdminsFBToken();
      //var ids = []string{fbTokens}
      sendNotificationWithFCM(fbTokens,data)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
    }

  }
}

func GetFeedback(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    //fmt.Println("creadMap : ",credMap)
    collection = setCollection("rozgar_db","feedback_collection")
    feedbackDetailsStruct := []StructConfig.FeedbackDetails{}
    err := collection.Find(bson.M{"feed_type":credMap["feedType"].(string)}).Select(bson.M{"feed_msg":1,"feed_by":1,"feed_on":1,"feed_reason":1,"feed_type":1}).All(&feedbackDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching user details list : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      //fmt.Println(userDetailsStruct)
      responder(w,[]StructConfig.FeedbackDetailsList{StructConfig.FeedbackDetailsList{Response:"true",FeedbackDetailsList:feedbackDetailsStruct,ErrInResponse:""}})
    }
  }
}

func GetCompanyEarningRecords(w http.ResponseWriter, r *http.Request,interfaceName string){
  if credMap,credErr := Underdog.StringToMap(interfaceName); credErr != nil {
    fmt.Println("Error while interface to map @GetCompanyEarningRecords : ",credErr)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  } else if len(credMap) > 0 {
    collection = setCollection("rozgar_db","userDetails_collection")
    //userDetailsStruct := StructConfig.UserDetails{}
    transactionHistoryStruct := []StructConfig.TransactionHistory{}
    var err error
    //Generate HRP Request HRP
    //fmt.Println("credMap : ",credMap)
    m := []bson.M{}
    if credMap["transactionType"].(string) == "Others" {
      err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Bank Details Update"}},bson.M{"$gt":[]interface{}{"$$t.units",0}}}}}}}}}).All(&m)
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Bank Details Update"}}}}}}}).All(&m)
    }else if credMap["transactionType"].(string) == "All" {
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":1}}}).All(&m)
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":[]interface{}{bson.M{"$and":[]bson.M{bson.M{"$ne":[]interface{}{"$$t.transaction_for","Generate HRP"}},bson.M{"$ne":[]interface{}{"$$t.transaction_for","Request HRP"}},bson.M{"$gt":[]interface{}{"$$t.units",0}}}}}}}}}}).All(&m)
      err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$gt":[]interface{}{"$$t.units",0}}}}}}}).All(&m)
    }else{
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":[]interface{}{bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for",credMap["transactionType"].(string)}},bson.M{"$gt":[]interface{}{"$$t.units",0}}}}}}}}}}).All(&m)
      err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for",credMap["transactionType"].(string)}},bson.M{"$gt":[]interface{}{"$$t.units",0}}}}}}}}}).All(&m)
    }

    /*b, errMar := json.Marshal(m[0]["transaction_history"])
    if errMar != nil {
      fmt.Println("error while marshal : ", errMar)
      //fmt.Fprintf(w, "%s", b)
    } else {
      errUnm := json.Unmarshal([]byte(b), &transactionHistoryStruct)
      if errUnm != nil {
        fmt.Println("Error while unmarshal Grand : ", errUnm)
      }
    }*/

    if err != nil {
      fmt.Println("Error while fetching transaction history : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      b, errMar := json.Marshal(m[0]["transaction_history"])
      if errMar != nil {
        fmt.Println("error while marshal : ", errMar)
        //fmt.Fprintf(w, "%s", b)
      } else {
        errUnm := json.Unmarshal([]byte(b), &transactionHistoryStruct)
        if errUnm != nil {
          fmt.Println("Error while unmarshal Grand : ", errUnm)
        }
      }
      responder(w,[]StructConfig.TransactionHistoryResponse{StructConfig.TransactionHistoryResponse{Response:"true",TransactionHistory:transactionHistoryStruct,ErrInResponse:""}})
    }
  }else{
    fmt.Println("Credential map is empty...")
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  }
}

func GetReports(w http.ResponseWriter, r *http.Request,interfaceName string){
  if credMap,credErr := Underdog.StringToMap(interfaceName); credErr != nil {
    fmt.Println("Error while interface to map @GetCompanyEarningRecords : ",credErr)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  } else if len(credMap) > 0 {
    collection = setCollection("rozgar_db","userDetails_collection")
    //userDetailsStruct := StructConfig.UserDetails{}
    grandTransactionHistoryStruct := []StructConfig.TransactionHistory{}
    transactionHistoryStruct := []StructConfig.TransactionHistory{}
    var err error
    //Generate HRP Request HRP
    //fmt.Println("credMap : ",credMap)
    m := []bson.M{}
    /*if credMap["transactionType"].(string) == "Others" {
      err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Bank Details Update"}},bson.M{"$gt":[]interface{}{"$$t.units",0}}}}}}}}}).All(&m)
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$eq":[]interface{}{"$$t.transaction_for","Bank Details Update"}}}}}}}).All(&m)
    }else if credMap["transactionType"].(string) == "FundTransfer" {
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":1}}}).All(&m)
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":[]interface{}{bson.M{"$and":[]bson.M{bson.M{"$ne":[]interface{}{"$$t.transaction_for","Generate HRP"}},bson.M{"$ne":[]interface{}{"$$t.transaction_for","Request HRP"}},bson.M{"$gt":[]interface{}{"$$t.units",0}}}}}}}}}}).All(&m)
      err = collection.Pipe([]bson.M{bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for","Fund Transfer"}},bson.M{"$gte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["startDate"].(string))}},bson.M{"$lte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["endDate"].(string))}}}}}}}}}).All(&m)
    }else{
      //err = collection.Pipe([]bson.M{bson.M{"$match":bson.M{"user_role":"company"}},bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":[]interface{}{bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for",credMap["transactionType"].(string)}},bson.M{"$gt":[]interface{}{"$$t.units",0}}}}}}}}}}).All(&m)
      err = collection.Pipe([]bson.M{bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for",credMap["transactionType"].(string)}},bson.M{"$gte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["startDate"].(string))}},bson.M{"$lte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["endDate"].(string))}}}}}}}}}).All(&m)
    }*/
    fmt.Println("Transaction Type : ",credMap["transactionType"].(string))
    err = collection.Pipe([]bson.M{bson.M{"$project":bson.M{"transaction_history":bson.M{"$filter":bson.M{"input":"$transaction_history","as":"t","cond":bson.M{"$and":[]bson.M{bson.M{"$eq":[]interface{}{"$$t.transaction_for",credMap["transactionType"].(string)}},bson.M{"$gte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["startDate"].(string))}},bson.M{"$lte":[]interface{}{"$$t.transaction_on",stringToInt64(credMap["endDate"].(string))}}}}}},"_id":0}}}).All(&m)
    if err != nil {
      fmt.Println("Error while fetching transaction history : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      for i := 0; i < len(m); i++ {
        //if m[i]["transaction_history"] != nil {
          //fmt.Println("Index : ",i,">>",m[i])
          b, errMar := json.Marshal(m[i]["transaction_history"])
          if errMar != nil {
            fmt.Println("error while marshal : ", errMar)
            //fmt.Fprintf(w, "%s", b)
          } else {
            errUnm := json.Unmarshal([]byte(b), &transactionHistoryStruct)
            if errUnm != nil {
              fmt.Println("Error while unmarshal Grand : ", errUnm)
            }else{
              for ind := 0; ind < len(transactionHistoryStruct); ind++ {
                grandTransactionHistoryStruct = append(grandTransactionHistoryStruct,transactionHistoryStruct[ind])
              }
            }
          }
        //}
      }
      //fmt.Println("transactionHistoryStruct : ",grandTransactionHistoryStruct)
      responder(w,[]StructConfig.TransactionHistoryResponse{StructConfig.TransactionHistoryResponse{Response:"true",TransactionHistory:grandTransactionHistoryStruct,ErrInResponse:""}})
    }
  }else{
    fmt.Println("Creadential map is empty...")
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
  }
}

func stringToInt64(s string)int64{
  i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
    		//panic(err)
        fmt.Println("Error while stringToInt64 : ",err)
        return 0
	}
  //fmt.Println(i)
  return i
}

func GetNewUserReports(w http.ResponseWriter, r *http.Request,interfaceName string){
  credMap,err := Underdog.StringToMap(interfaceName)
  if err != nil {
    fmt.Println("error while converting string to map : ",err)
    responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Error while converting string to map"}})
  } else {
    var err error
    collection = setCollection("rozgar_db","userDetails_collection")
    userDetailsStruct := []StructConfig.UserDetails{}
    err = collection.Find(bson.M{"user_role":"","user_added_on":bson.M{"$gte":stringToInt64(credMap["startDate"].(string)),"$lte":stringToInt64(credMap["endDate"].(string))}}).Select(bson.M{"user_id":1,"user_name":1,"sponsor_uname":1,"hrp":1,"account_status":1,"user_added_on":1,"personal_info.full_name":1,"identity_code":1,"direct_child_count":1}).All(&userDetailsStruct)
    if err != nil {
      fmt.Println("Error while fetching new user list reports : ",err)
      responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:"Something went wrong... try again"}})
    }else{
      responder(w,[]StructConfig.UsersListResponse{StructConfig.UsersListResponse{Response:"true",UsersList:userDetailsStruct,ErrInResponse:""}})
    }
  }
}


//Change user password

//Suspend User

//User Transaction history

// User Member History (LevelWise)

//Member Stats

//Mamber Earning levelk wise
