import { createRouter, createWebHistory } from 'vue-router'
import Auth from '../view/auth/Login.vue'
import Panel from '../view/panel/Panel.vue'
import Sub from '../view/sub/Sub.vue'

const router = createRouter({
  history: createWebHistory(window.__PANEL_PATH__ || '/'),
  routes: [
    {
      path: '/login',
      component: Auth,
      meta: { title: 'SLINX · 登录' }
    },
    {
      path: '/panel',
      component: Panel,
      children: [
        {
          path: '',
          redirect: '/panel/dashboard'
        },
        {
          path: 'dashboard',
          component: () => import('../view/panel/dashboard/Dashboard.vue'),
          meta: { title: 'SLINX · 仪表盘' }
        },
        {
          path: 'inbound',
          component: () => import('../view/panel/inbound/Inbound.vue'),
          meta: { title: 'SLINX · 入站管理' }
        },
        {
          path: 'user',
          component: () => import('../view/panel/user/User.vue'),
          meta: { title: 'SLINX · 用户管理' }
        },
        {
          path: 'board',
          component: () => import('../view/panel/board/Board.vue'),
          meta: { title: 'SLINX · 面板对接' }
        },
        {
          path: 'detect',
          component: () => import('../view/panel/detect/Detect.vue'),
          meta: { title: 'SLINX · IP检测' }
        },
        {
          path: 'core',
          component: () => import('../view/panel/core/Core.vue'),
          meta: { title: 'SLINX · 核心配置' }
        },
        {
          path: 'config',
          component: () => import('../view/panel/config/Config.vue'),
          meta: { title: 'SLINX · 面板配置' }
        },
      ]
    },
    {
      path: '/sub/:token',
      component: Sub,
      meta: { title: 'SLINX · 订阅' }
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/login'
    }
  ]
})

router.beforeEach((to) => {
  document.title = (to.meta.title as string) || 'SLINX'

  const token = localStorage.getItem('token')

  if (token && to.path === '/login') {
    return '/panel'
  }

  if (!token && to.path.startsWith('/panel')) {
    return '/login'
  }
})

export default router