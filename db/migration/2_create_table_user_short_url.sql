CREATE TABLE IF NOT EXISTS public.user_short_url (
       id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
       user_id int8 NOT NULL,
       url_id int8 NOT NULL,
       CONSTRAINT user_short_url_url_list_fk FOREIGN KEY (url_id) REFERENCES public.url_list(id) ON DELETE CASCADE ON UPDATE CASCADE,
       CONSTRAINT user_short_url_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
