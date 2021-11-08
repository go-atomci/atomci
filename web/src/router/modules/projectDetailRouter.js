import Layout from '@/layout'

export const projectDetailRouter = [
  { 
    path: '/project/detail',
    name: 'projectMenu',
    component: Layout,        
    meta: { title: '项目概览', noCache: true},
    redirect: '/project/detail/deashbord',
    children: [
        {
            path: '/project/detail/deashbord',
            name: 'projectSummary',
            component: () => import('@/views/project/ProjectDashboard.vue'),
            meta: { title: '项目概览', noCache: true},
        }
    ]
},
{
    path: '/project/detail/app',
    name: 'projectApp',
    component: Layout,
    meta: { title: '应用代码', noCache: true},
    children: [
        {
            path: '/project/detail/app',
            name: 'projectApp',
            meta: { title: '应用代码', noCache: true},
            iconCls: 'app',
            component: () => import('@/views/project/ProjectApp.vue'),
        }
    ]
},
{
    path: '/project/detail/ci',
    name: 'projectCI',
    component: Layout,
    meta: { title: '构建部署', noCache: true},
    children: [
        {
            path: '/project/detail/ci',
            name: 'projectCI',
            meta: { title: '构建部署', noCache: true},
            component: () => import('@/views/project/ProjectCICD.vue'),
        }
    ]
},
{
    path: '/project/detail/service',
    name: 'projectService',
    component: Layout,
    meta: { title: '应用服务', noCache: true},
    children: [
      {
        path: '/project/detail/service',
        meta: { title: '应用服务', noCache: true},
        name: 'projectService',
        iconCls: 'service',
        component: () => import('@/views/project/Service.vue'),
      },
      { 
        path: '/project/service/:clusterName/:namespace/:appName', 
        meta: { title: '应用详情', noCache: true }, 
        name: 'projectServiceDetail', 
        component: () => import('@/views/project/detail/ServiceDetail.vue'), 
        hidden: true 
    },
    ]
  },
  // {
  //   path: '/project/detail/projectStats',
  //   name: 'projectStatistics',
  //   component: Layout,
  //   meta: { title: '发布统计', noCache: true},
  //   children: [
  //     {
  //       path: '/project/detail/projectStats',
  //       name: 'projectStatistics',
  //       meta: { title: '发布统计', noCache: true },
  //       component: () => import('@/views/project/CICDStats.vue'),
  //     }
  //   ]
  // },
  {
    path: '/project/detail/projectSets',
    name: 'projectSetting',
    component: Layout,
    meta: { title: '项目设置', noCache: true },
    redirect: '/project/detail/projectInfo',
    children: [
      // { 
      //   path: '/project/detail/projectInfo', 
      //   name: 'projectInfo', 
      //   meta: { title: '成员管理', noCache: true },
      //   component: () => import('@/views/project/ProjectInfo.vue')
      // },
      { 
        path: '/project/detail/projectEnv', 
        name: 'projectEnv', 
        meta: { title: '项目环境', noCache: true },
        component: () => import('@/views/project/ProjectEnv.vue')
      },
      { 
        path: '/project/detail/projectSteps', 
        name: 'projectPipeline',
        meta: { title: '项目流程', noCache: true},
        component: () => import('@/views/project/ProjectPipeline.vue'), 
      },
      { 
        path: '/project/pipelines/:pipeId',
        name: 'pipelinesAdd', 
        component: () => import('@/views/project/detail/PipelineAdd.vue'),
        meta: { title: '流程详情', noCache: true,
        },
        hidden: true
      },
    ],
  }
]
