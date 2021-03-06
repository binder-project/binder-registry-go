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

func TestMongoRegisterTemplate(t *testing.T) {
    var MONGODB_URL string = "127.0.0.1"
    var MONGODB_DB string = "binder_registery_tests"
    var MONGODB_COL string = "templates"
    store := NewMongoStore(MONGODB_URL, MONGODB_DB, MONGODB_COL)

    tmpl := Template{
        Name: "MongoTest",
        ImageName: "jupyter/mongo_test",
    }

    registered, err := store.RegisterTemplate(tmpl)

    if (err != nil) {
        t.Error("Error when registering template: ", err)
    }

    if (registered == Template{}) {
        t.Error("Template was not registered properly")
    }

    equals(t, registered.Name, tmpl.Name)
    equals(t, registered.ImageName, tmpl.ImageName)
}

func TestMongoGetTemplate(t * testing.T) {
    var MONGODB_URL string = "127.0.0.1"
    var MONGODB_DB string = "binder_registery_tests"
    var MONGODB_COL string = "templates"
    store := NewMongoStore(MONGODB_URL, MONGODB_DB, MONGODB_COL)

    name := "MongoTest"
    tmpl, err := store.GetTemplate(name)

    if (err != nil) {
        t.Error("Error when getting template: ", err)
    }

    equals(t, tmpl.Name, name)
}

func TestMongoListTemplates(t *testing.T) {
    var MONGODB_URL string = "127.0.0.1"
    var MONGODB_DB string = "binder_registery_tests"
    var MONGODB_COL string = "templates"
    store := NewMongoStore(MONGODB_URL, MONGODB_DB, MONGODB_COL)

    results, err := store.ListTemplates()

    if (err != nil) {
        t.Error("Error when listing templates: ", err)
    }

    if (len(results) == 0) {
        t.Error("ListTemplates did not find any tempaltes!")
    }
}

func TestMongoDeleteTemplate(t *testing.T) {
    var MONGODB_URL string = "127.0.0.1"
    var MONGODB_DB string = "binder_registery_tests"
    var MONGODB_COL string = "templates"
    store := NewMongoStore(MONGODB_URL, MONGODB_DB, MONGODB_COL)

    tmpl := Template{
        Name: "DeleteMePlease",
        ImageName: "jupyter/delete_me",
    }

    _, err := store.RegisterTemplate(tmpl)

    if (err == nil) {
        err := store.DeleteTemplate(tmpl.Name)

        if (err != nil) {
            t.Error("Error when deleting template: ", err)
        } else {
            tmpl, err = store.GetTemplate(tmpl.Name)

            if (err == nil && tmpl != Template{}) {
                t.Error("Template was not successfully deleted")
            }
        }
    }
}

func TestMongoUpdateTemplate(t *testing.T) {
    var MONGODB_URL string = "127.0.0.1"
    var MONGODB_DB string = "binder_registery_tests"
    var MONGODB_COL string = "templates"
    store := NewMongoStore(MONGODB_URL, MONGODB_DB, MONGODB_COL)

    tmpl := Template{
        Name: "UpdateMePlease",
        ImageName: "jupyter/update_me",
    }

    _, err := store.RegisterTemplate(tmpl)

    if (err == nil) {
        updates := make(map[string]string)
        updates["name"] = "Updated!"
        err := store.UpdateTemplate(tmpl.Name, updates)

        if (err != nil) {
            t.Error("Error when updating template: ", err)
        } else {
            updated, err := store.GetTemplate("Updated!")

            if (err == nil && updated == Template{}) {
                t.Error("Template was not successfully updated!")
            }
        }
    }
}
