DELIMITER ;;

DROP PROCEDURE IF EXISTS `update_thread_po`;;
DROP PROCEDURE IF EXISTS `create_thread_po`;;
DROP PROCEDURE IF EXISTS `delete_thread_po`;;
DROP PROCEDURE IF EXISTS `insert_image_po`;;
DROP PROCEDURE IF EXISTS `delete_image_po`;;
DROP PROCEDURE IF EXISTS `insert_post_po`;;
DROP PROCEDURE IF EXISTS `delete_post_po`;;

DROP TRIGGER IF EXISTS `before_ins_po`;;
DROP TRIGGER IF EXISTS `after_ins_po`;;
DROP TRIGGER IF EXISTS `after_del_po`;;

CREATE PROCEDURE `update_thread_po` (ins INT, tnum INT, subnum INT, timestamp INT, media INT, email VARCHAR(100))
BEGIN
  UPDATE
    `po_threads` op
  SET
    op.time_last = IF((ins AND subnum = 0), GREATEST(timestamp, op.time_last), op.time_last),
    op.time_bump = IF((ins AND subnum = 0), GREATEST(timestamp, op.time_bump), op.time_bump),
    op.time_ghost = IF((ins AND subnum != 0), GREATEST(timestamp, COALESCE(op.time_ghost, 0)), op.time_ghost),
    op.time_ghost_bump = IF((ins AND subnum != 0 AND (email IS NULL OR email != 'sage')), GREATEST(timestamp, COALESCE(op.time_ghost_bump, 0)), op.time_ghost_bump),
    op.time_last_modified = GREATEST(timestamp, op.time_last_modified),
    op.nreplies = IF(ins, (op.nreplies + 1), (op.nreplies - 1)),
    op.nimages = IF(media, IF(ins, (op.nimages + 1), (op.nimages - 1)), op.nimages)
  WHERE op.thread_num = tnum;
END;;

CREATE PROCEDURE `create_thread_po` (num INT, timestamp INT)
BEGIN
  INSERT IGNORE INTO `po_threads` VALUES (num, timestamp, timestamp, timestamp, NULL, NULL, timestamp, 0, 0, 0, 0);
END;;

CREATE PROCEDURE `delete_thread_po` (tnum INT)
BEGIN
  DELETE FROM `po_threads` WHERE thread_num = tnum;
END;;

CREATE PROCEDURE `insert_image_po` (n_media_hash VARCHAR(25), n_media VARCHAR(20), n_preview VARCHAR(20), n_op INT)
BEGIN
  IF n_op = 1 THEN
    INSERT INTO `po_images` (media_hash, media, preview_op, total)
    VALUES (n_media_hash, n_media, n_preview, 1)
    ON DUPLICATE KEY UPDATE
      media_id = LAST_INSERT_ID(media_id),
      total = (total + 1),
      preview_op = COALESCE(preview_op, VALUES(preview_op)),
      media = COALESCE(media, VALUES(media));
  ELSE
    INSERT INTO `po_images` (media_hash, media, preview_reply, total)
    VALUES (n_media_hash, n_media, n_preview, 1)
    ON DUPLICATE KEY UPDATE
      media_id = LAST_INSERT_ID(media_id),
      total = (total + 1),
      preview_reply = COALESCE(preview_reply, VALUES(preview_reply)),
      media = COALESCE(media, VALUES(media));
  END IF;
END;;

CREATE PROCEDURE `delete_image_po` (n_media_id INT)
BEGIN
  UPDATE `po_images` SET total = (total - 1) WHERE media_id = n_media_id;
END;;

CREATE TRIGGER `before_ins_po` BEFORE INSERT ON `po`
FOR EACH ROW
BEGIN
  IF NEW.media_hash IS NOT NULL THEN
    CALL insert_image_po(NEW.media_hash, NEW.media_orig, NEW.preview_orig, NEW.op);
    SET NEW.media_id = LAST_INSERT_ID();
  END IF;
END;;

CREATE TRIGGER `after_ins_po` AFTER INSERT ON `po`
FOR EACH ROW
BEGIN
  IF NEW.op = 1 THEN
    CALL create_thread_po(NEW.num, NEW.timestamp);
  END IF;
  CALL update_thread_po(1, NEW.thread_num, NEW.subnum, NEW.timestamp, NEW.media_id, NEW.email);
END;;

CREATE TRIGGER `after_del_po` AFTER DELETE ON `po`
FOR EACH ROW
BEGIN
  CALL update_thread_po(0, OLD.thread_num, OLD.subnum, OLD.timestamp, OLD.media_id, OLD.email);
  IF OLD.op = 1 THEN
    CALL delete_thread_po(OLD.num);
  END IF;
  IF OLD.media_hash IS NOT NULL THEN
    CALL delete_image_po(OLD.media_id);
  END IF;
END;;

DELIMITER ;
