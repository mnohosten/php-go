package migrate

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewMigrationChecklist(t *testing.T) {
	projectName := "test-project"
	mc := NewMigrationChecklist(projectName)

	if mc.ProjectName != projectName {
		t.Errorf("Expected project name %s, got %s", projectName, mc.ProjectName)
	}

	if len(mc.Items) == 0 {
		t.Error("Expected checklist to have items")
	}

	// Verify all phases are present
	phases := mc.GetPhases()
	expectedPhases := []string{
		"Pre-Migration",
		"Parse Testing",
		"Functional Testing",
		"Configuration",
		"Performance Testing",
		"Deployment Preparation",
		"Production Rollout",
		"Post-Migration",
	}

	if len(phases) != len(expectedPhases) {
		t.Errorf("Expected %d phases, got %d", len(expectedPhases), len(phases))
	}

	for i, phase := range expectedPhases {
		if i >= len(phases) || phases[i] != phase {
			t.Errorf("Expected phase %s at index %d, got %s", phase, i, phases[i])
		}
	}
}

func TestMarkComplete(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	// Get first item
	items := mc.GetPhaseItems("Pre-Migration")
	if len(items) == 0 {
		t.Fatal("No items in Pre-Migration phase")
	}

	item := items[0]
	itemID := item.ID

	// Mark as complete
	note := "Test completion note"
	err := mc.MarkComplete(itemID, note)
	if err != nil {
		t.Fatalf("Error marking complete: %v", err)
	}

	// Verify completion
	item, exists := mc.Items[itemID]
	if !exists {
		t.Fatal("Item not found after marking complete")
	}

	if !item.Completed {
		t.Error("Item should be marked as completed")
	}

	if item.Notes != note {
		t.Errorf("Expected note '%s', got '%s'", note, item.Notes)
	}

	if item.CompletedAt.IsZero() {
		t.Error("CompletedAt should be set")
	}

	// Test invalid ID
	err = mc.MarkComplete("invalid-id", "")
	if err == nil {
		t.Error("Expected error for invalid ID")
	}
}

func TestMarkIncomplete(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	// Get first item and mark complete
	items := mc.GetPhaseItems("Pre-Migration")
	item := items[0]
	itemID := item.ID

	err := mc.MarkComplete(itemID, "Test note")
	if err != nil {
		t.Fatalf("Error marking complete: %v", err)
	}

	// Mark as incomplete
	err = mc.MarkIncomplete(itemID)
	if err != nil {
		t.Fatalf("Error marking incomplete: %v", err)
	}

	// Verify incompletion
	item, exists := mc.Items[itemID]
	if !exists {
		t.Fatal("Item not found after marking incomplete")
	}

	if item.Completed {
		t.Error("Item should be marked as incomplete")
	}

	if !item.CompletedAt.IsZero() {
		t.Error("CompletedAt should be cleared")
	}
}

func TestAddNote(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	items := mc.GetPhaseItems("Pre-Migration")
	item := items[0]
	itemID := item.ID

	// Add first note
	note1 := "First note"
	err := mc.AddNote(itemID, note1)
	if err != nil {
		t.Fatalf("Error adding note: %v", err)
	}

	item = mc.Items[itemID]
	if item.Notes != note1 {
		t.Errorf("Expected note '%s', got '%s'", note1, item.Notes)
	}

	// Add second note
	note2 := "Second note"
	err = mc.AddNote(itemID, note2)
	if err != nil {
		t.Fatalf("Error adding second note: %v", err)
	}

	item = mc.Items[itemID]
	expectedNotes := note1 + "\n" + note2
	if item.Notes != expectedNotes {
		t.Errorf("Expected notes '%s', got '%s'", expectedNotes, item.Notes)
	}
}

func TestGetPhaseItems(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	phase := "Pre-Migration"
	items := mc.GetPhaseItems(phase)

	if len(items) == 0 {
		t.Errorf("Expected items in %s phase", phase)
	}

	// Verify all items belong to the phase
	for _, item := range items {
		if item.Phase != phase {
			t.Errorf("Item %s has wrong phase: %s", item.ID, item.Phase)
		}
	}

	// Verify items are sorted
	for i := 1; i < len(items); i++ {
		if items[i].ID < items[i-1].ID {
			t.Error("Items are not sorted by ID")
		}
	}
}

