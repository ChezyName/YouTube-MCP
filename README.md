<p align="center">
  <b><font size="6">YouTube MCP</font></b>
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/ChezyName/YouTube-MCP"/>
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg"/>
    <img src="https://img.shields.io/badge/platform-Windows%20|%20Mac%20|%20Linux-orange">
</p>


<p align="center">
  <a href="#about"> About </a> |
  <a href="#getting-started"> Getting Started </a> |
  <a href="#tools"> Tools </a> |
  <a href="#roadmap"> Roadmap </a>
</p>

---


## About
YouTube MCP is a small program designed to be used by desktop AI clients such as [Gemini CLI](https://geminicli.com/) or [Claude Desktop](https://claude.com/download) to allow AI agents to **READ ONLY** the stats of your YouTube Channel including analytics - such as AVD, AVP, Views, Subscribers, and more.

> This program does not write to your channel, <br/>
> does not touch, read, write monetization data. <br/>
> All data is controlled by you, with everything working through your machine.

## Getting Started
### OAuth and Credentials
> These are required for both the Installer or Manual Installation.

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
6. Add Yourself as [Test User](https://console.cloud.google.com/auth/audience)
7. [Publish your app](https://console.cloud.google.com/auth/audience) as to not have your Token-Refresh not reset after 7 days
   1. Set it as Production
> **Note:** You should save the ClientID and the ClientSecret someplace secret, You can always create a new secret

1. Create [YouTube API Key](https://console.cloud.google.com/apis/credentials)
   1. **Name:** Whatever you want
   2. **API Restrictions:** YouTube Data API v3, YouTube Analytics API
   3. Create and save API Key
2. Run the installer script, it will ask you to write the Client ID, Client Secret, and YouTube API Key.
3. 
### Installer
> This is the recommended option; You can still download without it but requires more work with the token-grabber and the YouTube MCP executables.

Running the Installer requires Admin privileges due to how windows 'sees' it, It would say no published because its built in house, If you do not trust it, just manually install it.


### Manual Installation
> You need GoLang for this.

For this, its recommended to download the GitHub repo and manually run the token-grabber directory,

`go run .\cwd\token-getter\`

Which will ask you for ClientID and ClientSecret,
then it would auto-open the browser to setup the RefreshToken which is needed for the Private API,

Manually setting up the Config is simple JSON.
Its located in 

* **On Unix systems:** [$XDG_CONFIG_HOME](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) if non-empty, else $HOME/.config.
* **On Mac OS:** $HOME/Library/Application Support.
* **On Windows:** %AppData%.

Create a new folder for the YouTubeMCP in `HomeDirectory/YouTube-MCP/`,
Then create a file `config.json` which would have this schema:
``` json
{
  "YOUTUBE_API": "YOUTUBE_API_KEY",
  "YOUTUBE_REFRESH_TOKEN": "REFRESH FROM TOKEN-GRABBER OR INSTALLER",
  "ChannelHandle": "YOUTUBE HANDLE WITHOUT '@'",
  "YOUTUBE_CLIENT_ID": "CLIENT ID FROM GOOGLE DEVELOPER CONSOLE",
  "YOUTUBE_CLIENT_SECRET": "CLIENT SECRET FROM GOOGLE DEVELOPER CONSOLE"
}
```

Finally, install the version of the MCP you wish to have, an put the executable anywhere
as long as the AI agent can access it, recommended same folder as *config.json*

### Connecting to an AI
Each AI app is different but the main way to connect the MCP is just a path to the binary, for example

> Mac & Linux users: Run the installer script — it will configure the correct path for your system automatically.

**Claude Desktop @ Windows, <br/>
Gemini CLI @ Windows**
``` json
  "mcpServers": {
    "YouTube": {
      "command": "C:\\AppData\\Roaming\\YouTube-MCP\\YouTube-MCP.exe"
    }
  },
```


## Tools
These are the available tools that the MCP exposes for AI agents to use.
The assumption is that the AI would automatically use these tools based on your request.

### Examples
Use it as if AI has the context of your content with questions such as the following:

  * "Give me a detailed overview of my latest video's performance, especially the retention by the transcript."
  * "Compare my latest short to my top 3 best shorts"
  * "On my latest video, where are viewers dropping the video the most?"
  * "Can you compare my latest longform video to a popular video of the same niche?"

There will always be some prompts the AI cannot fulfil, but for what this tool is designed to do - read the data of your Channel, it should help creators analyze their videos and see where they need to improve.

### API
> 🌐 Public: Requires the YouTube API Key Only <br/>
> 🔒 Private: Request the YouTube OAuth and Refresh Token

> **Range Format:** "30" (last 30 days), "365" (last year), 
> "lifetime" (all time), or "YYYY-MM-DD/YYYY-MM-DD" for custom ranges. <br/>
> Applies to get_top_videos, get_channel_analytics, and get_video_analytics.


| Tool Name | Type | Description | Returns |
| :--- | :---: | :--- | :--- |
| `auth_check`  | 🔒 Private </br> 🌐 Public | Allows MCP to check if this Data API and Analytics API is valid| Auth Check Struct |
| `all_videos` | 🌐 Public | Returns a list of all public videos on the channel, including basic data. Can specify shortform, longform, or both. | List of Videos: id, title, description, thumbnail, published_at, content_type |
| `get_video` | 🌐 Public | Gets detailed information for a single video. | ID, Title, Description, Thumbnail, PublishedAt, Duration, Views, Likes, Dislikes, CommentCount |
| `get_video_comments` | 🌐 Public | Gets a limited number of comments from a specific video. | Author, Text, LikeCount, PublishedAt, UpdatedAt, Comment ID |
| `get_video_transcript` | 🌐 Public | Returns a structured list of the video transcript. | Array of transcript with text and timestamps |
| `get_top_videos` | 🌐 Public | Returns a list of top videos given a date range, limit and content_type (defaults to last 30 days, top 10, both long and short form content). | Basic: Name, Description, Likes · Detailed: adds Views, Duration, PublishedAt, CommentCount |
| `search_video` | 🌐 Public | Searches for videos given a query (supports basic or detailed view). | Array of Basic or Detailed Video Data |
| `get_channel` | 🌐 Public | Gets public information for a user's (self) channel. | Name, Handle, Banner, Icon, Description |
| `get_channel_analytics`| 🔒 Private | Gets channel analytics for a date range (defaults to 30 days). | Views, Watch time, AVD, CTR, Subscribers, Demographics, Traffic sources, Daily breakdown |
| `get_video_analytics`  | 🔒 Private | Gets video-specific analytics for a date range (defaults to last 30 days). | Views, Watch time, AVD, CTR, Demographics, Traffic sources, Daily breakdown, Retention Graph |

## Roadmap

- Search for any YouTube channel by handle or ID
- Get public stats for any channel (not just your own)
- Drop API Key requirement — OAuth only
- Caching for Performance Improvements < Not Fully Needed>

- Competitor tracking — add/remove competitors and compare basic public stats (must be enabled by config)
- Roadmaps -> Allow AI to edit `Context/<year>.md` with tooling to read and write (must be enabled by config)
  - Allows AI to get a context of your plans for each year

- TikTok integration (Channel stats, video analytics, top videos)

---
<p align="center">
© 2026 ChezyName
</p>
