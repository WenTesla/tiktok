package config

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type Comment struct {
	ID         int64     // 评论id
	UserId     int64     // 用户Id
	VideoId    int64     //视频Id
	Text       string    // 评论内容
	IsCancel   int64     // 是否取消
	CreateTime time.Time `gorm:"column:createTime"` // 创建时间
}

// CommentInfo 评论详细信息
type CommentInfo struct {
	Content    string   `json:"content"`     // 评论内容
	CreateDate string   `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64    `json:"id"`          // 评论id
	UserInfo   UserInfo `json:"user"`        // 评论用户的具体信息
}

// UserInfo  最详细的信息 不与数据库模型对应
type UserInfo struct {
	Id             int64  `json:"id,omitempty"`              //主键
	Name           string `json:"name,omitempty"`            //昵称
	FollowCount    int64  `json:"follow_count"`              //关注总数
	FollowerCount  int64  `json:"follower_count"`            //粉丝总数
	IsFollow       bool   `json:"is_follow"`                 //是否关注
	AvatarUrl      string `json:"avatar,omitempty"`          //用户的url
	TotalFavorited int64  `json:"total_favorited,omitempty"` //获赞数量
	WorkCount      int64  `json:"work_count,omitempty"`      //作品数量
	FavoriteCount  int64  `json:"favorite_count,omitempty"`  //点赞数量
}

func TestRedis(t *testing.T) {
	ctx := context.Background()
	_, err := InitRedisClient()
	if err != nil {
		fmt.Println(err)
	}
	_ = ctx
}

type guo struct {
	Name string
	Age  int
}

func TestAddRedis(t *testing.T) {
	ctx := context.Background()
	db, _ := InitRedisClient()
	comment := CommentInfo{
		Content:    "111",
		CreateDate: time.Now().String(),
		ID:         1,
		UserInfo: UserInfo{
			Id:             2,
			Name:           "22",
			FollowCount:    3,
			FollowerCount:  4,
			IsFollow:       false,
			AvatarUrl:      "132",
			TotalFavorited: 1,
			WorkCount:      2,
			FavoriteCount:  3,
		},
	}
	bytes, err := json.Marshal(comment)
	if err != nil {
		fmt.Println(err)
	}
	db.LPush(ctx, "commentInfo", bytes)
}

func TestQueryRedis(t *testing.T) {
	ctx := context.Background()
	db, _ := InitRedisClient()
	push := db.LRange(ctx, "commentInfo", 0, -1)
	result, err := push.Result()
	if err != nil {
		fmt.Println(err)
	}
	s := result[3]
	bytes := []byte(s)
	g := CommentInfo{}
	json.Unmarshal(bytes, &g)
	fmt.Printf("%v", g)
}

func TestTimeRedis(t *testing.T) {
	ctx := context.Background()
	db := RedisClient(10)
	db.Incr(ctx, "111")
}

func TestIs(t *testing.T) {
	ctx := context.Background()
	client := RedisClient(1)
	client.SAdd(ctx, "1", "2", "3", "4")
	member, _ := client.SIsMember(ctx, "1", "1").Result()
	fmt.Println(member)
}

func TestCount(t *testing.T) {
	ctx := context.Background()
	client := RedisClient(4)
	card, _ := client.SCard(ctx, "423").Result()
	fmt.Printf("%T", card)
}
