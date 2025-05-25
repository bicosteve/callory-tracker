# 🥗 Callory Tracker

**Callory Tracker** is a web application designed to help users monitor and manage their daily calorie intake. Users can register, log in, and log meals such as breakfast, lunch, dinner, or snacks, and track their nutritional consumption including calories, proteins, carbohydrates, and fats.

---

## 🚀 Features

- 🔐 **User Authentication**

  - Register
  - Login

- 🍽 **Meal Management**

  - Create a new food entry
  - Edit existing food entries
  - Delete food entries
  - Get a specific food entry by ID

- 📊 **Nutrition Analysis**
  - Calculates total daily nutritional consumption
  - Inputs: meal type (e.g., breakfast), food name, calories, protein, carbohydrates, fats
  - Output: nutritional summary with total calories and macros

---

## 🧰 Tech Stack

| Layer    | Technology  |
| -------- | ----------- |
| Backend  | Golang (Go) |
| Database | MySQL       |
| Frontend | HTML, CSS   |
| Hosting  | Heroku      |

---

## 📦 Project Structure

```bash
callory-tracker/
├── cmd/
│ └── web/
│ └── main.go
├── pkg/
│ ├── configs/
| ├── db/
| ├── forms/
| ├── helpers/
│ ├── models/
│ └── utils/
├── tables/
├── ui/
| ├── css/
| ├── html/
├── go.mod
├── go.sum
├── Procfile
└── README.md
```

---

## 🛠️ Installation and Setup

**Clone the repository**

```bash
git clone https://github.com/bicosteve/callory-tracker.git
cd callory-tracker

```

---

## 🛠️ Setting the app's db connection configs

1. DB_USER=your-db-username
2. DB_PASS=your-db-password
3. DB_HOST=your-db-host:3306
4. DB_NAME=your-db-name
5. SECRET=your-secret
6. PORT=4001

---

## 🛠️ Installing Dependancies

```bash
go mod tidy
```

---

## 🛠️ Run the application

```bash
go run ./cmd/web
```

## 🛠️ Deployment

```bash

cd /callory-tracker

# Build the binary
GOOS=linux GOARCH=amd64 go build -o callory-tracker ./cmd/web


heroku create clrytracker

git init
heroku git:remote -a myapp-name
git add .
git commit -m "Deploying callory-tracker"
git push heroku main
```
