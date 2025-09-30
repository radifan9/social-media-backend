-- public.user_profiles definition



CREATE TABLE public.user_profiles (
	user_id uuid,
	"name" varchar(50),
	avatar text,
	bio text,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT user_profiles_pkey PRIMARY KEY (user_id)
);


-- public.user_profiles foreign keys

ALTER TABLE public.user_profiles ADD CONSTRAINT user_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;