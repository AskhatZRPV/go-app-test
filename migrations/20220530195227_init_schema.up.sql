CREATE TABLE users (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  phone varchar not null unique,
  name varchar not null
);

CREATE TABLE doctors (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  name varchar not null,
  spec varchar not null,
  slots jsonb not null default '[]'::jsonb
);