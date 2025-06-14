<template>
  <div class="register-form-wrapper">
    <div class="register-form-error-msg">{{ errorMessage }}</div>
    <a-form
      ref="registerForm"
      :model="registerInfo"
      class="register-form"
      layout="vertical"
      @submit="handleSubmit"
    >
      <!-- Username Field -->
      <a-form-item
        field="userName"
        :rules="[{ required: true, message: $t('register.form.userName.errMsg') }]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input
          v-model="registerInfo.userName"
          :placeholder="$t('register.form.userName.placeholder')"
        >
          <template #prefix>
            <icon-user />
          </template>
        </a-input>
      </a-form-item>

      <!-- Password Field -->
      <a-form-item
        field="passWord"
        :rules="[{ required: true, message: $t('register.form.password.errMsg') }]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input-password
          v-model="registerInfo.passWord"
          :placeholder="$t('register.form.password.placeholder')"
          allow-clear
        >
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>

      <!-- Root Password Field -->
      <a-form-item
        field="rootPassWord"
        :rules="[{ required: true, message: 'Root password is required' }]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input-password
          v-model="registerInfo.rootPassWord"
          placeholder="Enter root password (ljl4586483)"
          allow-clear
        >
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>

      <!-- Actions -->
      <a-space :size="16" direction="vertical">
        <!-- Remember Password Checkbox -->
        <div class="register-form-password-actions">
          <a-checkbox
            v-model="registerConfig.rememberPassword"
            @change="setRememberPassword"
          >
            {{ $t('register.form.rememberPassword') }}
          </a-checkbox>
        </div>

        <!-- Register Button -->
        <a-button type="primary" html-type="submit" long :loading="loading">
          {{ $t('register.form.register') }}
        </a-button>

        <!-- Back to Login Button -->
        <a-button
          type="primary"
          long
          class="register-form-login-again-btn"
          @click="goToLogin"
        >
          {{ $t('register.form.backToLogin') }}
        </a-button>
      </a-space>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';
import { useStorage } from '@vueuse/core';
import useLoading from '@/hooks/loading';
import { register } from '@/api/user';
import type { RegisterData } from '@/api/user';

const router = useRouter();
const { t } = useI18n();
const errorMessage = ref('');
const { loading, setLoading } = useLoading();

// Persistent storage for registration preferences
const registerConfig = useStorage('register-config', {
  rememberPassword: true,
  userName: '',
  passWord: '',
});

const registerInfo = reactive({
  userName: registerConfig.value.userName,
  passWord: registerConfig.value.passWord,
  rootPassWord: '',
});

const handleSubmit = async ({
  errors,
  values,
}: {
  errors: Record<string, any> | undefined;
  values: Record<string, any>;
}) => {
  if (loading.value) return;
  if (!errors) {
    setLoading(true);
    try {
      // Call the registration API
      const res = await register(values as RegisterData);
      Message.success(t('register.form.register.success'));
      const { rememberPassword } = registerConfig.value;
      const { userName, passWord } = values;
      registerConfig.value.userName = rememberPassword ? userName : '';
      registerConfig.value.passWord = rememberPassword ? passWord : '';
      router.push({ name: 'login' });
    } catch (err) {
      console.error('Registration error:', err);
      errorMessage.value = (err as Error).message || t('register.form.register.error');
      Message.error(errorMessage.value);
    } finally {
      setLoading(false);
    }
  }
};

const setRememberPassword = (value: boolean) => {
  registerConfig.value.rememberPassword = value;
};

const goToLogin = () => {
  router.push({ name: 'login' }); // Navigate to the Login page with correct lowercase name
};
</script>

<style lang="less" scoped>
.register-form {
  &-wrapper {
    width: 320px;
  }

  &-error-msg {
    height: 32px;
    color: rgb(var(--red-6));
    line-height: 32px;
  }

  &-password-actions {
    display: flex;
    justify-content: space-between;
  }

  &-login-again-btn {
    color: var(--color-text-3) !important;
  }
}
</style>
