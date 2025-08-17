# Introduction
This is a Project for a simple OTP (One Time Password) I did for a job interview.

# Installation
Not much needed to be done except 

```bash
docker compose up
```
and then 
```bash
go run ./internal/cmd
```
```bash
go run ./UI

# Database
I chose PostgreSQL as the database because:

For a single-node setup, it is more than sufficient.

It supports transactions, making it a very good candidate for an OTP service where data integrity matters.

On startup, a test user named john_doe is automatically created.
# Usage 