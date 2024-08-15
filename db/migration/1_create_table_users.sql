CREATE TABLE IF NOT EXISTS public.users (
       id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
       name varchar(100) NOT NULL,
       login varchar(100) NOT NULL,
       password varchar(200) NOT NULL,
       created_at timestamp DEFAULT now() NOT NULL,
       deleted_at timestamp NULL,
       CONSTRAINT users_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS users_login_password_idx ON public.users (login,"password");
CREATE UNIQUE INDEX IF NOT EXISTS users_login_idx ON public.users (login) WHERE deleted_at IS NULL;