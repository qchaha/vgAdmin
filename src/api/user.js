import request from '@/utils/request'

export function getList() {
  return request({
    url: '/user',
    method: 'get'
  })
}

export function newUser(username, password, role, status, email) {
  return request({
    url: '/user/new',
    method: 'post',
    data: {
      username,
      password,
      role,
      status,
      email
    }
  })
}

export function editUser(username, password, role, status, email, origin_username) {
  return request({
    url: '/user/edit',
    method: 'post',
    data: {
      username,
      password,
      role,
      status,
      email,
      origin_username
    }
  })
}

export function deleteUser(username) {
  return request({
    url: '/user/delete',
    method: 'post',
    data: {
      username
    }
  })
}
