package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type SkillInput struct {
	Name string `yaml:"name"`
}

func readInputs(filePath string) (*SkillInput, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read inputs file: %v", err)
	}

	var input SkillInput
	if err := yaml.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("failed to parse inputs file: %v", err)
	}

	return &input, nil
}

func addChipToValidChips(filePath string, chipName string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Find ValidChips array
	contentStr := string(content)
	validChipsStart := strings.Index(contentStr, "ValidChips = []int{")
	if validChipsStart == -1 {
		return fmt.Errorf("ValidChips array not found")
	}

	// Find the end of the array
	validChipsEnd := strings.Index(contentStr[validChipsStart:], "}")
	if validChipsEnd == -1 {
		return fmt.Errorf("ValidChips array end not found")
	}
	validChipsEnd += validChipsStart

	// Add new chip ID before the closing brace
	newContent := contentStr[:validChipsEnd] + fmt.Sprintf("chip.ID%s,\n", chipName) + contentStr[validChipsEnd:]

	// Write the updated content back to the file
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated file: %v", err)
	}

	// Format the generated file
	if err := exec.Command("go", "fmt", filePath).Run(); err != nil {
		return fmt.Errorf("failed to format skill file: %v", err)
	}

	return nil
}

func generateSkillFile(skillName string) error {
	fileName := fmt.Sprintf("pkg/router/skill/skill_%s.go", strings.ToLower(skillName))
	const templateContent = `package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type {{.LowerName}} struct {
	ID   string
	Arg  Argument
	Core skillcore.SkillCore
}

func new{{.Name}}(arg Argument, core skillcore.SkillCore) *{{.LowerName}} {
	return &{{.LowerName}}{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core,
	}
}

func (p *{{.LowerName}}) Draw() {
	// nothing to do at router
}

func (p *{{.LowerName}}) Update() (bool, error) {
	return p.Core.Update()
}

func (p *{{.LowerName}}) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.Type{{.Name}},
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *{{.LowerName}}) StopByOwner() {
}
`

	t, err := template.New("skill_file").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// Create directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	return t.Execute(file, map[string]string{
		"Name":      skillName,
		"LowerName": strings.ToLower(skillName),
	})
}

func addAnimConst(skillName string) error {
	filePath := "pkg/router/anim/anim.go"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Find TypeMax const
	contentStr := string(content)
	typeMaxIndex := strings.Index(contentStr, "TypeMax")
	if typeMaxIndex == -1 {
		return fmt.Errorf("TypeMax not found")
	}

	// Find the line before TypeMax
	beforeTypeMax := contentStr[:typeMaxIndex]
	lastNewline := strings.LastIndex(beforeTypeMax, "\n")
	if lastNewline == -1 {
		return fmt.Errorf("failed to find position for new type")
	}

	// Insert new type before TypeMax
	newContent := contentStr[:lastNewline] + fmt.Sprintf("Type%s\n", skillName) + contentStr[lastNewline:]

	// Write the updated content back to the file
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated file: %v", err)
	}

	// Format the generated file
	if err := exec.Command("go", "fmt", filePath).Run(); err != nil {
		return fmt.Errorf("failed to format skill file: %v", err)
	}

	return nil
}

func addSkillToSwitch(skillName string) error {
	filePath := "pkg/router/skill/skill.go"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Find default case in switch statement
	contentStr := string(content)
	defaultIndex := strings.Index(contentStr, "default:")
	if defaultIndex == -1 {
		return fmt.Errorf("default case not found in switch statement")
	}

	// Find the line before default
	beforeDefault := contentStr[:defaultIndex]
	lastNewline := strings.LastIndex(beforeDefault, "\n")
	if lastNewline == -1 {
		return fmt.Errorf("failed to find position for new case")
	}

	// Insert new case before default
	newCase := fmt.Sprintf("case resources.Skill%s:\nreturn new%s(arg, core)", skillName, skillName)
	newContent := contentStr[:lastNewline] + "\n" + newCase + contentStr[lastNewline:]

	// Write the updated content back to the file
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated file: %v", err)
	}

	// Format the generated file
	if err := exec.Command("go", "fmt", filePath).Run(); err != nil {
		return fmt.Errorf("failed to format skill file: %v", err)
	}

	return nil
}

func addDrawAnimCase(skillName string) error {
	filePath := "pkg/app/game/netbattle/draw/draw_anim.go"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Find default case in switch statement
	contentStr := string(content)
	defaultIndex := strings.Index(contentStr, "default:")
	if defaultIndex == -1 {
		return fmt.Errorf("default case not found in switch statement")
	}

	// Find the line before default
	beforeDefault := contentStr[:defaultIndex]
	lastNewline := strings.LastIndex(beforeDefault, "\n")
	if lastNewline == -1 {
		return fmt.Errorf("failed to find position for new case")
	}

	// Insert new case before default
	newCase := fmt.Sprintf("case anim.Type%s:\n// no animation", skillName)
	newContent := contentStr[:lastNewline] + "\n" + newCase + contentStr[lastNewline:]

	// Write the updated content back to the file
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated file: %v", err)
	}

	// Format the generated file
	if err := exec.Command("go", "fmt", filePath).Run(); err != nil {
		return fmt.Errorf("failed to format file: %v", err)
	}

	return nil
}

func main() {
	inputs, err := readInputs("tools/skill-code-generator/network/inputs.yaml")
	if err != nil {
		fmt.Printf("Error reading inputs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generating network skill: %s\n", inputs.Name)

	netBattlePath := "pkg/app/game/netbattle/net_battle.go"
	if err := addChipToValidChips(netBattlePath, inputs.Name); err != nil {
		fmt.Printf("Error adding chip to ValidChips: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully added chip.ID%s to ValidChips\n", inputs.Name)

	if err := generateSkillFile(inputs.Name); err != nil {
		fmt.Printf("Error generating skill file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated skill file for %s\n", inputs.Name)

	if err := addAnimConst(inputs.Name); err != nil {
		fmt.Printf("Error adding animation constant: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully added Type%s to animation constants\n", inputs.Name)

	if err := addSkillToSwitch(inputs.Name); err != nil {
		fmt.Printf("Error adding skill to switch statement: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully added skill case to switch statement\n")

	if err := addDrawAnimCase(inputs.Name); err != nil {
		fmt.Printf("Error adding draw animation case: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully added draw animation case\n")

	// Format the files
	if err := exec.Command("go", "fmt", "pkg/router/anim/anim.go").Run(); err != nil {
		fmt.Printf("Error formatting anim.go: %v\n", err)
		os.Exit(1)
	}
}
