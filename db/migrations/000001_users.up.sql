-- public.users definition



CREATE TABLE public.users (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	email varchar(255) NOT NULL,
	"password" text NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);