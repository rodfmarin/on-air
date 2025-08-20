<img src="./resources/icon.svg" alt="on-air">

## A LIFX Smart Bulb Utility for reflecting your Free/Busy Google Calendar Status


This is a Go-based utility that automatically controls a LIFX smart bulb based on your Google Calendar's free/busy status. When your calendar shows you as busy, the utility will set your LIFX bulb to a busy state (e.g., red). When you are free, it will set the bulb to a free state (e.g., white).

## Features
- Monitors your Google Calendar for free/busy status
- Controls a LIFX bulb using the LIFX REST API
- Configuration via a simple `config.json` file
- Command-line flags for ad-hoc overrides

## Requirements

1. **Google API OAuth Client**
   - You must create a Google API OAuth client with the Google Calendar API enabled.
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Navigate to **APIs & Services → Library → Google Calendar API → Enable**
   - Download your OAuth client credentials as `credentials.json` and place it in the project directory.

2. **First Run Authorization**
   - On the first run, the app will open a loopback browser window to authorize your OAuth client with your Google account.
   - This will generate a `token.json` file for future runs.

3. **LIFX Bulb and Developer Token**
   - You need a LIFX smart bulb.
   - Obtain a developer token by following instructions at: [LIFX API Authentication](https://api.developer.lifx.com/reference/authentication)
   - Place your token in the `config.json` file under the `lifx_token` field.

4. **Configuration**
   - Update `config.json` with all required settings:
     - `credentials`: Path to your Google OAuth client JSON file
     - `token`: Path to your Google OAuth token file
     - `calendar`: Your Google Calendar ID (or `primary` for your main calendar)
     - `days`: How many days ahead to check for events
     - `lifx_token`: Your LIFX API token
     - `lifx_light_id`: The ID of the LIFX bulb to control
     - `lifx_light_label`: The label of the LIFX bulb to control
     - `lifx_busy_color`: The color to set when busy (e.g., "red saturation:0.8")
     - `lifx_free_color`: The color to set when free (e.g., "kelvin:3500")
     - `reload_interval_seconds`: How often to reload the calendar schedule (in seconds)

   Example `config.json`:
   ```json
   {
     "credentials": "path_to_oauth_credentials.json",
     "token": "path_to_oauth_token.json",
     "calendar": "your_calendar@gmail.com",
     "days": 7,
     "lifx_token": "YOUR_LIFX_API_TOKEN_HERE",
     "lifx_light_id": "YOUR_LIFX_LIGHT_ID_HERE",
     "lifx_light_label": "YOUR_LIFX_LIGHT_LABEL_HERE",
     "lifx_busy_color": "red saturation:0.8",
     "lifx_free_color": "kelvin:3500",
     "reload_interval_seconds": 120
   }
   ```

## Usage

```sh
go run main.go
```

You can override config values with command-line flags for ad-hoc things like this:

```sh
go run main.go -calendar="your_calendar_id" -lifx_token="your_token_here" -lifx_busy_color="blue saturation:1.0" -reload_interval_seconds=300
```

## Notes
- Make sure your LIFX bulb is online and connected to your account.
- The utility will continuously monitor your calendar and update the bulb state in real time.
- Keep your API tokens secure and do not share them publicly.
