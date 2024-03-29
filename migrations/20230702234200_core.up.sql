CREATE EXTENSION IF NOT EXISTS "pgcrypto" CASCADE;
CREATE EXTENSION IF NOT EXISTS "earthdistance" CASCADE;


CREATE FUNCTION "trigger_set_timestamp"()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = CURRENT_TIMESTAMP;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE "configs" (
	"key" VARCHAR(255) NOT NULL PRIMARY KEY,
	"value" JSONB NOT NULL DEFAULT 'null',
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX "configs_unique_idx" ON "configs" ("key");
CREATE TRIGGER "configs_set_timestamp" BEFORE UPDATE ON "configs" FOR EACH ROW EXECUTE PROCEDURE "trigger_set_timestamp"();

INSERT INTO "configs" ("key", "value") VALUES
	('host', '"0.0.0.0"'),
	('port', '3000'),
	('uri', '"http://localhost:3000"'),
	('ui.uri', '"http://localhost:5173"'),
	('dashboard.enabled', 'true'),
	('openapi.enabled', 'true'),
	('jwt.secret', to_jsonb(encode(public.gen_random_bytes(255), 'base64'))),
	('jwt.expires.seconds', '86400');


CREATE FUNCTION configs_notify() RETURNS trigger AS $$
DECLARE
  "key" VARCHAR(255);
  "value" JSONB;
BEGIN
	IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
		"key" = NEW."key";
	ELSE
		"key" = OLD."key";
	END IF;
	IF TG_OP != 'UPDATE' OR NEW."value" != OLD."value" THEN
		PERFORM pg_notify('configs_channel', json_build_object('table', TG_TABLE_NAME, 'key', "key", 'value', NEW."value", 'action_type', TG_OP)::text);
	END IF;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "configs_notify_update" AFTER UPDATE ON "configs" FOR EACH ROW EXECUTE PROCEDURE configs_notify();
CREATE TRIGGER "configs_notify_insert" AFTER INSERT ON "configs" FOR EACH ROW EXECUTE PROCEDURE configs_notify();
CREATE TRIGGER "configs_notify_delete" AFTER DELETE ON "configs" FOR EACH ROW EXECUTE PROCEDURE configs_notify();


CREATE TYPE SEX AS ENUM ('male', 'female');


CREATE TABLE "users"(
	"id" SERIAL PRIMARY KEY,
	"username" VARCHAR(255) NOT NULL,
	"email" VARCHAR(255) NOT NULL,
	"terms_of_service_acknowledge" BOOLEAN DEFAULT false NOT NULL,
	"birthdate" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP - INTERVAL '18 years' NOT NULL,
	"sex" SEX DEFAULT 'male' NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX "users_username_unique_idx" ON "users" ("username");
CREATE UNIQUE INDEX "users_email_unique_idx" ON "users" ("email");
CREATE TRIGGER "users_set_timestamp" BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE "trigger_set_timestamp"();


CREATE TABLE "places"(
	"id" SERIAL PRIMARY KEY,
	"creator_id" INT4 NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"street_address_line1" VARCHAR(255) NOT NULL,
	"street_address_line2" VARCHAR(255) NOT NULL DEFAULT '',
	"locality" VARCHAR(255) NOT NULL DEFAULT '',
	"region" VARCHAR(255) NOT NULL DEFAULT '',
	"postal_code" VARCHAR(255) NOT NULL DEFAULT '',
	"country" VARCHAR(255) NOT NULL DEFAULT '',
	"latitude" FLOAT8 NOT NULL,
	"longitude" FLOAT8 NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT "places_creator_id_fk" FOREIGN KEY("creator_id") REFERENCES "users"("id") ON DELETE CASCADE
);
CREATE TRIGGER "places_set_timestamp" BEFORE UPDATE ON "places" FOR EACH ROW EXECUTE PROCEDURE "trigger_set_timestamp"();


CREATE TYPE FEATURE AS ENUM (
	'full-basketball-court',
	'half-basketball-court',
	'full-football-field',
	'half-football-field',
	'mini-football-field',
	'full-american-football-field',
	'tennis-court',
	'tennis-wall',
	'ping-pong-table',
	'pool-table',
	'racquetball-court');


CREATE TYPE FEATURE_LOCATION AS ENUM (
	'indoor', 
	'outdoor', 
	'outdoor-covered');


CREATE TABLE "place_features"(
	"id" SERIAL PRIMARY KEY,
	"place_id" INT4 NOT NULL,
	"type" FEATURE NOT NULL,
	"location_type" FEATURE_LOCATION NOT NULL DEFAULT 'indoor',
	"latitude_offset" FLOAT4,
	"longitude_offset" FLOAT4,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT "place_features_place_id_fk" FOREIGN KEY("place_id") REFERENCES "places"("id") ON DELETE CASCADE
);