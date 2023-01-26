# 插入users表
INSERT INTO `tiktok`.`users` (`name`, `password`)
VALUES ('zhangbowen', '77a90868207689664f244ad398a871fc');
INSERT INTO `tiktok`.`users` (`name`, `password`)
VALUES ('zhangbowen1', '9e45c2b95fd89642f4219b28537c652c');
INSERT INTO `tiktok`.`users` (`name`, `password`)
VALUES ('bowenzhang', '77a90868207689664f244ad398a871fc');
INSERT INTO `tiktok`.`users` (`name`, `password`)
VALUES ('lichangyuan', '77a90868207689664f244ad398a871fc');
INSERT INTO `tiktok`.`users` (`name`, `password`)
VALUES ('sunshixin', 'de9ae573b41776f624526219666336d2');
INSERT INTO `tiktok`.`users` (`name`, `password`)
VALUES ('tandonghang', '77a90868207689664f244ad398a871fc');
INSERT INTO `tiktok`.`users` (`name`, `password`)
VALUES ('tuzhuangzhuang', '77a90868207689664f244ad398a871fc');
INSERT INTO `tiktok`.`users` (`name`, `password`)
VALUES ('zhangchangyueyan', '77a90868207689664f244ad398a871fc');

# 插入follows表

INSERT INTO `tiktok`.`follows` (`id`, `user_id`, `follower_id`, `cancel`, `createTime`)
VALUES (1, 1, 2, 0, '2023-01-26 18:56:28');
INSERT INTO `tiktok`.`follows` (`id`, `user_id`, `follower_id`, `cancel`, `createTime`)
VALUES (2, 1, 3, 0, '2023-01-26 18:56:47');
INSERT INTO `tiktok`.`follows` (`id`, `user_id`, `follower_id`, `cancel`, `createTime`)
VALUES (3, 1, 4, 0, '2023-01-26 18:57:02');
INSERT INTO `tiktok`.`follows` (`id`, `user_id`, `follower_id`, `cancel`, `createTime`)
VALUES (4, 1, 5, 0, '2023-01-26 18:57:27');
INSERT INTO `tiktok`.`follows` (`id`, `user_id`, `follower_id`, `cancel`, `createTime`)
VALUES (5, 2, 1, 0, '2023-01-26 18:57:51');
INSERT INTO `tiktok`.`follows` (`id`, `user_id`, `follower_id`, `cancel`, `createTime`)
VALUES (6, 2, 3, 0, '2023-01-26 18:58:01');

