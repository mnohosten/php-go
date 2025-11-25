package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// ChecklistItem represents a single task in the migration checklist
type ChecklistItem struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Phase       string    `json:"phase"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Notes       string    `json:"notes,omitempty"`
}

// MigrationChecklist manages the migration checklist
type MigrationChecklist struct {
	ProjectName string                    `json:"project_name"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
	Items       map[string]*ChecklistItem `json:"items"`
}

// NewMigrationChecklist creates a new migration checklist
func NewMigrationChecklist(projectName string) *MigrationChecklist {
	now := time.Now()
	mc := &MigrationChecklist{
		ProjectName: projectName,
		CreatedAt:   now,
		UpdatedAt:   now,
		Items:       make(map[string]*ChecklistItem),
	}
	mc.initializeItems()
	return mc
}

// initializeItems populates the checklist with all migration tasks
func (mc *MigrationChecklist) initializeItems() {
	phases := []struct {
		phase string
		items []string
	}{
		{
			phase: "Pre-Migration",
			items: []string{
				"Review migration guide completely",
				"Assess current PHP version (target PHP 8.2+)",
				"Inventory all PHP extensions used",
				"List all critical standard library functions used",
				"Check framework compatibility (WordPress/Laravel/Symfony)",
				"Review extension status section",
				"Identify blocking dependencies",
				"Set up test environment with PHP-Go",
				"Install PHP-Go build tools",
				"Create migration timeline and plan",
				"Get team buy-in and training scheduled",
				"Set up monitoring and logging",
				"Create rollback procedure document",
				"Backup production environment",
			},
		},
		{
			phase: "Parse Testing",
			items: []string{
				"Run parse test on entire codebase",
				"Document parse success rate",
				"Identify and categorize parse errors",
				"Fix critical parse errors",
				"Update code for PHP 8.4 compatibility",
				"Test inline HTML/PHP mixing",
				"Verify alternative syntax support",
				"Test with all application entry points",
			},
		},
		{
			phase: "Functional Testing",
			items: []string{
				"Set up parallel test environment (PHP and PHP-Go)",
				"Run unit test suite on PHP-Go",
				"Document test failures and incompatibilities",
				"Create polyfills for missing functions",
				"Implement workarounds for missing extensions",
				"Run integration tests",
				"Test critical user workflows",
				"Verify data consistency",
				"Test error handling and edge cases",
				"Validate API responses match exactly",
				"Check database interactions",
				"Test file uploads and downloads",
				"Verify session management",
				"Test authentication and authorization",
				"Validate email sending",
				"Check cron jobs and background tasks",
			},
		},
		{
			phase: "Configuration",
			items: []string{
				"Migrate php.ini to php-go.yaml",
				"Configure resource limits",
				"Set up logging and error reporting",
				"Configure extensions",
				"Set up environment variables",
				"Document configuration differences",
				"Test configuration in dev/staging/prod environments",
				"Configure monitoring and metrics",
				"Set up health checks",
				"Configure graceful shutdown",
			},
		},
		{
			phase: "Performance Testing",
			items: []string{
				"Run baseline performance tests on PHP",
				"Run same tests on PHP-Go",
				"Compare response times",
				"Compare memory usage",
				"Compare CPU usage",
				"Test concurrent request handling",
				"Run load tests",
				"Profile hot code paths",
				"Optimize based on profiling results",
				"Enable parallelization where beneficial",
				"Tune configuration for workload",
				"Verify performance meets requirements",
			},
		},
		{
			phase: "Deployment Preparation",
			items: []string{
				"Create deployment scripts",
				"Set up blue-green or canary deployment",
				"Prepare rollback procedure",
				"Document deployment steps",
				"Create monitoring dashboards",
				"Set up alerts for errors and performance",
				"Train operations team",
				"Schedule deployment window",
				"Notify stakeholders",
				"Prepare communication plan",
			},
		},
		{
			phase: "Production Rollout",
			items: []string{
				"Deploy to staging environment",
				"Run full test suite in staging",
				"Perform smoke tests",
				"Start with 1-5% traffic",
				"Monitor errors and performance",
				"Gradually increase to 10%",
				"Continue monitoring",
				"Increase to 25%",
				"Increase to 50%",
				"Monitor for 24-48 hours",
				"Increase to 100%",
				"Monitor for 1 week",
				"Document any issues encountered",
				"Optimize based on production metrics",
			},
		},
		{
			phase: "Post-Migration",
			items: []string{
				"Remove compatibility shims",
				"Optimize for PHP-Go features",
				"Enable parallelization",
				"Implement Go integration where beneficial",
				"Update documentation",
				"Share learnings with team",
				"Contribute improvements to PHP-Go",
				"Monitor long-term stability",
				"Plan next optimization phase",
			},
		},
	}

	for _, phase := range phases {
		for i, desc := range phase.items {
			id := fmt.Sprintf("%s-%d", strings.ToLower(strings.ReplaceAll(phase.phase, " ", "-")), i+1)
			mc.Items[id] = &ChecklistItem{
				ID:          id,
				Description: desc,
				Phase:       phase.phase,
				Completed:   false,
			}
		}
	}
}

