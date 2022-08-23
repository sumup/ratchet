package parser

import (
	"fmt"
	"strings"

	"github.com/sethvargo/ratchet/resolver"
	"gopkg.in/yaml.v3"
)

type Actions struct{}

// Parse pulls the GitHub Actions refs from the document.
func (a *Actions) Parse(m *yaml.Node) (*RefsList, error) {
	var refs RefsList

	if m == nil {
		return nil, nil
	}

	if m.Kind != yaml.DocumentNode {
		return nil, fmt.Errorf("expected document node, got %v", m.Kind)
	}

	// Top-level object map
	for _, docMap := range m.Content {
		if docMap.Kind != yaml.MappingNode {
			continue
		}

		// runs: keyword
		for i, runsMap := range docMap.Content {
			if runsMap.Value != "runs" {
				continue
			}

			// Individual steps names
			actions := docMap.Content[i+1]
			if actions.Kind != yaml.MappingNode {
				continue
			}
			// List of steps, iterate over each step and find the "uses" clause.
			for j, sub := range actions.Content {
			  if sub.Value == "steps" {
				  steps := actions.Content[j+1]
				  for _, step := range steps.Content {
					  if step.Kind != yaml.MappingNode {
						  continue
					  }
					  for k, property := range step.Content {
						  if property.Value == "uses" {
							  uses := step.Content[k+1]
							  switch {
							  case strings.HasPrefix(uses.Value, "docker://"):
								  ref := resolver.NormalizeContainerRef(uses.Value)
								  refs.Add(ref, uses)
							  case strings.Contains(uses.Value, "@"):
								  ref := resolver.NormalizeActionsRef(uses.Value)
								  refs.Add(ref, uses)
							  }
						  }
					  }
				  }
			  }
      }
 		}
 	}
	return &refs, nil
}
