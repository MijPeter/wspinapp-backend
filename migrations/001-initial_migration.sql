CREATE TABLE "walls"
(
    "id"                bigserial,
    "created_at"        timestamptz,
    "updated_at"        timestamptz,
    "deleted_at"        timestamptz,
    "image_url"         text,
    "image_preview_url" text,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_walls_deleted_at" ON "walls" ("deleted_at");
CREATE TABLE "routes"
(
    "id"         bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "wall_id"    bigint,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_routes_deleted_at" ON "routes" ("deleted_at");
CREATE TABLE "holds"
(
    "id"         bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "x"          decimal NOT NULL,
    "y"          decimal NOT NULL,
    "size"       decimal NOT NULL,
    "wall_id"    bigint  NOT NULL,
    "shape"      text    NOT NULL DEFAULT 'circle',
    "angle"      decimal,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_walls_holds" FOREIGN KEY ("wall_id") REFERENCES "walls" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_holds_deleted_at" ON "holds" ("deleted_at");
CREATE TABLE "route_top_hold"
(
    "route_id" bigint,
    "hold_id"  bigint,
    PRIMARY KEY ("route_id", "hold_id"),
    CONSTRAINT "fk_route_top_hold_route" FOREIGN KEY ("route_id") REFERENCES "routes" ("id"),
    CONSTRAINT "fk_route_top_hold_hold" FOREIGN KEY ("hold_id") REFERENCES "holds" ("id")
);
CREATE TABLE "route_start_holds"
(
    "route_id" bigint,
    "hold_id"  bigint,
    PRIMARY KEY ("route_id", "hold_id"),
    CONSTRAINT "fk_route_start_holds_hold" FOREIGN KEY ("hold_id") REFERENCES "holds" ("id"),
    CONSTRAINT "fk_route_start_holds_route" FOREIGN KEY ("route_id") REFERENCES "routes" ("id")
);
CREATE TABLE "route_holds"
(
    "route_id" bigint,
    "hold_id"  bigint,
    PRIMARY KEY ("route_id", "hold_id"),
    CONSTRAINT "fk_route_holds_route" FOREIGN KEY ("route_id") REFERENCES "routes" ("id"),
    CONSTRAINT "fk_route_holds_hold" FOREIGN KEY ("hold_id") REFERENCES "holds" ("id")
);
