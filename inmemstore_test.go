package registry

import "testing"

func TestNewInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()
	assert(t, store.templateMap != nil, "Template map initialized")
}

func TestGetTemplate(t *testing.T) {
	store := NewInMemoryStore()

	_, err := store.GetTemplate("nope")
	equals(t, err, UnavailableTemplateError)

	expected := Template{
		Name:      "Test1",
		ImageName: "jupyter/whoa",
	}

	_, err = store.RegisterTemplate(expected)
	ok(t, err)

	actual, err := store.GetTemplate("Test1")
	ok(t, err)
	equals(t, actual.Name, expected.Name)
	equals(t, actual.ImageName, expected.ImageName)

	// TODO: Introduce delete and show that getting fails

}

func TestRegisterTemplate(t *testing.T) {
	store := NewInMemoryStore()

	expected := Template{
		Name:      "Test1",
		ImageName: "jupyter/whoa",
	}

	registered, err := store.RegisterTemplate(expected)
	ok(t, err)

	registered.TimeCreated = expected.TimeCreated
	registered.TimeModified = expected.TimeModified
	equals(t, registered, expected)

	// Can't register twice
	_, err = store.RegisterTemplate(expected)
	equals(t, err, ExistingTemplateError)
}

func TestListTemplates(t *testing.T) {
	store := NewInMemoryStore()

	templates, err := store.ListTemplates()
	ok(t, err)

	// Should be empty
	equals(t, templates, []Template{})

	registered1, err := store.RegisterTemplate(Template{
		Name:      "Test1",
		ImageName: "jupyter/whoa",
	})
	ok(t, err)

	// Can assume ordering since only one element
	templates, err = store.ListTemplates()
	equals(t, templates, []Template{registered1})

	registered2, err := store.RegisterTemplate(Template{
		Name:      "Test2",
		ImageName: "jupyter/whoathere",
	})
	ok(t, err)
	templates, err = store.ListTemplates()

	equals(t, len(templates), 2)

	expectedTemplates := []Template{registered1, registered2}

	if templates[0].Name != expectedTemplates[0].Name {
		// Hot swap, can't assume ordering
		expectedTemplates[0], expectedTemplates[1] = expectedTemplates[1], expectedTemplates[0]
	}

	equals(t, templates, expectedTemplates)

}

func TestUpdateTemplate(t *testing.T) {
	store := NewInMemoryStore()
	_, err := store.UpdateTemplate(Template{Name: "NotHere", ImageName: "cool/app"})
	equals(t, err, UnavailableTemplateError)

	template := Template{
		Name:      "Test1",
		ImageName: "jupyter/whoa",
	}

	registeredTemplate, err := store.RegisterTemplate(template)
	ok(t, err)

	updatedTemplate := registeredTemplate

	updatedTemplate.ImageName = "nteract/poster"
	nequals(t, updatedTemplate, registeredTemplate)

	actualTemplate, err := store.UpdateTemplate(updatedTemplate)
	ok(t, err)

	// TODO: Mock time
	nequals(t, updatedTemplate.TimeModified, actualTemplate.TimeModified)

	updatedTemplate.TimeModified = actualTemplate.TimeModified
	equals(t, updatedTemplate, actualTemplate)

}
