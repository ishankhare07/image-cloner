package utils

import (
	"testing"
)

func TestRegistryInfo(t *testing.T) {
	inputs := []struct {
		Input    string
		Expected string
	}{
		{
			"nginx",
			"index.docker.io",
		},
		{
			"nginx:1.4.0",
			"index.docker.io",
		},
		{
			"ishankhare07/nginx:1.4.0",
			"index.docker.io/ishankhare07",
		},
		{
			"quay.io/ishankhare07/istio:1.4",
			"quay.io/ishankhare07",
		},
	}

	for _, input := range inputs {
		r := ExtractRegistryInfo(input.Input)
		if r.GetNameForClient() != input.Expected {
			t.Errorf("registry info wrong for %s, expected %s, got %s", input.Input, input.Expected, r.GetNameForClient())
		}
	}
}
