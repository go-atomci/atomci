import Layout from '@/layout'

const scmAppRouter = {
  path: '/scmapp',
  component: Layout,
  redirect: '/scmapp',
  name: 'scmapp',
  meta: {
    title: '我的应用',
    icon: 'chart'
  },
  children: [
    {
      path: '/scmapp',
      name: 'scmappIndex',
      component: () => import('@/views/scmapp/Scmapp.vue'),
      meta: { title: '我的应用', noCache: true }
    },
  ],
}

export default scmAppRouter
