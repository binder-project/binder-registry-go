package registry

import "testing"

func TestNewMongoStore(t *testing.T) {
    var MONGODB_URL string = "127.0.0.1"
    var MONGODB_DB string = "binder_registery_tests"
    var MONGODB_COL string = "templates"
    store := NewMongoStore(MONGODB_URL, MONGODB_DB, MONGODB_COL)

    if (store.err != nil) {
        t.Error("NewMongoStore raised error:", store.err);
    }

    if (store.connection.FullName != (MONGODB_DB + "." + MONGODB_COL)) {
       t.Error("Connected to wrong database:", store.connection.FullName)
    }
}
