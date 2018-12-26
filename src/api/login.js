// 所有的请求都引用request，所以请求头信息可以封装在request
import request from '@/utils/request'

export function login(username, password) {
  return request({
    url: '/user/login',
    method: 'post',
    data: {
      username,
      password
    }
  })
}

// http://127.0.0.1:8000/user/info?token=5c232c0b75f1c2308e7f8ef2ee26c6f5%0A
export function getInfo(token) {
  return request({
    url: '/user/info',
    method: 'get',
    params: { token }
  })
}

export function logout(token) {
  return request({
    url: '/user/logout',
    method: 'post',
    data: {
      token
    }
  })
}
