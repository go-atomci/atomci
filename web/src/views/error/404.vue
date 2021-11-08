<template>
  <div class="panel-error">
    <div class="panel-container">
      <div class="panel-content">
        <div class="error-404"></div>
        <p class="large-gray">404!</p>
        <p class="small-gray mt10">程序员被外星人抓走啦，没有找到你的页面</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .panel-error {
    position: absolute;
    width: 100%;
    height: 90%;
    height: calc(100% - 75px);
    padding: 20px;
    top: 70px;
  }
  .panel-container {
    height: 100%;
    background-color: #fff;
  }
  .panel-content {
    width: 500px;
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    text-align: center;
  }
  .error-404 {
    height: 225px;
  }
  .small-gray {
    font-size: 14px;
    font-family: PingFangSC-Regular,PingFang SC;
    font-weight: 400;
    color: rgba(144,147,153,1);
    line-height: 20px;
  }
  .large-gray {
    font-size: 32px;
    font-family: PingFangSC-Regular,PingFang SC;
    font-weight: 400;
    color: rgba(144,147,153,1);
    line-height: 32px;
  }
</style>
<script>
  import { MessageBox } from 'element-ui';
  export default {
    data() {
      return {};
    },
    methods: {
      logout() {
        MessageBox.confirm('确认退出吗?', this.$t('bm.infrast.tips'), {type: 'warning'})
        .then(() => {
          // 确认
          this.$router.push('/logout'); // 用go刷新
        }).catch(() => {});
      },
    },
  }
</script>