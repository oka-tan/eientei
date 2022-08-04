CREATE TABLE post (
    board TEXT NOT NULL,
    post_number BIGINT NOT NULL CHECK(post_number > 0),
    thread_number BIGINT NOT NULL CHECK(thread_number > 0),
    op BOOLEAN NOT NULL CHECK(op = (post_number = thread_number)),
    deleted BOOLEAN NOT NULL,
    hidden BOOLEAN NOT NULL,
    time_posted TIMESTAMP WITH TIME ZONE NOT NULL,
    last_modified TIMESTAMP WITH TIME ZONE NOT NULL CHECK(last_modified >= created_at),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL CHECK(created_at >= time_posted),
    name TEXT CHECK(name != '' AND name != 'Anonymous'),
    tripcode TEXT CHECK(tripcode != ''),
    capcode TEXT CHECK(capcode != ''),
    poster_id TEXT CHECK(poster_id != ''),
    country TEXT CHECK(country != ''),
    flag TEXT CHECK(flag != ''),
    email TEXT CHECK(email != ''),
    subject TEXT CHECK(subject != ''),
    comment TEXT CHECK(comment != ''),
    has_media BOOLEAN NOT NULL,
    media_deleted BOOLEAN,
    time_media_deleted TIMESTAMP WITH TIME ZONE CHECK(time_media_deleted >= created_at AND last_modified >= time_media_deleted),
    media_timestamp BIGINT CHECK(media_timestamp > 0),
    media_4chan_hash BYTEA CHECK(OCTET_LENGTH(media_4chan_hash) > 0),
    media_internal_hash BYTEA CHECK(OCTET_LENGTH(media_internal_hash) > 0),
    thumbnail_internal_hash BYTEA CHECK(OCTET_LENGTH(thumbnail_internal_hash) > 0),
    media_extension TEXT CHECK(media_extension != ''),
    media_file_name TEXT CHECK(media_file_name != ''),
    media_size INTEGER CHECK(media_size > 0),
    media_height SMALLINT CHECK(media_height > 0),
    media_width SMALLINT CHECK(media_width > 0),
    thumbnail_height SMALLINT CHECK(thumbnail_height > 0),
    thumbnail_width SMALLINT CHECK(thumbnail_width > 0),
    spoiler BOOLEAN,
    custom_spoiler SMALLINT CHECK(custom_spoiler > 0),
    sticky BOOLEAN,
    closed BOOLEAN,
    posters SMALLINT CHECK(posters > 0),
    replies SMALLINT CHECK(replies >= 0),
    since4pass SMALLINT CHECK(since4pass > 2011),
    oekaki_internal_hash BYTEA CHECK(OCTET_LENGTH(oekaki_internal_hash) > 0),
    
    CONSTRAINT replies_need_to_come_after_the_op CHECK(post_number >= thread_number),
    CONSTRAINT only_ops_have_subjects CHECK(op OR subject IS NULL),
    CONSTRAINT either_a_post_has_a_country_or_a_flag CHECK(country IS NULL OR flag IS NULL),
    CONSTRAINT only_posts_with_media_have_deleted_media CHECK(has_media OR media_deleted IS NULL),
    CONSTRAINT only_posts_with_media_have_media_timestamps CHECK(has_media OR media_timestamp IS NULL),
    CONSTRAINT only_posts_with_media_have_4chan_hashes CHECK(has_media OR media_4chan_hash IS NULL),
    CONSTRAINT only_posts_with_media_have_internal_hashes CHECK(has_media OR media_internal_hash IS NULL),
    CONSTRAINT only_posts_with_media_have_thumbnail_hashes CHECK(has_media OR thumbnail_internal_hash IS NULL),
    CONSTRAINT only_posts_with_media_have_extensions CHECK(has_media OR media_extension IS NULL),
    CONSTRAINT media_extensions_cannot_contain_dots CHECK(media_extension NOT LIKE '%.%'),
    CONSTRAINT only_posts_with_media_have_file_names CHECK(has_media OR media_file_name IS NULL),
    CONSTRAINT only_posts_with_media_have_media_sizes CHECK(has_media OR media_size IS NULL),
    CONSTRAINT only_posts_with_media_have_media_dimensions CHECK(has_media OR (media_height IS NULL AND media_width IS NULL)),
    CONSTRAINT either_both_media_dimensions_are_present_or_neither_is CHECK((media_height IS NULL) = (media_width IS NULL)),
    CONSTRAINT only_posts_with_media_have_thumbnail_dimensions CHECK(has_media OR (thumbnail_height IS NULL and thumbnail_width IS NULL)),
    CONSTRAINT either_both_thumbnail_dimensions_are_present_or_neither_is CHECK((thumbnail_height IS NULL) = (thumbnail_width IS NULL)),
    CONSTRAINT only_posts_with_media_have_spoilers CHECK(has_media OR spoiler IS NULL),
    CONSTRAINT only_posts_with_spoilers_have_custom_spoilers CHECK(spoiler IS TRUE OR custom_spoiler IS NULL),
    CONSTRAINT either_op_posts_are_stickies_or_they_arent CHECK(op = (sticky IS NOT NULL)),
    CONSTRAINT either_op_posts_are_closed_or_they_arent CHECK(op = (closed IS NOT NULL)),
    CONSTRAINT only_op_posts_have_posters CHECK(op OR posters IS NULL),
    CONSTRAINT only_op_posts_have_replies CHECK(op OR replies IS NULL),

    PRIMARY KEY(board, post_number)
) PARTITION BY LIST(board);

ALTER TABLE post ALTER COLUMN board SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN name SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN tripcode SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN capcode SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN poster_id SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN country SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN flag SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN email SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN subject SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN comment SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN media_4chan_hash SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN media_internal_hash SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN thumbnail_internal_hash SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN media_extension SET STORAGE EXTERNAL;
ALTER TABLE post ALTER COLUMN media_file_name SET STORAGE EXTERNAL;

CREATE INDEX post_thread_number_post_number_index ON post(thread_number, post_number);
CREATE INDEX post_op_post_number_index ON post(post_number) WHERE op;
CREATE INDEX post_last_modified_index ON post(last_modified, post_number);

CREATE TABLE media (hash BYTEA PRIMARY KEY);
ALTER TABLE media ALTER COLUMN hash SET STORAGE EXTERNAL;

CREATE TABLE index_tracker(
    board TEXT PRIMARY KEY,
    last_modified TIMESTAMP WITH TIME ZONE NOT NULL,
    post_number BIGINT NOT NULL
);
