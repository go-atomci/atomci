import Vue from 'vue';
import Router from 'vue-router';
import backend from '../api/backend';
import store from '@/store';
import Layout from '@/layout'
import projectRouter from './modules/projectRouter'
import { projectDetailRouter } from './modules/projectDetailRouter'

import { getToken } from '@/utils/auth' // get token from cookie

import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style

NProgress.configure({ showSpinner: false }) // NProgress Configuration

Vue.use(Router);

let constantRoutes = [
    // 登录
    {
        path: '/login',
        type: 'login',
        hiden: true,
        name: 'login',
        component: () => import('../views/Login.vue')
    },
    // 登出
    {
        path: '/logout',
        type: 'logout',
        hiden: true,
        name: 'logout',
        component: () => import('../views/Logout.vue')
    },
    projectRouter,
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
            },
            {
                path: '/project/addApp/:projectId',
                name: 'addApp',
                meta: { title: '新增应用', noCache: true },
                component: () => import('@/views/project/detail/ProjectAppAdd.vue'),
                hidden: true
            },
            {
                path: '/project/projectAppDetail/:projectId/:appId/:tabs',
                name: 'projectAppDetail',
                meta: { title: '代码仓库详情', noCache: true },
                component: () => import('@/views/project/detail/ProjectAppDetail.vue'),
                hidden: true
            },
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
            },
            {
                path: '/project/projectCIDetail/:projectId/:versionId',
                meta: { title: '构建部署详情', noCache: true },
                name: 'projectCIDetail',
                component: () => import('@/views/project/detail/ProjectCIDetail.vue'),
                hidden: true
            },
            {
                path: '/project/projectPubDetail/:projectId/:jobName/:runId/:stageId',
                rename: '发布详情',
                name: 'projectPubDetail',
                component: () => import('@/views/project/detail/ProjectPubDetail.vue'),
                hidden: true
            },
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
    //     path: '/project/detail/projectStats',
    //     name: 'projectStatistics',
    //     component: Layout,
    //     meta: { title: '发布统计', noCache: true},
    //     children: [
    //     {
    //         path: '/project/detail/projectStats',
    //         name: 'projectStatistics',
    //         meta: { title: '发布统计', noCache: true },
    //         component: () => import('@/views/project/CICDStats.vue'),
    //     }
    //     ]
    // },
    {
        path: '/project/detail/projectSets',
        name: 'project-setting',
        component: Layout,
        redirect: '/project/detail/projectInfo',
        meta: { title: '项目设置', noCache: true},
        children: [
            // { 
            //     path: '/project/detail/projectInfo', 
            //     name: 'projectInfo', 
            //     meta: { title: '项目信息', noCache: true },
            //     component: () => import('@/views/project/ProjectInfo.vue')
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
    
    },
];

const otherUrl = [
    {
        path: '/404',
        type: 'error',
        name: 'error404',
        hiden: true,
        component: Layout,
        children: [
            {
                path: '/404',
                name: 'error404',
                component: () => import('@/views/error/404.vue')
            },
        ]
    },
    {
        path: '/403',
        type: 'error',
        name: 'error403',
        hiden: true,
        component: Layout,
        children: [{
            path: '/403',
            name: 'error403',
            component: () => import('@/views/error/403.vue')
        },]
    },
];

export const asyncRoutes = [
    {
        path: '/settings',
        component: Layout,
        redirect: '/settings/integrate',
        name: 'settings',
        meta: { title: '系统设置', icon: 'chart'},
        children: [
            {
                path: '/settings/integrate',
                component: () => import('@/views/setting/IntegrateSetting.vue'),
                name: 'serviceIntegrate',
                meta: { title: '服务集成', noCache: true }
            },
            {
                path: '/settings/task',
                component: () => import('@/views/setting/Node.vue'),
                name: 'taskTemplate',
                meta: { title: '任务模板', noCache: true }
            },
            {
                path: '/settings/compile_env',
                component: () => import('@/views/setting/CompileEnv.vue'),
                name: 'taskTemplate',
                meta: { title: '编译环境', noCache: true }
            },
        ]
    },
    {
        path: '/users',
        component: Layout,
        redirect: '/sysusers',
        name: "users",
        meta: { title: '用户管理', noCache: true },
        children: [
            {
                path: '/sysusers',
                name: 'usermanage',
                component: () => import('@/views/user/UserManage.vue'),
                meta: { title: '用户列表', noCache: true }
            },
            { 
                path: '/sysusers/:user/detail', 
                meta: { title: '用户授权', noCache: true }, 
                name: 'managementUser', 
                component: () => import('@/views/user/UserDetail.vue'), 
                hidden: true 
            }, 
        ]
    },

    {
        path: '/roles',
        component: Layout,
        redirect: '/sysroles',
        name: "roles",
        meta: { title: '角色管理', noCache: true },
        children: [
            {
                path: '/sysroles',
                name: 'userrole',
                component: () => import('@/views/role/UserRole.vue'),
                meta: { title: '角色列表', noCache: true }
            },
            { 
                path: '/sysroles/:role', 
                meta: { title: '资源操作', noCache: true}, 
                name: 'listPermission', 
                component: () => import('@/views/role/detail/RoleResOpers.vue'), 
                hidden: true 
            }, // 角色权限操作
        ]
    },

    {
        path: '/audit',
        component: Layout,
        redirect: '/sysaudit',
        name: "audit",
        meta: { title: '操作审计', noCache: true },
        children: [
            {
                path: '/sysaudit',
                name: 'sysaudit',
                component: () => import('@/views/Audit.vue'),
                meta: { title: '操作审计', noCache: true }
            },
        ]
    },

    // 404 page must be placed at the end !!!
    { path: '*', redirect: '/404', hidden: true }
]


