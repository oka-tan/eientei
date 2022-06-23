/*
Basic search indexes.
*/
DROP INDEX post_resto_no_index;
DROP INDEX post_deleted_no_index;
DROP INDEX post_spoiler_no_index;
DROP INDEX post_op_no_index;
DROP INDEX post_image_no_index;
DROP INDEX post_capcode_no_index;
DROP INDEX post_since4pass_no_index;
DROP INDEX post_md5_no_index;
DROP INDEX post_trip_no_index;
DROP INDEX post_country_no_index;

/*
Search indexes by time.
*/
DROP INDEX post_time_no_index;
DROP INDEX post_deleted_time_no_index;
DROP INDEX post_spoiler_time_no_index;
DROP INDEX post_op_time_no_index;
DROP INDEX post_image_time_no_index;
DROP INDEX post_capcode_time_no_index;
DROP INDEX post_since4pass_time_no_index;

/*
Text search indexes.
The duplication SHOULD be cheap, but I might also be wrong and stupid.
*/
DROP INDEX post_com_index;
DROP INDEX post_name_index;
DROP INDEX post_filename_index;

DROP INDEX post_sub_index;
DROP INDEX post_op_com_index;
DROP INDEX post_op_name_index;
DROP INDEX post_op_filename_index;

DROP INDEX post_deleted_com_index;
DROP INDEX post_deleted_name_index;
DROP INDEX post_deleted_filename_index;

DROP INDEX post_image_com_index;
DROP INDEX post_image_name_index;
