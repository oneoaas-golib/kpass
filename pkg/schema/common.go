package schema

import (
	"encoding/json"

	"github.com/seccom/kpass/pkg/util"
	"github.com/tidwall/buntdb"
)

const (
	// KeyPrefixUser ...
	keyPrefixUser   = "U:"
	keyPrefixTeam   = "T:"
	keyPrefixEntry  = "E:"
	keyPrefixSecret = "S:"
	keyPrefixShare  = "SH:"
)

// InitIndex ...
func InitIndex(DB *buntdb.DB) {
	DB.CreateIndex("user_by_id", "U:*", buntdb.IndexJSON("id"))
	DB.CreateIndex("entry_by_team", "E:*", buntdb.IndexJSON("teamID"))
	DB.CreateIndex("team_by_user", "T:*", buntdb.IndexJSON("userID"))
	DB.CreateIndex("share_by_user", "SH:*", buntdb.IndexJSON("userID"))
	DB.CreateIndex("share_by_entry", "SH:*", buntdb.IndexJSON("entryID"))
	DB.CreateIndex("share_by_team", "SH:*", buntdb.IndexJSON("teamID"))
}

// UserKey returns the user's db key
func UserKey(name string) string {
	return keyPrefixUser + name
}

// TeamKey returns the team's db key
func TeamKey(id util.OID) string {
	return keyPrefixTeam + id.String()
}

// TeamIDFromKey returns team' ID from key
func TeamIDFromKey(key string) util.OID {
	val := key[len(keyPrefixTeam):]
	id, err := util.ParseOID(val)
	if err != nil {
		panic(err)
	}
	return id
}

// EntryKey returns the entry's db key
func EntryKey(id util.OID) string {
	return keyPrefixEntry + id.String()
}

// EntryIDFromKey returns entry' ID from key
func EntryIDFromKey(key string) util.OID {
	val := key[len(keyPrefixEntry):]
	id, err := util.ParseOID(val)
	if err != nil {
		panic(err)
	}
	return id
}

// SecretKey returns the secret's db key
func SecretKey(id util.OID) string {
	return keyPrefixSecret + id.String()
}

// ShareKey returns the share's db key
func ShareKey(id util.OID) string {
	return keyPrefixShare + id.String()
}

// ShareIDFromKey returns share' ID from key
func ShareIDFromKey(key string) util.OID {
	val := key[len(keyPrefixShare):]
	id, err := util.ParseOID(val)
	if err != nil {
		panic(err)
	}
	return id
}

func jsonMarshal(v interface{}) (str string) {
	if res, err := json.Marshal(v); err == nil {
		str = string(res)
	}
	return
}

// StringSlice ...
type StringSlice []string

// Has returns whether the str is in the slice.
func (s StringSlice) Has(str string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == str {
			return true
		}
	}
	return false
}

// Add adds the str to the slice.
func (s StringSlice) Add(str string) ([]string, bool) {
	if s.Has(str) {
		return s, false
	}
	return append(s, str), true
}

// Remove remove the str from the slice.
func (s StringSlice) Remove(str string) ([]string, bool) {
	offset := 0
	for i := 0; i < len(s); i++ {
		if s[i] != str {
			s[offset] = s[i]
			offset++
		}
	}
	if offset < len(s) {
		return s[0:offset], true
	}
	return s, false
}
