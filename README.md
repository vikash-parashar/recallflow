# RecallFlow

> Never lose a patient because of a missed call.

RecallFlow is a healthcare-focused SaaS platform that helps US clinics automatically recover missed patient calls using SMS automation and AI-powered conversation workflows.

## 🎯 Product Overview

RecallFlow acts as an operational automation layer that:
- Detects missed patient calls in real-time
- Sends automatic SMS responses
- Uses AI to classify patient intent
- Notifies clinic staff with actionable information
- Tracks recovery metrics and ROI

**Not an EHR replacement** — integrates alongside existing systems.

## 🏥 Target Market

- Dental clinics
- Med spas
- Chiropractors
- Physiotherapy clinics
- Cosmetic clinics

## 🛠 Tech Stack

| Category | Technology |
|----------|-----------|
| Backend | Golang |
| Frontend | Next.js + TypeScript |
| Database | PostgreSQL |
| Queue | Redis |
| Hosting | AWS (EC2, RDS) |
| SMS/Voice | Twilio |
| AI | OpenAI |
| Payments | Stripe |

## 📁 Project Structure

```
/backend          - Golang REST API
/frontend         - Next.js admin dashboard
/docs            - Documentation
/scripts         - Deployment and utility scripts
```

## 🚀 MVP Features

1. ✅ Missed call detection via Twilio webhooks
2. ✅ Automatic SMS replies
3. ✅ AI intent classification (OpenAI)
4. ✅ Staff notifications
5. ✅ Admin dashboard with analytics
6. ✅ Multi-location support
7. ✅ Stripe billing integration

## 💰 Pricing

- **Solo Clinic**: $99/month
- **Multi-Provider**: $299/month
- **Multi-Location**: $999/month

## 📦 Quick Start

### Backend Setup
```bash
cd backend
go mod download
cp .env.example .env
go run cmd/api/main.go
```

### Frontend Setup
```bash
cd frontend
npm install
cp .env.local.example .env.local
npm run dev
```

## 🏗 Development Philosophy

- **Ship fast** — MVP in 3 weeks
- **Solve one problem extremely well** — missed call recovery
- **Iterate from real feedback** — launch early, improve continuously
- **Production-ready from day one** — no technical debt

## 📈 Business Model

SaaS with recurring revenue. Target: Acquire 50 clinics in first 6 months for ~$15,000 MRR.

## 📄 License

Proprietary - All Rights Reserved

---

Built with ❤️ for healthcare clinics
