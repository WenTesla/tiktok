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
    PRIMARY KEY (`id`),
    KEY        `name-password` (`name`,`password`) USING BTREE COMMENT '用户名-密码索引\r\n'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb3 COMMENT='用户表';

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
    PRIMARY KEY (`id`) USING BTREE,
    KEY            `publish_time` (`publish_time`) COMMENT '发布时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8mb3 COMMENT='\r\n视频表';

DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes`
(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`    int      DEFAULT NULL COMMENT '点赞用户的id',
    `video_id`   int      DEFAULT NULL COMMENT '视频作者的id',
    `is_cancel`  tinyint  DEFAULT '0' COMMENT '是否点赞 0-点赞 1-未点赞',
    `createTime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY          `user_id` (`user_id`),
    KEY          `video_id` (`video_id`),
    KEY          `user_id_video_id` (`user_id`,`video_id`)
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='点赞列表';


DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id`    int          DEFAULT NULL COMMENT '用户的id',
    `video_id`   int          DEFAULT NULL COMMENT '视频的id',
    `text`       varchar(255) DEFAULT NULL COMMENT '评论的内容',
    `is_cancel`  int          DEFAULT '0' COMMENT '是否取消评论',
    `createTime` datetime     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY          `user_id` (`user_id`),
    KEY          `video_id` (`video_id`),
    KEY          `user_id_video_id` (`user_id`,`video_id`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评论表-用户评论视频';

DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows`
(
    `id`          bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id`     int      DEFAULT NULL COMMENT '用户的id',
    `follower_id` int      DEFAULT NULL COMMENT '关注的用户',
    `cancel`      tinyint  DEFAULT '0' COMMENT '是否关注',
    `createTime`  datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY           `user_id` (`user_id`),
    KEY           `follower_id` (`follower_id`),
    KEY           `user_id_follower_id` (`user_id`,`follower_id`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户关注列表';

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id`     bigint       NOT NULL COMMENT '用户的Id',
    `to_user_id`  bigint       NOT NULL COMMENT '接受消息的用户Id',
    `content`     varchar(256) NOT NULL COMMENT '消息内容',
    `is_withdraw` tinyint DEFAULT '0' COMMENT '是否撤回  0-不撤回，1-撤回',
    `createTime`  bigint  DEFAULT NULL COMMENT '时间戳',
    PRIMARY KEY (`id`),
    KEY           `用户索引` (`user_id`) USING BTREE COMMENT '发送信息的用户Id索引',
    KEY           `接受信息的用户索引` (`to_user_id`) USING BTREE COMMENT '接受用户的用户Id索引'
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='消息表';
