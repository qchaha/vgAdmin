import axios from 'axios'
import { Message, MessageBox } from 'element-ui'
import store from '../store'
import { getToken } from '@/utils/auth'

// 创建axios实例
const service = axios.create({
  baseURL: process.env.BASE_API, // api 的 base_url
  timeout: 5000 // 请求超时时间
})

// request拦截器
service.interceptors.request.use(
  config => {
    if (store.getters.token) {
      config.headers['X-Token'] = getToken() // 让每个请求携带自定义token 请根据实际情况自行修改
    }
    return config
  },
  error => {
    // Do something with request error
    console.log(error) // for debug
    Promise.reject(error)
  }
)

// response 拦截器
service.interceptors.response.use(
  response => {
    /**
     * code为非20000是抛错 可结合自己业务进行修改
     */
    const res = response.data
    if (res.code !== 20000) {
      Message({
        message: res.message,
        type: 'error',
        duration: 5 * 1000
      })
      // 50001: 身份认证失败
      if (res.code === 50001) {
        Message({
          message: '身份认证失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50002: 保存token失败
      if (res.code === 50002) {
        Message({
          message: '保存token失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50003: 查询用户信息失败
      if (res.code === 50003) {
        Message({
          message: '查询用户信息失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50004: 删除token失败
      if (res.code === 50004) {
        Message({
          message: '删除token失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50005: 该用户已经在其他地方登录，清除数据库token失败
      if (res.code === 50005) {
        Message({
          message: '该用户已经在其他地方登录，清除数据库token失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50006: 添加用户失败
      if (res.code === 50006) {
        Message({
          message: '添加失败...',
          type: 'error'
        })
      }

      // 50007: 该用户已经在其他地方登录，清除本地token失败
      if (res.code === 50007) {
        Message({
          message: '该用户已经在其他地方登录，清除本地token失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50008:非法的token
      if (res.code === 50008) {
        MessageBox.confirm(
          '你已被登出，可以取消继续留在该页面，或者重新登录',
          '确定登出',
          {
            confirmButtonText: '重新登录',
            cancelButtonText: '取消',
            type: 'warning'
          }
        ).then(() => {
          store.dispatch('FedLogOut').then(() => {
            location.reload() // 为了重新实例化vue-router对象 避免bug
          })
        })
      }

      // 50009: 编辑用户信息失败
      if (res.code === 50009) {
        Message({
          message: '编辑用户信息失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50010: 删除用户失败
      if (res.code === 50010) {
        Message({
          message: '删除用户失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50011: 查询用户状态失败
      if (res.code === 50011) {
        Message({
          message: '查询用户状态失败...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      // 50012: 用户被锁
      if (res.code === 50012) {
        Message({
          message: '用户被锁...',
          type: 'error',
          duration: 5 * 1000
        })
      }

      return Promise.reject('error')
    } else {
      return response.data
    }
  },
  error => {
    console.log('err' + error) // for debug
    Message({
      message: error.message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service
