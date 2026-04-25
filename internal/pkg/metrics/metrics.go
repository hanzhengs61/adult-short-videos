package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	UserRegistrationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_registrations_total",
		Help: "用户注册总数",
	})

	UserLoginsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "user_logins_total",
		Help: "用户登录总数",
	}, []string{"success"})

	VideoPlaysTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "video_plays_total",
		Help: "视频播放历史写入成功次数（含已存在记录的更新）",
	})

	FavoritesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "favorites_total",
		Help: "收藏操作总数",
	}, []string{"action"}) // action: add | remove

	SearchRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "search_requests_total",
		Help: "搜索请求总数",
	}, []string{"type"}) // type: video | actor | fanhao | advanced

	CommentsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "comments_total",
		Help: "评论操作总数",
	}, []string{"action"}) // action: add | like | unlike
)
