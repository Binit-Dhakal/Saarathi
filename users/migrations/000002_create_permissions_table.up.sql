create table if not exists permissions (
        id  serial primary key,
        code varchar(255) not null unique,
        description text
);

-- Role <-> Permission (many-to-many)
create table if not exists role_permissions(
        role_id int not null references roles(id) on delete cascade,
        permission_id int not null references permissions(id) on delete cascade,
        primary key (role_id, permission_id)
);

-- User <-> Role (many-to-many)
create table if not exists user_roles (
        user_id uuid not null references users(id) on delete cascade,
        role_id int not null references roles(id) on delete cascade,
        PRIMARY KEY(user_id, role_id)
);
