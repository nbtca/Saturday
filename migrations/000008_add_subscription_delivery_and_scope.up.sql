-- Add delivery method and scope to event_subscription table
ALTER TABLE public.event_subscription
ADD COLUMN delivery_method VARCHAR(20) DEFAULT 'webhook' NOT NULL,
ADD COLUMN email VARCHAR(255),
ADD COLUMN scope VARCHAR(20) DEFAULT 'related' NOT NULL;

-- Add check constraints
ALTER TABLE public.event_subscription
ADD CONSTRAINT check_delivery_method CHECK (delivery_method IN ('webhook', 'email', 'both'));

ALTER TABLE public.event_subscription
ADD CONSTRAINT check_scope CHECK (scope IN ('related', 'global'));

-- Ensure email is provided when delivery method includes email
ALTER TABLE public.event_subscription
ADD CONSTRAINT check_email_required CHECK (
    (delivery_method = 'webhook' AND email IS NULL) OR
    (delivery_method IN ('email', 'both') AND email IS NOT NULL)
);

-- Ensure callback_url is provided when delivery method includes webhook
ALTER TABLE public.event_subscription
ADD CONSTRAINT check_callback_required CHECK (
    (delivery_method = 'email' AND callback_url IS NULL) OR
    (delivery_method IN ('webhook', 'both') AND callback_url IS NOT NULL)
);

-- Make callback_url nullable since email-only subscriptions don't need it
ALTER TABLE public.event_subscription
ALTER COLUMN callback_url DROP NOT NULL;

-- Create index for email delivery
CREATE INDEX idx_subscription_delivery_method ON public.event_subscription(delivery_method) WHERE active = true;
CREATE INDEX idx_subscription_scope ON public.event_subscription(scope) WHERE active = true;
