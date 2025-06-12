<template>
    <div class="container">
      <div class="logo">
        <svg> ... </svg>
        <div class="logo-text">金铲铲</div>
      </div>
      <div class="content">
        <div class="content-inner">
          <RegisterForm />
        </div>
      </div>
    </div>
  </template>
  
  <script lang="ts" setup>
  import RegisterForm from './components/register-form.vue';
  // Import correct locale files but don't create a local i18n instance
  import '@/views/register/locale/en-US';
  import '@/views/register/locale/zh-CN';
  </script>
  
  <style lang="less" scoped>
  .container {
    display: flex;
    height: 100vh;
    background: linear-gradient(to right, #33ccff 0%, #ff99cc 100%);
  }
  
  .logo {
    position: fixed;
    top: 24px;
    left: 22px;
    display: inline-flex;
    align-items: center;
  
    &-text {
      margin-left: 4px;
      color: var(--color-fill-1);
      font-size: 20px;
    }
  }
  
  .content {
    display: flex;
    flex: 1;
    align-items: center;
    justify-content: center;
  }
  </style>
  