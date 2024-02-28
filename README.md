# YouTube Playlist Analyzer

YouTube Playlist Analyzer is a tool that retrieves information about a YouTube playlist, such as the total number of videos, average video length, and total playlist duration.

## Features

- Retrieve playlist information from a YouTube playlist URL
- Calculate total number of videos in the playlist
- Calculate average length of videos in the playlist
- Calculate total duration of the playlist
- Display duration at different playback speeds

### Prerequisites

Before running the application, you need to have the following:

- Go installed on your machine
- A YouTube Data API key. You can obtain one by following the instructions [here](https://developers.google.com/youtube/registering_an_application)

### Installation

Clone the repository:
```
git clone https://github.com/MISHRA-TUSHAR/YT-playlist-length-calculator.git
 ```

Create a .env file in the root directory and add your YouTube Data API key:
```
echo "YOUTUBE_API_KEY=your-api-key" > .env
```
Usage
Run the application:
```
run main.go
```

Enter the URL of the YouTube playlist when prompted.

View the playlist information displayed in the console.
