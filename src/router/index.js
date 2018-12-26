import Vue from 'vue'
import Router from 'vue-router'

// in development-env not use lazy-loading, because lazy-loading too many pages will cause webpack hot update too slow. so only in production use lazy-loading;
// detail: https://panjiachen.github.io/vue-element-admin-site/#/lazy-loading

Vue.use(Router)

/* Layout */
import Layout from '../views/layout/Layout'

/**
* hidden: true                   if `hidden:true` will not show in the sidebar(default is false)
* alwaysShow: true               if set true, will always show the root menu, whatever its child routes length
*                                if not set alwaysShow, only more than one route under the children
*                                it will becomes nested mode, otherwise not show the root menu
* redirect: noredirect           if `redirect:noredirect` will no redirect in the breadcrumb
* name:'router-name'             the name is used by <keep-alive> (must set!!!)
* meta : {
    title: 'title'               the name show in submenu and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar,
  }
**/
export const constantRouterMap = [
  { path: '/login', component: () => import('@/views/login/index'), hidden: true },
  { path: '/404', component: () => import('@/views/404'), hidden: true },

  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    name: 'Dashboard',
    hidden: true,
    children: [{
      path: 'dashboard',
      component: () => import('@/views/dashboard/index')
    }]
  }
]

// 前端根据用户角色显示不同的菜单，通过meta字段的roles来控制
export const asyncRouterMap = [
  {
    path: '/dashboard',
    component: Layout,
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        meta: { title: '仪表板', icon: 'dashboard' }
      }
    ]
  },

  {
    path: '/user',
    component: Layout,
    meta: { icon: 'user', roles: ['super'] },
    children: [
      {
        meta: { title: '人员管理', icon: 'user', roles: ['super'] },
        path: 'user',
        name: 'user',
        component: () => import('@/views/user/index/')
      }
    ]
  },

  {
    path: '/setting',
    component: Layout,
    meta: { title: '系统配置', icon: 'setting', roles: ['super'] },
    children: [
      {
        path: 'use1r',
        name: 'us11er',
        component: () => import('@/views/user/index/'),
        meta: { title: '333', icon: 'user', roles: ['super'] }
      },
      {
        path: 'r11rr',
        name: 'r11rr',
        component: () => import('@/views/setting/index'),
        meta: { title: 'rrr', icon: 'setting', roles: ['super'] }
      }
    ]
  },

  { path: '*', redirect: '/404', hidden: true }
]

export default new Router({
  // mode: 'history', //后端支持可开
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRouterMap
})
