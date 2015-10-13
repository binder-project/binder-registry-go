package registry

import "time"
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

func (store MongoStore) RegisterTemplate(tmpl Template) (Template, error) {
    result, err := store.GetTemplate(tmpl.Name)
    if (err != nil && result != Template{}) {
        // This template is already in the database
        return Template{}, ExistingTemplateError
    }

    tmpl.TimeModified = time.Now().UTC()
    tmpl.TimeCreated = tmpl.TimeModified

    error := store.connection.Insert(&tmpl)

    if (error != nil) {
        return Template{}, error
    }

    return tmpl, nil
}
