import Layout from '@/layout'

const scmAppRouter = {
  path: '/scmapps',
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
    {
      path: '/scmapp/register',
      name: 'addScmApp',
      meta: { title: '新增应用', noCache: true },
      component: () => import('@/views/scmapp/detail/ScmAppAdd.vue'),
      hidden: true
    },
    {
      path: '/scmapp/:appId',
      name: 'scmAppDetail',
      meta: { title: '代码仓库详情', noCache: true },
      component: () => import('@/views/scmapp/detail/AppDetail.vue'),
      hidden: true
    },
  ],
}

export default scmAppRouter
