import Vue from 'vue'

import Cookies from 'js-cookie'

import VueClipboard from 'vue-clipboard2';
import VueResource from 'vue-resource';
import VueCodeMirror from 'vue-codemirror';

import './style';
import '@/assets/iconfont.css';
import '@/assets/css/style.css';


import App from './App'
import i18n from './i18n';
import router from './router';
import store from './store';
import NoResult from './components/utils/NoResult';
import IconI from './components/icon-i';

Vue.config.productionTip = false

Vue.use(VueClipboard);
Vue.use(VueResource);
Vue.use(VueCodeMirror);

import Element from 'element-ui'
import './style/element-variables.scss'
// import enLang from 'element-ui/lib/locale/lang/en'// 如果使用中文语言包请默认支持，无需额外引入，请删除该依赖

import MessageBox from 'element-ui'
import Notification from 'element-ui'

Vue.prototype.$notify = Notification;
Vue.prototype.$confirm = MessageBox.confirm;

Vue.use(Element, {
  size: Cookies.get('size') || 'medium', // set element-ui default size
  // locale: enLang // 如果使用中文，无需设置，请删除
})

Vue.component('no-result', NoResult);
Vue.component('icon-i', IconI);
/* eslint-disable no-new */
new Vue({
  el: '#app',
  i18n,
  store,
  router,
  render: h => h(App)
})
