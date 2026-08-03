CREATE TABLE IF NOT EXISTS bookmarks
(
    id              varchar(36)         NOT NULL,
    user_id         varchar(36)         NOT NULL,
    code            varchar(10)         NOT NULL,
    url             varchar(2048)       NOT NULL,
    description     varchar(255),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE,
    CONSTRAINT      bookmark_pkey       PRIMARY KEY (id),
    CONSTRAINT      uni_code            UNIQUE (code),
    CONSTRAINT      fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
)