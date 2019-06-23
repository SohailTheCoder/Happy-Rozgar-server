package TokenManager

import (
    "fmt"
    jwt "github.com/dgrijalva/jwt-go"
    "packages/StructConfig"
)

func GenerateToken(userDetailsStruct StructConfig.UserDetails)(string,error){
  token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &StructConfig.JwtStruct{
    Username: userDetailsStruct.Username,
    SponsorUname: userDetailsStruct.SponsorUname,
    AccountStatus: userDetailsStruct.AccountStatus,
    FullName : userDetailsStruct.PersonalInfo.FullName,
    MobileNumber : userDetailsStruct.PersonalInfo.MobileNumber,
  })
  tokenstring, err := token.SignedString([]byte("write_some_secret_key_here"))
  if err != nil {
    fmt.Println("Error while generate token : ",err)
    return "",err
  } else {
    return tokenstring,nil
  }
  //return "",nil
}

func DecodeToken(tokenString string)(interface{},error){
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte("write_some_secret_key_here"), nil
  })
  if err != nil {
    fmt.Println("Error while DecodeToken : ",err)
    return "",err
  } else {
    return token.Claims,nil
  }
  //return make(map[string]interface{}),nil
}

func IsTokenValid(tokenString string)(bool,error){
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte("write_some_secret_key_here"), nil
  })
  // When using `Parse`, the result `Claims` would be a map.
  if err != nil {
    fmt.Println("Error while IsTokenValid : ",err)
    return false,err
  } else {
    return token.Valid,nil
  }
}
