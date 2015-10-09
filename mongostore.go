package registry

func NewMongoStore(server string, database string, collection string) Collection {
    session, error := mgo.Dial(server);
    if (error) {
        panic(error);
    }
    defer session.Close();

    connection = session.DB(database).C(collection)
    return collection
}