// TODO: comment tmp
console.log('current constantRoutes: ')
console.log(constantRoutes)

const createRouter = () => new Router({
    mode: 'history',
    // scrollBehavior: () => ({ y: 0 }),
    routes: constantRoutes
});

const router = createRouter()

const whiteList = ['/login'] // no redirect whitelist


router.beforeEach((to, from, next) => {
    NProgress.start()

    // determine whether the user has logged in
    const hasToken = getToken()
    const defaultRedirectPath = "/project";
    if (hasToken) {
        if (to.path === '/login'  || to.path === '/' ) {
            next({ path: defaultRedirectPath })
            NProgress.done() // hack: https://github.com/PanJiaChen/vue-element-admin/pull/2939
        } else {
            // determine whether the user has obtained his permission roles through getInfo
            const hasRoles = store.getters.roles && store.getters.roles.length > 0
            if (hasRoles) {
                next()
                NProgress.done()
            } else {
                try {
                    backend.getUserInfo((data) => {
                        if (data && data.user) {
                            store.dispatch('user/setUserInfo', data);

                            // TODO: dynamically add accessible routes
                            // generate accessible routes map based on roles
                            const accessRoutes = generateRoutes(data.admin)
                            
                            
                            // dynamically add accessible routes
                            router.addRoutes(accessRoutes)
                            
                            
                            next({ ...to, replace: true })
                        } else {
                            next(`/login?redirect=${to.path}`)
                        }
                    });
                } catch (error) {
                    console.log(error)
                    store.dispatch('user/resetToken')
                    Message.error(error || 'Has Error')
                    next(`/login?redirect=${to.path}`)
                    NProgress.done()
                }
            }
        }
    }
    else {
        /* has no token*/
    
        if (whiteList.indexOf(to.path) !== -1) {
          // in the free login whitelist, go directly
            next()
            NProgress.done()
        } else {
          // other pages that do not have permission to access are redirected to the login page.
            next(`/login?redirect=${to.path}`)
            NProgress.done()
        }
    }
});

router.onError((error) => {
    const pattern = /Loading chunk (\d)+ failed/g;
    const isChunkLoadFailed = error.message.match(pattern);
    const targetPath = router.history.pending.fullPath;
    if (isChunkLoadFailed) {
        console.info('router-error=====');
        router.replace(targetPath);
    }
});

export default router;

export function getUserSibeBarRoutes(routerPath) {
    let routers = []
    console.log('current path: ', routerPath)
    if (routerPath.startsWith('/settings/') || routerPath.startsWith('/sysusers') || routerPath.startsWith('/sysroles') || routerPath === '/sysaudit' || routerPath === "/environment" || routerPath === "/node" || routerPath === "/pipelines") {
        routers = asyncRoutes
    } else if (routerPath.startsWith('/project/')) {
        routers = projectDetailRouter
    } else if (routerPath === '/' || routerPath === '/project') {
        routers = [projectRouter]
    } else {
        console.log(routerPath)
    }
    return routers
}

function generateRoutes(isAdmin) {
    let accessedRoutes = []
    if (isAdmin === 1) {
        accessedRoutes = asyncRoutes || []
    }
    // else {
    //   accessedRoutes = filterAsyncRoutes(asyncRoutes, roles)
    // }
    return accessedRoutes
}
