package player

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	player := Player{
		HP:         1,
		ShotPower:  2,
		ChipFolder: [FolderSize]ChipInfo{},

		WinNum: 3,
	}
	expect := "development#0#1#2#3#"
	for i := 0; i < FolderSize; i++ {
		expect += "0#"
	}

	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	// Save with no encryption
	if err := player.Save(tmpfile.Name(), nil); err != nil {
		t.Errorf("Failed to save data: %+v", err)
	}

	tmpfile.Close()
	res, _ := ioutil.ReadFile(tmpfile.Name())
	if string(res) != expect {
		t.Errorf("Player info save failed. expect %s, but got %s", expect, string(res))
	}
}
