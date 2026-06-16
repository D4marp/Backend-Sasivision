# SasiVision Backend - Golang

REST API untuk Flutter app dengan MySQL.

## Quick Start

```bash
cp .env.example .env
docker-compose up --build
```

API: `http://localhost:8080`  
Storage: `http://localhost:8080/storage/`

## Demo Account

- Email: `demo@sasivision.com`
- Password: `Sasivision123`

## Seed Data

Migration `002_seed_data.sql` mengisi:

- 3 motif Sasirangan (Bintang Bahambur, Naga Balimbur, Kulat Karikit)
- 3 video edukasi
- Quiz kategori Post-Test + 4 soal
- Feature switches (AR active, Quizzes active, Vocabulary inactive)

## API Routes

```
POST   /api/auth/sign-in
POST   /api/auth/sign-up
POST   /api/auth/verify-token
POST   /api/auth/logout

GET    /api/content/markers
GET    /api/content/videos
GET    /api/features/switches/:feature

GET    /api/quiz/categories
GET    /api/quiz/questions/:category
POST   /api/quiz/attempts
GET    /api/quiz/history/:email
```

## Storage Files

```
storage/
  audio/          # motif narration MP3
  models/         # AR GLB models
  markers/        # motif thumbnail images
  descriptions/   # motif text descriptions
```
