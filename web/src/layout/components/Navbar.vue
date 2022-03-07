<style>

.nav-container {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 9;
  width: 100%;
  /* transition: width 0.28s; */
}

.container .user-img {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: inline-block;
}
.container .topbar-wrap .userinfo-inner .user-img-icon {
  width: 32px;
  height: 32px;
  background-size: 32px 32px;
  margin-left: 0;
  vertical-align: 5px;
}
.container .el-submenu__icon-arrow.el-icon-arrow-down {
  font-size: 14px;
}
.container .el-submenu__icon-arrow.el-icon-arrow-down::before {
  content: '\e790';
}

/* 下拉弹出框 */
.userinfo-menu-container {
  color: rgba(255, 255, 255, 1);
  font-size: 16px;
  line-height: 25px;
  flex: 1;
}
.el-menu--popup-bottom-start {
  padding: 0;
  border-radius: 4px 4px 0px 0px;
  min-width: 278px;
}
body > .el-menu--horizontal {
  min-width: 278px;
}
body > .el-menu--horizontal .el-menu-item:first-child {
  height: auto;
  /* 去除背景色 */
  /* background: linear-gradient(233deg, rgba(106, 17, 203, 1) 0%, rgba(37, 117, 252, 1) 100%); */
  border-radius: 4px 4px 0px 0px;
}
.userinfo-menu {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 6px;
}

.nav-button {
  float: left;
  background: #00364d!important;
  color: snow;
  border: 0px;
  margin: 5px 15px;
  font-size: 13px;
}

.user-icon {
  display: inline-block;
  width: 72px;
  height: 72px;
  border-radius: 50%;
}
.img-user-icon {
  width: 72px;
  height: 72px;
  background-size: 72px 72px;
}
body > .el-menu--horizontal .el-menu .el-menu-item {
  min-height: 40px;
  line-height: 40px;
}
body > .el-menu--horizontal .el-menu-item:not(.is-disabled):not(:first-child):hover {
  background: rgba(242, 246, 252, 1);
}

.el-submenu__title:hover {
    background-color: transparent !important;
    color: #fff !important;
}
</style>

<template>
<div class="nav-container">
  <el-row class="container">
    <!--头部-->
    <el-col :span="24"
            class="topbar-wrap">
      <div class="topbar-logos">
        <a href="/project"><span><img src="@/assets/logo.png"></span></a>
      </div>
      <div class="topbar-title">
            <el-menu :default-active="defaultActiveIndex" class="el-menu-demo" mode="horizontal" router>
              <el-menu-item index="/scmapp" >我的应用</el-menu-item>
              <el-menu-item index="/project" >我的项目</el-menu-item>
              <el-menu-item index="/settings" v-if="menuTrue">系统管理</el-menu-item>
            </el-menu>
      </div>
      <div class="topbar-account topbar-btn">
            <el-menu mode='horizontal' class="userinfo-inner">
              <a href="https://go-atomci.github.io/atomci-press/" target="_blank" class="nav-button" >
                帮助文档
              </a>
              <el-submenu index='4s'
                          v-show="false">
                <template slot='title'>{{currentLanguage.name}}</template>
                <el-menu-item v-for="(item, indx) in languageList"
                              :key="indx"
                              :index="`4-${indx + 1}`"
                              @click="changeLanuage(item)">
                  {{item.name}}
                </el-menu-item>
              </el-submenu>
              <el-submenu index='2s'>
                <template slot='title'>
                  {{nickname}}
                </template>
                <el-menu-item :index="`2-3`"
                              @click.native="logout">
                  <i class="el-icon-switch-button"></i>
                  {{$t('bm.dashboard.logOut')}}
                </el-menu-item>
              </el-submenu>
            </el-menu>
      </div>
    </el-col>
  </el-row>

</div>

</template>
<script>
import { MessageBox } from 'element-ui';
import { mapGetters } from 'vuex';

export default {
  name: 'home',
  data() {
    return {
      loading: false,
      // 支持的语言种类
      languageList: [
        {
          key: 'zh-CN',
          name: '简体中文',
        },
        {
          key: 'en',
          name: 'English',
        },
      ],
      // 当前选中的语言
      currentLanguage: '',
      menuTrue: false,
    };
  },
  components: {},
  created() {
    // 组件创建完后获取数据，
    const isAdmin = this.isAdmin;
    if (isAdmin == 1) {
      this.menuTrue = true;
    }
    
    // this.fetchNavData();
    // 设置默认语言
    const lang = window.localStorage.getItem('language');
    this.currentLanguage = JSON.parse(lang) || this.languageList[0];
  },
  computed: {
    ...mapGetters({
      sidebar: 'sidebar',
      isAdmin: 'isAdmin',
    }),
    nickname() {
      return this.$store.state.user.name
    },
    defaultActiveIndex() {
      if (this.$route.path.startsWith('/project')){
        return '/project'  
      }　else if (this.$route.path.startsWith('/scmapp')) {
        return '/scmapp'  
      } else {
        return '/settings'
      }
    },
  },
  methods: {
    // 添加切换语言功能
    changeLanuage(item) {
      window.loadLanguage(item.key);
      this.currentLanguage = item;
      window.localStorage.setItem('language', JSON.stringify(item));
      window.location.reload();
    },
    // 获取详情的地址
    initDetailNav(truePath, origin) {
      const path = origin.replace(/:[a-zA-Z0-9]+(\/)?/g, function(a, b, num, cont) {
        const len = cont.substr(0, num).split('/').length;
        return truePath.split('/')[len - 1] + (a[a.length - 1] === '/' ? '/' : '');
      });
      return truePath === path;
    },
    childData(children, curPath) {
      let flag = false;
      if (children) {
        for (let i = 0; i < children.length; i++) {
          if (children[i].path === curPath) {
            flag = true;
            break;
          }
          const sFlag = this.initDetailNav(curPath, children[i].path);
          if (sFlag) {
            flag = true;
            break;
          }
          if (children[i].children) {
            let child = children[i].children;
            let valFlag = this.childData(child, curPath);
            if (valFlag) {
              flag = true;
              break;
            }
          }
        }
      }
      return flag;
    },
    fetchNavData() {
      // 初始化菜单激活项
      const curPath = this.$route.path; // 获取当前路由
      const routers = this.$router.options.routes; // 获取路由对象
      let navType = '';
      let navName = '';
      for (let i = 0; i < routers.length; i++) {
        const children = routers[i].children;
        if (children) {
          var types = this.childData(children, curPath);
          if (types) {
            navType = routers[i].type;
            navName = routers[i].name;
            break;
          }
        }
      }
      this.$store.state.topNavState = navType; // 改变topNavState状态的值
      this.$store.state.leftNavState = navName; // 改变leftNavState状态的值

      // if (navType === 'cluster') {
      //   this.defaultActiveIndex = '/mycluster';
      // } else if (navType === 'projectInner') {
      //   this.defaultActiveIndex = '/myproject';
      // } else if (navType === 'project') {
      //   this.defaultActiveIndex = '/myproject';
      // } else {
      //   this.defaultActiveIndex = `/${navName}Manager`;
      // }
    },
    logout() {
      MessageBox.confirm('确认退出吗?', this.$t('bm.infrast.tips'), { type: 'warning' })
        .then(() => {
          // 确认
          this.$store.dispatch('user/logout')
          window.sessionStorage.clear()
          this.$router.push(`/login?redirect=${this.$route.fullPath}`)
        })
        .catch(() => {});
    },
  },
  watch: {
    // $route: 'fetchNavData', // 监听router值改变时，改变导航菜单激活项
  },
};
</script>
