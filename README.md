# Introduction
This is the BE for Youtube Video Sharing App. Features:
- User login using Google Oauth, you will be automatically registered at the first login
- Sharing YouTube videos
- Viewing a list of shared videos
- Real-time notifications for new video shares: When a user shares a new video, other logged-in users will receive a real-time notification about the newly shared video.

# Prerequisites
    Go 1.22.1
    Docker
    Docker compose
# Installation & Configuration
## Install dependencies
```
go mod download
```
## Configuration
You can change .env values if you want.
# Running the Application
## Not using Docker
```
cd infra
go run ../cmd/svc/main.go
```

To run unit tests
```
go test -v
```
## Using Docker
```
docker compose up -d
```
# Deployment
This is deployed in a Digital Ocean Droplet with these setup:
- Nginx is the reverse proxy
- Let's encrypt for SSL cert
- Docker (using Docker Compose as local)
# Demo
Go to my demo app https://funnymovies.thao.tech/.

It will redirect to `/login` if you have never logged in before. Your account will be automatically registered if you log in for the first time. You need a Google email because I use Google OAuth2.

After logging in, you will be redirected back to the homepage. Here, a list of videos shared by users (actually, only me) will appear. On the header, there is a button labeled `Share a movie`. When you click on it, parse one URL into the text box, then press OK, it will notify me what you are watching :joy:
# Troubleshooting
