# TUpload - Simple File Upload Service

A simple service that allows you to upload files and get temporary download links valid for 24 hours.

<div align="center">
  
![20250525_15h30m46s_grim](https://github.com/user-attachments/assets/2c7c917a-8446-4797-8cd9-5d81d365df7a)

</div>

## Features

- Simple file upload via curl
- Short, easy-to-share download links
- Files expire after 24 hours
- HTTPS support via Traefik
- Docker support

## Quick Start


## Configuration

Edit `config/config.yaml` to customize:
- Server port and host
- Storage path
- Domain name

```
server:
  port: 6060
  host: "0.0.0.0"
  trusted_proxies:
    - "127.0.0.1"
    - "::1"
    - "10.0.0.0/8"
    - "172.16.0.0/12"
    - "192.168.0.0/16"

storage:
  path: "./uploads"

domain: "your-domain.com"
```


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

### Example 

#### Upload a file:
```bash
curl https://your-domain.com -T your_file.txt
```
```
=========================

Uploaded 1 file, 1311 bytes

wget https://your-domain.com/download/2I9AV

=========================  
```

#### Download a file:
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

## License

MIT
