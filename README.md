# TUpload - Simple File Upload Service

A simple service that allows you to upload files and get temporary download links valid for 24 hours.

## Features

- Simple file upload via curl
- Short, easy-to-share download links
- Files expire after 24 hours
- HTTPS support via Traefik
- Docker support

## Quick Start

### Using Docker Compose

1. Configure your domain in `docker-compose.yml`:
```yaml
services:
  tupload:
    # ... other configuration ...
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.tupload.rule=Host(`your-domain.com`)"
      - "traefik.http.routers.tupload.entrypoints=websecure"
      - "traefik.http.routers.tupload.tls=true"
```

2. Build and start the service:
```bash
docker-compose up -d
```

3. Upload a file:
```bash
curl your-domain.com -F your_file.txt
```

4. Download a file:
```bash
wget https://your-domain.com/download/<file-id>
```

### Manual Build

1. Build the application:
```bash
go build
```

2. Run the application:
```bash
./tupload
```

## Example

Upload a file:
```bash
$ curl your-domain.com -F test.txt
=========================

Uploaded 1 file, 121 bytes

wget https://your-domain.com/download/Ab3Cd

=========================
```

Download the file:
```bash
$ wget https://your-domain.com/download/Ab3Cd
```

## Configuration

Edit `config/config.yaml` to customize:
- Server port and host
- Storage path
- Domain name

## License

MIT
