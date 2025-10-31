package sshconf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func keyPathsFor(key string) (priv, pub string) {
	home, _ := os.UserHomeDir()
	priv = filepath.Join(home, ".ssh", "id_ed25519_"+key)
	pub = priv + ".pub"
	return
}

// Assert private & public exist and .pub has ::keysej:: tag
func AssertKeysejKeyExists(key string) error {
	priv, pub := keyPathsFor(key)
	if _, err := os.Stat(priv); err != nil {
		return fmt.Errorf("private key not found: %s (generate with `keysej new %s` or use --force)", priv, key)
	}
	b, err := os.ReadFile(pub)
	if err != nil {
		return fmt.Errorf("public key not found: %s", pub)
	}
	line := string(b)
	if !strings.Contains(line, "::keysej::") {
		return errors.New("public key lacks ::keysej:: tag; use a key created by keysej or --force")
	}
	return nil
}
