<h1 align="center">YouTube MCP</h1>
<p align="center">
  <img src="https://img.shields.io/github/v/release/ChezyName/YouTube-MCP"/>
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg"/>
</p>

- [About](#about)
- [Getting Started](#getting-started)
  - [Installer](#installer)
  - [OAuth and Credentials](#oauth-and-credentials)
  - [Connecting to an AI](#connecting-to-an-ai)
- [Tools / Commands](#tools--commands)
- [Roadmap](#roadmap)


## About
YouTube MCP is a small program designed to be used by desktop AI clients such as [Gemini CLI](https://geminicli.com/) or [Claude Desktop](https://claude.com/download) to allow AI agents to **READ ONLY** the stats of your YouTube Channel including analytics - such as AVD, AVP, Views, Subscribers, and more.

> This program does not write to your channel, <br/>
> does not touch, read, write monetization data. <br/>
> All data is controlled by you, with everything working through your machine.

## Getting Started
### Installer
> This is the recommended option; You can still download without it but requires more work with the token-grabber and the YouTube MCP executables.

You can use the installer which would make it significantly easier to install this program by auto creating most of the items, and just asking you to type the Google OAuth ID and Secrets, as well as the Google API Key

### OAuth and Credentials
> These are required for the installer anyways, so better to have it setup now.
1. Create a new application on [Google Developer Console](https://console.cloud.google.com/)
2. Enable the [YouTube Data API](https://console.cloud.google.com/marketplace/product/google/youtube.googleapis.com)
3. Enable the [YouTube Analytics API](https://console.cloud.google.com/apis/library/youtubeanalytics.googleapis.com)
4. Create an [OAuth Screen](https://console.cloud.google.com/auth/overview)
   1. **App Information:** N/A
   2. **Audience:** External
   3. **Contact Information:** N/A
   4. **Finish:** Accept the agreement
5. Create an [OAuth Client](https://console.cloud.google.com/auth/clients)
   1. **Application Type:** Web Application
   2. **Name:** Whatever you want, ie: "YouTube MCP"
   3. **Authorized Redirect URIs:** `http://localhost:9999/callback`
> **Note:** You should save the ClientID and the ClientSecret someplace, You can always create a new secret

6. Create [YouTube API Key](https://console.cloud.google.com/apis/credentials)
   1. **Name:** Whatever you want
   2. **API Restrictions:** YouTube Data API v3, YouTube Analytics API
   3. Create and save API Key
7. Run the installer script, it will ask you to write the Client ID, Client Secret, and YouTube API Key.

> Note: If you are not using the Installer script, this program needs a YouTube OAuth Refresh Token in order to function properly

### Connecting to an AI
Each AI app is different but the main way to connect the MCP is just a path to the binary, for example

> Mac & Linux users: Run the installer script — it will configure the correct path for your system automatically.

Claude Desktop @ Windows, <br/>
Gemini CLI @ Windows,
``` json
  "mcpServers": {
    "YouTube": {
      "command": "C:\\AppData\\Roaming\\YouTube-MCP\\YouTube-MCP.exe"
    }
  },
```

## Tools / Commands
These are the available tools that the MCP exposes for AI agents to use.
> 🌐 Public: Requires the YouTube API Key Only <br/>
> 🔒 Private: Request the YouTube OAuth and Refresh Token

> **Range Format:** "30" (last 30 days), "365" (last year), 
> "lifetime" (all time), or "YYYY-MM-DD/YYYY-MM-DD" for custom ranges. <br/>
> Applies to get_top_videos, get_channel_analytics, and get_video_analytics.

| Tool Name | Type | Description | Returns |
| :--- | :---: | :--- | :--- |
| `all_videos` | 🌐 Public | Gets all public videos for a user. | ID, Title, Description, Thumbnail, PublishedAt |
| `get_video` | 🌐 Public | Gets detailed information for a single video. | ID, Title, Description, Thumbnail, PublishedAt, Duration, Views, Likes, Dislikes, CommentCount |
| `get_video_comments` | 🌐 Public | Gets a limited number of comments from a specific video. | Author, Text, LikeCount, PublishedAt, UpdatedAt, Comment ID |
| `get_video_transcript` | 🌐 Public | Returns a structured list of the video transcript. | Structured transcript text and timestamps |
| `get_top_videos` | 🌐 Public | Returns a list of top videos given a date range (defaults to last 90 days, top 10). | Basic: Name, Description, Likes · Detailed: adds Views, Duration, PublishedAt, CommentCount |
| `search_video` | 🌐 Public | Searches for videos given a query (supports basic or detailed view). | Basic or detailed video metadata |
| `get_channel` | 🌐 Public | Gets public information for a user's channel. | Name, Handle, Banner, Icon, Description |
| `get_channel_analytics`| 🔒 Private | Gets channel analytics for a date range (defaults to lifetime). | Views, Watch time, AVD, CTR, Subscribers, Demographics, Traffic sources, Daily breakdown |
| `get_video_analytics`  | 🔒 Private | Gets video-specific analytics for a date range (defaults to last 90 days). | Views, Watch time, AVD, CTR, Demographics, Traffic sources, Daily breakdown |

## Roadmap

**v1.4.0**
- Caching for Performance Improvements
- Search for any YouTube channel by handle or ID
- Get public stats for any channel (not just your own)
- Drop API Key requirement — OAuth only

**v1.5.0**
- Competitor tracking — add/remove competitors and compare basic public stats

**Future**
- TikTok integration (Channel stats, video analytics, top videos)

---
<p align="center">
© 2026 ChezyName
</p>
