-- RecallFlow Database Schema
-- PostgreSQL

-- Organizations table
CREATE TABLE organizations (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    plan_type VARCHAR(50) DEFAULT 'solo', -- solo, multi_provider, multi_location
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_organizations_email ON organizations(email);
CREATE INDEX idx_organizations_is_active ON organizations(is_active);

-- Locations table
CREATE TABLE locations (
    id VARCHAR(36) PRIMARY KEY,
    organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id),
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    address TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_locations_org_id ON locations(organization_id);
CREATE INDEX idx_locations_phone ON locations(phone_number);
CREATE INDEX idx_locations_is_active ON locations(is_active);

-- Users table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) DEFAULT 'staff', -- owner, admin, staff
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_org_id ON users(organization_id);
CREATE INDEX idx_users_is_active ON users(is_active);

-- Calls table
CREATE TABLE calls (
    id VARCHAR(36) PRIMARY KEY,
    location_id VARCHAR(36) NOT NULL REFERENCES locations(id),
    organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id),
    twilio_call_sid VARCHAR(255) UNIQUE,
    from_number VARCHAR(20) NOT NULL,
    to_number VARCHAR(20) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, missed, answered, voicemail
    duration INTEGER DEFAULT 0,
    call_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_calls_org_id ON calls(organization_id);
CREATE INDEX idx_calls_location_id ON calls(location_id);
CREATE INDEX idx_calls_twilio_sid ON calls(twilio_call_sid);
CREATE INDEX idx_calls_status ON calls(status);
CREATE INDEX idx_calls_call_time ON calls(call_time);
CREATE INDEX idx_calls_from_number ON calls(from_number);

-- Conversations table
CREATE TABLE conversations (
    id VARCHAR(36) PRIMARY KEY,
    call_id VARCHAR(36) NOT NULL REFERENCES calls(id),
    organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id),
    location_id VARCHAR(36) NOT NULL REFERENCES locations(id),
    patient_phone VARCHAR(20) NOT NULL,
    status VARCHAR(50) DEFAULT 'active', -- active, resolved, escalated
    intent VARCHAR(50), -- appointment, billing, emergency, insurance, prescription, hours, general
    summary TEXT,
    is_resolved BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_conversations_org_id ON conversations(organization_id);
CREATE INDEX idx_conversations_call_id ON conversations(call_id);
CREATE INDEX idx_conversations_patient_phone ON conversations(patient_phone);
CREATE INDEX idx_conversations_status ON conversations(status);
CREATE INDEX idx_conversations_created_at ON conversations(created_at);

-- SMS Messages table
CREATE TABLE sms_messages (
    id VARCHAR(36) PRIMARY KEY,
    conversation_id VARCHAR(36) NOT NULL REFERENCES conversations(id),
    twilio_message_sid VARCHAR(255),
    direction VARCHAR(20) NOT NULL, -- inbound, outbound
    from_number VARCHAR(20) NOT NULL,
    to_number VARCHAR(20) NOT NULL,
    body TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'sent',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sms_conversation_id ON sms_messages(conversation_id);
CREATE INDEX idx_sms_direction ON sms_messages(direction);
CREATE INDEX idx_sms_created_at ON sms_messages(created_at);

-- AI Classifications table
CREATE TABLE ai_classifications (
    id VARCHAR(36) PRIMARY KEY,
    conversation_id VARCHAR(36) NOT NULL REFERENCES conversations(id),
    message_id VARCHAR(36) REFERENCES sms_messages(id),
    intent VARCHAR(50) NOT NULL,
    confidence DECIMAL(3,2) DEFAULT 0.0,
    is_emergency BOOLEAN DEFAULT false,
    requires_staff BOOLEAN DEFAULT false,
    suggested_action TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ai_conversation_id ON ai_classifications(conversation_id);
CREATE INDEX idx_ai_is_emergency ON ai_classifications(is_emergency);

-- Notifications table
CREATE TABLE notifications (
    id VARCHAR(36) PRIMARY KEY,
    organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id),
    user_id VARCHAR(36) REFERENCES users(id),
    conversation_id VARCHAR(36) REFERENCES conversations(id),
    type VARCHAR(50) NOT NULL, -- missed_call, urgent, new_message
    title VARCHAR(255) NOT NULL,
    message TEXT,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_org_id ON notifications(organization_id);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Subscriptions table
CREATE TABLE subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id) UNIQUE,
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    plan_type VARCHAR(50) NOT NULL, -- solo, multi_provider, multi_location
    status VARCHAR(50) DEFAULT 'active', -- active, canceled, past_due
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_org_id ON subscriptions(organization_id);
CREATE INDEX idx_subscriptions_stripe_customer ON subscriptions(stripe_customer_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);

-- Analytics table (aggregated daily metrics)
CREATE TABLE analytics (
    id SERIAL PRIMARY KEY,
    organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id),
    date DATE NOT NULL,
    total_calls INTEGER DEFAULT 0,
    missed_calls INTEGER DEFAULT 0,
    recovered_leads INTEGER DEFAULT 0,
    sms_sent INTEGER DEFAULT 0,
    sms_responses INTEGER DEFAULT 0,
    appointments_booked INTEGER DEFAULT 0,
    estimated_revenue DECIMAL(10,2) DEFAULT 0.0,
    response_rate DECIMAL(5,2) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, date)
);

CREATE INDEX idx_analytics_org_id ON analytics(organization_id);
CREATE INDEX idx_analytics_date ON analytics(date);
