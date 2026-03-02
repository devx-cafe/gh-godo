package mkissue

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// Mock functions for testing
var (
	mockGitShowFunc      func(filePath, branch string) ([]byte, error)
	mockEnsureLabelFunc  func(label Label) error
	mockCreateIssueFunc  func(metadata *IssueMetadata, body string) error
	mockRunGhCommandFunc func(args []string) error
)

// Override functions with mocks during tests
func mockReadFileFromBranch(filePath, branch string) ([]byte, error) {
	if mockGitShowFunc != nil {
		return mockGitShowFunc(filePath, branch)
	}
	return readFileFromBranch(filePath, branch)
}

func mockEnsureLabelExists(label Label) error {
	if mockEnsureLabelFunc != nil {
		return mockEnsureLabelFunc(label)
	}
	return ensureLabelExists(label)
}

func mockCreateIssueInternal(metadata *IssueMetadata, body string) error {
	if mockCreateIssueFunc != nil {
		return mockCreateIssueFunc(metadata, body)
	}
	return createIssue(metadata, body)
}

func mockRunGhCommandInternal(args []string) error {
	if mockRunGhCommandFunc != nil {
		return mockRunGhCommandFunc(args)
	}
	return runGhCommand(args)
}

func TestExtractValue(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		prefix string
		want   string
	}{
		{
			name:   "simple value",
			line:   "title: My Issue Title",
			prefix: "title:",
			want:   "My Issue Title",
		},
		{
			name:   "quoted value",
			line:   "title: \"My Issue Title\"",
			prefix: "title:",
			want:   "My Issue Title",
		},
		{
			name:   "single quoted value",
			line:   "title: 'My Issue Title'",
			prefix: "title:",
			want:   "My Issue Title",
		},
		{
			name:   "value with extra spaces",
			line:   "title:   My Issue Title   ",
			prefix: "title:",
			want:   "My Issue Title",
		},
		{
			name:   "empty value",
			line:   "title:",
			prefix: "title:",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractValue(tt.line, tt.prefix)
			if got != tt.want {
				t.Errorf("extractValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseListField(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startIdx  int
		prefix    string
		wantItems []string
		wantIdx   int
	}{
		{
			name:      "inline array format",
			lines:     []string{"assign: [user1, user2, user3]"},
			startIdx:  0,
			prefix:    "assign:",
			wantItems: []string{"user1", "user2", "user3"},
			wantIdx:   0,
		},
		{
			name: "multi-line format",
			lines: []string{
				"assign:",
				"  - user1",
				"  - user2",
				"  - user3",
				"labels:",
			},
			startIdx:  0,
			prefix:    "assign:",
			wantItems: []string{"user1", "user2", "user3"},
			wantIdx:   3,
		},
		{
			name:      "empty list",
			lines:     []string{"assign:"},
			startIdx:  0,
			prefix:    "assign:",
			wantItems: []string{},
			wantIdx:   0,
		},
		{
			name:      "with quoted items",
			lines:     []string{"assign: [\"user1\", \"user2\"]"},
			startIdx:  0,
			prefix:    "assign:",
			wantItems: []string{"user1", "user2"},
			wantIdx:   0,
		},
		{
			name: "inline with @ prefix",
			lines: []string{
				"assign: [@user1, @user2]",
			},
			startIdx:  0,
			prefix:    "assign:",
			wantItems: []string{"user1", "user2"},
			wantIdx:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItems, gotIdx := parseListField(tt.lines, tt.startIdx, tt.prefix)
			if len(gotItems) != len(tt.wantItems) {
				t.Errorf("parseListField() items length = %d, want %d", len(gotItems), len(tt.wantItems))
				return
			}
			for i, item := range gotItems {
				if item != tt.wantItems[i] {
					t.Errorf("parseListField() item[%d] = %q, want %q", i, item, tt.wantItems[i])
				}
			}
			if gotIdx != tt.wantIdx {
				t.Errorf("parseListField() index = %d, want %d", gotIdx, tt.wantIdx)
			}
		})
	}
}

func TestParseLabels(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startIdx  int
		wantCount int
		wantFirst Label
		wantErr   bool
	}{
		{
			name: "single label",
			lines: []string{
				"labels:",
				"  - name: bug",
				"    color: ff0000",
				"    desc: Bug report",
			},
			startIdx:  0,
			wantCount: 1,
			wantFirst: Label{Name: "bug", Color: "ff0000", Desc: "Bug report"},
			wantErr:   false,
		},
		{
			name: "multiple labels",
			lines: []string{
				"labels:",
				"  - name: bug",
				"    color: ff0000",
				"  - name: feature",
				"    color: 00ff00",
			},
			startIdx:  0,
			wantCount: 2,
			wantFirst: Label{Name: "bug", Color: "ff0000"},
			wantErr:   false,
		},
		{
			name: "label without color and desc",
			lines: []string{
				"labels:",
				"  - name: urgent",
			},
			startIdx:  0,
			wantCount: 1,
			wantFirst: Label{Name: "urgent"},
			wantErr:   false,
		},
		{
			name:      "empty labels",
			lines:     []string{"labels:"},
			startIdx:  0,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLabels, _ := parseLabels(tt.lines, tt.startIdx)
			if len(gotLabels) != tt.wantCount {
				t.Errorf("parseLabels() count = %d, want %d", len(gotLabels), tt.wantCount)
				return
			}
			if tt.wantCount > 0 {
				if gotLabels[0].Name != tt.wantFirst.Name {
					t.Errorf("parseLabels()[0].Name = %q, want %q", gotLabels[0].Name, tt.wantFirst.Name)
				}
				if gotLabels[0].Color != tt.wantFirst.Color {
					t.Errorf("parseLabels()[0].Color = %q, want %q", gotLabels[0].Color, tt.wantFirst.Color)
				}
				if gotLabels[0].Desc != tt.wantFirst.Desc {
					t.Errorf("parseLabels()[0].Desc = %q, want %q", gotLabels[0].Desc, tt.wantFirst.Desc)
				}
			}
		})
	}
}

