CREATE TABLE IF NOT EXISTS public.url_list (
           id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
           short_url varchar(100) NOT NULL,
           url varchar(2000) NOT NULL,
           created_at timestamp DEFAULT now() NOT NULL,
           deleted_at timestamp NULL,
           CONSTRAINT url_pk PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS url_list_url_idx ON public.url_list (url) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS short_url_idx ON public.url_list USING btree (short_url)