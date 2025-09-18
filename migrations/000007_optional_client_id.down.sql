-- Rollback: Make client_id mandatory again in event table
-- This rollback will fail if there are any events with NULL client_id

-- Remove the indexes we added
DROP INDEX IF EXISTS public.event_anonymous_idx;
DROP INDEX IF EXISTS public.event_client_id_not_null_idx;

-- Drop the modified foreign key constraint
ALTER TABLE public.event DROP CONSTRAINT event_client_id_fkey;

-- Make client_id NOT NULL again (this will fail if there are NULL values)
ALTER TABLE public.event ALTER COLUMN client_id SET NOT NULL;

-- Recreate the original foreign key constraint
ALTER TABLE public.event ADD CONSTRAINT event_client_id_fkey 
    FOREIGN KEY (client_id) REFERENCES public.client(client_id) 
    ON UPDATE RESTRICT ON DELETE RESTRICT;