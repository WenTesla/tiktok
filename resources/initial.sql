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

DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes`
(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`    int      DEFAULT NULL COMMENT '点赞用户的id',
    `author_id`  int      DEFAULT NULL COMMENT '视频作者的id',
    `createTime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT '点赞列表';


DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id`    int          DEFAULT NULL COMMENT '用户的id',
    `video_id`   int          DEFAULT NULL COMMENT '视频的id',
    `text`       varchar(255) DEFAULT NULL COMMENT '评论的内容',
    `cancel`     int          DEFAULT NULL COMMENT '是否取消评论',
    `createTime` datetime     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='评论表-用户评论视频';

DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows`
(
    `id`          bigint NOT NULL COMMENT '自增id',
    `user_id`     int      DEFAULT NULL COMMENT '用户的id',
    `follower_id` int      DEFAULT NULL COMMENT '粉丝id',
    `cancel`      tinyint  DEFAULT NULL COMMENT '是否关注',
    `createTime`  datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='用户关注列表';


