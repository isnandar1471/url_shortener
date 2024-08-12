-- Adminer 4.8.1 PostgreSQL 16.3 dump

-- \connect "url_shortener";

CREATE SEQUENCE short_clicks_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."short_clicks" (
    "id" integer DEFAULT nextval('short_clicks_id_seq') NOT NULL,
    "short_id" integer NOT NULL,
    "clicked_at" integer NOT NULL,
    "ip_address" character varying NOT NULL,
    "user_agent" character varying NOT NULL,
    CONSTRAINT "short_clicks_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "short_clicks" ("id", "short_id", "clicked_at", "ip_address", "user_agent") VALUES
(1,	15,	1723175237,	'ipaddres',	'user agent'),
(2,	15,	1723175297,	'ipaddres',	'user agent'),
(3,	15,	1723175315,	'ipaddres',	'user agent'),
(4,	15,	1723175389,	'ipaddres',	'user agent'),
(5,	15,	1723175446,	'ipaddres',	'user agent'),
(6,	15,	1723175504,	'ipaddres',	'user agent'),
(7,	15,	1723175657,	'ipaddres',	'user agent'),
(8,	15,	1723175687,	'ipaddres',	'user agent'),
(9,	15,	1723175764,	'ipaddres',	'user agent'),
(10,	15,	1723175882,	'ipaddres',	'user agent'),
(11,	15,	1723175938,	'ipaddres',	'user agent'),
(12,	15,	1723176108,	'ipaddres',	'user agent'),
(13,	15,	1723176113,	'ipaddres',	'user agent'),
(14,	15,	1723176119,	'ipaddres',	'user agent'),
(15,	18,	1723495915,	'ipaddres',	'user agent'),
(16,	16,	1723496041,	'ipaddres',	'user agent'),
(17,	17,	1723496050,	'ipaddres',	'user agent');

CREATE SEQUENCE shorts_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."shorts" (
    "id" integer DEFAULT nextval('shorts_id_seq') NOT NULL,
    "name" character varying NOT NULL,
    "code" character varying NOT NULL,
    "destination_url" text NOT NULL,
    "user_id" integer,
    "created_at" integer,
    "click_count" integer DEFAULT '0' NOT NULL,
    CONSTRAINT "shorts_code" UNIQUE ("code"),
    CONSTRAINT "shorts_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "shorts" ("id", "name", "code", "destination_url", "user_id", "created_at", "click_count") VALUES
(4,	'namaganti',	'zzzzzzganti',	'urlganti',	3,	NULL,	0),
(15,	'halo',	'hai',	'https://www.google.com/',	3,	1723149870,	14),
(18,	'Isnandar Instagram',	'ig',	'https://instagram.com/isnandar_fajar_pangestu',	6,	1723495876,	1),
(16,	'Isnandar WhatsApp',	'wa',	'https://wa.me/6285814718596',	6,	1723495381,	1),
(17,	'Isnandar Github',	'gh',	'https://github.com/isnandar1471',	6,	1723495449,	1);

CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "username" character varying NOT NULL,
    "email" character varying NOT NULL,
    "password_hash" character varying NOT NULL,
    "created_at" integer NOT NULL,
    "shorts_count" integer DEFAULT '0' NOT NULL,
    CONSTRAINT "users_email_unique" UNIQUE ("email"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "users_username_unique" UNIQUE ("username")
) WITH (oids = false);

INSERT INTO "users" ("id", "username", "email", "password_hash", "created_at", "shorts_count") VALUES
(3,	'fajar',	'fajar',	'$2a$10$PyozqvO4h71/8YqaFZalyuc1NTNQH336cibmayNY5NT2QSUcC5t9S',	1723077920,	1),
(6,	'isnandar1',	'isnandar1',	'$2a$10$550V9VyRNRo7rXsIs3Rx2uG5KIYYA1UItLOqkWi20wyL/ApwnANaO',	1723492338,	3),
(7,	'guest1',	'guest1',	'$2a$10$qfyU1mhuW9B7ZAuZF9a6v.ccuJhldw.ck3s/RJ8VIsPVR6y4IQ2za',	1723496133,	0);

-- 2024-08-12 20:58:41.460369+00
