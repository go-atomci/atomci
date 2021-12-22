import Layout from '@/layout'
export function projectDetailRouter() {
  return [
    {
      path: '/project/:projectID/detail',
      name: 'projectMenu',
      component: Layout,
      meta: { title: '项目概览', noCache: true },
      redirect: '/project/:projectID/detail/deashbord',
      children: [
        {
          path: '/project/:projectID/detail/deashbord',  
          name: 'projectSummary',
          component: () => import('@/views/project/ProjectDashboard.vue'),
          meta: { title: '项目概览', noCache: true },
        }
      ]
    },
    {
      path: '/project/:projectID/detail/app',
      name: 'projectApp',
      component: Layout,
      meta: { title: '应用代码', noCache: true},
      children: [
        {
          path: '/project/:projectID/detail/app',
          name: 'projectApp',
          meta: { title: '应用代码', noCache: true},
          component: () => import('@/views/project/ProjectApp.vue'),
        },
        {
          path: '/project/addApp/:projectID',
          name: 'addApp',
          meta: { title: '新增应用', noCache: true },
          component: () => import('@/views/project/detail/ProjectAppAdd.vue'),
          hidden: true
        },
        {
          path: '/project/projectAppDetail/:projectID/:appId/:tabs',
          name: 'projectAppDetail',
          meta: { title: '代码仓库详情', noCache: true },
          component: () => import('@/views/project/detail/ProjectAppDetail.vue'),
          hidden: true
        },
      ]
    },
    {
      path: '/project/:projectID/detail/ci',
      name: 'projectCI',
      component: Layout,
      meta: { title: '构建部署', noCache: true },
      children: [
        {
          path: '/project/:projectID/detail/ci',
          name: 'projectCI',
          meta: { title: '构建部署', noCache: true },
          component: () => import('@/views/project/ProjectCICD.vue'),
        },
        {
          path: '/project/projectCIDetail/:projectID/:versionId',
          meta: { title: '构建部署详情', noCache: true },
          name: 'projectCIDetail',
          component: () => import('@/views/project/detail/ProjectCIDetail.vue'),
          hidden: true
        },
        {
          path: '/project/projectPubDetail/:projectID/:jobName/:runId/:stageId',
          rename: '发布详情',
          name: 'projectPubDetail',
          component: () => import('@/views/project/detail/ProjectPubDetail.vue'),
          hidden: true
        },
      ]
    },
    {
      path: '/project/:projectID/detail/service',
      name: 'projectService',
      component: Layout,
      meta: { title: '应用服务', noCache: true },
      children: [
        {
          path: '/project/:projectID/detail/service',
          meta: { title: '应用服务', noCache: true },
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
    {
      path: '/project/:projectID/detail/projectSets',
      name: 'projectSetting',
      component: Layout,
      meta: { title: '项目设置', noCache: true },
      redirect: '/project/:projectID/detail/projectEnv',
      children: [
        {
          path: '/project/:projectID/detail/projectEnv',
          name: 'projectEnv',
          meta: { title: '项目环境', noCache: true },
          component: () => import('@/views/project/ProjectEnv.vue')
        },
        {
          path: '/project/:projectID/detail/projectSteps',
          name: 'projectPipeline',
          meta: { title: '项目流程', noCache: true },
          component: () => import('@/views/project/ProjectPipeline.vue'),
        },
        {
          path: '/project/pipelines/:pipeId',
          name: 'pipelinesAdd',
          component: () => import('@/views/project/detail/PipelineAdd.vue'),
          meta: {
            title: '流程详情', noCache: true,
          },
          hidden: true
        },
      ],
    }
  ]
}