const getters = {
    sidebar: state => state.app.sidebar,
    size: state => state.app.size,
    device: state => state.app.device,
    visitedViews: state => state.tagsView.visitedViews,
    cachedViews: state => state.tagsView.cachedViews,
    token: state => state.user.token,
    name: state => state.user.name,
    isAdmin: state => state.user.admin,
    roles: state => state.user.roles,
    projectID: state => state.project.projectID
  }
  export default getters
  