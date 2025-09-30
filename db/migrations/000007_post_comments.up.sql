-- public.post_comments definition



CREATE TABLE public.post_comments (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	post_id uuid NOT NULL,
	user_id uuid NOT NULL,
	"comment" text NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT post_comments_pkey PRIMARY KEY (id)
);


-- public.post_comments foreign keys

ALTER TABLE public.post_comments ADD CONSTRAINT post_comments_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE;
ALTER TABLE public.post_comments ADD CONSTRAINT post_comments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;