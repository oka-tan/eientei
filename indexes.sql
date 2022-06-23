/*
Basic search indexes.
*/
CREATE UNIQUE INDEX post_resto_no_index ON post(board, resto, no);
CREATE UNIQUE INDEX post_deleted_no_index ON post (board, no) WHERE deleted;
CREATE UNIQUE INDEX post_spoiler_no_index ON post (board, no) WHERE spoiler;
CREATE UNIQUE INDEX post_op_no_index ON post (board, no) WHERE op;
CREATE UNIQUE INDEX post_image_no_index ON post(board, no) WHERE tim IS NOT NULL;
CREATE UNIQUE INDEX post_capcode_no_index ON post (board, capcode, no) WHERE capcode IS NOT NULL;
CREATE UNIQUE INDEX post_since4pass_no_index ON post (board, no) WHERE since4pass IS NOT NULL;
CREATE UNIQUE INDEX post_md5_no_index ON post (board, md5, no) WHERE md5 IS NOT NULL;
CREATE UNIQUE INDEX post_trip_no_index ON post (board, trip, no) WHERE trip IS NOT NULL;
CREATE UNIQUE INDEX post_country_no_index ON post (board, country, no) WHERE country IS NOT NULL;

/*
Search indexes by time.
*/
CREATE UNIQUE INDEX post_time_no_index ON post(board, time, no);
CREATE UNIQUE INDEX post_deleted_time_no_index ON post(board, time, no) WHERE deleted;
CREATE UNIQUE INDEX post_spoiler_time_no_index ON post(board, time, no) WHERE spoiler;
CREATE UNIQUE INDEX post_op_time_no_index ON post(board, time, no) WHERE op;
CREATE UNIQUE INDEX post_image_time_no_index ON post(board, time, no) WHERE tim IS NOT NULL;
CREATE UNIQUE INDEX post_capcode_time_no_index ON post(board, capcode, time, no) WHERE capcode IS NOT NULL;
CREATE UNIQUE INDEX post_since4pass_time_no_index ON post(board, time, no) WHERE since4pass IS NOT NULL;

/*
Text search indexes.
The duplication SHOULD be cheap, but I might also be wrong and stupid.
*/
CREATE INDEX post_com_index ON post USING GIN (com_tsvector);
CREATE INDEX post_name_index ON post USING GIN (name_tsvector);
CREATE INDEX post_filename_index ON post USING GIN (filename_tsvector);

CREATE INDEX post_sub_index ON post USING GIN (sub_tsvector);
CREATE INDEX post_op_com_index ON post USING GIN (com_tsvector) WHERE op;
CREATE INDEX post_op_name_index ON post USING GIN (name_tsvector) WHERE op;
CREATE INDEX post_op_filename_index ON post USING GIN (filename_tsvector) WHERE op;

CREATE INDEX post_deleted_com_index ON post USING GIN(com_tsvector) WHERE deleted;
CREATE INDEX post_deleted_name_index ON post USING GIN(name_tsvector) WHERE deleted;
CREATE INDEX post_deleted_filename_index ON post USING GIN(filename_tsvector) WHERE deleted;

CREATE INDEX post_image_com_index ON post USING GIN(com_tsvector) WHERE tim IS NOT NULL;
CREATE INDEX post_image_name_index ON post USING GIN(name_tsvector) WHERE tim IS NOT NULL;
