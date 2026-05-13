import request from './request'

const base = '/manage'

export const manageApi = {
    dashboard: () => request.get(`${base}/dashboard`),

    // 视频
    videoList:         p  => request.get(`${base}/videos`, { params: p }),
    videoStatus:       (id, status) => request.put(`${base}/videos/${id}/status`, { status }),
    videoDelete:       id => request.delete(`${base}/videos/${id}`),

    // 评论
    commentList:       p  => request.get(`${base}/comments`, { params: p }),
    commentDelete:     id => request.delete(`${base}/comments/${id}`),
    commentBatchDel:   ids => request.delete(`${base}/comments`, { data: { ids } }),

    // 用户
    userList:          p  => request.get(`${base}/users`, { params: p }),
    userUpdate:        (id, data) => request.put(`${base}/users/${id}`, data),

    // 广告 Banner
    bannerList:        () => request.get(`${base}/banners`),
    bannerCreate:      d  => request.post(`${base}/banners`, d),
    bannerUpdate:      (id, d) => request.put(`${base}/banners/${id}`, d),
    bannerDelete:      id => request.delete(`${base}/banners/${id}`),

    // 广告 App
    appList:           () => request.get(`${base}/apps`),
    appCreate:         d  => request.post(`${base}/apps`, d),
    appUpdate:         (id, d) => request.put(`${base}/apps/${id}`, d),
    appDelete:         id => request.delete(`${base}/apps/${id}`),

    // 吃瓜文章
    gossipList:        p  => request.get(`${base}/gossip`, { params: p }),
    gossipCreate:      d  => request.post(`${base}/gossip`, d),
    gossipUpdate:      (id, d) => request.put(`${base}/gossip/${id}`, d),
    gossipDelete:      id => request.delete(`${base}/gossip/${id}`),
}
