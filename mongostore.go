package registry

import "gopkg.in/mgo.v2"

func NewMongoStore(server string,
                    database string,
                    collection string) (*mgo.Collection, error) {
    session, error := mgo.Dial(server);
    if (error != nil) {
        return nil, error
    }
    defer session.Close();

    connection := session.DB(database).C(collection)
    return connection, nil
}
