
# telegram_notify
This repo contains a simple Go application that can be used for sending telegram messages containing an optional photo or video.

## Usage
1. Download the executable for your platform from the build directory within this repo. If the executable you want is not in the build directory, then clone the repo and build it for the desired platform.
2. Place the executable somewhere you can run it
3. Run one of the following commands:

* Get the chat ID of a Telegram chat that your Telegram Bot is in
    * Command:
        ```
        telegram_notify -botToken <my Telegram bot token here> -mode getChatID
        ```
    * Response (the desired chat ID is on the 15th line at result.message.chat.id):
        ```json
        {
            "ok": true,
            "result": [
                {
                    "update_id": 11111111,
                    "message": {
                        "message_id": 1,
                        "from": {
                            "id": 111111111,
                            "first_name": "MyFirstName",
                            "last_name": "MyLastName",
                            "username": "MyUsername"
                        },
                        "chat": {
                            "id": 186939052,
                            "first_name": "MyFirstName",
                            "last_name": "MyLastName",
                            "username": "MyUsername",
                            "type": "private"
                        },
                        "date": 1522189532,
                        "text": "Test"
                    }
                }
            ]
        }
        ```
* Send a message with text
    * Command: 
        ```
        telegram_notify -botToken <my Telegram bot token here> -chatID <my Telegram chat ID> -mode sendText -text "my message"
        ```
    * Response:
        ```
        Text sent
        ```
* Send a photo
    * Command: 
        ```
        telegram_notify -botToken <my Telegram bot token here> -chatID <my Telegram chat ID> -mode sendPhoto -filePath "path to photo file"
        ```
    * Response:
        ```
        Photo sent
        ```
* Send a photo with a caption
    * Command:
        ```
        telegram_notify -botToken <my Telegram bot token here> -chatID <my Telegram chat ID> -mode sendPhoto -filePath "path to photo file" -text "my optional caption"
        ```
    * Response:
        ```
        Photo sent
        ```
* Send a video
    * Command:
        ```
        telegram_notify -botToken <my Telegram bot token here> -chatID <my Telegram chat ID> -mode sendVideo -filePath "path to video file"
        ```
    * Response:
        ```
        Video sent
        ```
* Send a video with a caption
    * Command:
        ```
        telegram_notify -botToken <my Telegram bot token here> -chatID <my Telegram chat ID> -mode sendVideo -filePath "path to video file" -text "my optional caption"
        ```
    * Response:
        ```
        Video sent
        ```

### Manually Building Executable
When compiling for raspberry pi 2 use:

Windows:
```
set GOARCH=arm
set GOOS=linux
set GOARM=5
go build
```

Mac
```
env GOOS=linux GOARCH=arm GOARM=5 go build
```