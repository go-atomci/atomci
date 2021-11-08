import backend from '@/api/backend'
import { getToken, setToken, removeToken } from '@/utils/auth'
import router from '@/router'

const state = {
  token: getToken(),
  name: '',
  avatar: '',
  roles: [],
  admin: 0
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_USER: (state, name) => {
    state.name = name
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  },
  SET_ADMIN: (state, admin) => {
    state.admin = admin
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { username, password, login_type } = userInfo
    return new Promise((resolve, reject) => {
      backend.login({ username: username.trim(), password: password, login_type: login_type }, response => {
        const { token } = response
        commit('SET_TOKEN', token)
        setToken(token)
        resolve()
      })
    })
  },

  // get user info
  setUserInfo({ commit }, userInfo) {
    return new Promise(resolve => {
      commit('SET_ROLES', userInfo.roles)
      commit('SET_USER', userInfo.user)
      commit('SET_ADMIN', userInfo.admin)
      resolve()
    })
  },

  // user logout
  logout({ commit, state, dispatch }) {
    return new Promise(resolve => {
    commit('SET_TOKEN', '')
    commit('SET_ROLES', [])
    removeToken()
    // resetRouter()

    // reset visited views and cached views
    // to fixed https://github.com/PanJiaChen/vue-element-admin/issues/2485
    dispatch('tagsView/delAllViews', null, { root: true })

    resolve()
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      removeToken()
      resolve()
    })
  },

  // dynamically modify permissions
  async changeRoles({ commit, dispatch }, role) {
    const token = role + '-token'

    commit('SET_TOKEN', token)
    setToken(token)

    const { roles } = await dispatch('getInfo')

    // resetRouter()

    // generate accessible routes map based on roles
    const accessRoutes = await dispatch('permission/generateRoutes', roles, { root: true })
    // dynamically add accessible routes
    router.addRoutes(accessRoutes)

    // reset visited views and cached views
    dispatch('tagsView/delAllViews', null, { root: true })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
