-- create
-- database tiktok;

use
tiktok;
-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`       bigint       NOT NULL AUTO_INCREMENT COMMENT '用户id，自增主键',
    `name`     varchar(255) NOT NULL COMMENT '用户名',
    `password` varchar(255) NOT NULL COMMENT '用户密码',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COMMENT='用户表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`
(
    `id`           bigint       NOT NULL AUTO_INCREMENT COMMENT '自增主键，视频唯一id',
    `author_id`    bigint       NOT NULL COMMENT '视频作者id',
    `play_url`     varchar(255) NOT NULL COMMENT '播放url',
    `cover_url`    varchar(255) NOT NULL COMMENT '封面的url',
    `title`        varchar(255)          DEFAULT NULL COMMENT '视频名称',
    `publish_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间戳',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COMMENT='\r\n视频表';

DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes`
(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`    int      DEFAULT NULL COMMENT '点赞用户的id',
    `video_id`   int      DEFAULT NULL COMMENT '视频作者的id',
    `is_cancel`  tinyint  DEFAULT '0' COMMENT '是否点赞 0-点赞 1-未点赞',
    `createTime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='点赞列表';


DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id`    int          DEFAULT NULL COMMENT '用户的id',
    `video_id`   int          DEFAULT NULL COMMENT '视频的id',
    `text`       varchar(255) DEFAULT NULL COMMENT '评论的内容',
    `is_cancel`  int          DEFAULT '0' COMMENT '是否取消评论',
    `createTime` datetime     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评论表-用户评论视频';

DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows`
(
    `id`          bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id`     int      DEFAULT NULL COMMENT '用户的id',
    `follower_id` int      DEFAULT NULL COMMENT '关注的用户',
    `cancel`      tinyint  DEFAULT '0' COMMENT '是否关注',
    `createTime`  datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户关注列表';


