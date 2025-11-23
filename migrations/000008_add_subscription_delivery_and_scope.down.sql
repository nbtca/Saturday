-- Remove indexes
DROP INDEX IF EXISTS public.idx_subscription_delivery_method;
DROP INDEX IF EXISTS public.idx_subscription_scope;

-- Remove constraints
ALTER TABLE public.event_subscription
DROP CONSTRAINT IF EXISTS check_delivery_method;

ALTER TABLE public.event_subscription
DROP CONSTRAINT IF EXISTS check_scope;

ALTER TABLE public.event_subscription
DROP CONSTRAINT IF EXISTS check_email_required;

ALTER TABLE public.event_subscription
DROP CONSTRAINT IF EXISTS check_callback_required;

-- Restore callback_url NOT NULL constraint
ALTER TABLE public.event_subscription
ALTER COLUMN callback_url SET NOT NULL;

-- Remove columns
ALTER TABLE public.event_subscription
DROP COLUMN IF EXISTS delivery_method;

ALTER TABLE public.event_subscription
DROP COLUMN IF EXISTS email;

ALTER TABLE public.event_subscription
DROP COLUMN IF EXISTS scope;
