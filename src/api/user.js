import request from '@/utils/request'

// 获取用户信息表
export function getList() {
  return request({
    url: '/user',
    method: 'get'
  })
}

// 新增用户
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

// 编辑用户信息
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


// 删除用户
export function deleteUser(username) {
  return request({
    url: '/user/delete',
    method: 'post',
    data: {
      username
    }
  })
}
