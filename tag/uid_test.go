package tag_test

import (
	"testing"

	"github.com/dolmen-go/legodim/tag"
)

func TestUID(t *testing.T) {
	uid := tag.MustParseUID("0447710a524280") // Krusty
	pwd := uid.Pwd()
	t.Logf("%s => pwd %#x\n", uid, pwd)
	key := uid.Key()
	t.Logf("%s => key %s\n", uid, key)
	c := tag.Krusty
	data := key.EncryptCharacter(c)
	t.Logf("%d => [% X]\n", c, data)
	var data2 [8]byte
	key.Decrypt(data2[:], data)
	t.Logf("[% X] => [% X]\n", data, data2[:])
	t.Logf("[% X] => %d\n", data, key.DecryptCharacter(data))
}
