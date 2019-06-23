package main
import(
  "fmt"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2/txn"
)
var sess *mgo.Session
var collection *mgo.Collection

func setCollection(dbName string, collectionName string) *mgo.Collection {
	if sess == nil {
		fmt.Println("Not connected... Connecting to Mongo")
		sess = GetConnected()
	}
	collection = sess.DB(dbName).C(collectionName)
	return collection
}

func main(){
  fmt.Println("In main fun ")
  runner := txn.NewRunner(setCollection("dummydb","emp_col"))
  ops := []txn.Op{{
      C:      "emp_col",
      Id:     bson.ObjectId("5c76651a172224eb50237560"),
      Assert: bson.M{"name": bson.M{"$eq": "prashant"}},
      Update: bson.M{"$set": bson.M{"city": "Aurangabad"}},
    },
    {
      C:      "student_col",
      Id:     bson.ObjectId("5c766614172224eb50237562"),
      Assert: bson.M{"name": bson.M{"$eq": "sohil"}},
      Update: bson.M{"$set": bson.M{"city": "Aurangabad"}},
    },
  }
  id := bson.NewObjectId() // Optional
  runnerErr := runner.Run(ops, id, nil)
  if runnerErr != nil {
    fmt.Println("Error while runner : ",runnerErr)
  }
}
