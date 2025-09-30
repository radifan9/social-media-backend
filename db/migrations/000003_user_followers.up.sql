-- public.user_followers definition



CREATE TABLE public.user_followers (
	user_id uuid NOT NULL,
	follower_id uuid NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT check_not_self_follow CHECK ((user_id <> follower_id)),
	CONSTRAINT user_followers_pkey PRIMARY KEY (user_id, follower_id)
);


-- public.user_followers foreign keys

ALTER TABLE public.user_followers ADD CONSTRAINT user_followers_follower_id_fkey FOREIGN KEY (follower_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE public.user_followers ADD CONSTRAINT user_followers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;