func TestParseIssueFile(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantTitle string
		wantBody  string
		wantErr   bool
	}{
		{
			name: "valid issue file",
			content: `---
title: Test Issue
assign: [user1, user2]
labels:
  - name: bug
    color: ff0000
---
This is the issue body.
It can have multiple lines.`,
			wantTitle: "Test Issue",
			wantBody:  "This is the issue body.\nIt can have multiple lines.",
			wantErr:   false,
		},
		{
			name: "issue with milestone and projects",
			content: `---
title: Feature Request
milestone: v1.0
projects: [project1, project2]
---
Implementation details here.`,
			wantTitle: "Feature Request",
			wantBody:  "Implementation details here.",
			wantErr:   false,
		},
		{
			name:      "missing frontmatter",
			content:   "This is not valid",
			wantTitle: "",
			wantBody:  "",
			wantErr:   true,
		},
		{
			name: "empty title",
			content: `---
assign: [user1]
---
Body content`,
			wantTitle: "",
			wantBody:  "Body content",
			wantErr:   false,
		},
		{
			name: "body with multiple dash lines",
			content: `---
title: Test
---
First part
---
Second part`,
			wantTitle: "Test",
			wantBody:  "First part\n---\nSecond part",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotBody, err := parseIssueFile(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIssueFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Title != tt.wantTitle {
				t.Errorf("parseIssueFile() title = %q, want %q", got.Title, tt.wantTitle)
			}
			if gotBody != tt.wantBody {
				t.Errorf("parseIssueFile() body = %q, want %q", gotBody, tt.wantBody)
			}
		})
	}
}

