package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// decodeSubject decodes an encoded subject line
func decodeSubject(encoded string) (string, error) {
	parts := strings.Split(encoded, "?")
	if len(parts) != 5 {
		return "", fmt.Errorf("invalid encoded-word format")
	}

	encoding := parts[2]
	encodedText := parts[3]

	var decoded []byte
	var err error
	switch encoding {
	case "B":
		decoded, err = base64.StdEncoding.DecodeString(encodedText)
		if err != nil {
			return "", fmt.Errorf("base64 decode error: %v", err)
		}
	case "Q":
		decoded, err = decodeQ(encodedText)
		if err != nil {
			return "", fmt.Errorf("quoted-printable decode error: %v", err)
		}
	default:
		return "", fmt.Errorf("unsupported encoding: %v", encoding)
	}

	var decodedStr string
	switch parts[1] {
	case "iso-2022-jp":
		decodedStr, err = decodeISO2022JP(decoded)
		if err != nil {
			return "", fmt.Errorf("ISO-2022-JP decode error: %v", err)
		}
	default:
		decodedStr = string(decoded)
	}

	return decodedStr, nil
}

// decodeISO2022JP decodes ISO-2022-JP encoded bytes to a string
func decodeISO2022JP(encoded []byte) (string, error) {
	decoder := japanese.ISO2022JP.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(encoded), decoder)
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// decodeQ decodes a quoted-printable encoded string
func decodeQ(encoded string) ([]byte, error) {
	encoded = strings.ReplaceAll(encoded, "_", " ")
	encoded = strings.ReplaceAll(encoded, "=", "")
	return []byte(encoded), nil
}

func main() {
	// Connect to IMAP server
	c, err := client.DialTLS(SERVER, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// login to IMAP server
	if err := c.Login(ID, PASSWORD); err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// Serarch all unread messages from INBOX
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	if len(ids) == 0 {
		fmt.Println("No unread messages")
		return
	}

	// Fetch unread messages
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	section := &imap.BodySectionName{Peek: true}
	items := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	go func() {
		if err := c.Fetch(seqset, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Unread messages:")
	for msg := range messages {
		subject := msg.Envelope.Subject
		// fmt.Println("Original Subject:", subject) // print original subject
		// decode if it is necessary 
		latestNum, err := GetLatestNumber();
		if err != nil {
			log.Fatal(err);
		}
		if latestNum >= int(msg.SeqNum) {
			// No unread message
			fmt.Println("No unread message")
			return 
		}
		var pushMessage string
		pushMessage = "From: "
		for _, from := range msg.Envelope.From {
			pushMessage += fmt.Sprintf("%s\n", from.Address())
		}

		if strings.HasPrefix(subject, "=?") {
			decodedSubject, err := decodeSubject(subject)
			if err != nil {
				log.Printf("Failed to decode subject: %v", err)
				fmt.Println("Subject:", subject)
			} else {
				fmt.Println("Decoded Subject:", decodedSubject)				
				pushMessage += fmt.Sprintf("Subject: %s", decodedSubject)
				FetchToLine(pushMessage)
			}
		} else { /* if decoding is not necessary */
			fmt.Println("Subject:", subject)
			pushMessage += fmt.Sprintf("Subject: %s", subject)
			FetchToLine(pushMessage)
		}
		UpdateLatestNum(msg.SeqNum)
	}
}
