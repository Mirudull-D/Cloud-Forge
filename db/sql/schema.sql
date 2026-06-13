CREATE TABLE deployments (
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                             git_url TEXT NOT NULL,

                             status TEXT NOT NULL DEFAULT 'queued'
                                 CHECK (status IN (
                                                   'queued',
                                                   'building',
                                                   'running',
                                                   'failed',
                                                   'stopped'
                                     )),

                             image_name TEXT,
                             container_id TEXT,

                             port INTEGER UNIQUE,

                             error_message TEXT,

                             created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                             updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);