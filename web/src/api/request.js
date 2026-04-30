import axios from 'axios'

const request = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 15000,
})

request.interceptors.request.use(config => {
    const token = localStorage.getItem('access_token')
    if (token) config.headers.Authorization = `Bearer ${token}`
    return config
})

request.interceptors.response.use(
    res => {
        if (res.data.code !== 200) {
            return Promise.reject(new Error(res.data.msg || '请求失败'))
        }
        return res.data
    },
    err => {
        if (err.response?.status === 401) {
            localStorage.removeItem('access_token')
            window.location.href = '/login'
        }
        return Promise.reject(new Error(err.response?.data?.msg || err.message || '请求失败'))
    }
)

export default request
