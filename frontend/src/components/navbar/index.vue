<template>
  <div class="navbar">
    <div class="left-side">
      <a-space>
        <svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
          x="0px" y="0px" width="48" height="48" viewBox="0 0 200 200" enable-background="new 0 0 200 200"
          xml:space="preserve">
          <image id="image0" width="200" height="200" x="0" y="0"
            href="data:image/png;base64,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">
          </image>
        </svg> <a-typography-title :style="{ margin: 0, fontSize: '18px' }" :heading="5">
          金铲铲
        </a-typography-title>
        <icon-menu-fold v-if="!topMenu && appStore.device === 'mobile'" style="font-size: 22px; cursor: pointer"
          @click="toggleDrawerMenu" />
      </a-space>
    </div>
    <div class="center-side">
      <Menu v-if="topMenu" />
    </div>
    <ul class="right-side">
      <!-- <li>
        <a-tooltip :content="$t('settings.search')">
          <a-button class="nav-btn" type="outline" :shape="'circle'">
            <template #icon>
              <icon-search />
            </template>
</a-button>
</a-tooltip>
</li> -->
      <li>
        <a-tooltip :content="$t('settings.language')">
          <a-button class="nav-btn" type="outline" :shape="'circle'" @click="setDropDownVisible">
            <template #icon>
              <icon-language />
            </template>
          </a-button>
        </a-tooltip>
        <a-dropdown trigger="click" @select="changeLocale as any">
          <div ref="triggerBtn" class="trigger-btn"></div>
          <template #content>
            <a-doption v-for="item in locales" :key="item.value" :value="item.value">
              <template #icon>
                <icon-check v-show="item.value === currentLocale" />
              </template>
              {{ item.label }}
            </a-doption>
          </template>
        </a-dropdown>
      </li>
      <li>
        <a-tooltip :content="theme === 'light'
          ? $t('settings.navbar.theme.toDark')
          : $t('settings.navbar.theme.toLight')
          ">
          <a-button class="nav-btn" type="outline" :shape="'circle'" @click="handleToggleTheme">
            <template #icon>
              <icon-moon-fill v-if="theme === 'dark'" />
              <icon-sun-fill v-else />
            </template>
          </a-button>
        </a-tooltip>
      </li>
      <!-- <li>
        <a-tooltip :content="$t('settings.navbar.alerts')">
          <div class="message-box-trigger">
            <a-badge :count="9" dot>
              <a-button
                class="nav-btn"
                type="outline"
                :shape="'circle'"
                @click="setPopoverVisible"
              >
                <icon-notification />
              </a-button>
            </a-badge>
          </div>
        </a-tooltip>
        <a-popover
          trigger="click"
          :arrow-style="{ display: 'none' }"
          :content-style="{ padding: 0, minWidth: '400px' }"
          content-class="message-popover"
        >
          <div ref="refBtn" class="ref-btn"></div>
          <template #content>
            <message-box />
          </template>
        </a-popover>
      </li> -->
      <li>
        <a-tooltip :content="isFullscreen
          ? $t('settings.navbar.screen.toExit')
          : $t('settings.navbar.screen.toFull')
          ">
          <a-button class="nav-btn" type="outline" :shape="'circle'" @click="toggleFullScreen">
            <template #icon>
              <icon-fullscreen-exit v-if="isFullscreen" />
              <icon-fullscreen v-else />
            </template>
          </a-button>
        </a-tooltip>
      </li>
      <li>
        <a-tooltip :content="$t('settings.title')">
          <a-button class="nav-btn" type="outline" :shape="'circle'" @click="setVisible">
            <template #icon>
              <icon-settings />
            </template>
          </a-button>
        </a-tooltip>
      </li>
      <li>
        <a-dropdown trigger="click">
          <a-avatar :size="32" :style="{ marginRight: '8px', cursor: 'pointer' }">
            <img alt="avatar" src="../../assets//images/ks.png" />
          </a-avatar>
          <template #content>
            <!-- <a-doption>
              <a-space @click="switchRoles">
                <icon-tag />
                <span>
                  {{ $t('messageBox.switchRoles') }}
                </span>
              </a-space>
            </a-doption> -->
            <!-- <a-doption>
              <a-space @click="$router.push({ name: 'Info' })">
                <icon-user />
                <span>
                  {{ $t('messageBox.userCenter') }}
                </span>
              </a-space>
            </a-doption>
            <a-doption>
              <a-space @click="$router.push({ name: 'Setting' })">
                <icon-settings />
                <span>
                  {{ $t('messageBox.userSettings') }}
                </span>
              </a-space>
            </a-doption> -->
            <a-doption>
              <a-space @click="handleLogout">
                <icon-export />
                <span>
                  {{ $t('messageBox.logout') }}
                </span>
              </a-space>
            </a-doption>
          </template>
        </a-dropdown>
      </li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, inject } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useDark, useToggle, useFullscreen } from '@vueuse/core';
