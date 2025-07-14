package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ProjectMetadata holds the metadata extracted from PROGRESS.md
type ProjectMetadata struct {
	Version             string
	Date                string
	TestsTotal          int
	TestsCoverage       int
	EndpointsFunctional int
	EndpointsPending    int
	EndpointsTotal      int
	FeaturesImplemented int
	FeaturesIntegrated  int
	Database            string
	Architecture        string
	Status              string
	PriorityNext        string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run sync-docs.go <command>")
		fmt.Println("Commands:")
		fmt.Println("  sync     - Synchronize all documentation from PROGRESS.md")
		fmt.Println("  validate - Validate documentation consistency")
		fmt.Println("  check    - Check current state without changes")
		os.Exit(1)
	}

	command := os.Args[1]
	projectRoot := getProjectRoot()

	switch command {
	case "sync":
		syncAllDocumentation(projectRoot)
	case "validate":
		validateDocumentation(projectRoot)
	case "check":
		checkDocumentationState(projectRoot)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// getProjectRoot finds the project root directory
func getProjectRoot() string {
	// Assume we're running from tools/ directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	
	// Go up one level from tools/
	return filepath.Dir(wd)
}

// extractMetadata reads PROGRESS.md and extracts automation metadata
func extractMetadata(progressPath string) (*ProjectMetadata, error) {
	file, err := os.Open(progressPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PROGRESS.md: %w", err)
	}
	defer file.Close()

	metadata := &ProjectMetadata{}
	scanner := bufio.NewScanner(file)
	inMetadataSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if strings.Contains(line, "AUTOMATION_METADATA: START") {
			inMetadataSection = true
			continue
		}
		
		if strings.Contains(line, "AUTOMATION_METADATA: END") {
			break
		}
		
		if !inMetadataSection {
			continue
		}

		// Extract metadata values
		if strings.Contains(line, "VERSION:") {
			metadata.Version = extractValue(line)
		} else if strings.Contains(line, "DATE:") {
			metadata.Date = extractValue(line)
		} else if strings.Contains(line, "TESTS_TOTAL:") {
			metadata.TestsTotal, _ = strconv.Atoi(extractValue(line))
		} else if strings.Contains(line, "TESTS_COVERAGE:") {
			metadata.TestsCoverage, _ = strconv.Atoi(extractValue(line))
		} else if strings.Contains(line, "ENDPOINTS_FUNCTIONAL:") {
			metadata.EndpointsFunctional, _ = strconv.Atoi(extractValue(line))
		} else if strings.Contains(line, "ENDPOINTS_PENDING:") {
			metadata.EndpointsPending, _ = strconv.Atoi(extractValue(line))
		} else if strings.Contains(line, "ENDPOINTS_TOTAL:") {
			metadata.EndpointsTotal, _ = strconv.Atoi(extractValue(line))
		} else if strings.Contains(line, "FEATURES_IMPLEMENTED:") {
			metadata.FeaturesImplemented, _ = strconv.Atoi(extractValue(line))
		} else if strings.Contains(line, "FEATURES_INTEGRATED:") {
			metadata.FeaturesIntegrated, _ = strconv.Atoi(extractValue(line))
		} else if strings.Contains(line, "DATABASE:") {
			metadata.Database = extractValue(line)
		} else if strings.Contains(line, "ARCHITECTURE:") {
			metadata.Architecture = extractValue(line)
		} else if strings.Contains(line, "STATUS:") {
			metadata.Status = extractValue(line)
		} else if strings.Contains(line, "PRIORITY_NEXT:") {
			metadata.PriorityNext = extractValue(line)
		}
	}

	return metadata, scanner.Err()
}

