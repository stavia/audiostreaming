# audiostreaming
audiostreaming is a Go library for getting spotify, youtube, deezer, etc URI from an artist + song query.

## Usage ##

Inside ```cmd``` folder there is an example command. Clone the repository and go to folder ```cmd/audiostreaming```
```console
cd cmd/audiostreaming
```

Copy .env-example to .env and set the data for your environment
```console
cp .env-example .env
```

Run the command
```console
go run main.com --artist "Arctic Monkeys" --track "Do I Wanna Know?"
```
