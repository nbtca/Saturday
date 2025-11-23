-- Create event_subscription table
CREATE TABLE IF NOT EXISTS public.event_subscription (
    subscription_id BIGSERIAL PRIMARY KEY,
    member_id VARCHAR(10),
    client_id BIGINT,
    event_types TEXT[] NOT NULL,
    callback_url VARCHAR(500) NOT NULL,
    secret VARCHAR(100) NOT NULL,
    filters JSONB,
    active BOOLEAN DEFAULT true NOT NULL,
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_subscription_member FOREIGN KEY (member_id) REFERENCES public.member(member_id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_client FOREIGN KEY (client_id) REFERENCES public.client(client_id) ON DELETE CASCADE,
    CONSTRAINT check_owner CHECK (member_id IS NOT NULL OR client_id IS NOT NULL)
);

-- Create indexes for better query performance
CREATE INDEX idx_subscription_active ON public.event_subscription(active) WHERE active = true;
CREATE INDEX idx_subscription_member ON public.event_subscription(member_id) WHERE member_id IS NOT NULL;
CREATE INDEX idx_subscription_client ON public.event_subscription(client_id) WHERE client_id IS NOT NULL;
CREATE INDEX idx_subscription_event_types ON public.event_subscription USING GIN(event_types);

-- Create subscription delivery log table for tracking webhook deliveries
CREATE TABLE IF NOT EXISTS public.event_subscription_delivery (
    delivery_id BIGSERIAL PRIMARY KEY,
    subscription_id BIGINT NOT NULL,
    event_id BIGINT,
    event_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL, -- pending, success, failed
    attempts INT DEFAULT 0 NOT NULL,
    last_attempt TIMESTAMP,
    response_code INT,
    response_body TEXT,
    error_message TEXT,
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_delivery_subscription FOREIGN KEY (subscription_id) REFERENCES public.event_subscription(subscription_id) ON DELETE CASCADE,
    CONSTRAINT fk_delivery_event FOREIGN KEY (event_id) REFERENCES public.event(event_id) ON DELETE SET NULL
);

-- Create indexes for delivery tracking
CREATE INDEX idx_delivery_subscription ON public.event_subscription_delivery(subscription_id);
CREATE INDEX idx_delivery_status ON public.event_subscription_delivery(status);
CREATE INDEX idx_delivery_event ON public.event_subscription_delivery(event_id) WHERE event_id IS NOT NULL;
CREATE INDEX idx_delivery_created ON public.event_subscription_delivery(gmt_create DESC);
