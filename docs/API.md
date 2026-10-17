# 接口文档（API）

基址：`http://<host>:<port>`，所有业务接口挂在 `/douyin` 分组下（除 `/healthz`）。统一响应包裹 `status_code` / `status_msg`。

鉴权方式：受保护接口需在**查询参数**带 `token`（GET）或**表单字段**带 `token`（POST）；服务端校验后写入上下文 `Id`。

## 基础
| 方法 | 路径 | 鉴权 | 说明 |
| --- | --- | --- | --- |
| GET | `/healthz` | 否 | 健康检查，返回 `{"status":"ok"}` |

## 用户
| 方法 | 路径 | 鉴权 | 参数 | 说明 |
| --- | --- | --- | --- | --- |
| POST | `/douyin/user/register/` | 否 | `username`, `password` | 注册（密码 bcrypt 哈希入库） |
| POST | `/douyin/user/login/` | 否 | `username`, `password` | 登录，返回 `token` |
| GET | `/douyin/user/` | 是 | `token`, `user_id` | 用户信息（关注数/粉丝数/作品数/获赞数/是否关注） |

## 视频
| 方法 | 路径 | 鉴权 | 参数 | 说明 |
| --- | --- | --- | --- | --- |
| GET | `/douyin/feed/` | 否* | `latest_time`, `token`(可选) | 视频流；登录态回填 `is_favorite` / 作者 `is_follow` |
| POST | `/douyin/publish/action/` | 是 | `token`, `data`(视频文件), `title` | 发布（仅 mp4，≤100MB；依赖 COS） |
| GET | `/douyin/publish/list/` | 否 | `user_id`, `token`(可选) | 某用户发布列表 |

> *未带 token 时按 IP 限速；带 token 时按用户限速。

## 互动
| 方法 | 路径 | 鉴权 | 参数 | 说明 |
| --- | --- | --- | --- | --- |
| POST | `/douyin/favorite/action/` | 是 | `token`, `video_id`, `action_type`(0赞/1取消) | 点赞 / 取消 |
| GET | `/douyin/favorite/list/` | 否 | `user_id`, `token`(可选) | 喜欢列表 |
| POST | `/douyin/comment/action/` | 是 | `token`, `video_id`, `action_type`(1发/2删), `comment_text`/`comment_id` | 评论 |
| GET | `/douyin/comment/list/` | 否 | `video_id`, `token`(可选) | 评论列表 |

## 关系
| 方法 | 路径 | 鉴权 | 参数 | 说明 |
| --- | --- | --- | --- | --- |
| POST | `/douyin/relation/action/` | 是 | `token`, `to_user_id`, `action_type`(1关注/2取关) | 关注 / 取关 |
| GET | `/douyin/relation/follow/list/` | 是 | `token`, `user_id` | 关注列表 |
| GET | `/douyin/relation/follower/list/` | 是 | `token`, `user_id` | 粉丝列表 |
| GET | `/douyin/relation/friend/list/` | 是 | `token`, `user_id` | 好友（互关）列表 |

## 消息
| 方法 | 路径 | 鉴权 | 参数 | 说明 |
| --- | --- | --- | --- | --- |
| GET | `/douyin/message/chat/` | 是 | `token`, `to_user_id`, `pre_msg_time` | 与某用户的聊天记录 |
| POST | `/douyin/message/action/` | 是 | `token`, `to_user_id`, `content` | 发送消息 |

## 调用示例
```bash
# 注册
curl -F "username=alice" -F "password=secret123" http://localhost:8000/douyin/user/register/
# 登录拿到 token
TOKEN=$(curl -s -F "username=alice" -F "password=secret123" http://localhost:8000/douyin/user/login/ | python -c "import sys,json;print(json.load(sys.stdin)['token'])")
# 带 token 拉视频流
curl "http://localhost:8000/douyin/feed/?token=$TOKEN"
```
