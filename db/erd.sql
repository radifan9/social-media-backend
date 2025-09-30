-- ============================================
-- SIMPLIFIED SOCIAL MEDIA DATABASE SCHEMA
-- No indexes except PRIMARY KEY and UNIQUE constraints
-- ============================================

-- ENUM untuk tipe notifikasi
CREATE TYPE "notification_type" AS ENUM (
  'follow',
  'like',
  'comment'
);

-- ============================================
-- USERS & PROFILES
-- ============================================

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar(255) UNIQUE NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamptz DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "user_profiles" (
  "user_id" uuid PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "avatar" text,
  "bio" text,
  "created_at" timestamptz DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamptz DEFAULT (CURRENT_TIMESTAMP)
);

-- ============================================
-- FOLLOWERS
-- ============================================

CREATE TABLE "user_followers" (
  "user_id" uuid NOT NULL,
  "follower_id" uuid NOT NULL,
  "created_at" timestamptz DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("user_id", "follower_id"),
  CONSTRAINT "check_not_self_follow" CHECK ("user_id" != "follower_id")
);

-- ============================================
-- POSTS & IMAGES
-- ============================================

CREATE TABLE "posts" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL,
  "text_content" text,
  "created_at" timestamptz DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamptz DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "post_images" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "post_id" uuid NOT NULL,
  "image_url" text NOT NULL,
  "position" smallint DEFAULT 0,
  "created_at" timestamptz DEFAULT (CURRENT_TIMESTAMP)
);

-- ============================================
-- INTERACTIONS
-- ============================================

CREATE TABLE "post_likes" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "post_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamptz DEFAULT (CURRENT_TIMESTAMP),
  UNIQUE ("post_id", "user_id")
);

CREATE TABLE "post_comments" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "post_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "comment" text NOT NULL,
  "created_at" timestamptz DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamptz DEFAULT (CURRENT_TIMESTAMP)
);

-- ============================================
-- NOTIFICATIONS
-- ============================================

CREATE TABLE "notifications" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "recipient_id" uuid NOT NULL,
  "actor_id" uuid NOT NULL,
  "type" notification_type NOT NULL,
  "post_id" uuid,
  "comment_id" uuid,
  "read_at" timestamptz,
  "created_at" timestamptz DEFAULT (CURRENT_TIMESTAMP),
  CONSTRAINT "check_not_self_notify" CHECK ("recipient_id" != "actor_id")
);

-- ============================================
-- FOREIGN KEYS
-- ============================================

ALTER TABLE "user_profiles" 
  ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "user_followers" 
  ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "user_followers" 
  ADD FOREIGN KEY ("follower_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "posts" 
  ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "post_images" 
  ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;

ALTER TABLE "post_likes" 
  ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;

ALTER TABLE "post_likes" 
  ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "post_comments" 
  ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;

ALTER TABLE "post_comments" 
  ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "notifications" 
  ADD FOREIGN KEY ("recipient_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "notifications" 
  ADD FOREIGN KEY ("actor_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "notifications" 
  ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;

ALTER TABLE "notifications" 
  ADD FOREIGN KEY ("comment_id") REFERENCES "post_comments" ("id") ON DELETE CASCADE;