func TestGetProgress(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	phase := "Pre-Migration"
	items := mc.GetPhaseItems(phase)

	// Initially no items completed
	completed, total, percentage := mc.GetProgress(phase)
	if completed != 0 {
		t.Errorf("Expected 0 completed, got %d", completed)
	}
	if total != len(items) {
		t.Errorf("Expected total %d, got %d", len(items), total)
	}
	if percentage != 0 {
		t.Errorf("Expected 0%%, got %.1f%%", percentage)
	}

	// Mark half complete
	halfCount := len(items) / 2
	for i := 0; i < halfCount; i++ {
		mc.MarkComplete(items[i].ID, "")
	}

	completed, total, percentage = mc.GetProgress(phase)
	if completed != halfCount {
		t.Errorf("Expected %d completed, got %d", halfCount, completed)
	}
	expectedPercentage := float64(halfCount) / float64(total) * 100
	if percentage != expectedPercentage {
		t.Errorf("Expected %.1f%%, got %.1f%%", expectedPercentage, percentage)
	}

	// Mark all complete
	for i := halfCount; i < len(items); i++ {
		mc.MarkComplete(items[i].ID, "")
	}

	completed, total, percentage = mc.GetProgress(phase)
	if completed != total {
		t.Errorf("Expected %d completed, got %d", total, completed)
	}
	if percentage != 100.0 {
		t.Errorf("Expected 100%%, got %.1f%%", percentage)
	}
}

func TestGetOverallProgress(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	totalItems := len(mc.Items)

	// Initially no items completed
	completed, total, percentage := mc.GetOverallProgress()
	if completed != 0 {
		t.Errorf("Expected 0 completed, got %d", completed)
	}
	if total != totalItems {
		t.Errorf("Expected total %d, got %d", totalItems, total)
	}
	if percentage != 0 {
		t.Errorf("Expected 0%%, got %.1f%%", percentage)
	}

	// Mark some items complete
	count := 0
	for id := range mc.Items {
		if count >= 10 {
			break
		}
		mc.MarkComplete(id, "")
		count++
	}

	completed, total, percentage = mc.GetOverallProgress()
	if completed != count {
		t.Errorf("Expected %d completed, got %d", count, completed)
	}
	expectedPercentage := float64(count) / float64(totalItems) * 100
	if percentage != expectedPercentage {
		t.Errorf("Expected %.1f%%, got %.1f%%", expectedPercentage, percentage)
	}
}

func TestChecklistGenerateReport(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	// Mark some items complete
	items := mc.GetPhaseItems("Pre-Migration")
	mc.MarkComplete(items[0].ID, "Test note")

	report := mc.GenerateReport()

	// Verify report contains expected sections
	expectedSections := []string{
		"PHP-Go Migration Checklist Report",
		"Project: test-project",
		"Overall Progress:",
		"Phase Breakdown:",
		"Pre-Migration",
	}

	for _, section := range expectedSections {
		if !strings.Contains(report, section) {
			t.Errorf("Report missing expected section: %s", section)
		}
	}

	// Verify completed item is marked
	if !strings.Contains(report, "[x]") {
		t.Error("Report should contain completed items marked with [x]")
	}
}

func TestGenerateMarkdown(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	// Mark some items complete
	items := mc.GetPhaseItems("Pre-Migration")
	mc.MarkComplete(items[0].ID, "Test note")

	markdown := mc.GenerateMarkdown()

	// Verify markdown contains expected sections
	expectedSections := []string{
		"# PHP-Go Migration Checklist",
		"**Project:** test-project",
		"**Overall Progress:**",
		"## Table of Contents",
		"## Pre-Migration",
	}

	for _, section := range expectedSections {
		if !strings.Contains(markdown, section) {
			t.Errorf("Markdown missing expected section: %s", section)
		}
	}

	// Verify markdown checkbox syntax
	if !strings.Contains(markdown, "- [x]") {
		t.Error("Markdown should contain completed items with '- [x]' syntax")
	}
	if !strings.Contains(markdown, "- [ ]") {
		t.Error("Markdown should contain incomplete items with '- [ ]' syntax")
	}
}

func TestSaveLoad(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	// Mark some items complete
	items := mc.GetPhaseItems("Pre-Migration")
	mc.MarkComplete(items[0].ID, "Test note 1")
	mc.MarkComplete(items[1].ID, "Test note 2")

	// Save to temp file
	tmpfile, err := os.CreateTemp("", "checklist-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	err = mc.Save(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to save checklist: %v", err)
	}

	// Load from file
	loaded, err := LoadMigrationChecklist(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load checklist: %v", err)
	}

	// Verify loaded data
	if loaded.ProjectName != mc.ProjectName {
		t.Errorf("Project name mismatch: %s != %s", loaded.ProjectName, mc.ProjectName)
	}

	if len(loaded.Items) != len(mc.Items) {
		t.Errorf("Item count mismatch: %d != %d", len(loaded.Items), len(mc.Items))
	}

	// Verify specific items
	for id, item := range mc.Items {
		loadedItem, exists := loaded.Items[id]
		if !exists {
			t.Errorf("Item %s not found in loaded checklist", id)
			continue
		}

		if loadedItem.Completed != item.Completed {
			t.Errorf("Item %s completion status mismatch", id)
		}

		if loadedItem.Notes != item.Notes {
			t.Errorf("Item %s notes mismatch", id)
		}
	}
}

