package registry

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

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

func (store *mgo.Collection) GetTemplate(name string) (Template, error) {
    result := Template{}
    error := store.Find(bson.M{"name": name}).One(&result)
    if (error != nil) {
       return Template{}, UnavailableTemplateError
    }
    return result, nil
}
