# RecallFlow Project Documentation

## Project Overview

RecallFlow is a healthcare SaaS platform that automatically recovers missed patient calls using SMS automation and AI-powered conversation handling.

## Tech Stack

### Backend
- **Language**: Golang
- **Framework**: Gorilla Mux
- **Database**: PostgreSQL
- **Cache/Queue**: Redis
- **Integrations**: Twilio (SMS/Voice), OpenAI (AI), Stripe (Payments)

### Frontend
- **Framework**: Next.js 14
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **API Client**: Axios

## Project Structure

```
/backend
  /cmd/api          - Application entry point
  /internal
    /api            - HTTP handlers and routes
    /config         - Configuration management
    /middleware     - Auth, logging, etc.
    /models         - Data models
    /repositories   - Database layer
    /services       - Business logic
  /db               - Database schemas and migrations

/frontend
  /src
    /pages          - Next.js pages
    /components     - React components
    /lib            - Utilities and API client
    /styles         - CSS files
```

## Getting Started

### Prerequisites
- Go 1.22+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (recommended)

### Backend Setup

1. **Start infrastructure**:
```bash
cd backend
make docker-up
```

2. **Set up environment**:
```bash
cp .env.example .env
# Edit .env with your actual credentials
```

3. **Run migrations**:
```bash
make migrate
```

4. **Start API server**:
```bash
make dev
```

Server runs on `http://localhost:8080`

### Frontend Setup

1. **Install dependencies**:
```bash
cd frontend
npm install
```

2. **Set up environment**:
```bash
cp .env.local.example .env.local
```

3. **Start development server**:
```bash
npm run dev
```

Frontend runs on `http://localhost:3000`

## Core Workflow

1. **Missed Call Detection**
   - Twilio webhook triggers on unanswered call
   - System creates call record
   - Status callback confirms missed call

2. **SMS Automation**
   - Automatic SMS sent to patient
   - Message: "Sorry we missed your call to [Clinic]. How can we help you today?"

3. **AI Classification**
   - Patient responds via SMS
   - OpenAI classifies intent (appointment, billing, emergency, etc.)
   - System generates appropriate response

4. **Staff Notification**
   - Dashboard updates in real-time
   - Notifications sent to clinic staff
   - Actionable insights provided

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new organization
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/users/me` - Get current user

### Webhooks
- `POST /api/v1/webhooks/twilio/voice` - Incoming call webhook
- `POST /api/v1/webhooks/twilio/sms` - Incoming SMS webhook
- `POST /api/v1/webhooks/twilio/status` - Call status callback

### Conversations
- `GET /api/v1/conversations` - List conversations
- `GET /api/v1/conversations/:id` - Get conversation
- `GET /api/v1/conversations/:id/messages` - Get messages
- `POST /api/v1/conversations/:id/resolve` - Mark resolved

### Dashboard
- `GET /api/v1/dashboard/stats` - Get dashboard statistics
- `GET /api/v1/dashboard/analytics` - Get analytics data

## Database Schema

Key tables:
- `organizations` - Clinic organizations
- `locations` - Clinic locations with phone numbers
- `users` - System users
- `calls` - Inbound call records
- `conversations` - SMS conversation threads
- `sms_messages` - Individual SMS messages
- `ai_classifications` - AI analysis results
- `subscriptions` - Billing subscriptions

## Environment Variables

### Backend
```bash
PORT=8080
DATABASE_URL=postgresql://...
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
TWILIO_ACCOUNT_SID=...
TWILIO_AUTH_TOKEN=...
TWILIO_PHONE_NUMBER=...
OPENAI_API_KEY=...
STRIPE_SECRET_KEY=...
```

### Frontend
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Development Commands

### Backend
```bash
make dev          # Run development server
make build        # Build binary
make test         # Run tests
make docker-up    # Start PostgreSQL + Redis
make docker-down  # Stop containers
make migrate      # Run database migrations
```

### Frontend
```bash
npm run dev       # Start development server
npm run build     # Build for production
npm run start     # Start production server
npm run lint      # Run linter
```

## Deployment

### Backend Deployment (AWS EC2)
1. Build binary: `make build`
2. Upload to EC2 instance
3. Set environment variables
4. Run with systemd or supervisor

### Frontend Deployment (Vercel)
1. Connect GitHub repository
2. Set environment variables
3. Deploy automatically on push

## Testing Webhooks Locally

Use ngrok to expose local server:
```bash
ngrok http 8080
```

Update Twilio webhook URLs to ngrok URL.

## MVP Features Checklist

- [x] Backend API structure
- [x] Database schema
- [x] Authentication system
- [x] Twilio webhook handlers
- [x] SMS sending
- [x] OpenAI integration
- [x] Conversation management
- [x] Dashboard with analytics
- [x] Frontend (Next.js)
- [x] Landing page
- [x] Login/Register
- [ ] Stripe billing integration
- [ ] Production deployment
- [ ] Worker queue for async tasks
- [ ] Email notifications
- [ ] Onboarding flow

## Next Steps

1. **Week 1**: Complete Twilio integration testing
2. **Week 2**: Add worker queues for async SMS
3. **Week 3**: Stripe billing + deployment
4. **Launch**: Start customer acquisition

## Support

For questions or issues:
- Check documentation
- Review code comments
- Consult master skill document

---

Built with ❤️ for healthcare clinics
