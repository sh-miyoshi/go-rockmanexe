package player

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	plyr := Player{
		HP:        100,
		ShotPower: 5,
	}

	tmpfile, _ := os.CreateTemp("", "")
	defer os.Remove(tmpfile.Name())

	// Save with no encryption
	if err := plyr.Save(tmpfile.Name(), nil); err != nil {
		t.Errorf("Failed to save data: %+v", err)
	}

	tmpfile.Close()

	var res SaveData
	rawRes, _ := os.ReadFile(tmpfile.Name())
	json.Unmarshal(rawRes, &res)
	if res.Player.HP != plyr.HP || res.Player.ShotPower != plyr.ShotPower {
		t.Errorf("Player info save failed. expect %+v, but got %+v", plyr, res)
	}
}

func TestNewWithSaveData(t *testing.T) {
	tt := []struct {
		caseName string
		input    string
		expectOK bool
	}{
		{
			"v0.3 save data",
			`{"player":{"hp":200,"shot_power":1,"zenny":0,"chip_folder":[{"id":1,"code":"b"},{"id":1,"code":"b"},{"id":1,"code":"c"},{"id":1,"code":"c"},{"id":2,"code":"d"},{"id":2,"code":"d"},{"id":44,"code":"l"},{"id":44,"code":"l"},{"id":44,"code":"*"},{"id":44,"code":"*"},{"id":54,"code":"s"},{"id":54,"code":"s"},{"id":54,"code":"s"},{"id":54,"code":"s"},{"id":55,"code":"s"},{"id":55,"code":"s"},{"id":109,"code":"l"},{"id":109,"code":"l"},{"id":109,"code":"*"},{"id":109,"code":"*"},{"id":110,"code":"l"},{"id":110,"code":"l"},{"id":5,"code":"b"},{"id":5,"code":"b"},{"id":5,"code":"d"},{"id":5,"code":"d"},{"id":8,"code":"n"},{"id":8,"code":"n"},{"id":8,"code":"m"},{"id":8,"code":"m"}],"win_num":0,"play_count":2101,"back_pack":[],"battle_histories":[{"opponent_id":"804f3fc5-328f-49c6-969d-ac559d95108a","date":"2021-10-16T19:29:09.1331744+09:00","is_win":true}]},"program_version":"v0.3"}`,
			false,
		},
		{
			"v0.13 save data",
			`{"player":{"hp":300,"shot_power":1,"charge_time":180,"zenny":0,"chip_folder":[{"id":1,"code":"b"},{"id":1,"code":"b"},{"id":1,"code":"c"},{"id":1,"code":"c"},{"id":111,"code":"*"},{"id":111,"code":"*"},{"id":111,"code":"*"},{"id":111,"code":"*"},{"id":229,"code":"s"},{"id":229,"code":"s"},{"id":229,"code":"b"},{"id":229,"code":"b"},{"id":8,"code":"n"},{"id":8,"code":"n"},{"id":54,"code":"s"},{"id":54,"code":"s"},{"id":55,"code":"s"},{"id":55,"code":"s"},{"id":5,"code":"b"},{"id":5,"code":"b"},{"id":18,"code":"m"},{"id":9,"code":"f"},{"id":21,"code":"f"},{"id":21,"code":"f"},{"id":69,"code":"m"},{"id":69,"code":"m"},{"id":75,"code":"n"},{"id":75,"code":"n"},{"id":106,"code":"*"},{"id":106,"code":"*"}],"win_num":1,"play_count":815,"back_pack":[{"id":2,"code":"d"},{"id":2,"code":"d"},{"id":3,"code":"g"},{"id":44,"code":"*"},{"id":44,"code":"*"},{"id":44,"code":"*"},{"id":44,"code":"*"},{"id":112,"code":"g"},{"id":114,"code":"h"},{"id":54,"code":"s"},{"id":54,"code":"s"},{"id":55,"code":"s"},{"id":55,"code":"s"},{"id":107,"code":"*"},{"id":107,"code":"*"},{"id":108,"code":"*"},{"id":108,"code":"*"},{"id":229,"code":"*"}],"battle_histories":[],"navi_custom_parts":[{"id":0,"is_set":false,"x":0,"y":0},{"id":1,"is_set":false,"x":0,"y":0},{"id":2,"is_set":false,"x":0,"y":0},{"id":3,"is_set":false,"x":0,"y":0},{"id":5,"is_set":false,"x":0,"y":0}]},"program_version":"v0.13"}`,
			true,
		},
		{
			"development save data",
			`{"player":{},"program_version":"development"}`,
			true,
		},
		{
			"invalid save data",
			`{"invalid": "test message"}`,
			false,
		},
	}

	tmpfile, _ := os.CreateTemp("", "")
	defer os.Remove(tmpfile.Name())

	for _, tc := range tt {
		os.WriteFile(tmpfile.Name(), []byte(tc.input), 0644)
		_, err := NewWithSaveData(tmpfile.Name(), nil)
		if tc.expectOK && err != nil {
			t.Errorf("Test case %s expects ok, but got %v", tc.caseName, err)
		}
		if !tc.expectOK && err == nil {
			t.Errorf("Test case %s expects error, but got nil", tc.caseName)
		}
	}
}
