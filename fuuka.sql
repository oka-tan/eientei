INSERT INTO post(
	board,
	no,
	resto,
	time,
	name,
	trip,
	capcode,
	country,
	since4pass,
	sub,
	com,
	tim,
	md5,
	filename,
	ext,
	fsize,
	w,
	h,
	tn_w,
	tn_h,
	deleted,
	file_deleted,
	spoiler,
	custom_spoiler,
	op,
	sticky
)
SELECT
	'3',
	num,
	CASE parent
		WHEN 0 THEN num
		ELSE parent
	END,
	TO_TIMESTAMP(CAST(timestamp as DOUBLE PRECISION)) + INTERVAL '4 hours',
	CASE name
		WHEN 'Anonymous' THEN NULL
		ELSE name
	END,
	CASE trip
		WHEN '' THEN NULL
		ELSE trip
	END,
	CASE capcode
		WHEN 'N' THEN NULL
		WHEN 'M' THEN CAST('mod' AS CAPCODE)
	END,
	null,
	null,
	CASE title
		WHEN '' THEN NULL
		ELSE title
	END,
	CASE comment
		WHEN '' THEN NULL
		ELSE REPLACE(REGEXP_REPLACE(REGEXP_REPLACE(REPLACE(REPLACE(REPLACE(comment, '&', '&amp;'), '>', '&gt;'), '<', '&lt;'), '&gt;&gt;(\d+)', '<a class="quotelink" href="#p\1">&gt;&gt;\1</a>', 'gn'), '^&gt;(.*)', '<span class="quote">&gt;\1</span>', 'gn'), E'\n', '<br>')
	END,
	CASE preview IS NULL
		WHEN TRUE THEN NULL
		ELSE CAST(SUBSTRING(preview, 0, LENGTH(preview) - 4) AS BIGINT)
	END,
	media_hash,
	CASE media IS NULL
		WHEN TRUE THEN NULL
		ELSE SUBSTRING(media, 0, LAST_POSITION(media, '.'))
	END,
	CASE media IS NULL
		WHEN TRUE THEN NULL
		ELSE CAST(SUBSTRING(media, LAST_POSITION(media, '.')) AS EXT)
	END,
	CASE media_size
		WHEN 0 THEN NULL
		ELSE media_size
	END,
	CASE media_w
		WHEN 0 THEN NULL
		ELSE media_w
	END,
	CASE media_h
		WHEN 0 THEN NULL
		ELSE media_h
	END,
	CASE preview_w
		WHEN 0 THEN NULL
		ELSE preview_w
	END,
	CASE preview_h
		WHEN 0 THEN NULL
		ELSE preview_h
	END,
	deleted,
	false,
	spoiler,
	null,
	parent IN (0, num),
	sticky
FROM "3"
WHERE subnum = 0;
