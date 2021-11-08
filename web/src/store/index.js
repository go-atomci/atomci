import Vue from 'vue';
import Vuex from 'vuex';
import getters from './getters'

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';

// https://webpack.js.org/guides/dependency-management/#requirecontext
const modulesFiles = require.context('./modules', true, /\.js$/)

// you do not need `import app from './modules/app'`
// it will auto require all vuex module from modules file
const modules = modulesFiles.keys().reduce((modules, modulePath) => {
  // set './app.js' => 'app'
  const moduleName = modulePath.replace(/^\.\/(.*)\.\w+$/, '$1')
  const value = modulesFiles(modulePath)
  modules[moduleName] = value.default
  return modules
}, {})

// export default new Vuex.Store({
//   modules,
//   // state: index.state,
//   // getters: index.getters,
//   // mutations: index.mutations,
//   // actions: index.actions,
//   getters
//   // strict: debug,
// });


const store = new Vuex.Store({
  modules,
  getters
})

export default store
