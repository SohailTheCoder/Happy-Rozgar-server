package StructConfig

import(
	jwt "github.com/dgrijalva/jwt-go"

  //"time"
  //"gopkg.in/mgo.v2/bson"
)

type ArgumentList struct {
	Function string `json:"function" bson:"function"`
	Arguments string `json:"arguments" bson:"arguments"`
}

type InfoStruct struct {
	TransBy string `json:"trans_by" bson:"trans_by"`
	TransFor string `json:"trans_for" bson:"trans_for"`
	TransOn string `json:"trans_on" bson:"trans_on"`
}

type UserInstance struct {
	Username string `json:"user_name" bson:"user_name"`
	Password string `json:"password" bson:"password"`
	MobileNo string `json:"mobile_no" bson:"mobile_no"`
}
//UserDetails start
type UserDetails struct {
	//Id string `json:"_id" bson: "_id"`
	UserId string `json:"user_id" bson:"user_id"`
	Username string `json:"user_name" bson:"user_name"`
	UserRole string `json:"user_role" bson:"user_role"`
	IdentityCode string `json:"identity_code" bson:"identity_code"`
	FirebaseToken string `json:"firebase_token" bson:"firebase_token"`
	SponsorUname string `json:"sponsor_uname" bson:"sponsor_uname"`
	HRP float64 `json:"hrp" bson:"hrp"`
	AccountStatus string `json:"account_status" bson:"account_status"`
	UserAddedOn string `json:"user_added_on" bson:"user_added_on"`
	PersonalInfo PersonalInfo `json:"personal_info" bson:"personal_info"`
	BankDetails BankDetails `json:"bank_details" bson:"bank_details"`
	DirectChildCount int `json:"direct_child_count" bson:"direct_child_count"`
	TransactionHistory []TransactionHistory `json:"transaction_history" bson:"transaction_history"`
}

type PersonalInfo struct {
	FullName string `json:"full_name" bson:"full_name"`
	MobileNumber string `json:"mobile_number" bson:"mobile_number"`
	Gender string `json:"gender" bson:"gender"`
	DOB string `json:"dob" bson:"dob"`
}

type BankDetails struct {
	AccountNumber string `json:"account_number" bson:"account_number"`
	//BankName string `json:"bank_name" bson:"bank_name"`
	IFSCode string `json:"ifs_code" bson:"ifs_code"`
	//BankAddress string `json:"bank_address" bson:"bank_address"`
}

type TransactionHistory struct {
	//Id string `json:"_id" bson:"_id"`
	TransId string `json:"trans_id" bson:"trans_id"`
	To string `json:"to" bson:"to"`
	From string `json:"from" bson:"from"`
	Units float64 `json:"units" bson:"units"`
	TransactionOn string `json:"transaction_on" bson:"transaction_on"`
	TransactionFor string `json:"transaction_for" bson:"transaction_for"`
	BankTransId string `json:"bank_trans_id" bson:"bank_trans_id"`
	Status string `json:"status" bson:"status"`
	Level string `json:"level" bson:"level"`
	//TypeOfTransaction string `json:"type_of_transaction" bson:"type_of_transaction"` //Credit/Withdrawal(in/out)
}
//UserDetails End ...
//Admin related structs
type AdminActivityLogs struct {
	ActivityBy string `json:"activity_by" bson:"activity_by"`
	ActivityFor string `json:"activity_for" bson:"activity_for"`
	ActivityOn string `json:"activity_on" bson:"activity_on"`
	ActivityStatus string `json:"activity_status" bson:"activity_status"`
	ActivityPerformedOn string `json:"activity_performed_on" bson:"activity_performed_on"`
	ActivityDetails ActivityDetails `json:"activity_details" bson:"activity_details"`
}

type ActivityDetails struct {
	To string `json:"to" bson:"to"`
	From string `json:"from" bson:"from"`
	Amount float64 `json:"amount" bson:"amount"`
}

//Admin related structs end
type BroadcastDetails struct {
	BroadMsg string `json:"broad_msg" bson:"broad_msg"`
	BroadBy string `json:"broad_by" bson:"broad_by"`
	BroadOn string `json:"broad_on" bson:"broad_on"`
	BroadReason string `json:"broad_reason" bson:"broad_reason"`
	BroadStatus string `json:"broad_status" bson:"broad_status"`
}
// jwt res struct
type JwtStruct struct {
	Username string `json:"user_name"`
	SponsorUname string `json:"sponsor_uname"`
	AccountStatus string `json:"account_status"`
	FullName string `json:"full_name"`
	MobileNumber string `json:"mobile_number"`
	jwt.StandardClaims
}

//Response structs

type SingleResponse struct {
	Response string `json:"response"`
	ErrInResponse string `json:"errInResponse"`
}

type LoginResponse struct {
	Response string `json:"response"`
	UserDetails UserDetails `json:"user_details"`
	TokenString string `json:"tokenString"`
	FollowerCounts FollowerCounts `json:"follower_counts"`
	TeamCounts TeamCounts `json:"team_counts"`
	Earned float64 `json:"earned" bson:"earned"`
	BroadcastDetails BroadcastDetails `json:"broadcast_details"`
	ErrInResponse string `json:"errInResponse"`
}

type LoginAdminResponse struct {
	Response string `json:"response"`
	UserDetails UserDetails `json:"user_details"`
	TokenString string `json:"tokenString"`
	AdminCounts AdminCounts `json:"admin_counts"`
	BroadcastDetails BroadcastDetails `json:"broadcast_details"`
	ErrInResponse string `json:"errInResponse"`
}

