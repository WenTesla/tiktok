-- create
-- database tiktok;


-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`       bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '用户id，自增主键',
    `name`     varchar(255) NOT NULL COMMENT '用户名',
    `password` varchar(255) NOT NULL COMMENT '用户密码',
    PRIMARY KEY (`id`)
--     KEY        `name_password_idx` (`name`,`password`) USING BTREE
) ENGINE = InnoDB
--   AUTO_INCREMENT = 20044
  DEFAULT CHARSET = utf8 COMMENT ='用户表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`
(
    `id`           bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '自增主键，视频唯一id',
    `author_id`    bigint(20)   NOT NULL COMMENT '视频作者id',
    `play_url`     varchar(255) NOT NULL COMMENT '播放url',
    `cover_url`    varchar(255) NOT NULL COMMENT '封面的url',
    `publish_time` datetime     NOT NULL COMMENT '发布时间戳',
    `title`        varchar(255) DEFAULT NULL COMMENT '视频名称',
    PRIMARY KEY (`id`)
--     KEY `time` (`publish_time`) USING BTREE,
--     KEY `author` (`author_id`) USING BTREE
) ENGINE = InnoDB
--   AUTO_INCREMENT = 115
  DEFAULT CHARSET = utf8 COMMENT ='\r\n视频表';


