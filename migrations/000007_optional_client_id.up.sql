-- Make client_id nullable in event table to allow anonymous event creation
-- Related to: https://github.com/nbtca/Saturday/issues/192

-- Drop existing foreign key constraint
ALTER TABLE public.event DROP CONSTRAINT event_client_id_fkey;

-- Make client_id nullable
ALTER TABLE public.event ALTER COLUMN client_id DROP NOT NULL;

-- Recreate foreign key constraint allowing NULL values
ALTER TABLE public.event ADD CONSTRAINT event_client_id_fkey 
    FOREIGN KEY (client_id) REFERENCES public.client(client_id) 
    ON UPDATE RESTRICT ON DELETE SET NULL;

-- Add index for efficient querying of anonymous events (where client_id IS NULL)
CREATE INDEX event_anonymous_idx ON public.event (client_id) WHERE client_id IS NULL;

-- Add index for better performance on non-null client_id queries
CREATE INDEX event_client_id_not_null_idx ON public.event (client_id) WHERE client_id IS NOT NULL;