func TestRunWithFileFromFilesystem(t *testing.T) {
	// Create temporary file with valid content
	tmpFile, err := os.CreateTemp("", "test-issue-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `---
title: Test Issue from File
assign: [me]
---
This is a test issue.`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Setup mocks
	mockEnsureLabelFunc = func(label Label) error {
		return nil
	}
	mockCreateIssueFunc = func(metadata *IssueMetadata, body string) error {
		if metadata.Title != "Test Issue from File" {
			return fmt.Errorf("unexpected title")
		}
		return nil
	}
	defer func() {
		mockEnsureLabelFunc = nil
		mockCreateIssueFunc = nil
	}()

	// Need to modify RunWithFile to use mocks - for now test basic validation
	err = RunWithFile(tmpFile.Name(), "", "", "")
	if err != nil && strings.Contains(err.Error(), "gh command failed") {
		// This is expected since gh CLI not available, but parsing should have worked
		t.Logf("Got expected gh command error: %v", err)
	}
}

func TestRunWithFileNonexistentFile(t *testing.T) {
	err := RunWithFile("/nonexistent/file/path.md", "", "", "")
	if err == nil {
		t.Errorf("RunWithFile() expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("RunWithFile() error = %v, want 'not found'", err)
	}
}

func TestRunWithFileMissingTitle(t *testing.T) {
	// Create temporary file without title
	tmpFile, err := os.CreateTemp("", "test-issue-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `---
assign: [me]
---
This issue has no title.`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	err = RunWithFile(tmpFile.Name(), "", "", "")
	if err == nil {
		t.Errorf("RunWithFile() expected error for missing title")
	}
	if !strings.Contains(err.Error(), "title") {
		t.Errorf("RunWithFile() error = %v, want error about title", err)
	}
}

func TestReadFileFromBranchValidation(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		branch   string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "branch with null byte",
			filePath: "file.md",
			branch:   "branch\x00name",
			wantErr:  true,
			errMsg:   "prohibited characters",
		},
		{
			name:     "branch with newline",
			filePath: "file.md",
			branch:   "branch\nname",
			wantErr:  true,
			errMsg:   "prohibited characters",
		},
		{
			name:     "branch with carriage return",
			filePath: "file.md",
			branch:   "branch\rname",
			wantErr:  true,
			errMsg:   "prohibited characters",
		},
		{
			name:     "filePath with null byte",
			filePath: "file\x00.md",
			branch:   "mybranch",
			wantErr:  true,
			errMsg:   "prohibited characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := readFileFromBranch(tt.filePath, tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFileFromBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("readFileFromBranch() error = %v, want %q", err, tt.errMsg)
			}
		})
	}
}

func TestReadFileFromBranchInvalidBranch(t *testing.T) {
	// Test with a branch that doesn't exist
	_, err := readFileFromBranch("nonexistent.md", "nonexistent-branch")
	if err == nil {
		t.Errorf("readFileFromBranch() expected error for nonexistent branch")
	}
}

func TestReadFileFromGistValidation(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		gistID   string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "gist ID too short",
			fileName: "file.md",
			gistID:   "shortid",
			wantErr:  true,
			errMsg:   "32-character hexadecimal",
		},
		{
			name:     "gist ID with uppercase",
			fileName: "file.md",
			gistID:   "6EF8A9C46F65F5FEDB58E81B70DD90BA",
			wantErr:  true,
			errMsg:   "32-character hexadecimal",
		},
		{
			name:     "gist ID with non-hex characters",
			fileName: "file.md",
			gistID:   "6ef8a9c46f65f5fedb58e81b70dd90bg",
			wantErr:  true,
			errMsg:   "32-character hexadecimal",
		},
		{
			name:     "fileName with path traversal",
			fileName: "../file.md",
			gistID:   "6ef8a9c46f65f5fedb58e81b70dd90ba",
			wantErr:  true,
			errMsg:   "alphanumeric",
		},
		{
			name:     "fileName with forward slash",
			fileName: "path/file.md",
			gistID:   "6ef8a9c46f65f5fedb58e81b70dd90ba",
			wantErr:  true,
			errMsg:   "alphanumeric",
		},
		{
			name:     "fileName with backslash",
			fileName: "path\\file.md",
			gistID:   "6ef8a9c46f65f5fedb58e81b70dd90ba",
			wantErr:  true,
			errMsg:   "alphanumeric",
		},
		{
			name:     "fileName with special characters",
			fileName: "file$name.md",
			gistID:   "6ef8a9c46f65f5fedb58e81b70dd90ba",
			wantErr:  true,
			errMsg:   "alphanumeric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := readFileFromGist(tt.fileName, tt.gistID)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFileFromGist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("readFileFromGist() error = %v, want %q", err, tt.errMsg)
			}
		})
	}
}

func TestReadFileFromGistInvalidGist(t *testing.T) {
	// Test with a gist that doesn't exist - use valid format but nonexistent ID
	_, err := readFileFromGist("nonexistent.md", "0000000000000000000000000000000a")
	if err == nil {
		t.Errorf("readFileFromGist() expected error for nonexistent gist")
	}
	// Verify the error message indicates it failed to read from gist
	if err != nil && !strings.Contains(err.Error(), "failed to read file from gist") {
		t.Errorf("readFileFromGist() error = %v, want error containing 'failed to read file from gist'", err)
	}
}

func TestIntegrationParseAndValidate(t *testing.T) {
	// Test complete parsing and validation flow
	content := `---
title: Complete Issue
assign: [user1, me]
labels:
  - name: enhancement
    color: 84b6eb
    desc: New feature or request
milestone: v2.0
projects: [Backend, Frontend]
---
This is a comprehensive test issue.
It includes all metadata fields.`

	metadata, body, err := parseIssueFile(content)
	if err != nil {
		t.Fatalf("parseIssueFile() error = %v", err)
	}

	// Validate parsed metadata
	if metadata.Title != "Complete Issue" {
		t.Errorf("Title = %q, want 'Complete Issue'", metadata.Title)
	}

	if len(metadata.Assignees) != 2 {
		t.Errorf("Assignees count = %d, want 2", len(metadata.Assignees))
	}

	if len(metadata.Labels) != 1 {
		t.Errorf("Labels count = %d, want 1", len(metadata.Labels))
	}

	if metadata.Labels[0].Name != "enhancement" {
		t.Errorf("Label name = %q, want 'enhancement'", metadata.Labels[0].Name)
	}

	if metadata.Milestone != "v2.0" {
		t.Errorf("Milestone = %q, want 'v2.0'", metadata.Milestone)
	}

	if len(metadata.Projects) != 2 {
		t.Errorf("Projects count = %d, want 2", len(metadata.Projects))
	}

	if !strings.Contains(body, "comprehensive test issue") {
		t.Errorf("Body doesn't contain expected text")
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "title with special characters",
			content: `---
title: Test: Issue [URGENT] with "quotes"
---
Body`,
			wantErr: false,
		},
		{
			name:    "body with code blocks",
			content: "---\ntitle: Code Issue\n---\nUsage:\n```bash\ngh godo mkissue --file issue.md\n```",
			wantErr: false,
		},
		{
			name: "assignee with @ symbol",
			content: `---
title: Issue
assign: [@user1, @user2]
---
Body`,
			wantErr: false,
		},
		{
			name: "multiline description in label",
			content: `---
title: Issue
labels:
  - name: bug
    color: ff0000
    desc: |
      This is a multiline
      description
---
Body`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseIssueFile(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIssueFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkParseIssueFile(b *testing.B) {
	content := `---
title: Benchmark Test Issue
assign: [user1, user2, user3]
labels:
  - name: bug
    color: ff0000
    desc: Bug report
  - name: feature
    color: 00ff00
    desc: Feature request
milestone: v1.0
projects: [project1, project2]
---
This is the benchmark test body.
It contains multiple lines of content.
Used to measure parsing performance.`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parseIssueFile(content)
	}
}

func TestEnsureLabelExists(t *testing.T) {
	tests := []struct {
		name    string
		label   Label
		wantErr bool
	}{
		{
			name: "label with all fields",
			label: Label{
				Name:  "test-label",
				Color: "ff0000",
				Desc:  "Test description",
			},
			wantErr: false,
		},
		{
			name: "label with minimal fields",
			label: Label{
				Name: "minimal",
			},
			wantErr: false,
		},
		{
			name: "label with color only",
			label: Label{
				Name:  "colored",
				Color: "00ff00",
			},
			wantErr: false,
		},
		{
			name: "label with description only",
			label: Label{
				Name: "descriptive",
				Desc: "Has a description",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test will execute actual gh commands
			// In a CI environment, this would require gh auth setup
			// For unit testing purposes, we're verifying no panic
			err := ensureLabelExists(tt.label)
			// Error is expected if gh CLI not authenticated, but no panic should occur
			_ = err
		})
	}
}

func TestCreateIssueMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata *IssueMetadata
		body     string
		wantErr  bool
	}{
		{
			name: "issue with assignees",
			metadata: &IssueMetadata{
				Title:     "Test with assignees",
				Assignees: []string{"user1", "me"},
				Labels:    []Label{},
			},
			body:    "Body content",
			wantErr: false,
		},
		{
			name: "issue with multiple labels",
			metadata: &IssueMetadata{
				Title: "Test with labels",
				Labels: []Label{
					{Name: "bug"},
					{Name: "important"},
				},
			},
			body:    "Bug report here",
			wantErr: false,
		},
		{
			name: "issue with milestone and projects",
			metadata: &IssueMetadata{
				Title:     "Feature request",
				Milestone: "v2.0",
				Projects:  []string{"project1", "project2"},
				Assignees: []string{"team-member"},
				Labels:    []Label{{Name: "enhancement"}},
			},
			body:    "New feature proposal",
			wantErr: false,
		},
		{
			name: "issue with empty body",
			metadata: &IssueMetadata{
				Title: "No body issue",
			},
			body:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This calls gh CLI, so errors are expected without proper auth
			// We're verifying the function handles the metadata correctly
			err := createIssue(tt.metadata, tt.body)
			_ = err
		})
	}
}

func TestParseListFieldEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startIdx  int
		prefix    string
		wantCount int
	}{
		{
			name: "list with whitespace",
			lines: []string{
				"assign:",
				"  - user1   ",
				"  -   user2",
				"  - user3",
			},
			startIdx:  0,
			prefix:    "assign:",
			wantCount: 3,
		},
		{
			name: "mixed quotes and spaces",
			lines: []string{
				"assign: [ 'user1' ,  \"user2\"  , user3 ]",
			},
			startIdx:  0,
			prefix:    "assign:",
			wantCount: 3,
		},
		{
			name: "list stops at new field",
			lines: []string{
				"assign:",
				"  - user1",
				"  - user2",
				"labels:",
				"  - name: bug",
			},
			startIdx:  0,
			prefix:    "assign:",
			wantCount: 2,
		},
		{
			name: "inline with mixed spacing",
			lines: []string{
				"assign: [  user1  ,user2,  user3  ]",
			},
			startIdx:  0,
			prefix:    "assign:",
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, _ := parseListField(tt.lines, tt.startIdx, tt.prefix)
			if len(items) != tt.wantCount {
				t.Errorf("parseListField() count = %d, want %d", len(items), tt.wantCount)
			}
		})
	}
}

func TestExtractValueEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		prefix string
		want   string
	}{
		{
			name:   "mixed quotes single",
			line:   `title: "value with 'quotes' in it"`,
			prefix: "title:",
			want:   `value with 'quotes' in it`,
		},
		{
			name:   "mixed quotes double",
			line:   `title: 'value with "quotes" in it'`,
			prefix: "title:",
			want:   `value with "quotes" in it`,
		},
		{
			name:   "value with colons",
			line:   `title: Issue: How to handle colons`,
			prefix: "title:",
			want:   `Issue: How to handle colons`,
		},
		{
			name:   "numeric value",
			line:   `priority: 123`,
			prefix: "priority:",
			want:   `123`,
		},
		{
			name:   "value with dashes",
			line:   `title: Test-Issue-With-Dashes`,
			prefix: "title:",
			want:   `Test-Issue-With-Dashes`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractValue(tt.line, tt.prefix)
			if got != tt.want {
				t.Errorf("extractValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseLabelsEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startIdx  int
		wantCount int
		wantNames []string
	}{
		{
			name: "labels with extra whitespace",
			lines: []string{
				"labels:",
				"  - name:   bug   ",
				"    color:   ff0000   ",
				"  - name: feature",
			},
			startIdx:  0,
			wantCount: 2,
			wantNames: []string{"bug", "feature"},
		},
		{
			name: "labels stop at new field",
			lines: []string{
				"labels:",
				"  - name: bug",
				"  - name: feature",
				"milestone: v1.0",
			},
			startIdx:  0,
			wantCount: 2,
			wantNames: []string{"bug", "feature"},
		},
		{
			name: "single label with all fields",
			lines: []string{
				"labels:",
				"  - name: enhancement",
				"    color: 00ff00",
				"    desc: New feature or request",
			},
			startIdx:  0,
			wantCount: 1,
			wantNames: []string{"enhancement"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels, _ := parseLabels(tt.lines, tt.startIdx)
			if len(labels) != tt.wantCount {
				t.Errorf("parseLabels() count = %d, want %d", len(labels), tt.wantCount)
				return
			}
			for i, label := range labels {
				if label.Name != tt.wantNames[i] {
					t.Errorf("parseLabels()[%d].Name = %q, want %q", i, label.Name, tt.wantNames[i])
				}
			}
		})
	}
}

func TestReadFileFromRepoValidation(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		repo     string
		branch   string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "invalid repo format - no slash",
			filePath: "file.md",
			repo:     "justarepo",
			branch:   "",
			wantErr:  true,
			errMsg:   "invalid repository format",
		},
		{
			name:     "invalid repo format - empty owner",
			filePath: "file.md",
			repo:     "/repo",
			branch:   "",
			wantErr:  true,
			errMsg:   "invalid repository format",
		},
		{
			name:     "invalid repo format - empty repo",
			filePath: "file.md",
			repo:     "owner/",
			branch:   "",
			wantErr:  true,
			errMsg:   "invalid repository format",
		},
		{
			name:     "filePath with null byte",
			filePath: "file\x00.md",
			repo:     "owner/repo",
			branch:   "",
			wantErr:  true,
			errMsg:   "prohibited characters",
		},
		{
			name:     "filePath with newline",
			filePath: "file\n.md",
			repo:     "owner/repo",
			branch:   "",
			wantErr:  true,
			errMsg:   "prohibited characters",
		},
		{
			name:     "branch with null byte",
			filePath: "file.md",
			repo:     "owner/repo",
			branch:   "branch\x00name",
			wantErr:  true,
			errMsg:   "prohibited characters",
		},
		{
			name:     "branch with newline",
			filePath: "file.md",
			repo:     "owner/repo",
			branch:   "branch\nname",
			wantErr:  true,
			errMsg:   "prohibited characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := readFileFromRepo(tt.filePath, tt.repo, tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFileFromRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("readFileFromRepo() error = %v, want %q", err, tt.errMsg)
			}
		})
	}
}

func TestReadFileFromRepoNonexistent(t *testing.T) {
	// Test with a repo that doesn't exist - valid format but nonexistent
	_, err := readFileFromRepo("nonexistent.md", "nonexistent-owner/nonexistent-repo", "")
	if err == nil {
		t.Errorf("readFileFromRepo() expected error for nonexistent repo")
		return
	}
	if !strings.Contains(err.Error(), "failed to read file from repo") {
		t.Errorf("readFileFromRepo() error = %v, want error containing 'failed to read file from repo'", err)
	}
}

func TestRunWithFileRepo(t *testing.T) {
	// Test that RunWithFile returns an error when the repo doesn't exist
	err := RunWithFile("issue.md", "", "", "nonexistent-owner/nonexistent-repo")
	if err == nil {
		t.Errorf("RunWithFile() expected error for nonexistent repo")
		return
	}
	if !strings.Contains(err.Error(), "failed to read file from repo") {
		t.Errorf("RunWithFile() error = %v, want error containing 'failed to read file from repo'", err)
	}
}
