<template>
  <div></div>
</template>

<script>
import backend from '../api/backend';
export default {
  created() {
    this.logout();
  },
  methods: {
    // 退出
    logout() {
      backend.logout();
    },
  },
};
</script>