import { useAppStore, useUserStore } from '@/store';
import { LOCALE_OPTIONS } from '@/locale';
import useLocale from '@/hooks/locale';
import useUser from '@/hooks/user';
import Menu from '@/components/menu/index.vue';
import MessageBox from '../message-box/index.vue';

const appStore = useAppStore();
const userStore = useUserStore();
const { logout } = useUser();
const { changeLocale, currentLocale } = useLocale();
const { isFullscreen, toggle: toggleFullScreen } = useFullscreen();
const locales = [...LOCALE_OPTIONS];
const avatar = computed(() => {
  return userStore.userName;
});
const theme = computed(() => {
  return appStore.theme;
});
const topMenu = computed(() => appStore.topMenu && appStore.menu);
const isDark = useDark({
  selector: 'body',
  attribute: 'arco-theme',
  valueDark: 'dark',
  valueLight: 'light',
  storageKey: 'arco-theme',
  onChanged(dark: boolean) {
    // overridden default behavior
    appStore.toggleTheme(dark);
  },
});
const toggleTheme = useToggle(isDark);
const handleToggleTheme = () => {
  toggleTheme();
};
const setVisible = () => {
  appStore.updateSettings({ globalSettings: true });
};
const refBtn = ref();
const triggerBtn = ref();
const setPopoverVisible = () => {
  const event = new MouseEvent('click', {
    view: window,
    bubbles: true,
    cancelable: true,
  });
  refBtn.value.dispatchEvent(event);
};
const handleLogout = () => {
  logout();
};
const setDropDownVisible = () => {
  const event = new MouseEvent('click', {
    view: window,
    bubbles: true,
    cancelable: true,
  });
  triggerBtn.value.dispatchEvent(event);
};

const switchRoles = async () => {
  const res = await userStore.switchRoles();
  Message.success(res as string);
};
const toggleDrawerMenu = inject('toggleDrawerMenu') as () => void;
</script>

<style scoped lang="less">
.navbar {
  display: flex;
  justify-content: space-between;
  height: 100%;
  background-color: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border);
}

.left-side {
  display: flex;
  align-items: center;
  padding-left: 20px;
}

.center-side {
  flex: 1;
}

.right-side {
  display: flex;
  padding-right: 20px;
  list-style: none;

  :deep(.locale-select) {
    border-radius: 20px;
  }

  li {
    display: flex;
    align-items: center;
    padding: 0 10px;
  }

  a {
    color: var(--color-text-1);
    text-decoration: none;
  }

  .nav-btn {
    border-color: rgb(var(--gray-2));
    color: rgb(var(--gray-8));
    font-size: 16px;
  }

  .trigger-btn,
  .ref-btn {
    position: absolute;
    bottom: 14px;
  }

  .trigger-btn {
    margin-left: 14px;
  }
}
</style>

<style lang="less">
.message-popover {
  .arco-popover-content {
    margin-top: 0;
  }
}
</style>
