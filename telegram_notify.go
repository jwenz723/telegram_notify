package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"github.com/fsnotify/fsnotify"
	"strings"
	"path/filepath"
)

// GetUpdates is a struct to be used with the Telegram /getUpdates endpoint
type GetUpdates struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

// SendMessage is a struct to be used with the Telegram /send* endpoints
type SendMessage struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

// An HTTP client to be used for sending all Telegram API requests
var client = &http.Client{}

// getChatID should be used to obtain the chat id of the chat's that your Bot has been invited to
func getChatID(token string) {
	nurl := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", token)

	req, err := http.NewRequest("GET", nurl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		fmt.Println("Failed to build request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		fmt.Println("Failed to perform request")
	}

	defer resp.Body.Close()

	var record GetUpdates
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	j, err := json.MarshalIndent(record, "", "    ")
	if err != nil {
		fmt.Printf("Failed to parse to json -> %s\n", err)
	}
	fmt.Printf("%v\n", string(j))
}

// sendText Sends a message to the specified Telegram chatID
func sendText(token, chatID, text string) (err error) {
	// The Bot now sends a text to chatID
	url := "https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + chatID + "&text=" + text
	nurl := fmt.Sprintf(url)

	req, err := http.NewRequest("POST", nurl, nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var newrecord SendMessage

	err = json.NewDecoder(resp.Body).Decode(&newrecord)

	return
}

// sendPhoto Sends a photo and optional caption to the specified Telegram chatID
// file should be the full path to the photo file to be sent.
func sendPhoto(token, chatID, file, caption string) (err error) {
	url := "https://api.telegram.org/bot" + token + "/sendPhoto"

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("photo", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}

	// Add the chat_id
	if fw, err = w.CreateFormField("chat_id"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(chatID)); err != nil {
		return
	}

	// Add the caption
	if fw, err = w.CreateFormField("caption"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(caption)); err != nil {
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var newrecord SendMessage

	err = json.NewDecoder(resp.Body).Decode(&newrecord)

	return
}

// sendVideo Sends a video and optional caption to the specified Telegram chatID.
// file should be the full path to the video file to be sent.
func sendVideo(token, chatID, file, caption string) (err error) {
	url := "https://api.telegram.org/bot" + token + "/sendVideo"

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("video", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}

	// Add the chat_id
	if fw, err = w.CreateFormField("chat_id"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(chatID)); err != nil {
		return
	}

	// Add the caption
	if fw, err = w.CreateFormField("caption"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(caption)); err != nil {
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var newrecord SendMessage

	err = json.NewDecoder(resp.Body).Decode(&newrecord)

	return
}

// sendMediaGroup is designed to send multiple photos in a single message
// This function is NOT functional
func sendMediaGroup(token, chatid, file, caption string) (err error) {
	url := "https://api.telegram.org/bot" + token + "/sendPhoto"

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("photo", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}

	// Add the chat_id
	if fw, err = w.CreateFormField("chat_id"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(chatid)); err != nil {
		return
	}

	// Add the caption
	if fw, err = w.CreateFormField("caption"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(caption)); err != nil {
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	var newrecord SendMessage

	if err := json.NewDecoder(resp.Body).Decode(&newrecord); err != nil {
		log.Println(err)
	}

	//fmt.Printf("%+v\n", newrecord)

	return
}

var (
	chatID   string
	filePath string
	mode     string
	text     string
	botToken string
)

func init() {
	flag.StringVar(&chatID, "chatID", "", "specify the id of the Telegram chat that messages should be sent to")
	flag.StringVar(&filePath, "filePath", "", "specify the path to the file to be uploaded. For mode 'watcher' this specifies the directory to watch.")
	flag.StringVar(&mode, "mode", "", "specify 'getChatID', 'sendText', 'sendPhoto', 'sendVideo', or 'watcher'")
	flag.StringVar(&text, "text", "", "specify text to send")
	flag.StringVar(&botToken, "botToken", "", "specify the telegram bot api token to be used for sending messages")
	flag.Parse()
}

func main() {
	if mode == "" {
		fmt.Printf("Must specify mode flag!\n")
		return
	}

	if botToken == "" {
		fmt.Printf("Must specify botToken flag!\n")
		return
	}

	if mode == "getChatID" {
		getChatID(botToken)
	} else {
		if chatID == "" {
			fmt.Printf("Must specify chatID flag!\n")
			return
		}

		if mode == "sendText" {
			if text == "" {
				fmt.Printf("Must specify text flag!\n")
				return
			}

			err := sendText(botToken, chatID, text)
			if err != nil {
				fmt.Printf("Error occurred in sendText() -> %s", err)
			} else {
				fmt.Println("Text sent")
			}
		} else if mode == "sendVideo" {
			if filePath == "" {
				fmt.Printf("Must specify filePath flag!\n")
				return
			}

			err := sendVideo(botToken, chatID, filePath, text)
			if err != nil {
				fmt.Printf("Error occurred in sendVideo() -> %s", err)
			} else {
				fmt.Println("Video sent")
			}
		} else if mode == "sendPhoto" {
			if filePath == "" {
				fmt.Printf("Must specify filePath flag!\n")
				return
			}

			err := sendPhoto(botToken, chatID, filePath, text)
			if err != nil {
				fmt.Printf("Error occurred in sendPhoto() -> %s", err)
			} else {
				fmt.Println("Photo sent")
			}
		} else if mode == "watcher" {
			if filePath == "" {
				fmt.Printf("Must specify filePath flag!\n")
				return
			}

			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal(err)
			}
			defer watcher.Close()

			done := make(chan struct{})
			go func() {
				armed := false

				for {
					select {
					case event := <-watcher.Events:
						if event.Op == fsnotify.Create {
							if strings.HasSuffix(event.Name, ".jpg") {
								if armed {
									log.Printf("file created: %s\n", event.Name)
									err := sendPhoto(botToken, chatID, event.Name, text)
									if err != nil {
										fmt.Printf("Error occurred in watcher sendPhoto() -> %s", err)
									} else {
										fmt.Println("Photo sent")
									}
								}
							} else if event.Name == filepath.Join(filePath, "motionstarted") {
								armed = true
								log.Println("armed")
							}
						} else if event.Op == fsnotify.Remove {
							if event.Name == filepath.Join(filePath, "motionstarted") {
								armed = false
								log.Println("disarmed")
							}
						}
					case err := <-watcher.Errors:
						log.Println("error:", err)
					}
				}
			}()

			err = watcher.Add(filePath)
			if err != nil {
				log.Fatal(err)
			}
			<-done
		} else {
			fmt.Println("invalid mode parameter. use -help to view available values.")
		}
	}
}
