import Layout from '@/layout'

const projectRouter = {
  path: '/project',
  component: Layout,
  redirect: '/project',
  name: 'projects',
  meta: {
    title: '我的项目',
    icon: 'chart'
  },
  children: [
    {
      path: '/project',
      name: 'projectIndex',
      component: () => import('@/views/project/Project.vue'),
      meta: { title: '我的项目', noCache: true }
    },
  ],
}

export default projectRouter

