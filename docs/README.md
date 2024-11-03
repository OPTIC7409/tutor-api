# Tutor API Documentation

## Authentication

### Register a new user

POST /api/auth/register

Request body:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "userType": "student"
}
```

### Login

POST /api/auth/login

Request body:
```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Tutors

### Create a new tutor

POST /api/tutors

Request body:
```json
{
  "userID": 1,
  "subject": "Mathematics",
  "yearsExperience": 5,
  "hourlyRate": 50,
  "location": "New York"
}
```

### Get all tutors

GET /api/tutors

### Get a specific tutor

GET /api/tutors/:id

### Update a tutor

PUT /api/tutors/:id

Request body:
```json
{
  "subject": "Advanced Mathematics",
  "yearsExperience": 6,
  "hourlyRate": 55,
  "location": "Los Angeles"
}
```

### Delete a tutor

DELETE /api/tutors/:id

## Students

### Create a new student

POST /api/students

Request body:
```json
{
  "userID": 2,
  "age": 18,
  "subjects": "Mathematics, Physics",
  "location": "Chicago"
}
```

### Get all students

GET /api/students

### Get a specific student

GET /api/students/:id

### Update a student

PUT /api/students/:id

Request body:
```json
{
  "age": 19,
  "subjects": "Mathematics, Physics, Chemistry",
  "location": "Boston"
}
```

### Delete a student

DELETE /api/students/:id

## Chats

### Get all chats

GET /api/chats

### Get a specific chat

GET /api/chats/:id

### Create a new chat

POST /api/chats

### Send a message in a chat

POST /api/chats/:id/messages

Request body:
```json
{
  "senderID": 1,
  "content": "Hello, this is a test message."
}
```
