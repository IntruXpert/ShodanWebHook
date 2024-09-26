# Shodan Telegram Notifier

A Go application that listens for Shodan webhooks, processes the received data, captures screenshots of discovered services, and sends them to a specified Telegram channel.

## Features

- **Shodan Webhook Listener**: Listens for incoming webhooks from Shodan triggers.
- **Data Storage**: Stores processed data in a SQLite database to avoid duplicate processing.
- **Screenshot Capture**: Uses `chromedp` to take screenshots of services over HTTP and HTTPS.
- **Telegram Integration**: Sends screenshots and service details to a Telegram channel.
- **Duplicate Detection**: Checks the database to prevent reprocessing of the same IP and port.

## Prerequisites

- Go installed on your system.
- A Shodan account with webhook capabilities.
- A Telegram bot and channel (bot must be an admin in the channel).
- SQLite3 installed if you want to inspect the database.

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/IntruXpert/ShodanWebHook.git
cd ShodanWebHook
```

### 2. Install Dependencies

Ensure you have Go modules enabled.

```bash
go mod tidy
```

### 3. Configure the Application

- Open `main.go` and replace placeholder values:

  ```go
  var channelID int64 = YOUR_TELEGRAM_CHANNEL_ID // Replace with your Telegram channel ID
  ```

  ```go
  bot, err = tgbotapi.NewBotAPI("YOUR_TELEGRAM_BOT_API_TOKEN") // Replace with your Telegram bot token
  ```

- Ensure your Telegram bot is an admin in the channel.

### 4. Set Up Shodan Webhook

Configure Shodan to send webhooks to your server:

```
http://yourserver.com:9080/updat1X73rj92
```

Replace `yourserver.com` with your server's IP or domain.

### 5. Run the Application

```bash
go run main.go
```

The server listens on port `9080`.

### 6. Configure Server Access

- Make sure port `9080` is accessible externally.
- Configure firewall or use port forwarding if necessary.

## How It Works

1. **Webhook Reception**: Receives data from Shodan when your trigger conditions are met.
2. **Data Parsing**: Extracts relevant information from the JSON payload.
3. **Database Check**: Queries the SQLite database to check if the IP and port have been processed.
4. **Data Insertion**: Stores new data into the database.
5. **Screenshot Capture**: Navigates to the service URL and captures screenshots using `chromedp`.
6. **Telegram Notification**: Sends the screenshot and service details to the specified Telegram channel.

## License

This project is licensed under the GNU General Public License v3.0.

## Contributing

Feel free to submit issues or pull requests.

## Notes

- Ensure that `chromedp` can run headless Chrome on your server environment.
- The server endpoint `/updat1X73rj92` is hardcoded; you may change it for security purposes.
- Modify `time.Sleep(3 * time.Second)` if you need to adjust the delay between sending messages.

## Support Me

If you like this extension, consider donating via cryptocurrency:

- **USDT (TRON)**: `TUMiYtejQjuCXZA9iA5H7zMa65vdHBtrRC`
- **BTC**: `bc1qhv4x0lt5lwk46gfr5lmpauucck7eaxegcasr6d`
