package Underdog
import(
  "fmt"
  "encoding/json"
  "strings"
  //"packages/StructConfig"
)

func StringReplacer(str string)string{
  return strings.Replace(str, "=", ":", -1)
}

func StringToMap(interfaceName string)(map[string]interface{},error){
  //fmt.Println("In interfaceToMap : ",interfaceName)
  customizedMap := make(map[string]interface{})
  /*b, errInMarshal := json.Marshal(interfaceName)
  if errInMarshal != nil {
    fmt.Println("Error in marshal : ",errInMarshal)
  }else{*/
    //fmt.Println("Json : ",string(b))
    err := json.Unmarshal([]byte(interfaceName), &customizedMap)
    if err != nil {
      fmt.Println("Err in unmarsal : ",err)
      return customizedMap,err
    }
  //}
  return customizedMap,nil
}

func StringArrToMap(interfaceName string)([]map[string]interface{},error){
  var eventMapArr []map[string]interface{}
  err := json.Unmarshal([]byte(interfaceName),&eventMapArr)
  if err != nil {
    fmt.Println("Err in unmarsal : ",err)
    return eventMapArr,err
  }
  //fmt.Println("Event map arr : ",eventMapArr)
  return eventMapArr,nil
}

func InterfaceToMap(interfaceName interface{})(map[string]interface{},error){
  customizedMap := make(map[string]interface{})
  b, errInMarshal := json.Marshal(interfaceName)
  if errInMarshal != nil {
    fmt.Println("Error in marshal : ",errInMarshal)
  }else{
    //fmt.Println("Json marshalled string : ",string(b))
    err := json.Unmarshal([]byte(b), &customizedMap)
    if err != nil {
      fmt.Println("Err in unmarsal : ",err)
      return customizedMap,err
    }
  }
  return customizedMap,nil
}

func InterfaceArrToMap(interfaceName interface{})([]map[string]interface{},error){
  var customizedMapArr []map[string]interface{}

  b, errInMarshal := json.Marshal(interfaceName)
  if errInMarshal != nil {
    fmt.Println("Error in marshal : ",errInMarshal)
  }else{
    //fmt.Println("Json marshalled string : ",string(b))
    err := json.Unmarshal([]byte(b), &customizedMapArr)
    if err != nil {
      fmt.Println("Err in unmarsal : ",err)
      return customizedMapArr,err
    }
  }
  return customizedMapArr,nil
}

/*func InterfaceToStruct(interfaceName interface{}, structObj StructConfig.Event, structName string) StructConfig.Event {
	var errUnm error
	//fmt.Println("Interface at itGrand : ", interfaceName)
	b, errMar := json.Marshal(interfaceName)
	if errMar != nil {
		fmt.Println("error while marshal : ", errMar)
		//fmt.Fprintf(w, "%s", b)
	} else {
		if structName == "Host" {
			errUnm = json.Unmarshal([]byte(b), &structObj.Host_details)
		} else if structName == "Moment" {
      errUnm = json.Unmarshal([]byte(b),&structObj.Moments)
    }
		if errUnm != nil {
			fmt.Println("Error while unmarshal Grand : ", errUnm)
		} else {
			//fmt.Println("Unmarshaled Grand Struct : ", structObj)
			return structObj
		}
	}
	return structObj
}*/

func InterfaceToArr(interfaceName interface{})([]string,error){
  var retArr []string
  b, errInMarshal := json.Marshal(interfaceName)
  if errInMarshal != nil {
    fmt.Println("Error in marshal : ",errInMarshal)
  }else{
    //fmt.Println("Json marshalled string : ",string(b))
    err := json.Unmarshal([]byte(b), &retArr)
    if err != nil {
      fmt.Println("Err in unmarsal : ",err)
      return retArr,err
    }
  }
  return retArr,nil
}

func validator(){

}

//Check empty struct or not (reflect package needed)
/*if reflect.DeepEqual(eventStructObj, []StructConfig.Event{}) {
  responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"false",ErrInResponse:""}})
}else {
  responder(w,[]StructConfig.SingleResponse{StructConfig.SingleResponse{Response:"true",ErrInResponse:""}})
}*/
