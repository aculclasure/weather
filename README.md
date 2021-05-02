# Go OpenWeather API Client #
Go package that provides a client for interacting with the [current weather API](https://openweathermap.org/current) hosted by [OpenWeather]. Also provides a CLI that displays basic measurements of current weather for a given location.

## Setup ##
You will need to create an [OpenWeather] account and API key. See their [getting started](https://openweathermap.org/appid) guide for help on how to do this. Once you have the API key, it should be set as the environment variable `OPENWEATHER_API_KEY`.

## CLI Usage ##

Mac/Linux
```
$ cd cmd/weather

$  go run main.go -h
USAGE: weather [-units={standard|metric|imperial}] <location>

  -units string
        the units to use, one of: standard, metric, imperial (default "imperial")

$ OPENWEATHER_API_KEY=<YOUR-API-KEY> go run main.go --units=metric london

overcast clouds, 9.21 C, humidity 46%
```

Windows Powershell
```
PS ${env:OPENWEATHER_API_KEY}=<YOUR-API-KEY>

PS cd cmd\weather

PS go run main.go --units=metric london

overcast clouds, 9.21 C, humidity 46%
```

[OpenWeather]: https://openweathermap.org/