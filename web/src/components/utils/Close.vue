<template>
  <div></div>
</template>
<script>
import { Message, MessageBox } from 'element-ui';
import backend from '../../api/backend';

export default {
  props: ['routerName', 'isRefresh'],
  methods: {
    doClose(fnName, ...key) {
      MessageBox.confirm(this.$t('bm.add.sureCloseProject'), this.$t('bm.infrast.tips'), {
        confirmButtonText: this.$t('bm.other.confirm'),
        type: 'warning',
      })
        .then(() => {
          backend[fnName](...key, () => {
            Message.success(this.$t('bm.add.optionSuc'));
            if (this.$props.routerName) {
              this.$router.push({
                name: this.$props.routerName,
                params: this.$props.routerParams,
                query: {
                  isRefresh: this.$props.isRefresh,
                },
              });
              return;
            }
            this.$emit('getlist');
          });
        })
        .catch(() => {});
    },
  },
};
</script>
