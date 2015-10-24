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

func (store MongoStore) ListTemplates() ([]Template, error) {
    var results []Template

    err := store.connection.Find(bson.M{}).All(&results)

    if (err != nil) {
        return results, err
    }

    return results, nil
}

func contains(list []string, item string) bool {
    for _, a := range list {
        if a == item {
            return true
        }
    }
    return false
}

func (store MongoStore) UpdateTemplate(name string,
                                    update map[string]interface {}) error {
    TEMPLATE_KEYS := []string{"name", "image-name", "command", "limits",
                                "time-created", "time-modified", "redirect-uri",
                                "container-ip", "bind-port"}
    // If the update parameter is not nil and the
    // keys in update are all valid then execute the update
    ok := true
    for key, _ := range update {
        if !contains(TEMPLATE_KEYS, key) {
            ok = false
        }
    }

    if (update != nil && ok) {
        updates := bson.M{"$set": update}
        filter := bson.M{"name": name}

        err := store.connection.Update(filter, updates)

        if (err != nil) {
            return err
        } else {
            return nil
        }
    } else {
        return InvalidParameterError
    }
}

func (store MongoStore) DeleteTemplate(name string) error {
    filter := bson.M{"name": name}

    err := store.connection.Remove(filter)

    if (err != nil) {
        return err
    }

    return nil
}
