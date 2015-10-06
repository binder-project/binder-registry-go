package registry

import "testing"

func TestNewInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()
	assert(t, store.templateMap != nil, "Template map initialized")
}

func TestGetTemplate(t *testing.T) {
	store := NewInMemoryStore()

	_, err := store.GetTemplate("nope")
	equals(t, err, unavailableTemplateError)

	expected := Template{
		Name:      "Test1",
		ImageName: "jupyter/whoa",
	}

	registered, err := store.RegisterTemplate(expected)
	ok(t, err)

	// Note that the times will be different, so we adjust
	// Registration is tested elsewhere
	registered.TimeCreated = expected.TimeCreated
	registered.TimeModified = expected.TimeModified
	equals(t, registered, expected)

	actual, err := store.GetTemplate("Test1")
	ok(t, err)
	equals(t, actual.Name, expected.Name)
	equals(t, actual.ImageName, expected.ImageName)

}

func TestRegisterTemplate(t *testing.T) {
}

func TestListTemplates(t *testing.T) {
}

func TestUpdateTemplate(t *testing.T) {
}
