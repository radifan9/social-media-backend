-- public.posts definition



CREATE TABLE public.posts (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	user_id uuid NOT NULL,
	text_content text,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT posts_pkey PRIMARY KEY (id)
);


-- public.posts foreign keys

ALTER TABLE public.posts ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;