// MarkComplete marks an item as complete
func (mc *MigrationChecklist) MarkComplete(id string, notes string) error {
	item, exists := mc.Items[id]
	if !exists {
		return fmt.Errorf("item not found: %s", id)
	}

	item.Completed = true
	item.CompletedAt = time.Now()
	if notes != "" {
		item.Notes = notes
	}
	mc.UpdatedAt = time.Now()

	return nil
}

// MarkIncomplete marks an item as incomplete
func (mc *MigrationChecklist) MarkIncomplete(id string) error {
	item, exists := mc.Items[id]
	if !exists {
		return fmt.Errorf("item not found: %s", id)
	}

	item.Completed = false
	item.CompletedAt = time.Time{}
	mc.UpdatedAt = time.Now()

	return nil
}

// AddNote adds a note to an item
func (mc *MigrationChecklist) AddNote(id string, note string) error {
	item, exists := mc.Items[id]
	if !exists {
		return fmt.Errorf("item not found: %s", id)
	}

	if item.Notes != "" {
		item.Notes += "\n" + note
	} else {
		item.Notes = note
	}
	mc.UpdatedAt = time.Now()

	return nil
}

// GetPhaseItems returns all items for a specific phase
func (mc *MigrationChecklist) GetPhaseItems(phase string) []*ChecklistItem {
	var items []*ChecklistItem
	for _, item := range mc.Items {
		if item.Phase == phase {
			items = append(items, item)
		}
	}

	// Sort by ID to maintain order
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return items
}

// GetPhases returns all unique phase names in order
func (mc *MigrationChecklist) GetPhases() []string {
	phases := []string{
		"Pre-Migration",
		"Parse Testing",
		"Functional Testing",
		"Configuration",
		"Performance Testing",
		"Deployment Preparation",
		"Production Rollout",
		"Post-Migration",
	}
	return phases
}

// GetProgress returns the completion percentage for a phase
func (mc *MigrationChecklist) GetProgress(phase string) (completed int, total int, percentage float64) {
	items := mc.GetPhaseItems(phase)
	total = len(items)

	for _, item := range items {
		if item.Completed {
			completed++
		}
	}

	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}

	return
}

// GetOverallProgress returns the overall completion percentage
func (mc *MigrationChecklist) GetOverallProgress() (completed int, total int, percentage float64) {
	total = len(mc.Items)

	for _, item := range mc.Items {
		if item.Completed {
			completed++
		}
	}

	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}

	return
}