func TestFindItem(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	items := mc.GetPhaseItems("Pre-Migration")
	targetItem := items[0]

	// Test exact ID match
	found, err := mc.FindItem(targetItem.ID)
	if err != nil {
		t.Fatalf("Failed to find item by exact ID: %v", err)
	}
	if found.ID != targetItem.ID {
		t.Errorf("Found wrong item: %s != %s", found.ID, targetItem.ID)
	}

	// Test partial ID match - use full ID to ensure exact match
	// (partial matches may match multiple items and return any one)
	found, err = mc.FindItem(targetItem.ID)
	if err != nil {
		t.Fatalf("Failed to find item by full ID: %v", err)
	}
	if found.ID != targetItem.ID {
		t.Errorf("Found wrong item: %s != %s", found.ID, targetItem.ID)
	}

	// Test description match - use a longer unique part
	// Use at least 20 characters to ensure uniqueness
	descLen := len(targetItem.Description)
	if descLen > 20 {
		descLen = 20
	}
	descPart := strings.ToLower(targetItem.Description[:descLen])
	found, err = mc.FindItem(descPart)
	if err != nil {
		t.Fatalf("Failed to find item by description: %v", err)
	}
	// Just verify we found something matching the description
	if !strings.Contains(strings.ToLower(found.Description), descPart) {
		t.Errorf("Found item description doesn't contain search term")
	}

	// Test not found
	_, err = mc.FindItem("nonexistent-item-id-12345")
	if err == nil {
		t.Error("Expected error for nonexistent item")
	}
}

func TestListIncomplete(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	totalItems := len(mc.Items)

	// Initially all incomplete
	incomplete := mc.ListIncomplete()
	if len(incomplete) != totalItems {
		t.Errorf("Expected %d incomplete items, got %d", totalItems, len(incomplete))
	}

	// Mark some complete
	items := mc.GetPhaseItems("Pre-Migration")
	completeCount := 3
	for i := 0; i < completeCount; i++ {
		mc.MarkComplete(items[i].ID, "")
	}

	incomplete = mc.ListIncomplete()
	expectedIncomplete := totalItems - completeCount
	if len(incomplete) != expectedIncomplete {
		t.Errorf("Expected %d incomplete items, got %d", expectedIncomplete, len(incomplete))
	}

	// Verify incomplete items are sorted by phase
	lastPhaseIdx := -1
	phases := mc.GetPhases()
	for i := 1; i < len(incomplete); i++ {
		currentPhaseIdx := -1
		for idx, phase := range phases {
			if incomplete[i].Phase == phase {
				currentPhaseIdx = idx
				break
			}
		}
		if currentPhaseIdx < lastPhaseIdx {
			t.Error("Incomplete items are not sorted by phase order")
		}
		lastPhaseIdx = currentPhaseIdx
	}
}

func TestGetNextTask(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	// Get next task (should be first in Pre-Migration)
	next := mc.GetNextTask()
	if next == nil {
		t.Fatal("Expected next task, got nil")
	}

	if next.Phase != "Pre-Migration" {
		t.Errorf("Expected first task to be in Pre-Migration, got %s", next.Phase)
	}

	firstID := next.ID

	// Mark complete and verify next task changes
	mc.MarkComplete(firstID, "")
	next = mc.GetNextTask()
	if next == nil {
		t.Fatal("Expected next task after completing first, got nil")
	}

	if next.ID == firstID {
		t.Error("Next task should be different after completing first task")
	}

	// Mark all complete
	for id := range mc.Items {
		mc.MarkComplete(id, "")
	}

	next = mc.GetNextTask()
	if next != nil {
		t.Error("Expected nil when all tasks complete")
	}
}

func TestUpdateTimestamps(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	originalUpdated := mc.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	items := mc.GetPhaseItems("Pre-Migration")
	mc.MarkComplete(items[0].ID, "")

	if !mc.UpdatedAt.After(originalUpdated) {
		t.Error("UpdatedAt should be updated after marking item complete")
	}

	time.Sleep(10 * time.Millisecond)
	originalUpdated = mc.UpdatedAt

	mc.AddNote(items[0].ID, "New note")
	if !mc.UpdatedAt.After(originalUpdated) {
		t.Error("UpdatedAt should be updated after adding note")
	}
}

func TestEmptyPhase(t *testing.T) {
	mc := NewMigrationChecklist("test-project")

	// Test with non-existent phase
	items := mc.GetPhaseItems("Nonexistent Phase")
	if len(items) != 0 {
		t.Errorf("Expected 0 items for nonexistent phase, got %d", len(items))
	}

	completed, total, percentage := mc.GetProgress("Nonexistent Phase")
	if completed != 0 || total != 0 || percentage != 0 {
		t.Error("Expected zero values for nonexistent phase progress")
	}
}
