-- Drop delivery log table first (due to foreign key dependency)
DROP TABLE IF EXISTS public.event_subscription_delivery;

-- Drop subscription table
DROP TABLE IF EXISTS public.event_subscription;
