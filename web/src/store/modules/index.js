const initializeState = {
  loading: false,
  // 弹出框中的loading
  popLoading: false,
  isNeedLoading: true,
  userInfo: {},
  clusterList: [],
  harborList: [],
};
const getters = {
  getLoading(state) {
    return state.loading;
  },
  getPopLoading(state) {
    return state.popLoading;
  },
  getNeedLoading(state) {
    return state.isNeedLoading;
  },
  getUserInfo(state) {
    return state.userInfo;
  },
  getClusterList(state) {
    return state.clusterList;
  },
  getHarborList(state) {
    return state.harborList;
  },
};
// actions
const actions = {
  setLoading({ commit }, loading) {
    commit('SET_LOADING', { loading });
  },
  setPopLoading({ commit }, loading) {
    commit('SET_POP_LOADING', { loading });
  },
  setNeedLoading({ commit }, loading) {
    commit('SET_NEED_LOADING', { loading });
  },
  setPageSize({ commit }, curPageSize) {
    commit('CHANGE_PAGE_SIZE', { curPageSize });
  },
  setClusterList({ commit }, obj) {
    commit('CHANGE_CLUSTER_LIST', { obj });
  },
  setHarborList({ commit }, obj) {
    commit('CHANGE_HARBOR_LIST', { obj });
  },
};
// mutations
const mutations = {
  SET_LOADING(state, { loading }) {
    state.loading = loading;
  },
  SET_POP_LOADING(state, { loading }) {
    state.popLoading = loading;
  },
  SET_NEED_LOADING(state, { loading }) {
    state.isNeedLoading = loading;
  },
  CHANGE_USER_INFO(state, { userInfo }) {
    state.userInfo = userInfo;
  },
  CHANGE_CLUSTER_LIST(state, { obj }) {
    state.clusterList = obj;
    if (state.curCluster === '' && state.clusterList.length) {
      state.curCluster = state.clusterList[0].name;
    }
  },
  CHANGE_HARBOR_LIST(state, { obj }) {
    state.harborList = obj;
  },
};
export default {
  state: initializeState,
  getters,
  actions,
  mutations,
};