CREATE TYPE board AS ENUM (
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'gif', 'h', 'hr', 'k', 'm', 'o', 'p', 'r', 's', 't', 'u', 'v', 'vg', 'vm', 'vr', 'vrpg', 'vst', 'w', 'wg', 'i', 'ic', 'r9k', 's4s', 'vip', 'qa', 'cm', 'hm', 'lgbt', 'y', '3', 'aco', 'adv', 'an', 'bant', 'biz', 'cgl', 'ck', 'co', 'diy', 'fa', 'fit', 'gd', 'hc', 'his', 'int', 'jp', 'lit', 'mlp', 'mu', 'news', 'out', 'po', 'pol', 'pw', 'qst', 'sci', 'soc', 'sp', 'tg', 'toy', 'trv', 'tv', 'vp', 'vt', 'wsg', 'wsr', 'x', 'xs'
);

CREATE TYPE ext AS ENUM (
	'.jpg', '.png', '.gif', '.swf', '.pdf', '.webm'
);

CREATE TYPE capcode AS ENUM (
	'mod', 'admin', 'admin_highlight', 'manager', 'developer', 'founder'
);

CREATE TABLE post (
	board board NOT NULL,
	no BIGINT NOT NULL,
	resto BIGINT NOT NULL,
	time TIMESTAMP WITHOUT TIME ZONE,
	name TEXT,
	trip TEXT,
	capcode CAPCODE,
	country TEXT,
	since4pass SMALLINT,
	sub TEXT,
	com TEXT,
	tim BIGINT,
	md5 TEXT,
	filename TEXT,
	ext EXT,
	fsize BIGINT,
	w SMALLINT,
	h SMALLINT,
	tn_w SMALLINT,
	tn_h SMALLINT,
	deleted BOOLEAN NOT NULL,
	file_deleted BOOLEAN NOT NULL,
	spoiler BOOLEAN NOT NULL,
	custom_spoiler SMALLINT,
	op BOOLEAN NOT NULL,
	sticky BOOLEAN NOT NULL,
	name_tsvector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', name)) STORED,
	sub_tsvector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', sub)) STORED,
	com_tsvector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', com)) STORED,
	filename_tsvector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', filename)) STORED,

	PRIMARY KEY(board, no)
) PARTITION BY LIST(board);

CREATE TABLE post_sci PARTITION OF post
FOR VALUES IN ('sci');

CREATE TABLE post_wsr PARTITION OF post
FOR VALUES IN ('wsr');

