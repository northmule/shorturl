package migrations

const Migrations01 = `CREATE TABLE IF NOT EXISTS public.url_list (
           id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
           short_url varchar(100) NOT NULL,
           url varchar(2000) NOT NULL,
           created_at timestamp DEFAULT now() NOT NULL,
           deleted_at timestamp NULL,
           CONSTRAINT url_pk PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS url_list_url_idx ON public.url_list (url) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS short_url_idx ON public.url_list USING btree (short_url)`

const Migrations02 = `CREATE TABLE IF NOT EXISTS public.users (
       id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
       name varchar(100) NOT NULL,
       login varchar(100) NOT NULL,
       password varchar(200) NOT NULL,
       created_at timestamp DEFAULT now() NOT NULL,
       deleted_at timestamp NULL,
       "uuid" uuid NULL,
       CONSTRAINT users_uuid_unique UNIQUE (uuid),
       CONSTRAINT users_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS users_login_password_idx ON public.users (login,"password");
CREATE UNIQUE INDEX IF NOT EXISTS users_login_idx ON public.users (login) WHERE deleted_at IS NULL;`

const Migrations03 = `CREATE TABLE IF NOT EXISTS public.user_short_url (
       id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
       user_id int8 NOT NULL,
       url_id int8 NOT NULL,
       CONSTRAINT user_short_url_url_list_fk FOREIGN KEY (url_id) REFERENCES public.url_list(id) ON DELETE CASCADE ON UPDATE CASCADE,
       CONSTRAINT user_short_url_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
`
