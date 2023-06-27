CREATE TABLE links
(
    code       CHAR(10) PRIMARY KEY NOT NULL,
    url        VARCHAR(100000)      NOT NULL,
    created_at TIMESTAMPTZ          NOT NULL DEFAULT NOW()
);
