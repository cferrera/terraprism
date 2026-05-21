package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/viewport"

	"github.com/CaptShanks/terraprism/internal/parser"
)

func threeResourcesWithFolds() []parser.Resource {
	return []parser.Resource{
		{Address: "a.one", Action: parser.ActionUpdate, RawLines: []string{
			`  ~ resource "a" "one" {`, `      x = 1`, `    }`,
		}},
		{Address: "b.two", Action: parser.ActionUpdate, RawLines: []string{
			`  ~ resource "b" "two" {`,
			`      ~ metadata = {`, `          ~ values = {`, `              nested = true`, `            }`, `        }`,
			`    }`,
		}},
		{Address: "c.three", Action: parser.ActionUpdate, RawLines: []string{
			`  ~ resource "c" "three" {`, `      x = 3`, `    }`,
		}},
	}
}

// `e` at root scopes to the highlighted item: expand cursor's resource and its
// sub-folds, but leave siblings untouched.
func TestExpandAllScopesToHighlightedItem(t *testing.T) {
	resources := threeResourcesWithFolds()
	m := Model{
		plan:         &parser.Plan{Resources: resources},
		expanded:     map[int]bool{},
		foldedBlocks: make(map[string]bool),
		blockCursor:  -1,
		viewport:     viewport.New(120, 40),
		cursor:       1, // b.two
	}

	updated, _, _ := handleKeyExpandAll(m)

	if !updated.expanded[1] {
		t.Fatalf("cursor's resource (b.two) should be expanded; expanded=%v", updated.expanded)
	}
	if updated.expanded[0] || updated.expanded[2] {
		t.Fatalf("siblings should NOT have been expanded by scoped `e`; expanded=%v", updated.expanded)
	}
	for _, block := range findFoldBlocks(resources[1], resources[1].RawLines[1:]) {
		if updated.foldedBlocks[block.Key] {
			t.Fatalf("sub-fold %q of cursor's resource should be expanded after `e`", block.Key)
		}
	}
}

// `c` at root collapses just the cursor's item and its sub-folds.
func TestCollapseAllScopesToHighlightedItem(t *testing.T) {
	resources := threeResourcesWithFolds()
	m := Model{
		plan:         &parser.Plan{Resources: resources},
		expanded:     map[int]bool{0: true, 1: true, 2: true},
		foldedBlocks: make(map[string]bool),
		blockCursor:  -1,
		viewport:     viewport.New(120, 40),
		cursor:       1,
	}

	updated, _, _ := handleKeyCollapseAll(m)

	if updated.expanded[1] {
		t.Fatalf("cursor's resource should be collapsed; expanded=%v", updated.expanded)
	}
	if !updated.expanded[0] || !updated.expanded[2] {
		t.Fatalf("siblings should remain expanded; expanded=%v", updated.expanded)
	}
}

// Shift+E (handleKeyExpandEverything) is the global expand-all.
func TestExpandEverythingIsGlobal(t *testing.T) {
	resources := threeResourcesWithFolds()
	m := Model{
		plan:         &parser.Plan{Resources: resources},
		expanded:     map[int]bool{},
		foldedBlocks: make(map[string]bool),
		blockCursor:  -1,
		viewport:     viewport.New(120, 40),
		cursor:       1,
	}

	updated, _, _ := handleKeyExpandEverything(m)

	for idx := range resources {
		if !updated.expanded[idx] {
			t.Fatalf("resource %d should be expanded by `E`; expanded=%v", idx, updated.expanded)
		}
	}
}

// Shift+C is the global collapse-all.
func TestCollapseEverythingIsGlobal(t *testing.T) {
	resources := threeResourcesWithFolds()
	m := Model{
		plan:         &parser.Plan{Resources: resources},
		expanded:     map[int]bool{0: true, 1: true, 2: true},
		foldedBlocks: make(map[string]bool),
		blockCursor:  -1,
		viewport:     viewport.New(120, 40),
		cursor:       1,
	}

	updated, _, _ := handleKeyCollapseEverything(m)

	for idx := range resources {
		if updated.expanded[idx] {
			t.Fatalf("resource %d should be collapsed by `C`; expanded=%v", idx, updated.expanded)
		}
	}
}

// Inside a sub-fold, `e` still scopes to that fold.
func TestExpandAllInsideSubBlockKeepsScope(t *testing.T) {
	resources := threeResourcesWithFolds()
	m := Model{
		plan:         &parser.Plan{Resources: resources},
		expanded:     map[int]bool{1: true},
		foldedBlocks: make(map[string]bool),
		blockCursor:  0,
		viewport:     viewport.New(120, 40),
		cursor:       1,
	}

	updated, _, _ := handleKeyExpandAll(m)

	if updated.expanded[0] {
		t.Fatalf("sibling 0 should NOT expand from sub-block-scoped `e`; expanded=%v", updated.expanded)
	}
}
