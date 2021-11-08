const state = {
  projectID: 0
}

const mutations = {
  SET_PROJECTID: (state, projectID) => {
    state.projectID = projectID
  },
}

const actions = {
  // get user info
  setProjectID({ commit }, projectID) {
    commit('SET_PROJECTID', projectID)
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
