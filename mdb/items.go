package mdb

type Account struct {
	Id_      string `bson:"_id"`
	Password string `bson:"password"`
	Cookie   string `bson:"cookie"`
	Status   string `bson:"status"`
}

type Tweets struct {
	Id_             string `bson:"_id"`
	WeiboUrl        string `bson:"weibo_url"`
	CreatedAt       string `bson:"created_at"`
	LikeNum         int32  `bson:"like_num"`
	RepostNum       int32  `bson:"repost_num"`
	CommentNum      int32  `bson:"comment_num"`
	Content         string `bson:"content"`
	UserId          string `bson:"user_id"`
	Tool            string `bson:"tool"`
	ImageUrl        string `bson:"image_url"`
	VideoUrl        string `bson:"video_url"`
	OriginWeibo     string `bson:"origin_weibo"`
	LocationMapInfo string `bson:"location_map_info"`
	CrawlTime       int32  `bson:"crawl_time"`
}

type Information struct {
	Id_               string `bson:"_id"`
	Nickname          string `bson:"nick_name"`
	Gender            string `bson:"gender"`
	Province          string `bson:"province"`
	City              string `bson:"city"`
	BriefIntroduction string `bson:"brief_introduction"`
	Birthday          string `bson:"birthday"`
	TweetsNum         int32  `bson:"tweets_num"`
	FollowsNum        int32  `bson:"follows_num"`
	FansNum           int32  `bson:"fans_num"`
	SexOrientation    string `bson:"sex_orientation"`
	Sentiment         string `bson:"sentiment"`
	VipLevel          string `bson:"vip_level"`
	Authentication    string `bson:"authentication"`
	Labels            string `bson:"labels"`
	CrawlTime         int32  `bson:"crawl_time"`
}

type Relationships struct {
	Id_        string `bson:"_id"`
	FanId      string `bson:"fan_id"`
	FollowedId string `bson:"followed_id"`
	CrawlTime  int32  `bson:"crawl_time"`
}

type Comment struct {
	Id_           string `bson:"_id"`
	CommentUserId string `bson:"comment_user_id"`
	Content       string `bson:"content"`
	WeiboUrl      string `bson:"weibo_url"`
	CreatedAt     string `bson:"created_at"`
	LikeNum       int32 `bson:"like_num"`
	CrawlTime     int32 `bson:"crawl_time"`
}
