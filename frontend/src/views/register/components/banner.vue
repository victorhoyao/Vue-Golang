<template>
    <div class="banner">
      <div class="banner-inner">
        <a-carousel class="carousel" animation-name="fade">
          <a-carousel-item v-for="item in carouselItem" :key="item.slogan">
            <div :key="item.slogan" class="carousel-item">
              <img class="carousel-image" :src="item.image" />
            </div>
          </a-carousel-item>
        </a-carousel>
      </div>
    </div>
  </template>
  
  <script lang="ts" setup>
  import { computed } from 'vue';
  import { useI18n } from 'vue-i18n';
  import bannerImage from '../chaoren.png';
  
  const { t } = useI18n();
  const carouselItem = computed(() => [
    {
      slogan: t('register.banner.slogan1'),
      subSlogan: t('register.banner.subSlogan1'),
      image: bannerImage,
    },
    {
      slogan: t('register.banner.slogan2'),
      subSlogan: t('register.banner.subSlogan2'),
      image: bannerImage,
    },
    {
      slogan: t('register.banner.slogan3'),
      subSlogan: t('register.banner.subSlogan3'),
      image: bannerImage,
    },
  ]);
  </script>
  
  <style lang="less" scoped>
  .banner {
    display: flex;
    align-items: center;
    justify-content: center;
  
    &-inner {
      flex: 1;
      height: 100%;
    }
  }
  
  .carousel {
    height: 100%;
  
    &-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
    }
  
    &-image {
      width: 100%;
      margin-top: 30px;
    }
  }
  </style>
  