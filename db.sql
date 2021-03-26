CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE realms (
    id bytea NOT NULL UNIQUE,
    name VARCHAR(256) NOT NULL,
    slug VARCHAR(256) NOT NULL UNIQUE,
    website_url VARCHAR(1024),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TRIGGER update_realms_update_at BEFORE UPDATE
ON realms FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at();

CREATE TABLE users (
    id bytea NOT NULL UNIQUE,
    realm_id bytea NOT NULL,
    name VARCHAR(256) NOT NULL,
    login VARCHAR(256) NOT NULL,
    email VARCHAR(256),
    password VARCHAR(256) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (realm_id) REFERENCES realms (id)
);

CREATE TRIGGER update_users_update_at BEFORE UPDATE
ON users FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at();

CREATE TABLE tasks (
    id bytea NOT NULL UNIQUE,
    short_id VARCHAR(16) NOT NULL UNIQUE,
    title VARCHAR(256) NOT NULL,
    content TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TRIGGER update_tasks_update_at BEFORE UPDATE
ON tasks FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at();

CREATE TABLE completed_tasks (
    task_id bytea NOT NULL,
    user_id VARCHAR(256) NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks (id)
);

CREATE TABLE groups (
    id bytea NOT NULL UNIQUE,
    realm_id bytea NOT NULL,
    usable BOOLEAN NOT NULL DEFAULT true,
    slug VARCHAR(256) NOT NULL,
    name VARCHAR(256) NOT NULL,
    parent_id bytea,
    active BOOLEAN NOT NULL DEFAULT false,
    active_at TIMESTAMP,
    archived BOOLEAN NOT NULL DEFAULT false,
    archived_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (realm_id) REFERENCES realms (id),
    FOREIGN KEY (parent_id) REFERENCES groups (id)
);

CREATE TRIGGER update_groups_update_at BEFORE UPDATE
ON groups FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at();

CREATE TABLE admins (
    id bytea NOT NULL UNIQUE,
    realm_id bytea NOT NULL,
    login VARCHAR(256) NOT NULL,
    password VARCHAR(128) NOT NULL,
    name VARCHAR(256) NOT NULL,
    email VARCHAR(256) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (realm_id) REFERENCES realms (id)
);

CREATE TRIGGER update_admins_update_at BEFORE UPDATE
ON admins FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at();

CREATE TABLE group_users (
    group_id bytea NOT NULL,
    user_id bytea NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Enable function `unaccent` (removes accents)
CREATE EXTENSION unaccent;