type AdminCounts struct {
	AdminActiveUsers int `json:"admin_active_users"`
	AdminInactiveUsers int `json:"admin_inactive_users"`
	AdminActivitiesCount int `json:"admin_activities_count"`
	AdminTotalEarnings float64 `json:"admin_total_earnings"`
	CompanyActiveUsers int `json:"company_active_users"`
	CompanyInactiveUsers int `json:"company_inactive_users"`
	AllWithdrawalRequest int `json:"all_withdrawal_request"`
	CompanyNewJoineeMonth int `json:"company_new_joinee_month"`
	CompanyTotalEarnings float64 `json:"company_total_earnings"`
	CompanyBalance float64 `json:"company_balance"`
	CompanyFollowerCounts FollowerCounts `json:"company_follower_counts"`
}

type RefreshResponse struct {
	Response string `json:"response"`
	UserDetails UserDetails `json:"user_details"`
	FollowerCounts FollowerCounts `json:"follower_counts"`
	TeamCounts TeamCounts `json:"team_counts"`
	Earned float64 `json:"earned" bson:"earned"`
	BroadcastDetails BroadcastDetails `json:"broadcast_details"`
	ErrInResponse string `json:"errInResponse"`
}

type RefreshAdminResponse struct {
	Response string `json:"response"`
	UserDetails UserDetails `json:"user_details"`
	AdminCounts AdminCounts `json:"admin_counts"`
	BroadcastDetails BroadcastDetails `json:"broadcast_details"`
	ErrInResponse string `json:"errInResponse"`
}

type SpecificUserResponse struct {
	Response string `json:"response"`
	UserDetails UserDetails `json:"user_details"`
	SponsorFullName string `json:"sponsor_full_name"`
	FollowerCounts FollowerCounts `json:"follower_counts"`
	TeamCounts TeamCounts `json:"team_counts"`
	Earnings Earnings `json:"earnings" bson:"earnings"`
	ErrInResponse string `json:"errInResponse"`
}

type Earnings struct {
	Level1Earnings float64 `json:"level1_earnings" bson:"level1_earnings"`
	Level2Earnings float64 `json:"level2_earnings" bson:"level2_earnings"`
	Level3Earnings float64 `json:"level3_earnings" bson:"level3_earnings"`
	Level4Earnings float64 `json:"level4_earnings" bson:"level4_earnings"`
	Level5Earnings float64 `json:"level5_earnings" bson:"level5_earnings"`
	TotalEarnings float64 `json:"total_earnings" bson:"total_earnings"`
}


type FollowerCounts struct {
	Level1Count int `json:"level1_count" bson:"level1_count"`
	Level2Count int `json:"level2_count" bson:"level2_count"`
	Level3Count int `json:"level3_count" bson:"level3_count"`
	Level4Count int `json:"level4_count" bson:"level4_count"`
	Level5Count int `json:"level5_count" bson:"level5_count"`
}

type TeamCounts struct {
	ActiveMembersCount int `json:"active_members_count" bson:"active_members_count"`
	NonActiveMembersCount int `json:"non_active_members_count" bson:"non_active_members_count"`
	TotalMembersCount int `json:"total_members_count" bson:"total_members_count"`
}

type SponsorListResponse struct {
	Response string `json:"response"`
	SponsorList []UserDetails `json:"sponsorList" bson:"sponsorList"`
	ErrInResponse string `json:"errInResponse"`
}

type FollowersListResponse struct {
	Response string `json:"response"`
	FollowerList []UserDetails `json:"followerList" bson:"followerList"`
	ErrInResponse string `json:"errInResponse"`
}

type BankDetailsResponse struct {
	Response string `json:"response"`
	BankDetails BankDetails `json:"bankDetails" bson:"bankDetails"`
	ErrInResponse string `json:"errInResponse"`
}

type TransactionHistoryResponse struct {
	Response string `json:"response"`
	TransactionHistory []TransactionHistory `json:"transactionHistory" bson:"transactionHistory"`
	ErrInResponse string `json:"errInResponse"`
}

type UsersListResponse struct {
	Response string `json:"response"`
	UsersList []UserDetails `json:"users_list" bson:"users_list"`
	ErrInResponse string `json:"errInResponse"`
}

type UserDetailsResponse struct {
	Response string `json:"response"`
	UserDetails UserDetails `json:"user_details" bson:"user_details"`
	ErrInResponse string `json:"errInResponse"`
}

type AdminLogListResponse struct {
	Response string `json:"response"`
	AdminActivityLogs []AdminActivityLogs `json:"admin_activity_logs" bson:"admin_activity_logs"`
	ErrInResponse string `json:"errInResponse"`
}

type Uplinks struct { //Set sponsor ids by there level
	Level1Id string `json:"level_1_id" bson:"level_1_id"`
	Level2Id string `json:"level_2_id" bson:"level_2_id"`
	Level3Id string `json:"level_3_id" bson:"level_3_id"`
	Level4Id string `json:"level_4_id" bson:"level_4_id"`
	Level5Id string `json:"level_5_id" bson:"level_5_id"`
}

//LevelDoc started ...
type LevelDoc struct {
	//Id string `json:"_id" bson: "_id"`
	Username string `json:"user_name" bson:"user_name"`
	Level1 []Levels `json:"level1" bson:"level1"`
	Level2 []Levels `json:"level2" bson:"level2"`
	Level3 []Levels `json:"level3" bson:"level3"`
	Level4 []Levels `json:"level4" bson:"level4"`
	Level5 []Levels `json:"level5" bson:"level5"`
}

type Levels struct {
	Uname string `json:"u_name" bson:"u_name"`
	FullName string `json:"full_name" bson:"full_name"`
	UserAddedOn string `json:"user_added_on" bson:"user_added_on"`
	UserStatus string `json:"user_status" bson:"user_status"`
	SponsorUName string `json:"sponsor_uname" bson:"sponsor_uname"`
}
//LevelDoc Ended ...
