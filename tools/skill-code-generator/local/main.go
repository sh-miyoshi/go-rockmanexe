package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

// Chip represents a battle chip
type Chip struct {
	ID            int    `yaml:"id"`
	Name          string `yaml:"name"`
	Power         int    `yaml:"power"`
	Type          int    `yaml:"type"`
	PlayerAct     int    `yaml:"player_act"`
	Description   string `yaml:"description"`
	IsImplemented bool   `yaml:"is_implemented"`
	IconIndex     int    `yaml:"icon_index,omitempty"`
	KeepCnt       int    `yaml:"keep_cnt,omitempty"`
	ForMe         bool   `yaml:"for_me,omitempty"`
}

// ChipList represents the list of chips from chipList.yaml
type ChipList []Chip

func generateConstName(name string) string {
	table := map[string]string{
		"フルカスタム":  "FullCustom",
		"エアホッケー1": "AirHockey1",
		"エアホッケー2": "AirHockey2",
		"エアホッケー3": "AirHockey3",
	}

	return table[name]
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func updateResourcesFile(skillName string) error {
	constName := generateConstName(skillName)
	if constName == "" {
		return fmt.Errorf("failed to generate const name for %s", skillName)
	}

	resourcePath := filepath.Join("pkg", "app", "resources", "skill.go")
	content, err := os.ReadFile(resourcePath)
	if err != nil {
		return fmt.Errorf("failed to read skill.go: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" && strings.TrimSpace(lines[i+1]) == "SkillFailed" {
			// Insert new skill before SkillFailed with blank line
			tmp := make([]string, len(lines[i:]))
			copy(tmp, lines[i:])
			lines = append(lines[:i], fmt.Sprintf("Skill%s", constName))
			lines = append(lines, "") // Keep blank line
			lines = append(lines, tmp...)
			break
		}
	}

	err = os.WriteFile(resourcePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to skill.go: %v", err)
	}

	// Format the file
	if err := exec.Command("go", "fmt", resourcePath).Run(); err != nil {
		return fmt.Errorf("failed to format skill.go: %v", err)
	}

	return nil
}

func updateManagerFile(skillName string) error {
	constName := generateConstName(skillName)
	if constName == "" {
		return fmt.Errorf("failed to generate const name for %s", skillName)
	}

	managerPath := filepath.Join("pkg", "app", "skillcore", "manager", "manager.go")
	content, err := os.ReadFile(managerPath)
	if err != nil {
		return fmt.Errorf("failed to read manager.go: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.Contains(line, "default:") {
			// Insert new case before default
			newCase := fmt.Sprintf("case resources.Skill%s:\nreturn &processor.%s{Arg: arg}", constName, constName)
			tmp := make([]string, len(lines[i:]))
			copy(tmp, lines[i:])
			lines = append(lines[:i], newCase)
			lines = append(lines, tmp...)
			break
		}
	}

	err = os.WriteFile(managerPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to manager.go: %v", err)
	}

	// Format the file
	if err := exec.Command("go", "fmt", managerPath).Run(); err != nil {
		return fmt.Errorf("failed to format manager file: %v", err)
	}

	return nil
}

func generateProcessorTemplate(skillName string) error {
	const processorTemplate = `package processor

import (
"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type {{.Name}} struct {
Arg skillcore.Argument

count int
}

func (p *{{.Name}}) Update() (bool, error) {
p.count++
return true, nil
}

func (p *{{.Name}}) GetCount() int {
return p.count
}`

	constName := generateConstName(skillName)
	if constName == "" {
		return fmt.Errorf("failed to generate const name for %s", skillName)
	}

	t, err := template.New("processor_file").Parse(processorTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	filePath := filepath.Join("pkg", "app", "skillcore", "processor", "skill_"+toSnakeCase(constName)+".go")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	err = t.Execute(file, map[string]string{
		"Name": constName,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Format the generated file
	if err := exec.Command("go", "fmt", filePath).Run(); err != nil {
		return fmt.Errorf("failed to format processor file: %v", err)
	}

	return nil
}

func generateSkillTemplate(skillName string) error {
	const templateContent = `package skill

import (
"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type {{.Name}} struct {
ID      string
Arg     skillcore.Argument
Core    *processor.{{.Name}}
animMgr *manager.Manager
}

func new{{.Name}}(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *{{.Name}} {
return &{{.Name}}{
ID:      objID,
Arg:     arg,
Core:    core.(*processor.{{.Name}}),
animMgr: animMgr,
}
}

func (p *{{.Name}}) Draw() {
// TODO: implement draw method
}

func (p *{{.Name}}) Update() (bool, error) {
return p.Core.Update()
}

func (p *{{.Name}}) GetParam() anim.Param {
return anim.Param{
ObjID: p.ID,
}
}

func (p *{{.Name}}) StopByOwner() {
p.animMgr.AnimDelete(p.ID)
}`

	constName := generateConstName(skillName)
	if constName == "" {
		return fmt.Errorf("failed to generate const name for %s", skillName)
	}

	t, err := template.New("skill_file").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	filePath := filepath.Join("pkg", "app", "game", "battle", "skill", "skill_"+toSnakeCase(constName)+".go")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	err = t.Execute(file, map[string]string{
		"Name": constName,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Format the generated file
	if err := exec.Command("go", "fmt", filePath).Run(); err != nil {
		return fmt.Errorf("failed to format skill file: %v", err)
	}

	return nil
}

func updateChipGoFile(skill Chip) error {
	chipGoPath := "pkg/app/chip/chip.go"
	content, err := os.ReadFile(chipGoPath)
	if err != nil {
		return fmt.Errorf("failed to read chip.go: %v", err)
	}

	// 新しい定数の定義を生成
	constName := generateConstName(skill.Name)
	newConst := fmt.Sprintf("ID%s   = %d", constName, skill.ID)

	// ファイルの内容を文字列として処理
	lines := strings.Split(string(content), "\n")
	constInserted := false

	// 適切な位置に新しい定数を挿入
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], fmt.Sprintf("ID%s ", constName)) {
			// 既に定義されている場合はスキップ
			return nil
		}

		// IDが昇順になるように適切な位置に挿入
		if strings.Contains(lines[i], "// Program Advance") {
			// Insert before Program Advance section
			tmp := make([]string, len(lines[i:]))
			copy(tmp, lines[i:])
			newLines := append(lines[:i], newConst)
			newLines = append(newLines, tmp...)
			lines = newLines
			constInserted = true
			break
		} else if strings.Contains(lines[i], " = ") {
			parts := strings.Split(strings.TrimRight(lines[i], "\r\n"), " = ")
			if len(parts) == 2 {
				id, err := strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("failed to parse ID: %v", err)
				}
				if id > skill.ID {
					fmt.Printf("Inserting new constant before %s\n", lines[i])
					tmp := make([]string, len(lines[i:]))
					copy(tmp, lines[i:])
					newLines := append(lines[:i], newConst)
					newLines = append(newLines, tmp...)
					lines = newLines
					constInserted = true
					break
				}
			}
		}
	}

	if !constInserted {
		return fmt.Errorf("failed to find appropriate position to insert constant")
	}

	// ファイルに書き戻す
	err = os.WriteFile(chipGoPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to chip.go: %v", err)
	}
	// Format the generated file
	if err := exec.Command("go", "fmt", chipGoPath).Run(); err != nil {
		return fmt.Errorf("failed to format skill file: %v", err)
	}

	return nil
}

func main() {
	// Read input skill file
	yamlFile, err := os.ReadFile("tools/skill-code-generator/local/inputs.yaml")
	if err != nil {
		log.Fatalf("Failed to read input YAML file: %v", err)
	}

	// Parse input skill
	var skill Chip
	err = yaml.Unmarshal(yamlFile, &skill)
	if err != nil {
		log.Fatalf("Failed to parse input skill: %v", err)
	}

	if generateConstName(skill.Name) == "" {
		log.Fatalf("Failed to generate const name for %s", skill.Name)
	}

	// Read chipList.yaml
	chipListFile, err := os.ReadFile("data/chipList.yaml")
	if err != nil {
		log.Fatalf("Failed to read chipList.yaml: %v", err)
	}

	// Parse chipList
	var chipList ChipList
	err = yaml.Unmarshal(chipListFile, &chipList)
	if err != nil {
		log.Fatalf("Failed to parse chipList: %v", err)
	}

	// Update matching chip
	updated := false
	for i, chip := range chipList {
		if chip.ID == skill.ID {
			chipList[i] = skill
			chipList[i].IsImplemented = true // 自動的に実装済みとしてマーク
			updated = true
			fmt.Printf("Updated chip ID %d: %s\n", skill.ID, skill.Name)
			break
		}
	}

	if !updated {
		fmt.Printf("Warning: Chip ID %d not found in chipList.yaml\n", skill.ID)
		return
	}

	// Write updated chipList back to file
	output, err := yaml.Marshal(chipList)
	if err != nil {
		log.Fatalf("Failed to marshal updated chipList: %v", err)
	}

	err = os.WriteFile("data/chipList.yaml", output, 0644)
	if err != nil {
		log.Fatalf("Failed to write updated chipList: %v", err)
	}

	fmt.Println("Successfully updated chipList.yaml")

	// Update chip.go with new constant
	err = updateChipGoFile(skill)
	if err != nil {
		log.Fatalf("Warning: Failed to update chip.go: %v", err)
	} else {
		fmt.Println("Successfully updated pkg/app/chip/chip.go")
	}

	// Generate processor template
	err = generateProcessorTemplate(skill.Name)
	if err != nil {
		log.Fatalf("Warning: Failed to generate processor template: %v", err)
	} else {
		fmt.Printf("Successfully generated processor template\n")
	}

	// Generate skill template
	err = generateSkillTemplate(skill.Name)
	if err != nil {
		log.Fatalf("Warning: Failed to generate skill template: %v", err)
	} else {
		fmt.Printf("Successfully generated skill template\n")
	}

	// Update resources/skill.go with new skill constant
	err = updateResourcesFile(skill.Name)
	if err != nil {
		log.Fatalf("Warning: Failed to update resources/skill.go: %v", err)
	} else {
		fmt.Printf("Successfully updated resources/skill.go\n")
	}

	// Update manager.go with new processor case
	err = updateManagerFile(skill.Name)
	if err != nil {
		log.Fatalf("Warning: Failed to update manager.go: %v", err)
	} else {
		fmt.Printf("Successfully updated manager.go\n")
	}

	// WIP: Update pkg\app\skillcore\skill.go with new GetIDByChipID case
	// WIP: Update pkg\app\game\battle\skill\skill.go Get method
}
