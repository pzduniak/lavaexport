package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dghubble/sling"
	"github.com/lavab/api/routes"
)

const API_URL = "https://api.lavaboom.com/"

func main() {
	// Prepare a sling
	sb := sling.New().Base(API_URL)

	// Ask for username and password
	username, err := prompt("Your username: ", false)
	if err != nil {
		log.Fatal(err)
	}
	password, err := prompt("Your password: ", true)
	if err != nil {
		log.Fatal(err)
	}
	outputPath, err := prompt("Output directory: ", false)
	if err != nil {
		log.Fatal(err)
	}

	// Hash the password
	phash := sha256.Sum256([]byte(password))
	hashedPassword := hex.EncodeToString(phash[:])

	// Acquire a new token
	var authResponse *routes.TokensCreateResponse
	_, err = sb.New().Post("tokens").BodyJSON(&routes.TokensCreateRequest{
		Username: username,
		Password: hashedPassword,
		Type:     "auth",
	}).ReceiveSuccess(&authResponse)
	if err != nil {
		log.Fatal(err)
	}
	if !authResponse.Success {
		log.Fatal(authResponse.Message)
	}
	log.Printf("Acquired a new auth token: %s", authResponse.Token.ID)

	// Authorized sling
	ab := sb.New().Set("Authorization", "Bearer "+authResponse.Token.ID)

	// Prepare the output directory
	fi, err := os.Stat(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(outputPath, 0777); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		if !fi.IsDir() {
			log.Fatalf("%v is not a directory", outputPath)
		}

		if _, err := os.Stat(filepath.Join(outputPath, "account.json")); err != nil {
			if err == nil {
				log.Fatal("%v is already a Lavaboom account export", outputPath)
			} else {
				if !os.IsNotExist(err) {
					log.Fatal(err)
				}
			}
		}
	}

	// Account
	{
		var resp *routes.AccountsGetResponse
		_, err := ab.New().Get("accounts/me").ReceiveSuccess(&resp)
		if err != nil {
			log.Fatal(err)
		}
		if err := saveStruct(filepath.Join(outputPath, "account.json"), resp.Account); err != nil {
			log.Fatal(err)
		}
		log.Printf("Saved account %s", resp.Account.ID)
	}

	// Addresses
	{
		var resp *routes.AddressesListResponse
		_, err := ab.New().Get("addresses").ReceiveSuccess(&resp)
		if err != nil {
			log.Fatal(err)
		}
		if err := saveStruct(filepath.Join(outputPath, "addresses.json"), resp.Addresses); err != nil {
			log.Fatal(err)
		}
		log.Print("Saved addresses")
	}

	// Contacts
	{
		var resp *routes.ContactsListResponse
		_, err := ab.New().Get("contacts").ReceiveSuccess(&resp)
		if err != nil {
			log.Fatal(err)
		}

		os.Mkdir(filepath.Join(outputPath, "contacts"), 0777)

		for _, contact := range *resp.Contacts {
			if err := saveStruct(filepath.Join(outputPath, "contacts", contact.ID+".json"), contact); err != nil {
				log.Fatal(err)
			}

			log.Printf("Saved contact %s", contact.ID)
		}
	}

	// Emails
	fileIDs := map[string]struct{}{}
	{
		var resp *routes.EmailsListResponse
		_, err := ab.New().Get("emails").ReceiveSuccess(&resp)
		if err != nil {
			log.Fatal(err)
		}

		os.Mkdir(filepath.Join(outputPath, "emails"), 0777)

		for _, email := range *resp.Emails {
			if err := saveStruct(filepath.Join(outputPath, "emails", email.ID+".json"), email); err != nil {
				log.Fatal(err)
			}

			if email.Files != nil && len(email.Files) > 0 {
				for _, file := range email.Files {
					fileIDs[file] = struct{}{}
				}
			}

			log.Printf("Saved email %s", email.ID)
		}
	}

	// Files aren't that easy to fetch, but we have their IDs from emails
	{
		os.Mkdir(filepath.Join(outputPath, "files"), 0777)

		for file, _ := range fileIDs {
			var resp *routes.FilesGetResponse
			_, err := ab.New().Get("files/" + file).ReceiveSuccess(&resp)
			if err != nil {
				log.Fatal(err)
			}

			if err := saveStruct(filepath.Join(outputPath, "files", resp.File.ID+".json"), resp.File); err != nil {
				log.Fatal(err)
			}

			log.Printf("Saved file %s", resp.File.ID)
		}
	}

	// Fetch public keys
	{
		var resp *routes.KeysListResponse
		_, err := ab.New().Get("keys?user=" + username).ReceiveSuccess(&resp)
		if err != nil {
			log.Fatal(err)
		}

		os.Mkdir(filepath.Join(outputPath, "public_keys"), 0777)

		for _, key := range *resp.Keys {
			if err := saveStruct(filepath.Join(outputPath, "public_keys", key.ID+".json"), key); err != nil {
				log.Fatal(err)
			}

			log.Printf("Saved public key %s", key.ID)
		}
	}

	// Labels
	{
		var resp *routes.LabelsListResponse
		_, err := ab.New().Get("labels").ReceiveSuccess(&resp)
		if err != nil {
			log.Fatal(err)
		}

		os.Mkdir(filepath.Join(outputPath, "labels"), 0777)

		for _, label := range *resp.Labels {
			if err := saveStruct(filepath.Join(outputPath, "labels", label.ID+".json"), label); err != nil {
				log.Fatal(err)
			}

			log.Printf("Saved label %s", label.ID)
		}
	}

	// Threads
	{
		var resp *routes.ThreadsListResponse
		_, err := ab.New().Get("threads").ReceiveSuccess(&resp)
		if err != nil {
			log.Fatal(err)
		}

		os.Mkdir(filepath.Join(outputPath, "threads"), 0777)

		for _, thread := range *resp.Threads {
			if err := saveStruct(filepath.Join(outputPath, "threads", thread.ID+".json"), thread); err != nil {
				log.Fatal(err)
			}

			log.Printf("Saved thread %s", thread.ID)
		}
	}

	log.Print("Exporting complete. Press any key to continue.")
	fmt.Scanln()
}
