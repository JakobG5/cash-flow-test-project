CREATE TYPE merchant_status AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE api_key_status AS ENUM ('active', 'inactive', 'expired');
CREATE TYPE currency_type AS ENUM ('ETB', 'USD');
CREATE TYPE payment_status AS ENUM ('pending', 'processing', 'success', 'failed', 'cancelled');
CREATE TYPE transaction_status AS ENUM ('pending', 'success', 'failed');
CREATE TYPE payment_method_type AS ENUM ('card', 'bank_transfer', 'mobile_money');
CREATE TYPE event_type AS ENUM ('created', 'processing', 'completed', 'failed', 'cancelled');

CREATE TABLE merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    status merchant_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE merchant_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    api_key VARCHAR(255) UNIQUE NOT NULL,
    secret_key VARCHAR(255) NOT NULL,
    status api_key_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT fk_merchant_api_keys_merchant FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);

CREATE TABLE payment_intents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id VARCHAR(255) UNIQUE NOT NULL,
    merchant_id UUID NOT NULL REFERENCES merchants(id),
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    currency currency_type NOT NULL,
    reference VARCHAR(255) UNIQUE NOT NULL,
    callback_url VARCHAR(500) NOT NULL,
    status payment_status DEFAULT 'pending',
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + INTERVAL '24 hours'),
    CONSTRAINT fk_payment_intents_merchant FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);

CREATE TABLE payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id VARCHAR(255) UNIQUE NOT NULL,
    payment_intent_id UUID NOT NULL REFERENCES payment_intents(id),
    merchant_id UUID NOT NULL REFERENCES merchants(id),
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    currency currency_type NOT NULL,
    status transaction_status DEFAULT 'pending',
    third_party_reference VARCHAR(255) UNIQUE,
    payment_method payment_method_type,
    fee_amount DECIMAL(10,2) DEFAULT 0 CHECK (fee_amount >= 0),
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_payment_transactions_intent FOREIGN KEY (payment_intent_id) REFERENCES payment_intents(id),
    CONSTRAINT fk_payment_transactions_merchant FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);

CREATE TABLE merchant_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id),
    currency currency_type NOT NULL,
    available_balance DECIMAL(15,2) DEFAULT 0,
    pending_balance DECIMAL(15,2) DEFAULT 0,
    total_processed DECIMAL(15,2) DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_merchant_balances_merchant FOREIGN KEY (merchant_id) REFERENCES merchants(id),
    UNIQUE(merchant_id, currency)
);

CREATE INDEX idx_merchants_email ON merchants(email);
CREATE INDEX idx_merchants_status ON merchants(status);
CREATE INDEX idx_merchant_api_keys_merchant ON merchant_api_keys(merchant_id);
CREATE INDEX idx_merchant_api_keys_api_key ON merchant_api_keys(api_key);
CREATE INDEX idx_merchant_api_keys_status ON merchant_api_keys(status);
CREATE INDEX idx_payment_intents_merchant ON payment_intents(merchant_id);
CREATE INDEX idx_payment_intents_transaction_id ON payment_intents(transaction_id);
CREATE INDEX idx_payment_intents_reference ON payment_intents(reference);
CREATE INDEX idx_payment_intents_status ON payment_intents(status);
CREATE INDEX idx_payment_intents_expires_at ON payment_intents(expires_at);
CREATE INDEX idx_payment_transactions_transaction_id ON payment_transactions(transaction_id);
CREATE INDEX idx_payment_transactions_intent_id ON payment_transactions(payment_intent_id);
CREATE INDEX idx_payment_transactions_merchant ON payment_transactions(merchant_id);
CREATE INDEX idx_payment_transactions_status ON payment_transactions(status);
CREATE INDEX idx_payment_transactions_third_party_ref ON payment_transactions(third_party_reference);
CREATE INDEX idx_merchant_balances_merchant ON merchant_balances(merchant_id);
CREATE INDEX idx_merchant_balances_currency ON merchant_balances(currency);
