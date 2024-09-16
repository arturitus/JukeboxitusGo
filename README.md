# JukeboxitusGo

JukeboxitusGo is a Discord music bot that utilizes Lavalink for streaming, from Youtube, Soundcloud, and other platforms..

## Features

* Stream music from YouTube, SoundCloud, and more
  
* Add YouTube links for playback.
  
* Play YouTube public playlists.
  
* Skip songs.

## Usage

### Without Docker

You can use this bot without Docker by setting up the necessary configuration file. Create a file with the following content and ensure it includes the proper `Token`:

```yaml
Token: "Bot-Token"

Lavalink:
  Name: "test"
  Hostname: "lava-v4.ajieblogs.eu.org"
  Port: 80
  Password: "https://dsc.gg/ajidevserver"
  Secured: false
```

*Note: This file is only an example for showcase purposes. You must replace "Bot-Token" with your actual Discord bot token and adjust the Lavalink settings to match your server configuration.*

### With Docker
You can also use this bot with Docker. A Dockerfile is provided to help with the setup.

Here's an example of the container hosted on a Raspberry Pi, as shown in Portainer:

![image](https://github.com/user-attachments/assets/9e4057eb-b9a0-4681-91a1-0a6b521e22c5)


## License
This project is licensed under the MIT License. It is not for sale and is intended solely for academic purposes and skill demonstration.