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


CREATE TABLE post_a PARTITION OF post
FOR VALUES IN ('a');

CREATE TABLE post_c PARTITION OF post
FOR VALUES IN ('c');

CREATE TABLE post_d PARTITION OF post
FOR VALUES IN ('d');

CREATE TABLE post_e PARTITION OF post
FOR VALUES IN ('e');

CREATE TABLE post_g PARTITION OF post
FOR VALUES IN ('g');

CREATE TABLE post_gif PARTITION OF post
FOR VALUES IN ('gif');

CREATE TABLE post_h PARTITION OF post
FOR VALUES IN ('h');

CREATE TABLE post_hr PARTITION OF post
FOR VALUES IN ('hr');

CREATE TABLE post_k PARTITION OF post
FOR VALUES IN ('k');

CREATE TABLE post_m PARTITION OF post
FOR VALUES IN ('m');

CREATE TABLE post_o PARTITION OF post
FOR VALUES IN ('o');

CREATE TABLE post_p PARTITION OF post
FOR VALUES IN ('p');

CREATE TABLE post_r PARTITION OF post
FOR VALUES IN ('r');

CREATE TABLE post_s PARTITION OF post
FOR VALUES IN ('s');

CREATE TABLE post_t PARTITION OF post
FOR VALUES IN ('t');

CREATE TABLE post_u PARTITION OF post
FOR VALUES IN ('u');

CREATE TABLE post_v PARTITION OF post
FOR VALUES IN ('v');

CREATE TABLE post_vg PARTITION OF post
FOR VALUES IN ('vg');

CREATE TABLE post_vm PARTITION OF post
FOR VALUES IN ('vm');

CREATE TABLE post_vr PARTITION OF post
FOR VALUES IN ('vr');

CREATE TABLE post_vrpg PARTITION OF post
FOR VALUES IN ('vrpg');

CREATE TABLE post_vst PARTITION OF post
FOR VALUES IN ('vst');

CREATE TABLE post_w PARTITION OF post
FOR VALUES IN ('w');

CREATE TABLE post_wg PARTITION OF post
FOR VALUES IN ('wg');

CREATE TABLE post_i PARTITION OF post
FOR VALUES IN ('i');

CREATE TABLE post_ic PARTITION OF post
FOR VALUES IN ('ic');

CREATE TABLE post_r9k PARTITION OF post
FOR VALUES IN ('r9k');

CREATE TABLE post_s4s PARTITION OF post
FOR VALUES IN ('s4s');

CREATE TABLE post_vip PARTITION OF post
FOR VALUES IN ('vip');

CREATE TABLE post_qa PARTITION OF post
FOR VALUES IN ('qa');

CREATE TABLE post_cm PARTITION OF post
FOR VALUES IN ('cm');

CREATE TABLE post_hm PARTITION OF post
FOR VALUES IN ('hm');

CREATE TABLE post_lgbt PARTITION OF post
FOR VALUES IN ('lgbt');

CREATE TABLE post_y PARTITION OF post
FOR VALUES IN ('y');

CREATE TABLE post_3 PARTITION OF post
FOR VALUES IN ('3');

CREATE TABLE post_aco PARTITION OF post
FOR VALUES IN ('aco');

CREATE TABLE post_adv PARTITION OF post
FOR VALUES IN ('adv');

CREATE TABLE post_an PARTITION OF post
FOR VALUES IN ('an');

CREATE TABLE post_biz PARTITION OF post
FOR VALUES IN ('biz');

CREATE TABLE post_cgl PARTITION OF post
FOR VALUES IN ('cgl');

CREATE TABLE post_ck PARTITION OF post
FOR VALUES IN ('ck');

CREATE TABLE post_co PARTITION OF post
FOR VALUES IN ('co');

CREATE TABLE post_diy PARTITION OF post
FOR VALUES IN ('diy');

CREATE TABLE post_fa PARTITION OF post
FOR VALUES IN ('fa');

CREATE TABLE post_fit PARTITION OF post
FOR VALUES IN ('fit');

CREATE TABLE post_gd PARTITION OF post
FOR VALUES IN ('gd');

CREATE TABLE post_hc PARTITION OF post
FOR VALUES IN ('hc');

CREATE TABLE post_his PARTITION OF post
FOR VALUES IN ('his');

CREATE TABLE post_int PARTITION OF post
FOR VALUES IN ('int');

CREATE TABLE post_jp PARTITION OF post
FOR VALUES IN ('jp');

CREATE TABLE post_lit PARTITION OF post
FOR VALUES IN ('lit');

CREATE TABLE post_mlp PARTITION OF post
FOR VALUES IN ('mlp');

CREATE TABLE post_mu PARTITION OF post
FOR VALUES IN ('mu');

CREATE TABLE post_news PARTITION OF post
FOR VALUES IN ('news');

CREATE TABLE post_out PARTITION OF post
FOR VALUES IN ('out');

CREATE TABLE post_po PARTITION OF post
FOR VALUES IN ('po');

CREATE TABLE post_pol PARTITION OF post
FOR VALUES IN ('pol');

CREATE TABLE post_pw PARTITION OF post
FOR VALUES IN ('pw');

CREATE TABLE post_qst PARTITION OF post
FOR VALUES IN ('qst');

CREATE TABLE post_sci PARTITION OF post
FOR VALUES IN ('sci');

CREATE TABLE post_soc PARTITION OF post
FOR VALUES IN ('soc');

CREATE TABLE post_sp PARTITION OF post
FOR VALUES IN ('sp');

CREATE TABLE post_tg PARTITION OF post
FOR VALUES IN ('tg');

CREATE TABLE post_toy PARTITION OF post
FOR VALUES IN ('toy');

CREATE TABLE post_trv PARTITION OF post
FOR VALUES IN ('trv');

CREATE TABLE post_tv PARTITION OF post
FOR VALUES IN ('tv');

CREATE TABLE post_vp PARTITION OF post
FOR VALUES IN ('vp');

CREATE TABLE post_vt PARTITION OF post
FOR VALUES IN ('vt');

CREATE TABLE post_wsg PARTITION OF post
FOR VALUES IN ('wsg');

CREATE TABLE post_wsr PARTITION OF post
FOR VALUES IN ('wsr');

CREATE TABLE post_x PARTITION OF post
FOR VALUES IN ('x');

CREATE TABLE post_xs PARTITION OF post
FOR VALUES IN ('xs');

