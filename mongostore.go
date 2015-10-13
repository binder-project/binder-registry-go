package registry

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

type MongoStore struct {
    connection *mgo.Collection
    err error
}

func NewMongoStore(server string,
                    database string,
                    collection string) MongoStore {
    session, error := mgo.Dial(server);

    if (error != nil) {
        return MongoStore{connection: nil, err: error}
    }
    defer session.Close();

    connection := session.DB(database).C(collection)
    return MongoStore{connection: connection, err: nil}
}

func (store MongoStore) GetTemplate(name string) (Template, error) {
    result := Template{}
    error := store.connection.Find(bson.M{"name": name}).One(&result)
    if (error != nil) {
       return Template{}, UnavailableTemplateError
    }
    return result, nil
}
