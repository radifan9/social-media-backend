-- public.post_likes definition



CREATE TABLE public.post_likes (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	post_id uuid NOT NULL,
	user_id uuid NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT post_likes_pkey PRIMARY KEY (id),
	CONSTRAINT post_likes_post_id_user_id_key UNIQUE (post_id, user_id)
);


-- public.post_likes foreign keys

ALTER TABLE public.post_likes ADD CONSTRAINT post_likes_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE;
ALTER TABLE public.post_likes ADD CONSTRAINT post_likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;