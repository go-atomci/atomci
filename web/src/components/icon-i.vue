<template>
  <i class="font_family" :class="`icon-${type}`"></i>
</template>

<script>
export default {
  props: {
    type: String,
    resetClass: String,
  },
};
</script>

<style>
@import url('//at.alicdn.com/t/font_1568297_wvz48tcg1o.css');
</style>
