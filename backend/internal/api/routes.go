package api

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"github.com/recallflow/backend/internal/api/handlers"
	"github.com/recallflow/backend/internal/config"
	"github.com/recallflow/backend/internal/middleware"
	"github.com/recallflow/backend/internal/repositories"
	"github.com/recallflow/backend/internal/services"
)

func RegisterRoutes(router *mux.Router, db *sql.DB, redisClient *redis.Client, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	orgRepo := repositories.NewOrganizationRepository(db)
	locationRepo := repositories.NewLocationRepository(db)
	callRepo := repositories.NewCallRepository(db)
	conversationRepo := repositories.NewConversationRepository(db)
	smsRepo := repositories.NewSMSRepository(db)
	subRepo := repositories.NewSubscriptionRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	twilioService := services.NewTwilioService(cfg)
	openaiService := services.NewOpenAIService(cfg.OpenAIAPIKey)
	stripeService := services.NewStripeService(cfg)
	conversationService := services.NewConversationService(conversationRepo, smsRepo, openaiService, twilioService)
	callService := services.NewCallService(callRepo, conversationService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, userRepo, orgRepo)
	twilioHandler := handlers.NewTwilioHandler(callService, twilioService, locationRepo)
	conversationHandler := handlers.NewConversationHandler(conversationService)
	dashboardHandler := handlers.NewDashboardHandler(db)
	billingHandler := handlers.NewBillingHandler(stripeService, subRepo, orgRepo, cfg.StripeWebhookSecret)

	// Public routes (no auth required)
	router.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Twilio webhooks (Twilio validates these differently)

	// Stripe webhooks (no auth required, validated by signature)
	router.HandleFunc("/webhooks/stripe", billingHandler.HandleStripeWebhook).Methods("POST")
	router.HandleFunc("/webhooks/twilio/voice", twilioHandler.HandleVoiceWebhook).Methods("POST")
	router.HandleFunc("/webhooks/twilio/sms", twilioHandler.HandleSMSWebhook).Methods("POST")
	router.HandleFunc("/webhooks/twilio/status", twilioHandler.HandleStatusCallback).Methods("POST")

	// Protected routes (require authentication)
	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Organization routes
	protected.HandleFunc("/organizations/{id}", handlers.NotImplemented).Methods("GET")
	protected.HandleFunc("/organizations/{id}", handlers.NotImplemented).Methods("PUT")

	// Location routes
	protected.HandleFunc("/locations", handlers.NotImplemented).Methods("GET")
	protected.HandleFunc("/locations", handlers.NotImplemented).Methods("POST")
	protected.HandleFunc("/locations/{id}", handlers.NotImplemented).Methods("GET")
	protected.HandleFunc("/locations/{id}", handlers.NotImplemented).Methods("PUT")

	// Call routes
	protected.HandleFunc("/calls", handlers.NotImplemented).Methods("GET")
	protected.HandleFunc("/calls/{id}", handlers.NotImplemented).Methods("GET")

	// Conversation routes
	protected.HandleFunc("/conversations", conversationHandler.ListConversations).Methods("GET")
	protected.HandleFunc("/conversations/{id}", conversationHandler.GetConversation).Methods("GET")
	protected.HandleFunc("/conversations/{id}/messages", conversationHandler.GetMessages).Methods("GET")
	protected.HandleFunc("/conversations/{id}/resolve", conversationHandler.ResolveConversation).Methods("POST")

	// Dashboard/Analytics routes
	protected.HandleFunc("/dashboard/stats", dashboardHandler.GetStats).Methods("GET")
	protected.HandleFunc("/dashboard/analytics", dashboardHandler.GetAnalytics).Methods("GET")


	// Billing routes
	protected.HandleFunc("/billing/subscription", billingHandler.GetSubscription).Methods("GET")
	protected.HandleFunc("/billing/subscription", billingHandler.CreateSubscription).Methods("POST")
	protected.HandleFunc("/billing/subscription/cancel", billingHandler.CancelSubscription).Methods("POST")
	// User routes
	protected.HandleFunc("/users/me", authHandler.GetMe).Methods("GET")
	protected.HandleFunc("/users", handlers.NotImplemented).Methods("GET")
}
