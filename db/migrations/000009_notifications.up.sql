-- public.notifications definition


CREATE TABLE public.notifications (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	recipient_id uuid NOT NULL,
	actor_id uuid NOT NULL,
	"type" public."notification_type" NOT NULL,
	post_id uuid NULL,
	comment_id uuid NULL,
	read_at timestamptz NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT check_not_self_notify CHECK ((recipient_id <> actor_id)),
	CONSTRAINT notifications_pkey PRIMARY KEY (id)
);


-- public.notifications foreign keys

ALTER TABLE public.notifications ADD CONSTRAINT notifications_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE public.notifications ADD CONSTRAINT notifications_comment_id_fkey FOREIGN KEY (comment_id) REFERENCES public.post_comments(id) ON DELETE CASCADE;
ALTER TABLE public.notifications ADD CONSTRAINT notifications_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE;
ALTER TABLE public.notifications ADD CONSTRAINT notifications_recipient_id_fkey FOREIGN KEY (recipient_id) REFERENCES public.users(id) ON DELETE CASCADE;