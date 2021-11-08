<template>
  <el-button :plain="true" type="success" @click="doRefresh">
    <i class='icon-refresh'></i>
    {{$t('bm.other.refresh')}}</el-button>
</template>
<script>
export default {
  data() {
    return {
      instance: '',
    };
  },
  methods: {
    doRefresh() {
      this.$emit('getlist', true);
      if (this.instance) this.instance.close();
      this.instance = this.$notify({
        title: this.$t('bm.add.success'),
        message: this.$t('bm.add.refreshSuc'),
        type: 'success',
        duration: '1000',
      });
    },
  },
};
</script>
