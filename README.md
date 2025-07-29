# env-cracker-web 🔓

This app lets you upload a proprietary `.env` file and download a ZIP archive containing the extracted embedded files.

![Image of the app](./app.png)

## 📝 Project Purpose

Many internal tools or legacy systems store multiple files inside custom binary formats. `env-cracker-web` is designed to extract those embedded files from `.env` files in a user-friendly way, without needing to reverse engineer manually. This can be useful for internal forensics, debugging, or documentation purposes.

## 📦 Input & Output Format

- **Input**: A single binary file with `.env` extension, containing embedded files in a custom format.
- **Output**: A `.zip` archive containing all the extracted files with correct filenames and extensions.

## 🔗 Live Demo

[![Visit Live Demo](https://img.shields.io/badge/Visit-Demo-blue?style=for-the-badge)](https://env-cracker-web-production.up.railway.app/)

## 🐳 Run project with Docker

1. Ensure Docker is installed by running `docker --version`, which should return a message like `Docker version ...`.
2. Run: `docker pull nicolasbellanich/env-cracker-web:latest`
3. Run: `docker run -p 8080:8080 nicolasbellanich/env-cracker-web:latest`

## 💻 Run project locally

1. Clone this repo.
2. Navigate to the root directory and run: `go run ./...`