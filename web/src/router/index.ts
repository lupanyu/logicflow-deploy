import { createRouter, createWebHistory } from 'vue-router'
import LF from '@/components/LF.vue'
import ExecutionDetail from '@/components/ExecutionDetail.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: LF,
      meta: {
        title: '流程列表',
      }
    },
    
    {
      path: '/flows',
      name: 'FlowList',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('@/components/FlowList.vue'),
      meta: {
        title: '流程列表',
      }
    },
    {
      path: '/flows/:name',
      name: 'FlowDetail',
      component: LF,
    },
    {
      path: '/status/:id',
      name: 'status',
      component: ExecutionDetail,
      meta: {
        title: '执行详情',
      }
    },
    {
      path: '/newflow',
      name: 'newFlow',
      component: () => import('@/components/NewFlow.vue'),
      meta: {
        title: '新建流程',
      }
    }
  ],
})

export default router
