// Copyright jean-françois PHILIPPE 2014-2018
// Messages send to Collecte Through a domain socket

package pipe

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

const (
	// SocketName default socket name
	SocketName string = "/tmp/collecte.sock"
	// Secret default secret key
	Secret string = "secret_key_01"
	// CmdStop commande pour arreter
	CmdStop string = "stop"
	// CmdLogRotate force la rotation des logs
	CmdLogRotate string = "log_rotate"
	// DeltaMin intervalle de temps de validite d une commande.
	DeltaMin int64 = 2
)

// Msg Message transmited through pipe.
type Msg struct {
	Version string    `json:"version"`
	When    time.Time `json:"when"`
	Cmde    string    `json:"cmde"`
	Secure  string    `json:"hash"`
}

// ComputeHMAC  Compute an unique signature of a Message
func (m *Msg) ComputeHMAC() string {
	// Generate a key
	keyHash := md5.New()
	io.WriteString(keyHash, Secret)
	key := keyHash.Sum(nil)

	sig := hmac.New(sha256.New, key)
	io.WriteString(sig, m.Version)
	io.WriteString(sig, m.Cmde)
	io.WriteString(sig, m.When.Format(time.RFC3339))
	return hex.EncodeToString(sig.Sum(nil))
}

// New  New Message
func New(cmde string) (*Msg, error) {
	msg := &Msg{"1", time.Now(), cmde, ""}
	msg.Secure = msg.ComputeHMAC()
	return msg, nil
}

// Validate Check wether the message is valid or not.
//  return nil if OK.
func (m *Msg) Validate() error {
	// First Version
	if m.Version != "1" {
		return errors.New("Unsupported Message Version")
	}
	// Then outdated Msg
	dateRef := time.Now().Add(-time.Duration(DeltaMin) * time.Minute)
	if dateRef.After(m.When) {
		return errors.New("Outdated Message")
	}
	// Future Message
	dateRef = time.Now().Add(time.Duration(DeltaMin) * time.Minute)
	if dateRef.Before(m.When) {
		return errors.New("Future Message")
	}

	// Finally Check HMAC
	if m.Secure != m.ComputeHMAC() {
		return errors.New("Wrong Signature")
	}

	return nil
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
