Website showing weekly chart of Crypto/Stable ratio

Frontend with chart
- Frontend hosting solution - Github Pages?
- SPA in React + graph library + react-query
- Lambda to query S3 CSV in Go

Backend
- CSV in S3 to store weekly ratio
- Seed CSV in local
- Take seed CSV + API response to create initial S3 CSV
- Weekly cron job to update CSV

Infra
- Github pages
- CICD pipeline for frontend using Git Actions and Github Pages
- ?Terraform in CICD pipeline for backend

https://github.dev/nzoschke/gofaas/tree/master/handlers

`cd terraform && terraform apply` => Deploy AWS Lambda function