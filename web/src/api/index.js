import request from './request'

// 用户
export const userApi = {
    register: data => request.post('/user/register', data),
    login: data => request.post('/user/login', data),
    logout: () => request.post('/user/logout'),
    info: () => request.get('/user/info'),
    updateProfile: data => request.put('/user/profile', data),
    uploadAvatar: file => {
      const fd = new FormData()
      fd.append('file', file)
      return request.post('/user/avatar', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
    },
    creators: (limit = 20) => request.get('/user/creators', { params: { limit } }),
}

// 视频
export const videoApi = {
    list: params => request.get('/video/list', {params}),
    detail: id => request.get(`/video/detail/${id}`),
    popular: params => request.get('/video/popular', {params}),
}

// 搜索
export const searchApi = {
    videos: params => request.get('/search/videos', {params}),
    advanced: params => request.get('/search/advanced', {params}),
    byFanhao: fanhao => request.get(`/search/fanhao/${fanhao}`),
}

// 收藏
export const favoriteApi = {
    add: videoId => request.post('/favorite/add', {video_id: videoId}),
    remove: videoId => request.delete('/favorite/remove', {data: {video_id: videoId}}),
    list: params => request.get('/favorite/list', {params}),
    check: videoId => request.get('/favorite/check', {params: {video_id: videoId}}),
}

// 播放历史
export const playApi = {
    record: videoId => request.post('/play/record', {video_id: videoId}),
    history: params => request.get('/play/history', {params}),
}

// 关注
export const followApi = {
    follow: author => request.post('/follow', { author }),
    unfollow: author => request.delete('/follow', { data: { author } }),
    check: author => request.get('/follow/check', { params: { author } }),
    list: () => request.get('/follow/list'),
    feed: params => request.get('/follow/feed', { params }),
}

// 评论
export const commentApi = {
    list: params => request.get('/comment/list', {params}),
    add: data => request.post('/comment/add', data),
    like: commentId => request.post('/comment/like', {comment_id: commentId}),
    unlike: commentId => request.delete(`/comment/like/${commentId}`),
    delete: id => request.delete(`/comment/delete/${id}`),
}
