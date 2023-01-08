# GogoApps NASA

Application to collect images URL from NASA's APOD API.

## How to run

Start with cloning the repository

```bash
git clone https://github.com/Vrangz/S3J6eXN6dG9mIFN6dWxjIEdvZ29BcHBzIE5BU0E-.git
```

### Docker way

You can use already prepared scripts in `/script` directory. To run the application call

```bash
./scripts/run.sh
```

and if you want to stop the container then execute

```bash
./scripts/stop.sh
```

Alternatively, you can build and run using your own docker commands i.e.:

```bash
docker image build -t url-collector -f ./deployment/url-collector.Dockerfile .
```

and 

```bash
docker run -d -p 8080:8080 --name url-collector url-collector
```

If you want to reconfigure application you can use environment variables

```bash
docker run -e MY_VAR=MY_VALUE -d -p 8080:8080 --name url-collector url-collector
```

### Golang way

You can also run the application using source code calling

```bash
go run ./url-collector/cmd/main/main.go
```

and if you want to reconfigure the application you can also provide the environment variables

```bash
MY_VAR=MY_VALUE ./url-collector/cmd/main/main.go
```

or you can modify the `config.yaml` file.

## How to use

By default the server will start on http://localhost:8080/ and there you have one route available

```
- /api/v1/nasa/pictures
```

with two query params which allow to query for many urls from certain time

```
- from: defines start date in "2006-01-02" format
- to  : defines end date in "2006-01-02" format
```

Example `curl` calls:

The following example will return todays picture url of the day
```bash
curl "http://localhost:8080/api/v1/nasa/pictures"
```

The following example will return all the pictures url from days 2023-01-01 to 2023-01-05
```bash
curl "http://localhost:8080/api/v1/nasa/pictures?from=2023-01-01&to=2023-01-05"
```
