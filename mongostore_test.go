package registry

import "testing"

func TestNewMongoStore(t *testing.T) {
    var MONGODB_URL string = "http://localhost:27017"
    var MONGODB_DB string = "binder_registery_tests"
    var MONGODB_COL string = "templates"
    connection, error := NewMongoStore(MONGODB_URL, MONGODB_DB, MONGODB_COL)
    
    if (error != nil) {
        t.Error("NewMongoStore raised error:", error);
    }

    if (connection.FullName != (MONGODB_DB + "." + MONGODB_COL)) {
       t.Error("Connected to wrong database:", connection.FullName) 
    }
}
