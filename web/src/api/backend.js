import Vue from 'vue';
import { Message } from 'element-ui';
import store from '../store/index';
import i18n from '../i18n/index';

const language = JSON.parse(window.localStorage.getItem('language'));
const langs = (language && language.key) || 'zh-CN';

const Package = {
  // 400 提示错误信息
  notifyErr400(ErrMsg) {
    if (ErrMsg && ErrMsg.ErrCode) {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.tipWarn,
        dangerouslyUseHTMLString: true,
        message: `<div><div style="padding-bottom:10px;">${ErrMsg.ErrCode} : ${ErrMsg.ErrMsg}</div><div>${ErrMsg.ErrDetail}</div></div>`,
        type: 'error',
        duration: 15000,
      });
    } else {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.tipWarn,
        message: `${i18n.messages[langs].bm.add.paramsError}${ErrMsg}`,
        type: 'error',
        duration: 3000,
      });
    }
  },
  // 401 用户认证失败或用户登陆已失效
  notifyErr401(ErrMsg) {
    if (ErrMsg && ErrMsg.ErrCode) {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.tipWarn,
        dangerouslyUseHTMLString: true,
        message: `<div><div style="padding-bottom:10px;">${ErrMsg.ErrCode} : ${ErrMsg.ErrMsg}</div><div>${ErrMsg.ErrDetail}</div></div>`,
        type: 'error',
        duration: 15000,
      });
    } else {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.tipWarn,
        message: i18n.messages[langs].bm.add.errorMsg_401,
        type: 'error',
        duration: 3000,
      });
    }
  },
  // 403 权限拒绝
  notifyErr403(ErrMsg) {
    if (ErrMsg && ErrMsg.ErrCode) {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.tipWarn,
        dangerouslyUseHTMLString: true,
        message: `<div><div style="padding-bottom:10px;">${ErrMsg.ErrCode} : ${ErrMsg.ErrMsg}</div><div>${ErrMsg.ErrDetail}</div></div>`,
        type: 'error',
        duration: 15000,
      });
    } else {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.tipWarn,
        message: i18n.messages[langs].bm.add.errorMsg_403,
        type: 'error',
        duration: 3000,
      });
    }
  },
  // 500 错误信息
  notifyErr500(ErrMsg) {
    if (ErrMsg && ErrMsg.ErrCode) {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.serviceErr,
        dangerouslyUseHTMLString: true,
        message: `<div><div style="padding-bottom:10px;">${ErrMsg.ErrCode} : ${ErrMsg.ErrMsg}</div><div>${ErrMsg.ErrDetail}</div></div>`,
        type: 'error',
        duration: 3000,
      });
    } else {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.serviceErr,
        message: `${ErrMsg}`,
        type: 'error',
        duration: 3000,
      });
    }
  },
  // 502 错误信息
  notifyErr502(ErrMsg) {
    if (ErrMsg && ErrMsg.ErrCode) {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.tipWarn,
        dangerouslyUseHTMLString: true,
        message: `<div><div style="padding-bottom:10px;">${ErrMsg.ErrCode} : ${ErrMsg.ErrMsg}</div><div>${ErrMsg.ErrDetail}</div></div>`,
        type: 'error',
        duration: 3000,
      });
    } else {
      Message({
        showClose: true,
        title: i18n.messages[langs].bm.add.tipWarn,
        message: i18n.messages[langs].bm.add.responseTimeOut,
        type: 'error',
        duration: 3000,
      });
    }
  },
  notifyErr503() {
    Message({
      showClose: true,
      title: i18n.messages[langs].bm.add.tipWarn,
      message: '服务不可用，请稍后重试！',
      type: 'error',
      duration: 3000,
    });
  },
  httpMethods(method, url, cb, body = null, errCb = null, showMsg = true) {
    let loadingName;
    let headers = { headers: { 'Authorization': 'Bearer ' + backendAPI.getCookie('Authorization') } }
    let arrArgs = [url, headers];
    switch (method) {
      case 'get':
        if (store.getters.getNeedLoading) {
          loadingName = 'setLoading';
          store.dispatch(loadingName, true);
        }
        break;
      case 'delete':
        // 不要重复单一提交
        if (store.getters.getLoading) {
          console.warn('重复DELETE提交');
          return;
        }
        loadingName = 'setLoading';
        store.dispatch(loadingName, true);
        if (body) {
          // deleteGroupUserRole, 401 delete 带参数 401错误
          arrArgs = [url, { body }, headers];
        } else {
          arrArgs = [url, headers]
        }
        break;
      case 'post':
        if (
          !~url.indexOf('/apps/list')
          && !~url.indexOf('/ingress')
          && !~url.indexOf('/template')
          && !~url.indexOf('/projects')
          && !~url.indexOf('/rules/list')
          && !~url.indexOf('/issues/')
          && !~url.indexOf('/metrics/')
          && !~url.indexOf('/status?kind=')
          && !~url.indexOf('/status/')
          && !~url.indexOf('/duration/')
          && !~url.indexOf('/ops/')
          && !~url.indexOf('/successrate/')
          && !~url.indexOf('/status/node')
          && store.getters.getPopLoading
        ) {
          console.warn(url);
          console.warn('重复POST提交');
          return;
        }
        loadingName = 'setPopLoading';
        store.dispatch(loadingName, true);
        arrArgs = [url, body, headers];
        break;
      case 'put':
        if (store.getters.getPopLoading) {
          console.warn('重复PUT提交');
          return;
        }
        loadingName = 'setPopLoading';
        store.dispatch(loadingName, true);
        arrArgs = [url, body, headers];
        break;
      case 'patch':
        arrArgs = [url, body, headers];
        break;
    }
    Vue.http[method](...arrArgs).then(
      (response) => {
        if (response.body.IsSuccess) {
          loadingName && store.dispatch(loadingName, false);
          cb && cb(response.body.Data);
        } else if(response.body.loginUrl) {
          cb && cb(response.body);
        } else {
          if (response.status === 400) {
            showMsg && Package.notifyErr400(response.body); // .ErrMsg
          } else if (response.status === 401) {
            window.sessionStorage.clear();
            backendAPI.setCookie('Authorization', '');
            showMsg && Package.notifyErr401();
          } else if (response.status === 403) {
            showMsg && Package.notifyErr403();
          } else if (response.status === 502) {
            showMsg && Package.notifyErr502();
          } else if(response.status === 503 || response.status === 504) {
            showMsg && Package.notifyErr503();
          } else {
            showMsg && Package.notifyErr500(response.body);
          }
          loadingName && store.dispatch(loadingName, false);
          errCb && errCb();
        }
      },
      (response) => {
        setTimeout(() => {
          loadingName && store.dispatch(loadingName, false);
        }, 100);
        if (response.status === 400) {
          showMsg && Package.notifyErr400(response.body); // .ErrMsg
        } else if (response.status === 401) {
          window.sessionStorage.clear();
          backendAPI.setCookie('Authorization', '');
          showMsg && Package.notifyErr401();
        } else if (response.status === 403) {
          showMsg && Package.notifyErr403();
        } else if (response.status === 502) {
          showMsg && Package.notifyErr502();
        } else if(response.status === 503 || response.status === 504) {
          showMsg && Package.notifyErr503();
        } else {
          showMsg && Package.notifyErr500(response.body);
        }
        errCb && errCb();
      }
    );
  },
};
const backendAPI = {
  setCookie(name, value, min) {
    const exp = new Date();
    exp.setTime(exp.getTime() + min * 60 * 1000);
    document.cookie = `${name}=${value}; path=/;expires=${exp.toGMTString()}`;
  },
  getCookie(objname) {
    const cookie = document.cookie;
    if (!cookie){
      return '';
    }
    const arr = cookie.match(new RegExp(`(^| )${objname}=([^;]*)(;|$)`));
    if (arr != null) return unescape(arr[2]);
    return '';
  },
  updateLangCookie() {
    const lang = JSON.parse(window.localStorage.getItem('language')) || { key: 'zh-CN', name: '简体中文' };
    if (lang) {
      this.setCookie('lang', lang.key);
    }
  },
  // login
  login(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/login', cb, body)
  },

  // 查询项目流水线列表
  projectReleaseList(projectId, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectId}/publishes`, cb, body);
  },
  // 回退拿数据
  goregression(projectId, publishId, stageId, action, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${projectId}/publishes/${publishId}/stages/${stageId}/${action}`, cb);
  },
  // 流水线列表关闭接口
  goclose(projectId, publishId, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${projectId}/publishes/${publishId}`, cb);
  },

  // 获取jenkins地址
  getJenkinsServer(cb, stageId) {
    Package.httpMethods('get', `/atomci/api/v1/pipelines/stages/${stageId}/jenkins-config`, cb);
  },
  // 流水线删除接口
  getDeletionPublish(projectId, publishId, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${projectId}/publishes/${publishId}`, cb);
  },

  // 流水线操作
  publishListOperation(projectId, publishId, stageId, stepName, body, cb, errCb) {
    Package.httpMethods('post', `/atomci/api/v1/pipelines/${projectId}/publishes/${publishId}/stages/${stageId}/steps/${stepName}`, cb, body, errCb);
  },

  // 返回描述信息
  getDescription(projectId, publishId, stageId, cb) {
    Package.httpMethods('get', `/atomci/api/v1/pipelines/${projectId}/publishes/${publishId}/stages/${stageId}/steps/publish-audit`, cb);
  },
  // 流水线回退/流转至下一阶段
  gocontinue(projectId, publishId, stageId, action, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectId}/publishes/${publishId}/stages/${stageId}/${action}`, cb, body);
  },

  // 关闭功能的完善
  closeProject(projectId, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectId}/checkProjectOwner`, cb);
  },

  // 查询流水线
  getPipline(pid, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${pid}/pipeline/`, cb);
  },

  // 获得集群列表
  getClusterList(cb) {
    Package.httpMethods('get', '/atomci/api/v1/integrate/clusters', cb);
  },

  // 查询应用编排资源空间 项目内使用
  getAppArrange_namespace(projectId, arrange_env, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${projectId}/arrange_env/${arrange_env}/namespaces`, cb);
  },

  // 获得用户信息
  getUserInfo(cb, errCb) {
    Package.httpMethods('get', '/atomci/api/v1/getCurrentUser', cb, null, errCb);
  },


  // 获取服务详情事件记录
  getAppEventRecord(cluster, namespace, appname, cb) {
    Package.httpMethods('get', `/atomci/api/v1/clusters/${cluster}/namespaces/${namespace}/apps/${appname}/event`, cb);
  },

  // 获取流水线详情数据
  getListdetail(projectId, publishid, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${projectId}/publishes/${publishid}`, cb);
  },
  
  // 删除项目
  delProject(key, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${key}`, cb);
  },
  // 获取项目列表
  getProjectList(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/projects', cb, body);
  },
  // 获取项目详情
  getProjectDetail(id, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${id}`, cb);
  },
  // 获取项目应用模块列表
  getProjectApp(id, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${id}/apps`, cb);
  },
  // 获得应用仓库列表
  getWarehouse(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/apps', cb, body);
  },
  // 关联应用
  postWarehouse(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/apps/create', cb, body);
  },
  // 获得应用所有分支列表
  getBranches(appID, cb) {
    Package.httpMethods('get', `/atomci/api/v1/apps/${appID}/branches`, cb);
  },
  // 设置应用编排信息
  setArrangement(scmAppId, env, group, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/apps/${scmAppId}/${env}/${group}/arrange`, cb, body);
  },
  // 获得应用编排详情
  getArrangement(scmAppId, env, group, cb) {
    Package.httpMethods('get', `/atomci/api/v1/apps/${scmAppId}/${env}/${group}arrange/`, cb, null, null, false);
  },
  // 获得用户列表数据
  getUserList(cb) {
    Package.httpMethods('get', '/atomci/api/v1/users', cb);
  },
  // 创建用户
  addUser(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/users', cb, body);
  },
  // 创建用户组
  addGroup(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/groups', cb, body);
  },
  // 创建角色
  addRole(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/roles', cb, body);
  },
  // 创建应用
  addApp(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/applications', cb, body);
  },
  // 获得用户列表数据
  getUserViewList(user, cb) {
    Package.httpMethods('get', `/atomci/api/v1/users/${user}`, cb);
  },
  // 添加用户详情-角色
  addUserRoles(user, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/users/${user}/roles`, cb, body);
  },
  // 添加用户详情-约束
  addUserConstraints(user, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/users/${user}/constraints`, cb, body);
  },
  // 更新用户详情-约束
  putUserConstraints(user, constraints, body, cb) {
    Package.httpMethods(
      'put',
      `/atomci/api/v1/users/${user}/constraints/${constraints}`,
      cb,
      body
    );
  },
  // 获得用户组列表数据
  getGroupViewList(group, cb) {
    Package.httpMethods('get', `/atomci/api/v1/groups/${group}`, cb);
  },
  // 添加角色详情-资源操作
  addRolesOperations(role, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/roles/${role}/operations`, cb, body);
  },
  // 删除角色详情列表数据
  delRolePers(role, permissions, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/roles/${role}/permissions/${permissions}`, cb);
  },
  // 更新用户
  updateUser(user, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/users/${user}`, cb, body);
  },

  // 删除用户
  delUser(user, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/users/${user}`, cb);
  },
  // 删除用户权限 - 角色
  delUserRole(user, role, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/users/${user}/roles/${role}`, cb);
  },
  // 删除用户权限 - 约束
  delUserConstraints(user, constraints, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/users/${user}/constraints/${constraints}`, cb);
  },
  // 删除用户组权限 - 用户
  delGroupUser(group, user, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/groups/${group}/users/${user}`, cb);
  },
  // 删除用户组权限 - 角色
  delGroupConstraints(group, constraints, cb) {
    Package.httpMethods(
      'delete',
      `/atomci/api/v1/groups/${group}/constraints/${constraints}`,
      cb
    );
  },
  // 添加用户组权限 - 用户
  addGroupUser(group, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/groups/${group}/users`, cb, body);
  },
  // 添加用户组权限 - 角色
  addGroupRole(group, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/groups/${group}/roles`, cb, body);
  },
  // 添加用户组详情-约束
  addGroupConstraints(group, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/groups/${group}/constraints`, cb, body);
  },
  // 更新用户组详情-约束
  putGroupConstraints(group, constraints, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/groups/${group}/constraints/${constraints}`, cb, body);
  },
  // 操作审计列表
  getAudit(cb) {
    Package.httpMethods('get', '/atomci/api/v1/audit', cb);
  },

    
  // 获得权限列表数据
  getPermissionList(cb) {
    Package.httpMethods('get', '/atomci/api/v1/permissions', cb);
  },
  // 获得角色列表数据
  getRoleList(cb) {
    Package.httpMethods('get', '/atomci/api/v1/roles', cb);
  },
  // 删除列表角色
  delRole(role, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/roles/${role}`, cb);
  },
  // 创建项目
  createProject(body, cb, errCb) {
    Package.httpMethods('post', '/atomci/api/v1/projects/create/', cb, body, errCb);
  },
  // 更新项目
  updateProject(id, body, cb, errCb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${id}`, cb, body, errCb);
  },
  // 添加项目应用模块
  addProjectApp(id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${id}/apps`, cb, body);
  },
  // 删除项目应用模块
  delProjectApp(id, projectAppId, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${id}/apps/${projectAppId}`, cb);
  },
  // 更新应用模块
  updateProjectApp(id, projectAppId, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${id}/apps/${projectAppId}`, cb, body);
  },
  // 同步远程分支
  synBranch(appID, cb) {
    Package.httpMethods('post', `/atomci/api/v1/apps/${appID}/syncBranches`, cb);
  },

  // 获得部门列表
  getBus(cb) {
    Package.httpMethods('get', '/atomci/api/v1/groups', cb);
  },
  // 删除用户组
  delGroup(group, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/groups/${group}`, cb);
  },

  // pod实例状态列表
  getPodStatusViews(clusterName, namespace, appname, podname, cb) {
    Package.httpMethods(
      'get',
      `/atomci/api/v1/clusters/${clusterName}/namespaces/${namespace}/apps/${appname}/pods/${podname}/status`,
      cb
    );
  },

  // 移除应用
  removeService(clusterName, namespace, serviceID, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/clusters/${clusterName}/namespaces/${namespace}/apps/${serviceID}`, cb);
  },
  // 重启应用
  reStart(clusterName, namespace, appname, cb) {
    Package.httpMethods(
      'post',
      `/atomci/api/v1/clusters/${clusterName}/namespaces/${namespace}/apps/${appname}/restart`,
      cb
    );
  },
  // 获得应用详情
  getServiceInspect(clusterName, namespace, serviceID, cb) {
    Package.httpMethods(
      'get',
      `/atomci/api/v1/clusters/${clusterName}/namespaces/${namespace}/apps/${serviceID}`,
      cb
    );
  },
  // 水平扩展实例应用信息
  scaleService(clusterName, namespace, serviceID, replicas, cb) {
    Package.httpMethods(
      'post',
      `/atomci/api/v1/clusters/${clusterName}/namespaces/${namespace}/apps/${serviceID}/scale`
        + `?scaleBy=${replicas}`,
      cb
    );
  },
  // 获得应用日志信息
  getServiceLog(clusterName, namespace, serviceID, podName, containerName, cb) {
    Package.httpMethods('get', 
      `/atomci/api/v1/clusters/${clusterName}/namespaces/${namespace}/apps/${serviceID}/log` + `?podName=${podName}&containerName=${containerName}`, 
      cb
    );
  },
  // 获得控制台输出
  getShellResponse(clusterName, namespace, appName, podName, containerName, command, cb, errCb) {
    Package.httpMethods(
      'get',
      `/atomci/api/v1/clusters/${clusterName}/namespaces/${namespace}/apps/${appName}/pods/${podName}/containernames/${containerName}/terminal?shellCommand=${command}`,
      cb,
      null,
      errCb
    );
  },

  logout(cb, errCb){
    Package.httpMethods('get', '/atomci/api/v1/logout', cb, null, errCb);
  },
  
  // 查询角色列表
  getGroupRoleList(cb) {
    Package.httpMethods('get', `/atomci/api/v1/roles`, cb);
  },
  // 更新组角色
  updateGroupRole(role, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/roles/${role}`, cb, body);
  },

  // 删除用户组权限 - 角色
  delGroupRole(group, role, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/roles/${role}`, cb);
  },

  // 资源操作
  getResourcesOperation(cb) {
    Package.httpMethods('get', '/atomci/api/v1/resources-operations', cb);
  },
    
  // 创建权限策略
  addPolicies(group, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/groups/${group}/policies`, cb, body);
  },
  // 查看权限策略
  getPoliciesDetail(group, policy, cb) {
    Package.httpMethods('get', `/atomci/api/v1/groups/${group}/policies/${policy}`, cb);
  },

  // 查询用户组成员角色
  getGroupUserRole(group, user, cb) {
    Package.httpMethods('get', `/atomci/api/v1/groups/${group}/users/${user}/roles`, cb);
  },
  // 添加用户组成员角色
  addGroupUserRole(group, user, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/groups/${group}/users/${user}/roles`, cb, body);
  },
  // 删除用户组成员角色
  deleteGroupUserRole(role, group, user, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/groups/${group}/users/${user}/roles/${role}`, cb);
  },
  // 查询用户组成员
  getGroupUserList(group, cb) {
    Package.httpMethods('get', `/atomci/api/v1/groups/${group}/users`, cb);
  },
  // 查询角色详情
  getGroupRoleDetail(role, cb) {
    Package.httpMethods('get', `/atomci/api/v1/roles/${role}`, cb);
  },

  // 查询角色操作
  getRoleOperations(role, cb) {
    Package.httpMethods('get', `/atomci/api/v1/roles/${role}/operations`, cb);
  },

  // 删除角色操作
  deleteRoleOperation(role, operationID, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/roles/${role}/operations/${operationID}`, cb);
  },

  // 添加角色绑定用户
  addRoleBindUser(group, role, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/groups/${group}/roles/${role}/bundling`, cb, body);
  },
  // 删除用户绑定角色
  deleteRoleBindUser(body, group, role, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/groups/${group}/roles/${role}/bundling`, cb, body);
  },
  // 查询资源类型
  getResourceTypeList(cb) {
    Package.httpMethods('get', '/atomci/api/v1/resources', cb);
  },
  // 创建资源类型
  addResourceType(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/resources', cb, body);
  },
  // 更新资源类型
  updateResourceType(resourceType, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/resources/${resourceType}`, cb, body);
  },
  // 删除资源类型
  deleteResourceType(resourceType, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/resources/${resourceType}`, cb);
  },
  // 获取资源详情
  getResourceTypeDetail(type, cb) {
    Package.httpMethods('get', `/atomci/api/v1/resources/${type}`, cb);
  },
  // 创建资源操作
  addResourceOperations(resourceType, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/resources/${resourceType}/operations`, cb, body);
  },
  // 更新资源操作
  updateResourceOperations(resourceType, resourceOperation, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/resources/${resourceType}/operations/${resourceOperation}`, cb, body);
  },
  // 删除资源操作
  deleteResourceOperations(resourceType, resourceOperation, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/resources/${resourceType}/operations/${resourceOperation}`, cb);
  },
  // 创建资源约束
  addResourceConstraints(resourceType, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/resources/${resourceType}/constraints`, cb, body);
  },
  // 更新资源操作
  updateResourceConstraints(resourceType, resourceConstraints, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/resources/${resourceType}/constraints/${resourceConstraints}`, cb, body);
  },
  // 删除资源操作
  deleteResourceConstraints(resourceType, resourceConstraints, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/resources/${resourceType}/constraints/${resourceConstraints}`, cb);
  },

  // 用户约束列表
  getConstraintsList(group, user, cb) {
    Package.httpMethods('get', `/atomci/api/v1/groups/${group}/users/${user}/constraints`, cb);
  },
  // 添加用户约束
  postUserConstraints(group, user, constraints, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/groups/${group}/users/${user}/constraints/${constraints}/values`, cb, body);
  },
  // 更新用户约束
  putUserConstraintsItem(group, user, constraints, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/groups/${group}/users/${user}/constraints/${constraints}/values`, cb, body);
  },
  // 删除用户约束
  deleteUserConstraints(group = "system", user, constraints, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/groups/${group}/users/${user}/constraints/${constraints}`, cb);
  },

  // Integrate Setting
  getIntegrateSettings(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/integrate/settings', cb, body);
  },
  getStagesAll(cb) {
    Package.httpMethods('get', '/atomci/api/v1/integrate/settings', cb);
  },
  AddIntegrateSetting(body,cb) {
    Package.httpMethods('post', '/atomci/api/v1/integrate/settings/create', cb, body);
  },
  VerifyIntegrateSetting(body,cb) {
    Package.httpMethods('post', `/atomci/api/v1/integrate/settings/verify`, cb, body);
  },
  editIntegrateSetting(idStep, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/integrate/settings/${idStep}`, cb, body);
  },
  delIntegrateSetting(idStep, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/integrate/settings/${idStep}`, cb);
  },
  // 创建项目环境-所需要
  getAllIntegrateSettings(cb) {
    Package.httpMethods('get', '/atomci/api/v1/integrate/settings', cb);
  },
    
  //阶段列表
  getProjectEnvs(projectID, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectID}/envs`, cb, body);
  },
  getProjectEnvsAll(projectID, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${projectID}/envs`, cb);
  },
  AddProjectEnv(projectID, body,cb) {
    Package.httpMethods('post',  `/atomci/api/v1/projects/${projectID}/envs/create`, cb, body);
  },
  editProjectEnv(projectID, idStep, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${projectID}/envs/${idStep}`, cb, body);
  },
  delProjectEnv(projectID, idStep, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${projectID}/envs/${idStep}`, cb);
  },
  
  //步骤列表
  getStep(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/pipelines/flow/steps', cb, body);
  },
  getStepAll(cb) {
    Package.httpMethods('get', '/atomci/api/v1/pipelines/flow/steps', cb);
  },
  getPipeRow(stepId, cb) {
    Package.httpMethods('get', `/atomci/api/v1/pipelines?step_id=${stepId}`, cb);
  },
  AddStep(body,cb) {
    Package.httpMethods('post', '/atomci/api/v1/pipelines/flow/steps/create', cb, body);
  },
  editStep(idStep, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/pipelines/flow/steps/${idStep}`, cb, body);
  },
  delStep(idStep, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/pipelines/flow/steps/${idStep}`, cb);
  },
  getStepComponent(cb) {
    Package.httpMethods('get', '/atomci/api/v1/pipelines/flow/components', cb);
  },

  // 编译环境
  getCompileEnv(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/integrate/compile_envs', cb, body);
  },
  getCompileEnvAll(cb) {
    Package.httpMethods('get', '/atomci/api/v1/integrate/compile_envs', cb);
  },
  AddCompileEnv(body,cb) {
    Package.httpMethods('post', '/atomci/api/v1/integrate/compile_envs/create', cb, body);
  },
  editCompileEnv(idStep, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/integrate/compile_envs/${idStep}`, cb, body);
  },
  delCompileEnv(idStep, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/integrate/compile_envs/${idStep}`, cb);
  },

  getProjectPipeline(projectID, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectID}/pipelines`, cb, body);
  },
  addProjectPipe(projectID, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectID}/pipelines/create`, cb, body);
  },
  editProjectPipe(projectID, pipelineId, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${projectID}/pipelines/${pipelineId}`, cb, body);
  },
  delProjectPipe(projectID, pipelineId, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${projectID}/pipelines/${pipelineId}`, cb);
  },
  //流水线详情
  getProjectPipeDetail(projectID, pipelineId, cb) {
      Package.httpMethods('get', `/atomci/api/v1/projects/${projectID}/pipelines/${pipelineId}`, cb);
  },
    
  //流程图
  getPipe(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/pipelines', cb, body);
  },
  getPipeAll(cb) {
    Package.httpMethods('get', '/atomci/api/v1/pipelines', cb);
  },
  addPipe(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/pipelines/create', cb, body);
  },
  editPipe(pipelineId, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/pipelines/${pipelineId}`, cb, body);
  },
  delPipe(pipelineId, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/pipelines/${pipelineId}`, cb);
  },
  //流水线详情
  getPipeBase(pipelineId, cb) {
    Package.httpMethods('get', `/atomci/api/v1/pipelines/${pipelineId}/setup`, cb);
  },
  viewPipeBase(project_id, pipeline_id, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${project_id}/pipelines/${pipeline_id}`, cb);
  },
  updatePipeBase(pipelineId, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/pipelines/${pipelineId}/setup`, cb, body);
  },
  //发布中心--项目
  getProject(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/projects', cb, body);
  },
  getProjectCreate(body, cb) {
    Package.httpMethods('post', '/atomci/api/v1/projects/create', cb, body);
  },
  updateNewProject(projectId, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${projectId}`, cb, body);
  },
  delNewProject(projectId, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${projectId}`, cb);
  },
  getProjectMember(projectId, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${projectId}/members`, cb);
  },
  updateProjectMember(projectId, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${projectId}/members`, cb, body);
  },
  delProjectMember(projectId, memberId, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${projectId}/members/${memberId}`, cb);
  },
  getProjectPipe(projectId, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${projectId}/pipelines`, cb);
  },
  updateProjectPipe(projectId, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${projectId}/pipelines`, cb, body);
  },
  delProjectPipe(projectId, stepId, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${projectId}/pipelines/${stepId}`, cb);
  },
  // 应用列表
  getApp(projectId, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectId}/apps`, cb, body);
  },
  getAppAll(projectId, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${projectId}/apps`, cb);
  },
  appDialogBranch(repo_id, code_id, cb, body, errCb) {
    Package.httpMethods('get', `/atomci/api/v1/repos/${repo_id}/projects/${code_id}/branches`, cb, body, errCb);
  },
  // 获取单个代码仓库详情
  getAppDetail(projectId, appID, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${projectId}/apps/${appID}`, cb);
  },
  updateAppInfo(projectId, appId, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${projectId}/apps/${appId}`, cb, body);
  },
  // 更新应用模块
  changeBranch(id, projectAppId, body, cb) {
    Package.httpMethods('patch', `/atomci/api/v1/projects/${id}/apps/${projectAppId}`, cb, body);
  },
  getProjectBranch(project_id, app_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/apps/${app_id}/branches`, cb, body);
  },
  createBranch(project_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/apps/branches`, cb, body);
  },
  asyncBranch(project_id, app_id, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/apps/${app_id}/syncBranches`, cb);
  },
  //新增应用模块
  getRepos(project_id, cb) {
    Package.httpMethods('get', `/atomci/api/v1/repos?project_id=${project_id}`, cb);
  },
  getReposList(repo_id, project_id, body, cb, errcb) {
    Package.httpMethods('post', `/atomci/api/v1/repos/${repo_id}/projects?project_id=${project_id}`, cb, body, errcb);
  },
  //新增
  addAppPro(project_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/apps/create`, cb, body);
  },
  //应用编排
  getProjectArrange(project_id, app_id, arrange_env, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${project_id}/apps/${app_id}/${arrange_env}/arrange`, cb);
  },
  setProjectArrange(project_id, app_id, arrange_env, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/apps/${app_id}/${arrange_env}/arrange`, cb, body);
  },
  parseYamlImages(body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/arrange/yaml/parser`, cb, body);
  },
  //CI/CD
  getProjectCI(projectId, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectId}/publishes`, cb, body);
  },
  addProjectCI(project_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/publishes/create`, cb, body);
  },
  updateProjectCI(project_id, publish_id, body, cb) {
    Package.httpMethods('put', `/atomci/api/v1/projects/${project_id}/publishes/${publish_id}`, cb, body);
  },

  removeApp(project_id, publish_id, publish_app_id, cb) {
    Package.httpMethods('delete', `/atomci/api/v1/projects/${project_id}/publishes/${publish_id}/apps/${publish_app_id}`, cb);
  },
  versionApp(project_id, publish_id, cb) {
    Package.httpMethods('get', `/atomci/api/v1/projects/${project_id}/publishes/${publish_id}/apps/can_added`, cb);
  },
  addVersionApp(project_id, publish_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/publishes/${publish_id}/apps/create`, cb, body);
  },
  getMark(project_id, publish_id, stage_id, cb) {
    Package.httpMethods('get', `/atomci/api/v1/pipelines/${project_id}/publishes/${publish_id}/stages/${stage_id}/steps/manual`, cb);
  },
  setMark(project_id, publish_id, stage_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/pipelines/${project_id}/publishes/${publish_id}/stages/${stage_id}/steps/manual`, cb, body);
  },
  getBuildMerge(project_id, publish_id, stage_id, cb) {
    Package.httpMethods('get', `/atomci/api/v1/pipelines/${project_id}/publishes/${publish_id}/stages/${stage_id}/steps/build`, cb);
  },
  setBuildMerge(project_id, publish_id, stage_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/pipelines/${project_id}/publishes/${publish_id}/stages/${stage_id}/steps/build`, cb, body);
  },
  getDeploy(project_id, publish_id, stage_id, cb) {
    Package.httpMethods('get', `/atomci/api/v1/pipelines/${project_id}/publishes/${publish_id}/stages/${stage_id}/steps/deploy`, cb);
  },
  setDeploy(project_id, publish_id, stage_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/pipelines/${project_id}/publishes/${publish_id}/stages/${stage_id}/steps/deploy`, cb, body);
  },
  // 操作日志
  getOperationLog(projectId, publishId, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${projectId}/publishes/${publishId}/audits`, cb, body);
  },
  // 统计
  getStats(project_id, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/publish/stats`, cb, body);
  },
  // 项目下应用列表
  getProjectServiceList(project_id, cluster, body, cb) {
    Package.httpMethods('post', `/atomci/api/v1/projects/${project_id}/clusters/${cluster}/apps`, cb, body);
  },
  // 项目设置-角色
  getProjectUser(group, role, cb) {
    Package.httpMethods('get', `/atomci/api/v1/groups/${group}/roles/${role}`, cb);
  },
};
export default backendAPI;
