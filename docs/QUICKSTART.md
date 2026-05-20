# RecallFlow - Quick Start Guide

## ⚡️ 5-Minute Setup

### Prerequisites
- Go 1.22+
- Node.js 18+
- Docker & Docker Compose
- Git

### Step 1: Clone and Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/recallflow.git
cd recallflow

# Run setup script
chmod +x setup.sh
./setup.sh
```

### Step 2: Configure Environment

Edit `backend/.env`:
```bash
cd backend
nano .env
```

Minimum required settings:
```env
PORT=8080
DATABASE_URL=postgresql://recallflow:recallflow_dev_password@localhost:5432/recallflow?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-random-secret-key-change-this
TWILIO_ACCOUNT_SID=your_twilio_sid
TWILIO_AUTH_TOKEN=your_twilio_token
TWILIO_PHONE_NUMBER=+1234567890
OPENAI_API_KEY=sk-your-openai-key
```

### Step 3: Run Database Migrations

```bash
cd backend
make migrate
```

### Step 4: Start Backend

```bash
cd backend
make dev
```

✅ API running at `http://localhost:8080`

### Step 5: Start Frontend (New Terminal)

```bash
cd frontend
npm run dev
```

✅ Dashboard at `http://localhost:3000`

## 🧪 Testing the System

### 1. Create an Account
Visit `http://localhost:3000/register` and create a test clinic account.

### 2. Test Twilio Webhooks (Local)

Install ngrok:
```bash
# Mac
brew install ngrok

# Or download from https://ngrok.com
```

Expose your local API:
```bash
ngrok http 8080
```

Configure Twilio webhooks with your ngrok URL:
- Voice: `https://your-ngrok-url.ngrok.io/api/v1/webhooks/twilio/voice`
- SMS: `https://your-ngrok-url.ngrok.io/api/v1/webhooks/twilio/sms`

### 3. Simulate a Missed Call

Use Twilio Console to make a test call to your Twilio number. When the call is missed, the system will:
1. Create a call record
2. Send an SMS to the caller
3. Show up in your dashboard

## 📱 Using the Dashboard

### Dashboard Overview
- Total calls in last 30 days
- Missed calls count
- Recovered leads
- Response rate
- Estimated revenue

### Managing Conversations
1. Navigate to **Conversations**
2. Click on any conversation to view details
3. See SMS history
4. Mark as resolved when done

## 🔧 Common Commands

### Backend
```bash
# Start development server
make dev

# Build production binary
make build

# Run tests
make test

# Start database & Redis
make docker-up

# Stop containers
make docker-down

# View logs
make docker-logs

# Run migrations
make migrate
```

### Frontend
```bash
# Development server
npm run dev

# Production build
npm run build

# Start production server
npm start

# Type check
npm run type-check

# Lint code
npm run lint
```

## 🐛 Troubleshooting

### Backend won't start
```bash
# Check if PostgreSQL is running
docker ps

# Check logs
docker-compose logs postgres

# Restart database
make docker-down && make docker-up
```

### Database connection error
```bash
# Verify DATABASE_URL in .env
# Check PostgreSQL is accessible
psql postgresql://recallflow:recallflow_dev_password@localhost:5432/recallflow
```

### Frontend API errors
```bash
# Verify backend is running
curl http://localhost:8080/health

# Check NEXT_PUBLIC_API_URL in frontend/.env.local
```

### Twilio webhooks not working
- Ensure ngrok is running
- Verify webhook URLs in Twilio console
- Check backend logs for webhook requests
- Verify Twilio credentials in .env

## 📊 Database GUI (Optional)

Install a PostgreSQL client:

**pgAdmin** (GUI):
- Download from https://www.pgadmin.org
- Connect to `localhost:5432`
- Username: `recallflow`
- Password: `recallflow_dev_password`

**psql** (CLI):
```bash
psql postgresql://recallflow:recallflow_dev_password@localhost:5432/recallflow
```

## 🔑 Getting API Keys

### Twilio
1. Sign up at https://www.twilio.com
2. Get free trial credits
3. Copy Account SID and Auth Token
4. Buy a phone number with Voice + SMS

### OpenAI
1. Sign up at https://platform.openai.com
2. Create API key
3. Add $5-10 credits

### Stripe (for billing)
1. Sign up at https://stripe.com
2. Get test API keys
3. Configure webhooks

## 📚 Next Steps

1. **Customize SMS messages** - Edit in `backend/internal/services/twilio_service.go`
2. **Add more AI intents** - Modify `backend/internal/services/openai_service.go`
3. **Customize dashboard** - Edit `frontend/src/pages/dashboard/index.tsx`
4. **Add Stripe billing** - Follow Stripe documentation
5. **Deploy to production** - See `docs/DEPLOYMENT.md`

## 💡 Pro Tips

- Use `.env.local` for local frontend overrides
- Check `docs/README.md` for full documentation
- Run `make help` to see all backend commands
- Enable auto-save in your editor
- Use VS Code with Go and TypeScript extensions

## 🆘 Getting Help

- Check API logs: Backend will print all requests
- Use browser DevTools to debug frontend
- Check database with SQL queries
- Review Twilio webhook logs in console

## ✅ Verification Checklist

- [ ] Backend running on port 8080
- [ ] Frontend running on port 3000
- [ ] PostgreSQL accessible
- [ ] Redis accessible
- [ ] Can register new account
- [ ] Can login successfully
- [ ] Dashboard loads
- [ ] Twilio webhooks configured (for testing calls)

---

**Ready to build?** You now have RecallFlow running locally! 🚀