// GenerateReport generates a text report of the checklist
func (mc *MigrationChecklist) GenerateReport() string {
	var sb strings.Builder

	sb.WriteString("PHP-Go Migration Checklist Report\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	sb.WriteString(fmt.Sprintf("Project: %s\n", mc.ProjectName))
	sb.WriteString(fmt.Sprintf("Created: %s\n", mc.CreatedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Updated: %s\n\n", mc.UpdatedAt.Format("2006-01-02 15:04:05")))

	// Overall progress
	completed, total, percentage := mc.GetOverallProgress()
	sb.WriteString(fmt.Sprintf("Overall Progress: %d/%d (%.1f%%)\n\n", completed, total, percentage))

	// Progress bar
	barWidth := 50
	filled := int(float64(barWidth) * percentage / 100)
	sb.WriteString("[")
	sb.WriteString(strings.Repeat("=", filled))
	sb.WriteString(strings.Repeat(" ", barWidth-filled))
	sb.WriteString("]\n\n")

	// Phase breakdown
	sb.WriteString("Phase Breakdown:\n")
	sb.WriteString(strings.Repeat("-", 80) + "\n\n")

	for _, phase := range mc.GetPhases() {
		completed, total, percentage := mc.GetProgress(phase)
		status := "⬜"
		if percentage == 100 {
			status = "✅"
		} else if percentage > 0 {
			status = "🟡"
		}

		sb.WriteString(fmt.Sprintf("%s %s: %d/%d (%.1f%%)\n", status, phase, completed, total, percentage))

		items := mc.GetPhaseItems(phase)
		for _, item := range items {
			checkbox := "[ ]"
			if item.Completed {
				checkbox = "[x]"
			}
			sb.WriteString(fmt.Sprintf("  %s %s\n", checkbox, item.Description))

			if item.Notes != "" {
				for _, line := range strings.Split(item.Notes, "\n") {
					sb.WriteString(fmt.Sprintf("      Note: %s\n", line))
				}
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// GenerateMarkdown generates a markdown report of the checklist
func (mc *MigrationChecklist) GenerateMarkdown() string {
	var sb strings.Builder

	sb.WriteString("# PHP-Go Migration Checklist\n\n")

	sb.WriteString(fmt.Sprintf("**Project:** %s  \n", mc.ProjectName))
	sb.WriteString(fmt.Sprintf("**Created:** %s  \n", mc.CreatedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("**Updated:** %s  \n\n", mc.UpdatedAt.Format("2006-01-02 15:04:05")))

	// Overall progress
	completed, total, percentage := mc.GetOverallProgress()
	sb.WriteString(fmt.Sprintf("**Overall Progress:** %d/%d (%.1f%%)\n\n", completed, total, percentage))

	// Table of contents
	sb.WriteString("## Table of Contents\n\n")
	for i, phase := range mc.GetPhases() {
		phaseID := strings.ToLower(strings.ReplaceAll(phase, " ", "-"))
		completed, total, _ := mc.GetProgress(phase)
		sb.WriteString(fmt.Sprintf("%d. [%s](#%s) (%d/%d)\n", i+1, phase, phaseID, completed, total))
	}
	sb.WriteString("\n")

	// Phase details
	for _, phase := range mc.GetPhases() {
		completed, total, percentage := mc.GetProgress(phase)

		status := "⬜ Not Started"
		if percentage == 100 {
			status = "✅ Complete"
		} else if percentage > 0 {
			status = fmt.Sprintf("🟡 In Progress (%.1f%%)", percentage)
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n", phase))
		sb.WriteString(fmt.Sprintf("**Status:** %s  \n", status))
		sb.WriteString(fmt.Sprintf("**Progress:** %d/%d items\n\n", completed, total))

		items := mc.GetPhaseItems(phase)
		for _, item := range items {
			checkbox := "- [ ]"
			if item.Completed {
				checkbox = "- [x]"
			}
			sb.WriteString(fmt.Sprintf("%s %s\n", checkbox, item.Description))

			if item.Notes != "" {
				sb.WriteString(fmt.Sprintf("  - *Note:* %s\n", strings.ReplaceAll(item.Notes, "\n", " ")))
			}
			if item.Completed {
				sb.WriteString(fmt.Sprintf("  - *Completed:* %s\n", item.CompletedAt.Format("2006-01-02")))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// Save saves the checklist to a JSON file
func (mc *MigrationChecklist) Save(filename string) error {
	data, err := json.MarshalIndent(mc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checklist: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write checklist file: %w", err)
	}

	return nil
}

// Load loads a checklist from a JSON file
func LoadMigrationChecklist(filename string) (*MigrationChecklist, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read checklist file: %w", err)
	}

	var mc MigrationChecklist
	err = json.Unmarshal(data, &mc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal checklist: %w", err)
	}

	return &mc, nil
}

// FindItem finds an item by partial ID or description match
func (mc *MigrationChecklist) FindItem(query string) (*ChecklistItem, error) {
	query = strings.ToLower(query)

	// Try exact ID match first
	if item, exists := mc.Items[query]; exists {
		return item, nil
	}

	// Try partial ID match
	for id, item := range mc.Items {
		if strings.Contains(strings.ToLower(id), query) {
			return item, nil
		}
	}

	// Try description match
	for _, item := range mc.Items {
		if strings.Contains(strings.ToLower(item.Description), query) {
			return item, nil
		}
	}

	return nil, fmt.Errorf("no item found matching: %s", query)
}

// ListIncomplete returns all incomplete items
func (mc *MigrationChecklist) ListIncomplete() []*ChecklistItem {
	var items []*ChecklistItem
	for _, item := range mc.Items {
		if !item.Completed {
			items = append(items, item)
		}
	}

	// Sort by phase and ID
	sort.Slice(items, func(i, j int) bool {
		if items[i].Phase != items[j].Phase {
			// Use phase order
			phases := mc.GetPhases()
			iIdx, jIdx := -1, -1
			for idx, phase := range phases {
				if items[i].Phase == phase {
					iIdx = idx
				}
				if items[j].Phase == phase {
					jIdx = idx
				}
			}
			return iIdx < jIdx
		}
		return items[i].ID < items[j].ID
	})

	return items
}

// GetNextTask returns the next incomplete task
func (mc *MigrationChecklist) GetNextTask() *ChecklistItem {
	incomplete := mc.ListIncomplete()
	if len(incomplete) == 0 {
		return nil
	}
	return incomplete[0]
}
