-- public.post_images definition



CREATE TABLE public.post_images (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	post_id uuid NOT NULL,
	image_url text NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT post_images_pkey PRIMARY KEY (id)
);


-- public.post_images foreign keys

ALTER TABLE public.post_images ADD CONSTRAINT post_images_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE;