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

---

## 🚀 CI/CD Pipeline & Contabo VM Deployment

A fully automated CI/CD pipeline has been set up using **GitHub Actions**.

Whenever code is pushed to the `main` or `master` branches, the pipeline:

1. **Runs Tests**: Executes all unit and handler tests using Go's test runner.
2. **Builds & Pushes Docker Image**: Packages the application into a Docker image using the multi-stage `Dockerfile` and pushes it to your **Docker Hub** repository.
3. **Deploys to Contabo VM**: Logs into your remote Contabo virtual machine via SSH, pulls the newly built image, stops/removes any existing container, and starts the container with the required environment variables.

### GitHub Repository Secrets Configuration

To run the pipeline successfully, configure the following secrets under **Settings > Secrets and variables > Actions** in your GitHub repository:

| Secret Name          | Description                                        | Example Value                            |
| -------------------- | -------------------------------------------------- | ---------------------------------------- |
| `DOCKERHUB_USERNAME` | Your Docker Hub user name                          | `my_docker_user`                         |
| `DOCKERHUB_TOKEN`    | A personal access token generated from Docker Hub  | `dckr_pat_...`                           |
| `CONTABO_HOST`       | Public IP or domain name of your Contabo VM        | `192.168.1.100` or `vm.example.com`      |
| `CONTABO_USER`       | SSH user used to log in to the Contabo VM          | `root` or `ubuntu`                       |
| `CONTABO_SSH_KEY`    | Private SSH key matching the public key on your VM | `-----BEGIN OPENSSH PRIVATE KEY-----...` |
| `CONTABO_SSH_PORT`   | Port for SSH connection                            | `22`                                     |
| `APP_PORT`           | The port the app runs on (defaults to `4001`)      | `4001`                                   |
| `DBUSER`             | Production MySQL Database User                     | `production_db_user`                     |
| `DBPASSWORD`         | Production MySQL Database Password                 | `production_secure_pass`                 |
| `DBHOST`             | Production MySQL Hostname/IP                       | `production-mysql-db.example.com`        |
| `DBPORT`             | Production MySQL Database Port                     | `3306`                                   |
| `DBNAME`             | Production MySQL Database Name                     | `calorie_tracker`                        |
| `SESSION`            | A 32-byte secret key for session encryption        | `s3cr3t_s3ss10n_k3y_g03s_h3r3_123`      
