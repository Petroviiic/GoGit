package commands

import (
	"fmt"
	"log"

	"github.com/Petroviiic/GoGit/internal/core"
)

func TestFunc(repo *core.Repository) {
	fmt.Println("begin test")

	content := []byte("cao svima ja sam pera i ovo je content ovog filea.")
	b := core.NewBlob(content)
	log.Printf("normal object \n	type: %s\n	content: %s \n", b.Type, string(b.Content))

	serializedContent, _ := b.Serialize()

	deserialized, err := core.Deserialize(serializedContent)

	if err != nil {
		fmt.Println(err)
		return
	}

	if deserialized == nil {
		_ = fmt.Errorf("deserialized object is nil")
		return
	}

	log.Printf("deserialized object \n	type: %s\n	content: %s", deserialized.GetType(), string(deserialized.GetContent()))
}
