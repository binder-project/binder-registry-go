package registry

type MongoStore struct {
}

func NewMongoStore() MongoStore {
}

func (store MongoStore) GetTemplate(name string) (Template, error) {
}

func (store MongoStore) RegisterTemplate(tmpl Template) (Template, error) {
}

func (store MongoStore) ListTemplates() ([]Template, error) {
}

func (store MongoStore) UpdateTemplate(tmpl Template) (Template, error) {
}
