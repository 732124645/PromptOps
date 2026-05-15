<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NDialogProvider, darkTheme } from 'naive-ui'
import AppShell from './AppShell.vue'
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <n-message-provider>
      <n-dialog-provider>
        <app-shell />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>