// extractValue extracts the value from a metadata comment line
func extractValue(line string) string {
	re := regexp.MustCompile(`<!-- .+:\s*(.+)\s*-->`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// syncAllDocumentation synchronizes all documentation files
func syncAllDocumentation(projectRoot string) {
	fmt.Println("🔄 Synchronizing documentation from PROGRESS.md...")
	
	progressPath := filepath.Join(projectRoot, "PROGRESS.md")
	metadata, err := extractMetadata(progressPath)
	if err != nil {
		fmt.Printf("❌ Error extracting metadata: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📊 Extracted metadata: Version %s, %d tests, %d functional endpoints\n", 
		metadata.Version, metadata.TestsTotal, metadata.EndpointsFunctional)

	// Sync CLAUDE.md
	if err := syncClaudeMD(projectRoot, metadata); err != nil {
		fmt.Printf("❌ Error syncing CLAUDE.md: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ CLAUDE.md synchronized")

	// Sync .claude/commands/README.md
	if err := syncCommandsReadme(projectRoot, metadata); err != nil {
		fmt.Printf("❌ Error syncing commands README: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ .claude/commands/README.md synchronized")

	// Sync DEVELOPMENT.md
	if err := syncDevelopmentMD(projectRoot, metadata); err != nil {
		fmt.Printf("❌ Error syncing DEVELOPMENT.md: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ DEVELOPMENT.md synchronized")

	// Sync .claude/COMMAND_EXAMPLES.md
	if err := syncCommandExamples(projectRoot, metadata); err != nil {
		fmt.Printf("❌ Error syncing COMMAND_EXAMPLES.md: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ .claude/COMMAND_EXAMPLES.md synchronized")

	fmt.Println("🎉 All documentation synchronized successfully!")
}

// syncClaudeMD updates CLAUDE.md with current metadata
func syncClaudeMD(projectRoot string, metadata *ProjectMetadata) error {
	claudePath := filepath.Join(projectRoot, "CLAUDE.md")
	
	content, err := os.ReadFile(claudePath)
	if err != nil {
		return err
	}

	text := string(content)

	// Update version
	text = regexp.MustCompile(`\*\*Versión:\*\* v[\d\.]+-[\w-]+`).ReplaceAllString(text, 
		fmt.Sprintf("**Versión:** %s", metadata.Version))

	// Update date
	text = regexp.MustCompile(`\*\*Fecha:\*\* \d{4}-\d{2}-\d{2}`).ReplaceAllString(text, 
		fmt.Sprintf("**Fecha:** %s", metadata.Date))

	// Update test count
	text = regexp.MustCompile(`\*\*Cobertura Tests:\*\* \d+%\+ promedio \(\d+ tests\)`).ReplaceAllString(text,
		fmt.Sprintf("**Cobertura Tests:** %d%%+ promedio (%d tests)", metadata.TestsCoverage, metadata.TestsTotal))

	// Update endpoint count
	text = regexp.MustCompile(`\*\*Funcionalidades:\*\* CRUD completo \+ PostgreSQL FTS \+ Sistema de Imágenes \+ Cache LRU`).ReplaceAllString(text,
		fmt.Sprintf("**Funcionalidades:** %d endpoints funcionales + %d pendientes integración", metadata.EndpointsFunctional, metadata.EndpointsPending))

	// Update endpoint claims
	text = regexp.MustCompile(`- \*\*CRUD completo:\*\* \d+ endpoints API funcionales`).ReplaceAllString(text,
		fmt.Sprintf("- **CRUD completo:** %d endpoints API funcionales", metadata.EndpointsFunctional))

	// Update test count in completadas section
	text = regexp.MustCompile(`- \*\*Testing comprehensivo:\*\* \d+ tests con \d+%\+ cobertura`).ReplaceAllString(text,
		fmt.Sprintf("- **Testing comprehensivo:** %d tests con %d%%+ cobertura", metadata.TestsTotal, metadata.TestsCoverage))

	return os.WriteFile(claudePath, []byte(text), 0644)
}

// syncCommandsReadme updates .claude/commands/README.md
func syncCommandsReadme(projectRoot string, metadata *ProjectMetadata) error {
	readmePath := filepath.Join(projectRoot, ".claude", "commands", "README.md")
	
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	text := string(content)

	// Update version
	text = regexp.MustCompile(`- \*\*Versión:\*\* v[\d\.]+-[\w-]+`).ReplaceAllString(text,
		fmt.Sprintf("- **Versión:** %s", metadata.Version))

	// Update test count
	text = regexp.MustCompile(`- \*\*\d+ tests\*\* con \d+%\+ cobertura`).ReplaceAllString(text,
		fmt.Sprintf("- **%d tests** con %d%%+ cobertura", metadata.TestsTotal, metadata.TestsCoverage))

	// Update endpoint count
	text = regexp.MustCompile(`- \*\*\d+\+ endpoints\*\* API funcionales`).ReplaceAllString(text,
		fmt.Sprintf("- **%d endpoints funcionales + %d pendientes** integración", metadata.EndpointsFunctional, metadata.EndpointsPending))

	return os.WriteFile(readmePath, []byte(text), 0644)
}

// syncDevelopmentMD updates DEVELOPMENT.md
func syncDevelopmentMD(projectRoot string, metadata *ProjectMetadata) error {
	devPath := filepath.Join(projectRoot, "DEVELOPMENT.md")
	
	content, err := os.ReadFile(devPath)
	if err != nil {
		return err
	}

	text := string(content)

	// Update date in status section
	text = regexp.MustCompile(`### Estado Actual \(\d{4}-\d{2}-\d{2}\)`).ReplaceAllString(text,
		fmt.Sprintf("### Estado Actual (%s)", metadata.Date))

	// Update test count
	text = regexp.MustCompile(`- ✅ Testing comprehensivo \(\d+ tests, \d+%\+ cobertura\)`).ReplaceAllString(text,
		fmt.Sprintf("- ✅ Testing comprehensivo (%d tests, %d%%+ cobertura)", metadata.TestsTotal, metadata.TestsCoverage))

	return os.WriteFile(devPath, []byte(text), 0644)
}

// validateDocumentation checks for inconsistencies
func validateDocumentation(projectRoot string) {
	fmt.Println("🔍 Validating documentation consistency...")
	
	progressPath := filepath.Join(projectRoot, "PROGRESS.md")
	metadata, err := extractMetadata(progressPath)
	if err != nil {
		fmt.Printf("❌ Error extracting metadata: %v\n", err)
		os.Exit(1)
	}

	// Check CLAUDE.md
	if err := validateClaudeMD(projectRoot, metadata); err != nil {
		fmt.Printf("❌ CLAUDE.md inconsistency: %v\n", err)
		os.Exit(1)
	}

	// Check commands README
	if err := validateCommandsReadme(projectRoot, metadata); err != nil {
		fmt.Printf("❌ Commands README inconsistency: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ All documentation is consistent!")
}

// validateClaudeMD checks if CLAUDE.md is consistent with metadata
func validateClaudeMD(projectRoot string, metadata *ProjectMetadata) error {
	claudePath := filepath.Join(projectRoot, "CLAUDE.md")
	content, err := os.ReadFile(claudePath)
	if err != nil {
		return err
	}

	text := string(content)

	// Check version
	if !strings.Contains(text, metadata.Version) {
		return fmt.Errorf("version mismatch: expected %s", metadata.Version)
	}

	// Check test count
	testCountStr := fmt.Sprintf("%d tests", metadata.TestsTotal)
	if !strings.Contains(text, testCountStr) {
		return fmt.Errorf("test count mismatch: expected %s", testCountStr)
	}

	return nil
}

// validateCommandsReadme checks if commands README is consistent
func validateCommandsReadme(projectRoot string, metadata *ProjectMetadata) error {
	readmePath := filepath.Join(projectRoot, ".claude", "commands", "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	text := string(content)

	// Check version
	if !strings.Contains(text, metadata.Version) {
		return fmt.Errorf("version mismatch: expected %s", metadata.Version)
	}

	return nil
}

// checkDocumentationState shows current state without making changes
func checkDocumentationState(projectRoot string) {
	fmt.Println("📋 Current documentation state:")
	
	progressPath := filepath.Join(projectRoot, "PROGRESS.md")
	metadata, err := extractMetadata(progressPath)
	if err != nil {
		fmt.Printf("❌ Error extracting metadata: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Version: %s\n", metadata.Version)
	fmt.Printf("Date: %s\n", metadata.Date)
	fmt.Printf("Tests: %d (Coverage: %d%%)\n", metadata.TestsTotal, metadata.TestsCoverage)
	fmt.Printf("Endpoints: %d functional + %d pending = %d total\n", 
		metadata.EndpointsFunctional, metadata.EndpointsPending, metadata.EndpointsTotal)
	fmt.Printf("Features: %d implemented, %d integrated\n", 
		metadata.FeaturesImplemented, metadata.FeaturesIntegrated)
	fmt.Printf("Status: %s\n", metadata.Status)
	fmt.Printf("Next Priority: %s\n", metadata.PriorityNext)
}

// syncCommandExamples updates .claude/COMMAND_EXAMPLES.md
func syncCommandExamples(projectRoot string, metadata *ProjectMetadata) error {
	examplesPath := filepath.Join(projectRoot, ".claude", "COMMAND_EXAMPLES.md")
	
	content, err := os.ReadFile(examplesPath)
	if err != nil {
		return err
	}

	text := string(content)

	// Update examples context section if it exists
	contextPattern := `## 🏠 \*\*CONTEXTO DEL PROYECTO\*\*[^#]+`
	newContext := fmt.Sprintf(`## 🏠 **CONTEXTO DEL PROYECTO**

- **Versión:** %s
- **%d tests** con %d%%+ cobertura
- **%d endpoints funcionales + %d pendientes** integración
- **Sistema completo** de imágenes con cache LRU
- **PostgreSQL FTS** en español
- **Validaciones Ecuador** integradas

---`, metadata.Version, metadata.TestsTotal, metadata.TestsCoverage, 
	metadata.EndpointsFunctional, metadata.EndpointsPending)

	re := regexp.MustCompile(contextPattern)
	if re.MatchString(text) {
		text = re.ReplaceAllString(text, newContext)
	}

	return os.WriteFile(examplesPath, []byte(text), 0644)
}