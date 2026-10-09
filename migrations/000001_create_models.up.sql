CREATE TABLE sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT NOT NULL UNIQUE,
    status INTEGER NOT NULL,
    consecutive_fails INTEGER NOT NULL DEFAULT 0,
    interval INTEGER NOT NULL,
    last_check_at TIMESTAMPTZ,
    next_check_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id UUID NOT NULL,
    status INTEGER NOT NULL,
    status_code INTEGER NOT NULL,
    latency INTEGER NOT NULL,
    error TEXT,
    checked_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_checks_site
        FOREIGN KEY (site_id)
            REFERENCES sites(id)
            ON DELETE CASCADE
);