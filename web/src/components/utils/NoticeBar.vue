<template>
  <div class="notice-bar" ref="noticebar">
    <div class="scroll_box" @mouseover="onOver" @mouseleave="onLeave">
      <span class="scroll_begin">
        <span class="item">{{text}}</span>
      </span>
      <span class="scroll_end"></span>
    </div>
  </div>
</template>
<script>
export default {
  name: 'notice-bar',
  props: ['msg'],
  data() {
    return {
      text: '',
      timer: null,
    };
  },
  methods: {
    move() {
      const speed = 20;
      const el = this.$refs.noticebar;
      const scrollBegin = el.querySelector('.scroll_begin');
      const scrollEnd = el.querySelector('.scroll_end');
      const scrollBox = el.querySelector('.scroll_box');
      const len = Math.floor(scrollBox.offsetWidth / scrollBegin.offsetWidth);
      scrollEnd.innerHTML = scrollBegin.innerHTML;

      for (let i = 0; i < len; i++) {
        scrollBegin.appendChild(scrollBegin.querySelector('.item').cloneNode(true));
      }

      const Marquee = () => {
        scrollEnd.offsetWidth - scrollBox.scrollLeft <= 0
          ? scrollBox.scrollLeft -= scrollBegin.offsetWidth
          : scrollBox.scrollLeft += 1;
      };
      this.timer = setInterval(Marquee, speed);
    },
    onOver() {
      clearInterval(this.timer);
    },
    onLeave() {
      this.move();
    },
  },
  mounted() {
    this.text = this.msg;
  },
  updated() {
    this.move();
  },
};
</script>
<style lang="scss">
.notice-bar {
  position: relative;
  overflow: hidden;
  background: #fef0f0;
  border: dashed 1px #f3baba;
  border-radius: 4px;
  height: 30px;
  line-height: 30px;
  margin: 1px 0 10px;
  .notice-icon {
    position: absolute;
    height: 100%;
    z-index: 2;
    background: #fef0f0;
    font-size: 20px;
    float: left;
    color: #ce5a5a;
    padding: 0 5px;
  }
  .scroll_box {
    width: calc(100% - 70px);
    padding-left: 50px;
    white-space: nowrap;
    overflow: hidden;
  }
  .scroll_begin,
  .scroll_end {
    display: inline-block;
    height: 100%;
    .item {
      display: inline-block;
      height: 100%;
      margin-right: 50px;
    }
  }
}
</style>
