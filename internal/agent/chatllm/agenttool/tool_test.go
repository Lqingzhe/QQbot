package agenttool

import "testing"

func Test_NewTool(t *testing.T) {
	tool, err := NewTool()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tool)
}
