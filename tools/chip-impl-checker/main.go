package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle"
	"github.com/stretchr/stew/slice"
	"gopkg.in/yaml.v2"
)

type Chip struct {
	ID            int    `yaml:"id"`
	Name          string `yaml:"name"`
	IsImplemented bool   `yaml:"is_implemented"`
}

func main() {
	// Read chip list
	chipList, err := loadChips("data/chipList.yaml")
	if err != nil {
		log.Fatalf("failed to load chips: %v", err)
	}

	outFile, err := os.Create("tmp/docs/未実装チップ.md")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	outFile.WriteString("# 未実装チップ\n\n")
	for _, c := range chipList {
		if c.IsImplemented {
			continue
		}
		if c.Name == "" {
			outFile.WriteString(fmt.Sprintf("- %dのチップ\n", c.ID))
		} else {
			outFile.WriteString(fmt.Sprintf("- %s\n", c.Name))
		}
	}

	outNetFile, err := os.Create("tmp/docs/未実装ネットワーク_チップ.md")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outNetFile.Close()

	outNetFile.WriteString("# 通信対戦時の未実装チップ\n\n")
	for _, c := range chipList {
		if !c.IsImplemented {
			continue
		}
		if slice.Contains(netbattle.ValidChips, c.ID) {
			continue
		}

		outNetFile.WriteString(fmt.Sprintf("- %s\n", c.Name))
	}
}

func loadChips(fname string) ([]Chip, error) {
	// Read chip list
	buf, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read chip list: %w", err)
	}

	var chips []Chip
	if err := yaml.Unmarshal(buf, &chips); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chip list: %w", err)
	}

	return chips, nil
}
