import Vue from 'vue';
import VueI18n from 'vue-i18n';
import locale from 'element-ui/lib/locale';
import zhLocale from 'element-ui/lib/locale/lang/zh-CN';
import zh from './lang/zh-CN';


Vue.use(VueI18n);

const language = JSON.parse(window.localStorage.getItem('language'));
// 初始化时，从localStorage中读取缓存的语言种类
const i18n = new VueI18n({
  locale: (language && language.key) || 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': Object.assign(zh, zhLocale)
  },
});

// 为了实现element插件的多语言切换
locale.i18n((key, value) => i18n.t(key, value));


const loadedLanguages = ['zh-CN'];

function setI18nLanguage(lang) {
  i18n.locale = lang;
  return lang;
}

export function loadLanguage(lang) {
  if (i18n.locale !== lang) {
    if (!loadedLanguages.includes(lang)) {
      return import(/* webpackChunkName: 'lang-[request]' */ `./lang/${lang}`).then((msgs) => {
        // 增加element-ui 国际化组件的动态引入
        return import(`element-ui/lib/locale/lang/${lang}`).then((data) => {
          i18n.setLocaleMessage(lang, Object.assign(msgs.default, data.default));
          loadedLanguages.push(lang);
          return setI18nLanguage(lang);
        });
      });
    }
    return Promise.resolve(setI18nLanguage(lang));
  }
  return Promise.resolve(lang);
}

window.loadLanguage = loadLanguage;


export default i18n;
