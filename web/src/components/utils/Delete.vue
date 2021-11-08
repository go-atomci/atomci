<template>
  <div></div>
</template>
<script>
import { Message, MessageBox } from 'element-ui';
import backend from '../../api/backend';

export default {
  props: ['routerName', 'isRefresh'],
  methods: {
    doDelete(fnName, ...key) {
      MessageBox.confirm(this.$t('bm.add.sureDelete'), this.$t('bm.infrast.tips'), { type: 'warning' })
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
    // 成员角色删除
    doDeleteBody(fnName, body, group, user) {
      MessageBox.confirm(this.$t('bm.add.sureDelete'), this.$t('bm.infrast.tips'), { type: 'warning' })
        .then(() => {
          backend[fnName](body, group, user, () => {
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
