CREATE TABLE deployments (
                             id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             git_url       TEXT NOT NULL,
                             status        TEXT NOT NULL DEFAULT 'queued',

                             container_id  TEXT,
                             image_name    TEXT,

                             logs          TEXT,

                             created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
                             updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);