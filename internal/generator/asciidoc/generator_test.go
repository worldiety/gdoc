package asciidoc

import (
	"godocgenerator/internal/api"
	"testing"
)

func TestCreateModuleTemplate(t *testing.T) {
	// Create a sample module
	module := &api.Module{
		Name: "test",
		Packages: map[api.ImportPath]*api.Package{
			"mypackage": {
				Name: "mypackage",
				Consts: map[string]*api.Constant{
					"c1": {Name: "c1", Comment: "testComment1"},
					"c2": {Name: "c2", Comment: "testComment2"},
				},
				Vars: map[string]*api.Variable{
					"v1": {Name: "v1", Comment: "var test comment 1"},
					"v2": {Name: "v2", Comment: "var test comment 2"},
				},
				Functions: map[string]*api.Function{
					"f1": {Name: "f1", Comment: "Function 1", Parameters: map[string]*api.Parameter{"p1": {Name: "p1", TypeDefinition: 3}}},
					"f2": {Name: "f2", Comment: "Function 2", Parameters: map[string]*api.Parameter{"p1": {Name: "p1", TypeDefinition: 4}},
						Results: map[string]*api.Parameter{"r1": {Name: "r1", TypeDefinition: 5}}},
				},
			},
		},
	}

	// Call the function
	buffer, err := CreateModuleTemplate(module)

	// Check for errors
	if err != nil {
		t.Errorf("CreateModuleTemplate returned an error: %v", err)
	}

	// Check if the buffer is not nil
	if buffer == nil {
		t.Errorf("CreateModuleTemplate returned a nil buffer")
	}

	// Check if the buffer contains the expected content
	expectedContent := "= Module test\n:toc:\n\n\n== Package mypackage\n=== Constants\n\n* c1\n** __Comment__: testComment1\n* c2\n** __Comment__: testComment2\n=== Variables\n\nv1 * Comment: var test comment 1\nv2 * Comment: var test comment 2\n\n=== Functions\n\nf1(p1 , )\n\n* Parameters\n\n** p1 ()\n\n* Comment: Function 1\n\nf2(p1 , )r1 , \n\n* Parameters\n\n** p1 ()\n\n* Comment: Function 2\n\n"
	if buffer.String() != expectedContent {
		t.Errorf("CreateModuleTemplate returned an unexpected buffer content: %v", buffer.String())
	}
}
