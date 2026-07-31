CREATE TABLE IF NOT EXISTS users
(
    id           varchar(36)                                NOT NULL,
    display_name varchar(255)                               NOT NULL,       
    username     varchar(255)                               NOT NULL,
    email        varchar(255)                               NOT NULL,
    password     varchar(2048)                              NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT uni_username UNIQUE (username),
    CONSTRAINT uni_email UNIQUE (email)
);